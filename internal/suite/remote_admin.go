package suite

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// A narrowly scoped, bounded owner continuation lets an existing update page
// reconnect after the manager loses its in-memory sign-in sessions on restart.
// Only hashes are stored, protected by Windows encryption and private ACLs.
type adminContinuation struct {
	Hash, Owner string
	Expires     time.Time
}

const adminCookie = "workspace_update_owner"

func (p *Portal) remoteUpdateOwner(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != "GET" && r.Header.Get("Origin") != p.rt.Config.Origin {
		http.Error(w, "Origin rejected", 403)
		return false
	}
	p.adminMu.Lock()
	defer p.adminMu.Unlock()
	path := filepath.Join(p.rt.Config.Data, "maintenance", "remote-owner.dpapi")
	var stored []adminContinuation
	_ = platform.LoadSecret(path, &stored)
	now := time.Now()
	active := make([]adminContinuation, 0, len(stored))
	valid := false
	c, _ := r.Cookie(adminCookie)
	hash := ""
	if c != nil && len(c.Value) == 64 {
		sum := sha256.Sum256([]byte(c.Value))
		hash = hex.EncodeToString(sum[:])
	}
	for _, s := range stored {
		if s.Owner == p.rt.Config.Owner && now.Before(s.Expires) {
			active = append(active, s)
			if s.Hash == hash {
				valid = true
			}
		}
	}
	user, signed := p.user(r)
	if signed && user.Subject != p.rt.Config.Owner {
		http.Error(w, "Workspace owner access required", 403)
		return false
	}
	if signed && !valid {
		root := filepath.Dir(path)
		if e := os.MkdirAll(root, 0700); e != nil {
			http.Error(w, "Could not protect remote administration", 500)
			return false
		}
		if e := platform.ProtectMaintenanceDir(root, p.rt.Config.Data); e != nil {
			http.Error(w, "Could not protect remote administration", 500)
			return false
		}
		token := platform.Token()
		sum := sha256.Sum256([]byte(token))
		if len(active) >= 8 {
			active = active[len(active)-7:]
		}
		active = append(active, adminContinuation{Hash: hex.EncodeToString(sum[:]), Owner: user.Subject, Expires: now.Add(2 * time.Hour)})
		if e := platform.SaveSecret(path, active); e != nil {
			http.Error(w, "Could not save update continuation", 500)
			return false
		}
		http.SetCookie(w, &http.Cookie{Name: adminCookie, Value: token, Path: "/updates", HttpOnly: true, Secure: p.secure, SameSite: http.SameSiteStrictMode, MaxAge: 7200})
		valid = true
	}
	if !valid {
		if r.Method == "GET" && (r.URL.Path == "/updates" || r.URL.Path == "/updates/") {
			http.Redirect(w, r, "/auth/login?returnTo=%2Fupdates", 303)
		} else {
			http.Error(w, "Sign in as the Workspace owner to continue", 401)
		}
		return false
	}
	return true
}

func (p *Portal) revokeRemoteAdmin(w http.ResponseWriter, r *http.Request) bool {
	// Logout revokes all update continuations for this installation. A cookie
	// scoped to /updates is deliberately not sent to /auth/logout.
	p.adminMu.Lock()
	defer p.adminMu.Unlock()
	if e := os.Remove(filepath.Join(p.rt.Config.Data, "maintenance", "remote-owner.dpapi")); e != nil && !os.IsNotExist(e) {
		http.Error(w, "Could not revoke update access. Retry sign-out.", 500)
		return false
	}
	http.SetCookie(w, &http.Cookie{Name: adminCookie, Value: "", Path: "/updates", MaxAge: -1, HttpOnly: true, Secure: p.secure, SameSite: http.SameSiteStrictMode})
	return true
}
