//go:build linux

// SPDX-License-Identifier: Apache-2.0
package portablesecure

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type DataVolume struct{ Root, FileSystem string }

// Private state requires native Unix ownership and permissions. A mounted
// Windows/FAT/network volume is suitable for imports, not service credentials.
func CheckDataVolume(dir string) (DataVolume, error) {
	p, err := filepath.Abs(dir)
	if err != nil {
		return DataVolume{}, err
	}
	for {
		if _, err = os.Stat(p); err == nil {
			break
		}
		if !os.IsNotExist(err) || filepath.Dir(p) == p {
			return DataVolume{}, err
		}
		p = filepath.Dir(p)
	}
	var stat unix.Statfs_t
	if err = unix.Statfs(p, &stat); err != nil {
		return DataVolume{}, err
	}
	names := map[int64]string{0xef53: "ext4", 0x58465342: "XFS", 0x9123683e: "Btrfs", 0x2fc12fc1: "ZFS", 0xf2f52010: "F2FS", 0x794c7630: "overlayfs"}
	name, ok := names[int64(stat.Type)]
	if !ok {
		return DataVolume{}, fmt.Errorf("private app data needs a local Linux filesystem with Unix permissions (ext4, XFS or Btrfs); choose a folder on a Linux drive, filesystem type %x", stat.Type)
	}
	return DataVolume{Root: p, FileSystem: name}, nil
}

func RejectReparse(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for p := abs; ; p = filepath.Dir(p) {
		st, e := os.Lstat(p)
		if e != nil {
			return e
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("private paths must not contain symbolic links: %s", p)
		}
		if filepath.Dir(p) == p {
			return nil
		}
	}
}

func ValidateOwner(dir string) error {
	if _, err := os.Lstat(filepath.Join(dir, ".owner-sid")); err == nil {
		return errors.New("Windows private data cannot be opened on Linux; import an application export into a new Linux installation")
	} else if !os.IsNotExist(err) {
		return err
	}
	b, err := os.ReadFile(filepath.Join(dir, ".owner-uid"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(b)) != strconv.Itoa(os.Geteuid()) {
		return errors.New("this app data belongs to another Linux user; start Workspace as its installation owner")
	}
	return nil
}

func ProtectDir(dir string) error {
	if _, err := CheckDataVolume(dir); err != nil {
		return err
	}
	if err := RejectReparse(dir); err != nil {
		return err
	}
	if err := ValidateOwner(dir); err != nil {
		return err
	}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		st, e := d.Info()
		if e != nil {
			return e
		}
		if !d.IsDir() && !st.Mode().IsRegular() {
			return fmt.Errorf("private storage contains a link or special file: %s", p)
		}
		raw, ok := st.Sys().(*syscall.Stat_t)
		if !ok || raw.Uid != uint32(os.Geteuid()) || (!d.IsDir() && raw.Nlink != 1) {
			return fmt.Errorf("private storage has another owner or a hard link: %s", p)
		}
		mode := fs.FileMode(0600)
		if d.IsDir() || st.Mode().Perm()&0111 != 0 {
			mode = 0700
		}
		return os.Chmod(p, mode)
	})
	if err != nil {
		return err
	}
	return AtomicWrite(filepath.Join(dir, ".owner-uid"), []byte(strconv.Itoa(os.Geteuid())))
}

const secretMagic = "OWSLINUX1\x00"

