package suite

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// Cloudflared's local readiness endpoint reports a registered edge connection,
// including when Workspace's origin is temporarily stopped.
func cloudflareMetricsAddress() (string, error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}

func cloudflareCommand(c Config, s Secrets, executable, config, metrics string) *exec.Cmd {
	cmd := exec.Command(executable, "tunnel", "--config", config, "--no-autoupdate", "--metrics", metrics, "run")
	cmd.Env = append(os.Environ(), "TUNNEL_TOKEN="+s.TunnelToken, "TEMP="+runtimeTemp(c, "cloudflare"), "TMP="+runtimeTemp(c, "cloudflare"))
	cmd.Dir = c.Data
	return cmd
}

func waitCloudflareReady(ctx context.Context, metrics string, exited <-chan error) error {
	return WaitHTTPDetailed(ctx, "http://"+metrics+"/ready", 90*time.Second, exited, nil)
}

// The updater owns this verified replica in a different Windows job. The
// runtime can stop its own connector and replace its installed executable.
func keepUpdateTunnel(ctx context.Context, c Config, s Secrets, report func(string)) (stop func(), result error) {
	if !c.Tunnel {
		return func() {}, nil
	}
	if s.TunnelToken == "" {
		return nil, errors.New("the installed Cloudflare connector has no credentials")
	}
	if err := prepareRuntimeTemp(c, "cloudflare"); err != nil {
		return nil, err
	}
	connector, ok := c.Catalog.Packages["cloudflared"]
	if !ok {
		return nil, errors.New("the installed Cloudflare package is missing")
	}
	source, err := platform.Inside(c.Programs, connector.File)
	if err != nil {
		return nil, err
	}
	heldSource, err := platform.LockedBinary(source, connector.SHA256)
	if err != nil {
		return nil, err
	}
	defer heldSource.Close()

	root := filepath.Join(c.Data, "updates")
	if err = cleanupPath(root); err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp(root, "connector-")
	if err != nil {
		return nil, err
	}
	executable := filepath.Join(dir, "cloudflared.exe")
	config := filepath.Join(dir, "cloudflared.yml")
	var job *platform.Job
	var finished chan struct{}
	closeReplica := func() {
		job.Close()
		if finished != nil {
			<-finished
		}
		_ = os.Remove(executable)
		_ = os.Remove(config)
		_ = os.Remove(dir)
	}
	defer func() {
		if result != nil {
			closeReplica()
		}
	}()
	if err = platform.ProtectMaintenanceDir(dir, c.Data); err != nil {
		return nil, err
	}
	if err = platform.CopyFile(source, executable); err != nil {
		return nil, err
	}
	held, err := platform.LockedBinary(executable, connector.SHA256)
	if err != nil {
		return nil, err
	}
	if err = platform.AtomicWrite(config, []byte("loglevel: error\n")); err != nil {
		held.Close()
		return nil, err
	}
	metrics, err := cloudflareMetricsAddress()
	if err != nil {
		held.Close()
		return nil, err
	}
	job, err = platform.NewJob()
	if err != nil {
		held.Close()
		return nil, err
	}
	cmd := cloudflareCommand(c, s, executable, config, metrics)
	log := &platform.Log{Path: filepath.Join(root, "connector.log")}
	cmd.Stdout, cmd.Stderr = log, log
	if err = job.Start(cmd); err != nil {
		held.Close()
		return nil, err
	}
	exited := make(chan error, 1)
	finished = make(chan struct{})
	go func() {
		defer close(finished)
		defer held.Close()
		exited <- cmd.Wait()
	}()
	report("Waiting for a second Cloudflare connector before stopping Workspace")
	if err = waitCloudflareReady(ctx, metrics, exited); err != nil {
		return nil, fmt.Errorf("the backup Cloudflare connector did not connect: %w", err)
	}
	report("The backup Cloudflare connector is connected")
	return closeReplica, nil
}
