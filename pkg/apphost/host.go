// Package apphost is One Workspace's native application identity and storage contract.
package apphost

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"golang.org/x/oauth2"
	_ "modernc.org/sqlite"
)

type Config struct{ ID, Name, Version, Data, Listen, Origin, Issuer, ClientID, ClientSecret, Owner, WorkspaceURL, RegistryURL, RegistryToken string }
type User struct {
	ID      string `json:"id"`
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
}
type Session struct {
	User    User      `json:"user"`
	CSRF    string    `json:"csrf"`
	Expires time.Time `json:"expires"`
}
type flow struct {
	Verifier, Nonce, Return string
	Expires                 time.Time
}
type loginAttempt struct {
	Count int
	Since time.Time
}
type Host struct {
	Config   Config
	DB       *sql.DB
	Mux      *http.ServeMux
	Base     string
	Freeze   sync.RWMutex
	Log      *slog.Logger
	lock     *os.File
	logFile  *os.File
	ui       fs.FS
	authMu   sync.Mutex
	flows    map[string]flow
	attempts map[string]loginAttempt
	oauth    oauth2.Config
	verifier *oidc.IDTokenVerifier
	setup    string
}
type contextKey struct{}

func Current(r *http.Request) Session { s, _ := r.Context().Value(contextKey{}).(Session); return s }
func ID() string {
	b := make([]byte, 24)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func Hash(s string) string                   { b := sha256.Sum256([]byte(s)); return hex.EncodeToString(b[:]) }
func Identity(issuer, subject string) string { return Hash(issuer + "\x00" + subject) }
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func Problem(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}
func Decode(w http.ResponseWriter, r *http.Request, v any) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	d.DisallowUnknownFields()
	if d.Decode(v) != nil {
		Problem(w, 400, "Check the submitted fields and try again")
		return false
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		Problem(w, 400, "Only one JSON object is accepted")
		return false
	}
	return true
}
func FromEnv(c Config) Config {
	c.Issuer = os.Getenv("WORKSPACE_OIDC_ISSUER")
	c.ClientID = os.Getenv("WORKSPACE_OIDC_CLIENT_ID")
	c.ClientSecret = os.Getenv("WORKSPACE_OIDC_CLIENT_SECRET")
	c.Owner = os.Getenv("WORKSPACE_OWNER_SUBJECT")
	c.WorkspaceURL = os.Getenv("WORKSPACE_HOME_URL")
	c.RegistryURL = os.Getenv("WORKSPACE_REGISTRY_URL")
	c.RegistryToken = os.Getenv("WORKSPACE_REGISTRY_TOKEN")
	return c
}
func (c Config) Validate() error {
	if !regexp.MustCompile(`^[a-z][a-z0-9-]{0,39}$`).MatchString(c.ID) {
		return errors.New("invalid application ID")
	}
	if c.ID == "" || c.Name == "" || !filepath.IsAbs(c.Data) || filepath.Dir(filepath.Clean(c.Data)) == filepath.Clean(c.Data) {
		return errors.New("choose a dedicated absolute data folder")
	}
	host, port, e := net.SplitHostPort(c.Listen)
	if e != nil || port == "" || !(host == "127.0.0.1" || host == "::1" || host == "localhost") {
		return errors.New("listen must use a loopback address and port; use Workspace for public access")
	}
	u, e := url.Parse(c.Origin)
	if e != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || strings.Contains(u.Path, "..") || strings.ContainsAny(c.Origin, "\r\n\\") {
		return errors.New("origin must be the actual public app URL, without credentials, query or fragment")
	}
	ip := net.ParseIP(u.Hostname())
	if u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "localhost" || ip != nil && ip.IsLoopback())) {
		return errors.New("public application addresses require HTTPS")
	}
	if c.Issuer != "" {
		issuer, e := url.Parse(c.Issuer)
		if e != nil || issuer.Host == "" || issuer.User != nil || issuer.RawQuery != "" || issuer.Fragment != "" {
			return errors.New("invalid MyAccounts issuer")
		}
		ip := net.ParseIP(issuer.Hostname())
		if issuer.Scheme != "https" && !(issuer.Scheme == "http" && (issuer.Hostname() == "localhost" || ip != nil && ip.IsLoopback())) {
			return errors.New("MyAccounts requires HTTPS except on loopback")
		}
		if c.ClientID == "" || c.ClientSecret == "" {
			return errors.New("configure the MyAccounts client ID and server-side secret")
		}
	}
	return nil
}
func PrepareData(root, id string) error {
	if _, e := platform.CheckDataVolume(root); e != nil {
		return e
	}
	for p := filepath.Clean(root); ; p = filepath.Dir(p) {
		if _, e := os.Lstat(p); e == nil {
			if e = platform.RejectReparse(p); e != nil {
				return e
			}
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	if e := os.MkdirAll(root, 0700); e != nil {
		return e
	}
	marker := filepath.Join(root, "app.json")
	var saved struct {
		Schema int    `json:"schema"`
		ID     string `json:"id"`
	}
	if e := platform.LoadJSON(marker, &saved); e == nil {
		if saved.Schema != 1 || saved.ID != id {
			return errors.New("this data folder belongs to a different application")
		}
	} else if !os.IsNotExist(e) {
		return e
	} else {
		entries, e := os.ReadDir(root)
		if e != nil {
			return e
		}
		if len(entries) > 0 {
			return errors.New("choose an empty app data folder or an existing folder for this app")
		}
	}
	if e := validateOwner(root); e != nil {
		return e
	}
	if e := filepath.WalkDir(root, func(path string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		return platform.RejectReparse(path)
	}); e != nil {
		return e
	}
	if e := platform.ProtectDir(root); e != nil {
		return e
	}
	return platform.SaveJSON(marker, map[string]any{"schema": 1, "id": id})
}
func New(c Config, ui fs.FS) (h *Host, result error) {
	c.Origin = strings.TrimRight(c.Origin, "/")
	if e := c.Validate(); e != nil {
		return nil, e
	}
	if e := PrepareData(c.Data, c.ID); e != nil {
		return nil, e
	}
	lock, e := platform.Lock(filepath.Join(c.Data, "app.lock"))
	if e != nil {
		return nil, errors.New("this app is already running for the selected data folder")
	}
	db, e := sql.Open("sqlite", filepath.Join(c.Data, "state.sqlite"))
	if e != nil {
		lock.Close()
		return nil, e
	}
	db.SetMaxOpenConns(1)
	defer func() {
		if result != nil {
			db.Close()
			lock.Close()
		}
	}()
	if _, e = db.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000;
 CREATE TABLE IF NOT EXISTS host_users(id TEXT PRIMARY KEY,issuer TEXT NOT NULL,subject TEXT NOT NULL,name TEXT NOT NULL,email TEXT NOT NULL DEFAULT '',role TEXT NOT NULL DEFAULT 'member',password BLOB,UNIQUE(issuer,subject));
 CREATE TABLE IF NOT EXISTS host_sessions(id TEXT PRIMARY KEY,user_id TEXT NOT NULL REFERENCES host_users(id) ON DELETE CASCADE,csrf TEXT NOT NULL,expires TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS host_migrations(version INTEGER PRIMARY KEY,applied TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS host_audit(id INTEGER PRIMARY KEY,at TEXT NOT NULL,user_id TEXT NOT NULL,action TEXT NOT NULL,object TEXT NOT NULL);`); e != nil {
		return nil, e
	}
	u, _ := url.Parse(c.Origin)
	logFile, e := os.OpenFile(filepath.Join(c.Data, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if e != nil {
		return nil, e
	}
	h = &Host{Config: c, DB: db, Mux: http.NewServeMux(), Base: strings.TrimRight(u.Path, "/"), lock: lock, logFile: logFile, ui: ui, flows: map[string]flow{}, attempts: map[string]loginAttempt{}, Log: slog.New(slog.NewJSONHandler(logFile, nil))}
	var n int
	if e = db.QueryRow("SELECT count(*) FROM host_users").Scan(&n); e != nil {
		logFile.Close()
		return nil, e
	}
	if n == 0 && c.Issuer == "" {
		h.setup = ID() + ID()
	}
	h.routes()
	return h, nil
}
func (h *Host) Close() error {
	_, _ = h.DB.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	e := h.DB.Close()
	if h.logFile != nil {
		h.logFile.Close()
	}
	h.lock.Close()
	return e
}
func (h *Host) Migrate(version int, schema string) error {
	tx, e := h.DB.Begin()
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var current int
	if e = tx.QueryRow("SELECT coalesce(max(version),0) FROM host_migrations").Scan(&current); e != nil {
		return e
	}
	if current > version {
		return errors.New("app data needs a newer executable; restore a compatible backup to downgrade")
	}
	if current == version {
		return nil
	}
	if current+1 != version {
		return errors.New("missing application migration")
	}
	if _, e = tx.Exec(schema); e != nil {
		return e
	}
	if _, e = tx.Exec("INSERT INTO host_migrations VALUES(?,?)", version, time.Now().UTC().Format(time.RFC3339Nano)); e != nil {
		return e
	}
	return tx.Commit()
}
func (h *Host) Audit(u, action, object string) {
	_, _ = h.DB.Exec("INSERT INTO host_audit(at,user_id,action,object)VALUES(?,?,?,?)", time.Now().UTC().Format(time.RFC3339Nano), u, action, object)
}
func (h *Host) Require(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, ok := h.session(r)
		if !ok {
			Problem(w, 401, "Your app session needs reconnecting. Continue with MyAccounts.")
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			if !h.sameOrigin(r) || !Equal(r.Header.Get("X-CSRF-Token"), s.CSRF) {
				Problem(w, 403, "Refresh the app before submitting this change")
				return
			}
		}
		h.Freeze.RLock()
		defer h.Freeze.RUnlock()
		fn(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, s)))
	}
}
func (h *Host) Owner(fn http.HandlerFunc) http.HandlerFunc {
	return h.Require(func(w http.ResponseWriter, r *http.Request) {
		if Current(r).User.Role != "owner" {
			Problem(w, 403, "Only the Workspace owner can manage infrastructure settings")
			return
		}
		fn(w, r)
	})
}
func (h *Host) routes() {
	h.Mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		var n int
		if e := h.DB.QueryRow("SELECT 1").Scan(&n); e != nil {
			Problem(w, 503, "Database is unavailable")
			return
		}
		JSON(w, 200, map[string]any{"status": "ok", "app": h.Config.ID, "version": h.Config.Version})
	})
	h.Mux.HandleFunc("GET /auth/login", h.login)
	h.Mux.HandleFunc("GET /auth/callback", h.callback)
	h.Mux.HandleFunc("POST /auth/logout", h.logout)
	h.Mux.HandleFunc("POST /api/login", h.localLogin)
	h.Mux.HandleFunc("POST /api/setup", h.setupUser)
	h.Mux.HandleFunc("GET /api/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		h.authMu.Lock()
		needsSetup := h.setup != ""
		h.authMu.Unlock()
		JSON(w, 200, map[string]any{"name": h.Config.Name, "id": h.Config.ID, "version": h.Config.Version, "oidc": h.Config.Issuer != "", "needsSetup": needsSetup, "workspace": h.Config.WorkspaceURL, "base": h.Base})
	})
	h.Mux.HandleFunc("GET /api/me", h.Require(func(w http.ResponseWriter, r *http.Request) { JSON(w, 200, Current(r)) }))
	h.Mux.HandleFunc("POST /api/local-users", h.Owner(h.localUser))
	h.Mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if h.ui == nil {
			http.NotFound(w, r)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		f, e := h.ui.Open(path)
		if e != nil {
			http.NotFound(w, r)
			return
		}
		f.Close()
		http.FileServer(http.FS(h.ui)).ServeHTTP(w, r)
	})
}
func (h *Host) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				h.Log.Error("request panicked", "route", r.Pattern)
				Problem(w, 500, "The request could not finish; please retry")
			}
		}()
		u, _ := url.Parse(h.Config.Origin)
		host := r.Host
		_, port, _ := net.SplitHostPort(h.Config.Listen)
		allowed := host == u.Host || host == h.Config.Listen || host == "localhost:"+port || host == "127.0.0.1:"+port || host == "[::1]:"+port
		if !allowed {
			Problem(w, 403, "Open the configured application address")
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		scripts := "'self'"
		if h.Config.ID == "pdf" {
			scripts += " 'wasm-unsafe-eval'"
		}
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src "+scripts+"; style-src 'self' 'unsafe-inline'; font-src 'self' blob: data:; img-src 'self' data: blob:; worker-src 'self' blob:; frame-src 'self' blob:; object-src 'none'; base-uri 'self'; frame-ancestors 'self'")
		if h.Base != "" && (r.URL.Path == h.Base || strings.HasPrefix(r.URL.Path, h.Base+"/")) {
			if r.URL.Path == h.Base {
				http.Redirect(w, r, h.Base+"/", 303)
				return
			}
			r = r.Clone(r.Context())
			r.URL.Path = strings.TrimPrefix(r.URL.Path, h.Base)
		}
		h.Mux.ServeHTTP(w, r)
	})
}
func (h *Host) Run(ctx context.Context) error {
	l, e := net.Listen("tcp", h.Config.Listen)
	if e != nil {
		return e
	}
	srv := &http.Server{Handler: h.Handler(), ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	done := make(chan error, 1)
	go func() { done <- srv.Serve(l) }()
	h.Log.Info("application started", "app", h.Config.ID, "version", h.Config.Version)
	select {
	case e = <-done:
		if !errors.Is(e, http.ErrServerClosed) {
			return e
		}
	case <-ctx.Done():
	}
	stop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(stop)
}
func (h *Host) SetupURL() string {
	h.authMu.Lock()
	defer h.authMu.Unlock()
	if h.setup == "" {
		return ""
	}
	return h.Config.Origin + "/#setup/" + h.setup
}
func (h *Host) BlobPath(id string) (string, error) {
	if len(id) != 64 {
		return "", errors.New("invalid object identifier")
	}
	if _, e := hex.DecodeString(id); e != nil {
		return "", e
	}
	return platform.Inside(h.Config.Data, "objects/"+id)
}
func (h *Host) PutBlob(r io.Reader, max int64) (string, int64, error) {
	tmpdir := filepath.Join(h.Config.Data, "tmp")
	if e := os.MkdirAll(tmpdir, 0700); e != nil {
		return "", 0, e
	}
	f, e := os.CreateTemp(tmpdir, "upload-")
	if e != nil {
		return "", 0, e
	}
	name := f.Name()
	defer os.Remove(name)
	sum := sha256.New()
	n, e := io.Copy(io.MultiWriter(f, sum), io.LimitReader(r, max+1))
	if e == nil {
		e = f.Sync()
	}
	f.Close()
	if e != nil {
		return "", 0, e
	}
	if n > max {
		return "", 0, errors.New("file exceeds the application size limit")
	}
	id := hex.EncodeToString(sum.Sum(nil))
	path, e := h.BlobPath(id)
	if e != nil {
		return "", 0, e
	}
	if e = os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		return "", 0, e
	}
	if _, e = os.Stat(path); e == nil {
		return id, n, nil
	}
	if e = os.Rename(name, path); e != nil {
		return "", 0, e
	}
	return id, n, nil
}
func (h *Host) DBError(w http.ResponseWriter, e error) {
	h.Log.Error("database operation failed", "kind", fmt.Sprintf("%T", e))
	Problem(w, 500, "Could not save this change. Retry in a moment.")
}