func secretCipher(path string, create bool) (cipher.AEAD, error) {
	dir := filepath.Dir(path)
	if create {
		if e := os.MkdirAll(dir, 0700); e != nil {
			return nil, e
		}
	}
	if _, e := CheckDataVolume(dir); e != nil {
		return nil, e
	}
	if e := RejectReparse(dir); e != nil {
		return nil, e
	}
	st, e := os.Stat(dir)
	if e != nil {
		return nil, e
	}
	raw, ok := st.Sys().(*syscall.Stat_t)
	if !ok || raw.Uid != uint32(os.Geteuid()) || st.Mode().Perm()&0077 != 0 {
		return nil, errors.New("secret directory must belong to this Linux user and have mode 0700")
	}
	keyPath := filepath.Join(dir, ".workspace-secret-key")
	if create {
		f, e := os.OpenFile(keyPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e == nil {
			key := make([]byte, 32)
			_, e = rand.Read(key)
			if e == nil {
				_, e = f.Write(key)
			}
			if e == nil {
				e = f.Sync()
			}
			ce := f.Close()
			if e != nil {
				_ = os.Remove(keyPath)
				return nil, e
			}
			if ce != nil {
				return nil, ce
			}
		} else if !os.IsExist(e) {
			return nil, e
		}
	}
	fd, e := unix.Open(keyPath, unix.O_RDONLY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if e != nil {
		return nil, e
	}
	f := os.NewFile(uintptr(fd), keyPath)
	defer f.Close()
	ks, e := f.Stat()
	if e != nil {
		return nil, e
	}
	kr, ok := ks.Sys().(*syscall.Stat_t)
	if !ok || !ks.Mode().IsRegular() || ks.Size() != 32 || ks.Mode().Perm()&0077 != 0 || kr.Uid != uint32(os.Geteuid()) || kr.Nlink != 1 {
		return nil, errors.New("invalid or unprotected Linux secret key")
	}
	key, e := io.ReadAll(f)
	if e != nil {
		return nil, e
	}
	block, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	return cipher.NewGCM(block)
}

func SealSecret(path string, b []byte) ([]byte, error) {
	a, e := secretCipher(path, true)
	if e != nil {
		return nil, e
	}
	nonce := make([]byte, a.NonceSize())
	if _, e = rand.Read(nonce); e != nil {
		return nil, e
	}
	out := append([]byte(secretMagic), nonce...)
	return a.Seal(out, nonce, b, []byte(filepath.Base(path))), nil
}
func OpenSecret(path string, b []byte) ([]byte, error) {
	if !strings.HasPrefix(string(b), secretMagic) {
		return nil, errors.New("secret is not a Linux Workspace envelope; Windows DPAPI state requires reconfiguration")
	}
	a, e := secretCipher(path, false)
	if e != nil {
		return nil, e
	}
	b = b[len(secretMagic):]
	if len(b) < a.NonceSize() {
		return nil, errors.New("truncated secret")
	}
	return a.Open(nil, b[:a.NonceSize()], b[a.NonceSize():], []byte(filepath.Base(path)))
}
func protectSecret(path string, b []byte) ([]byte, error)   { return SealSecret(path, b) }
func unprotectSecret(path string, b []byte) ([]byte, error) { return OpenSecret(path, b) }

// Process groups cover normal shutdown. Keep the spawning OS thread alive so
// Linux's thread-scoped parent-death signal is not triggered by Go's scheduler.
// A systemd user service provides descendant cleanup after a manager crash.
type Job struct {
	mu       sync.Mutex
	groups   []int
	releases []chan struct{}
	closed   bool
}

func NewJob() (*Job, error) { return &Job{}, nil }
func (j *Job) Start(cmd *exec.Cmd) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return errors.New("process group is closed")
	}
	result := make(chan error, 1)
	release := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGKILL}
		e := cmd.Start()
		result <- e
		if e == nil {
			<-release
		}
	}()
	if e := <-result; e != nil {
		return e
	}
	j.groups = append(j.groups, cmd.Process.Pid)
	j.releases = append(j.releases, release)
	return nil
}
func (j *Job) Stop(ctx context.Context) error {
	if j == nil {
		return nil
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	for _, id := range j.groups {
		if e := unix.Kill(-id, unix.SIGKILL); e != nil && e != unix.ESRCH {
			return e
		}
	}
	// The caller reaps its direct children with Cmd.Wait. Signal delivery to
	// the entire group happens before returning; the supervisor waits for children.
	return ctx.Err()
}
func (j *Job) Close() {
	if j == nil {
		return
	}
	_ = j.Stop(context.Background())
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return
	}
	j.closed = true
	for _, release := range j.releases {
		close(release)
	}
}

