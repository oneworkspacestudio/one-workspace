package apphost

import (
	"archive/zip"
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"golang.org/x/crypto/argon2"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const backupMagic = "OWBACKUP1\x00"
const maxBackupBytes = int64(2 << 30)

type backupManifest struct {
	Schema  int       `json:"schema"`
	App     string    `json:"app"`
	Version string    `json:"version"`
	Created time.Time `json:"created"`
}

func archiveKey(password string, salt []byte) (cipher.AEAD, error) {
	if len(password) < 8 || len(password) > 128 {
		return nil, errors.New("use a backup password of 8–128 bytes")
	}
	key := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 32)
	block, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	return cipher.NewGCM(block)
}
func encryptArchive(dst io.Writer, src io.Reader, password string) error {
	salt := make([]byte, 16)
	if _, e := rand.Read(salt); e != nil {
		return e
	}
	a, e := archiveKey(password, salt)
	if e != nil {
		return e
	}
	if _, e = io.WriteString(dst, backupMagic); e != nil {
		return e
	}
	if _, e = dst.Write(salt); e != nil {
		return e
	}
	buf := make([]byte, 1<<20)
	var seq uint64
	for {
		n, readErr := io.ReadFull(src, buf)
		if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
			return readErr
		}
		nonce := make([]byte, a.NonceSize())
		if _, e = rand.Read(nonce); e != nil {
			return e
		}
		aad := make([]byte, 9)
		binary.BigEndian.PutUint64(aad, seq)
		if readErr != nil {
			aad[8] = 1
		}
		sealed := a.Seal(nil, nonce, buf[:n], aad)
		if e = binary.Write(dst, binary.BigEndian, uint32(len(sealed))); e != nil {
			return e
		}
		if _, e = dst.Write(aad); e != nil {
			return e
		}
		if _, e = dst.Write(nonce); e != nil {
			return e
		}
		if _, e = dst.Write(sealed); e != nil {
			return e
		}
		seq++
		if readErr != nil {
			return nil
		}
	}
}
func decryptArchive(dst io.Writer, src io.Reader, password string) error {
	r := bufio.NewReader(src)
	magic := make([]byte, len(backupMagic))
	if _, e := io.ReadFull(r, magic); e != nil || string(magic) != backupMagic {
		return errors.New("not a supported Workspace backup")
	}
	salt := make([]byte, 16)
	if _, e := io.ReadFull(r, salt); e != nil {
		return e
	}
	a, e := archiveKey(password, salt)
	if e != nil {
		return e
	}
	var seq uint64
	var total int64
	for {
		var size uint32
		if e = binary.Read(r, binary.BigEndian, &size); e != nil {
			return errors.New("backup is incomplete")
		}
		if size > 1<<20+uint32(a.Overhead()) || size < uint32(a.Overhead()) {
			return errors.New("invalid backup frame")
		}
		aad := make([]byte, 9)
		nonce := make([]byte, a.NonceSize())
		sealed := make([]byte, size)
		if _, e = io.ReadFull(r, aad); e != nil {
			return e
		}
		if binary.BigEndian.Uint64(aad) != seq || aad[8] > 1 {
			return errors.New("backup frames are out of order")
		}
		if _, e = io.ReadFull(r, nonce); e != nil {
			return e
		}
		if _, e = io.ReadFull(r, sealed); e != nil {
			return e
		}
		plain, e := a.Open(nil, nonce, sealed, aad)
		if e != nil {
			return errors.New("backup password or integrity check failed")
		}
		total += int64(len(plain))
		if total > maxBackupBytes {
			return errors.New("backup exceeds the restore size limit")
		}
		if _, e = dst.Write(plain); e != nil {
			return e
		}
		seq++
		if aad[8] == 1 {
			if _, e = r.ReadByte(); e != io.EOF {
				return errors.New("unexpected trailing backup data")
			}
			return nil
		}
	}
}

