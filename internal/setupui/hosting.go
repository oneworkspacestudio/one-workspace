package setupui

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fahimc/one-workspace/internal/platform"
	"github.com/fahimc/one-workspace/internal/suite"
)

func (s *Server) connectHosting(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Domain string `json:"domain"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&in) != nil {
		http.Error(w, "Enter a domain", 400)
		return
	}
	domain, e := suite.NormalizeDomain(in.Domain)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	s.mu.Lock()
	if s.connecting || s.busy || s.ready {
		s.mu.Unlock()
		http.Error(w, "Setup is already running", 409)
		return
	}
	s.connecting = true
	s.hosting = nil
	s.hostingError = ""
	s.hostingPhase = "Connecting Cloudflare. Complete the authorization in your browser if prompted."
	s.mu.Unlock()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
		defer cancel()
		path, e := suite.CloudflareLoginPath()
		if e == nil {
			if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
				e = s.providerLogin(ctx)
			} else if statErr != nil {
				e = errors.New("cannot read the saved Cloudflare connection")
			}
		}
		var hosting suite.Request
		if e == nil {
			hosting, e = suite.ResolveHosting(ctx, domain, path)
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		s.connecting = false
		if e != nil {
			s.hostingError = e.Error()
			s.hostingPhase = "Connection needs attention"
			return
		}
		s.hosting = &hosting
		s.hostingPhase = "Connected. Addresses and secure hosting will be created during installation."
	}()
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) providerLogin(ctx context.Context) error {
	p, ok := s.cat.Packages["cloudflared"]
	if !ok {
		return errors.New("this installer does not include the Cloudflare connector")
	}
	folder, e := os.MkdirTemp("", "workspace-connect-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(folder)
	if e = platform.ProtectDir(folder); e != nil {
		return e
	}
	executable := filepath.Join(folder, "cloudflared.exe")
	found := false
	for _, entry := range s.z.File {
		if entry.Name != p.File {
			continue
		}
		reader, e := entry.Open()
		if e != nil {
			return e
		}
		file, e := os.OpenFile(executable, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0700)
		if e != nil {
			reader.Close()
			return e
		}
		n, copyErr := io.Copy(file, io.LimitReader(reader, p.Size+1))
		closeErr := file.Close()
		reader.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		if n != p.Size {
			return errors.New("connector package size mismatch")
		}
		found = true
		break
	}
	if !found {
		return errors.New("Cloudflare connector is missing from the bundle")
	}
	held, e := platform.LockedBinary(executable, p.SHA256)
	if e != nil {
		return e
	}
	defer held.Close()
	job, e := platform.NewJob()
	if e != nil {
		return e
	}
	defer job.Close()
	// The official login flow opens Cloudflare's browser authorization. Its
	// credential is written to the standard location only when absent.
	cmd := exec.CommandContext(ctx, executable, "tunnel", "login")
	cmd.Dir = folder
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if e = job.Start(cmd); e != nil {
		return errors.New("could not start Cloudflare authorization")
	}
	if e = cmd.Wait(); e != nil {
		return errors.New("Cloudflare authorization was cancelled or timed out; click Connect Cloudflare to retry")
	}
	return nil
}