func Lock(path string) (*os.File, error) {
	fd, e := unix.Open(path, unix.O_CREAT|unix.O_RDWR|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0600)
	if e != nil {
		return nil, e
	}
	f := os.NewFile(uintptr(fd), path)
	if e = unix.Flock(fd, unix.LOCK_EX|unix.LOCK_NB); e != nil {
		f.Close()
		return nil, errors.New("another Workspace process owns this installation lock")
	}
	return f, nil
}
func LockedBinary(path, expected string) (*os.File, error) {
	fd, e := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0)
	if e != nil {
		return nil, e
	}
	f := os.NewFile(uintptr(fd), path)
	fail := func(e error) (*os.File, error) { f.Close(); return nil, e }
	st, e := f.Stat()
	if e != nil {
		return fail(e)
	}
	if !st.Mode().IsRegular() || st.Mode().Perm()&0022 != 0 || st.Mode().Perm()&0111 == 0 {
		return fail(errors.New("application must be executable and not writable by other users"))
	}
	if e = unix.Flock(fd, unix.LOCK_SH|unix.LOCK_NB); e != nil {
		return fail(e)
	}
	h := sha256.New()
	if _, e = io.Copy(h, f); e != nil {
		return fail(e)
	}
	if hex.EncodeToString(h.Sum(nil)) != expected {
		return fail(errors.New("application checksum does not match protected catalog"))
	}
	_, e = f.Seek(0, 0)
	if e != nil {
		return fail(e)
	}
	return f, nil
}
func OpenBrowser(url string) {
	cmd := exec.Command("xdg-open", url)
	if cmd.Start() == nil {
		go cmd.Wait()
	}
}
func StartDetached(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	return cmd.Start()
}

// Each installation has its own user service; changing one never disables another.
func AutostartCommand(name, exe string, args []string, enabled bool) error {
	root, e := os.UserConfigDir()
	if e != nil {
		return e
	}
	dir := filepath.Join(root, "systemd", "user")
	if e = os.MkdirAll(dir, 0700); e != nil {
		return e
	}
	if strings.ContainsAny(name, "/\\\r\n") {
		return errors.New("invalid service name")
	}
	unit := name + ".service"
	path := filepath.Join(dir, unit)
	run := func(args ...string) error {
		out, e := exec.Command("systemctl", append([]string{"--user"}, args...)...).CombinedOutput()
		if e != nil {
			return fmt.Errorf("user service manager: %w: %s", e, strings.TrimSpace(string(out)))
		}
		return nil
	}
	if !enabled {
		if _, e = os.Stat(path); os.IsNotExist(e) {
			return nil
		}
		if e = run("disable", unit); e != nil {
			return e
		}
		if e = os.Remove(path); e != nil {
			return e
		}
		return run("daemon-reload")
	}
	quote := func(s string) string {
		return "\"" + strings.NewReplacer("\\", "\\\\", "\"", "\\\"", "%", "%%", "$", "$$").Replace(s) + "\""
	}
	command := quote(exe)
	for _, arg := range args {
		if strings.ContainsAny(arg, "\x00\r\n") {
			return errors.New("invalid service argument")
		}
		command += " " + quote(arg)
	}
	if strings.ContainsAny(exe, "\x00\r\n") {
		return errors.New("invalid service executable")
	}
	body := "[Unit]\nDescription=One Workspace\nAfter=network-online.target\n\n[Service]\nType=simple\nExecStart=" + command + "\nRestart=on-failure\nRestartSec=5\nKillMode=control-group\nTimeoutStopSec=30\nUMask=0077\n\n[Install]\nWantedBy=default.target\n"
	if e = AtomicWrite(path, []byte(body)); e != nil {
		return e
	}
	if e = run("daemon-reload"); e != nil {
		return e
	}
	return run("enable", unit)
}
func Autostart(exe, data string, on bool) error {
	sum := sha256.Sum256([]byte(data))
	return AutostartCommand("one-workspace-"+hex.EncodeToString(sum[:8]), exe, []string{"--data", data, "run"}, on)
}

func WaitLock(ctx context.Context, path string) (func(), error) {
	for {
		f, e := Lock(path)
		if e == nil {
			return func() { f.Close() }, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("wait for installation lock: %w", ctx.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}