// Snapshot creates an online-consistent SQLite snapshot plus immutable objects.
// All application mutations and object finalization must hold Freeze.RLock.
func (h *Host) Snapshot(dst io.Writer, password string) error {
	if len(password) < 8 || len(password) > 128 {
		return errors.New("use a backup password of 8–128 bytes")
	}
	h.Freeze.Lock()
	defer h.Freeze.Unlock()
	tmp, e := os.MkdirTemp(h.Config.Data, ".backup-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(tmp)
	dbpath := filepath.Join(tmp, "state.sqlite")
	if _, e = h.DB.Exec("VACUUM main INTO ?", dbpath); e != nil {
		return e
	}
	archive, e := os.Create(filepath.Join(tmp, "snapshot.zip"))
	if e != nil {
		return e
	}
	defer archive.Close()
	z := zip.NewWriter(archive)
	var total int64
	count := 0
	add := func(path, name string) error {
		if e := platform.RejectReparse(path); e != nil {
			return e
		}
		st, e := os.Stat(path)
		if e != nil {
			return e
		}
		if !st.Mode().IsRegular() {
			return errors.New("backup only accepts ordinary files")
		}
		total += st.Size()
		count++
		if total > maxBackupBytes-(10<<20) || count > 50000 {
			return errors.New("app backup exceeds its configured size or file-count limit")
		}
		entry, e := z.Create(name)
		if e != nil {
			return e
		}
		f, e := os.Open(path)
		if e != nil {
			return e
		}
		defer f.Close()
		if strings.HasPrefix(name, "objects/") {
			sum := sha256.New()
			_, e = io.Copy(io.MultiWriter(entry, sum), f)
			if e == nil && hex.EncodeToString(sum.Sum(nil)) != strings.TrimPrefix(name, "objects/") {
				return errors.New("a stored object failed checksum verification")
			}
		} else {
			_, e = io.Copy(entry, f)
		}
		return e
	}
	if e = add(dbpath, "state.sqlite"); e != nil {
		return e
	}
	for _, name := range []string{"objects", "secret.key"} {
		path := filepath.Join(h.Config.Data, name)
		if _, e = os.Stat(path); os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return e
		}
		if e = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if err = platform.RejectReparse(p); err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(h.Config.Data, p)
			if err != nil {
				return err
			}
			return add(p, filepath.ToSlash(rel))
		}); e != nil {
			return e
		}
	}
	manifest, e := z.Create("manifest.json")
	if e != nil {
		return e
	}
	if e = json.NewEncoder(manifest).Encode(backupManifest{1, h.Config.ID, h.Config.Version, time.Now().UTC()}); e != nil {
		return e
	}
	if e = z.Close(); e != nil {
		return e
	}
	if _, e = archive.Seek(0, 0); e != nil {
		return e
	}
	return encryptArchive(dst, archive, password)
}
func Restore(root, app, archive, password string) (result error) {
	if _, e := os.Lstat(root); !os.IsNotExist(e) {
		return errors.New("restore requires a new, empty destination path; existing data is preserved")
	}
	parent := filepath.Dir(root)
	if _, e := platform.CheckDataVolume(parent); e != nil {
		return e
	}
	for p := parent; ; p = filepath.Dir(p) {
		if _, e := os.Stat(p); e == nil {
			if e = platform.RejectReparse(p); e != nil {
				return e
			}
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	if e := os.MkdirAll(parent, 0700); e != nil {
		return e
	}
	stage, e := os.MkdirTemp(parent, ".restore-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(stage)
	if e = platform.ProtectDir(stage); e != nil {
		return e
	}
	in, e := os.Open(archive)
	if e != nil {
		return e
	}
	defer in.Close()
	raw, e := os.Create(filepath.Join(stage, "snapshot.zip"))
	if e != nil {
		return e
	}
	if e = decryptArchive(raw, in, password); e != nil {
		raw.Close()
		return e
	}
	raw.Close()
	z, e := zip.OpenReader(filepath.Join(stage, "snapshot.zip"))
	if e != nil {
		return e
	}
	defer z.Close()
	var manifest backupManifest
	seen := map[string]bool{}
	var total uint64
	if len(z.File) > 50000 {
		return errors.New("too many backup files")
	}
	outdir := filepath.Join(stage, "data")
	if e = os.Mkdir(outdir, 0700); e != nil {
		return e
	}
	for _, f := range z.File {
		key := strings.ToLower(f.Name)
		if seen[key] || f.Mode()&os.ModeSymlink != 0 || f.FileInfo().IsDir() {
			return errors.New("invalid duplicate or linked backup entry")
		}
		seen[key] = true
		total += f.UncompressedSize64
		if total > uint64(maxBackupBytes) || f.UncompressedSize64 > 512<<20 {
			return errors.New("backup decompression limit exceeded")
		}
		if f.Name != "manifest.json" && f.Name != "state.sqlite" && f.Name != "secret.key" && !strings.HasPrefix(f.Name, "objects/") {
			return errors.New("unknown backup entry")
		}
		if strings.HasPrefix(f.Name, "objects/") {
			hash := strings.TrimPrefix(f.Name, "objects/")
			if len(hash) != 64 || strings.ToLower(hash) != hash {
				return errors.New("invalid backup object name")
			}
			if _, e := hex.DecodeString(hash); e != nil {
				return errors.New("invalid backup object name")
			}
		}
		path, e := platform.Inside(outdir, f.Name)
		if e != nil {
			return e
		}
		r, e := f.Open()
		if e != nil {
			return e
		}
		if f.Name == "manifest.json" {
			e = json.NewDecoder(io.LimitReader(r, 8192)).Decode(&manifest)
			r.Close()
			if e != nil {
				return e
			}
			continue
		}
		if e = os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			r.Close()
			return e
		}
		dest, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			r.Close()
			return e
		}
		sum := sha256.New()
		n, e := io.Copy(io.MultiWriter(dest, sum), io.LimitReader(r, int64(f.UncompressedSize64)+1))
		dest.Close()
		r.Close()
		if e != nil || n != int64(f.UncompressedSize64) {
			return errors.New("backup entry failed integrity verification")
		}
		if strings.HasPrefix(f.Name, "objects/") && hex.EncodeToString(sum.Sum(nil)) != strings.TrimPrefix(f.Name, "objects/") {
			return errors.New("backup object checksum mismatch")
		}
	}
	if manifest.Schema != 1 || manifest.App != app || !seen["state.sqlite"] {
		return errors.New("backup does not belong to this application")
	}
	if e = platform.SaveJSON(filepath.Join(outdir, "app.json"), map[string]any{"schema": 1, "id": app}); e != nil {
		return e
	}
	if e = platform.ProtectDir(outdir); e != nil {
		return e
	}
	db, e := sql.Open("sqlite", filepath.Join(outdir, "state.sqlite"))
	if e != nil {
		return e
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var integrity string
	e = db.QueryRowContext(ctx, "PRAGMA quick_check").Scan(&integrity)
	if e == nil && integrity != "ok" {
		e = errors.New("restored database failed integrity validation")
	}
	if e == nil {
		_, e = db.ExecContext(ctx, "PRAGMA trusted_schema=OFF; DELETE FROM host_sessions")
	}
	if e == nil {
		_, e = db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	}
	closeErr := db.Close()
	if e != nil {
		return e
	}
	if closeErr != nil {
		return closeErr
	}
	if e = os.Rename(outdir, root); e != nil {
		return e
	}
	return nil
}
func (h *Host) BackupRoute() {
	h.Mux.HandleFunc("POST /api/backups", func(w http.ResponseWriter, r *http.Request) {
		s, ok := h.session(r)
		if !ok {
			Problem(w, 401, "Sign in first")
			return
		}
		if s.User.Role != "owner" || !h.sameOrigin(r) || !Equal(r.Header.Get("X-CSRF-Token"), s.CSRF) {
			Problem(w, 403, "Only the app owner can export encrypted backups")
			return
		}
		var in struct {
			Password string `json:"password"`
		}
		if !Decode(w, r, &in) {
			return
		}
		tmp, e := os.CreateTemp(h.Config.Data, ".download-backup-")
		if e != nil {
			h.DBError(w, e)
			return
		}
		defer os.Remove(tmp.Name())
		defer tmp.Close()
		if e = h.Snapshot(tmp, in.Password); e != nil {
			Problem(w, 400, e.Error())
			return
		}
		if _, e = tmp.Seek(0, 0); e != nil {
			Problem(w, 500, "Could not open backup")
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+h.Config.ID+".owbackup")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = io.Copy(w, tmp)
	})
}
