package suite

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	_ "modernc.org/sqlite"
)

const myaiProjectsSource = "Projects"

var myaiProjectID = regexp.MustCompile(`^[a-f0-9]{48}$`)

type myaiProject struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	GitDir       string `json:"gitDir"`
	Folder       string `json:"folder"`
	Repository   string `json:"repository"`
	RepositoryID int64  `json:"repositoryId"`
	GitUser      string `json:"gitUser"`
	URL          string `json:"url"`
	DriveURL     string `json:"driveUrl"`
	Branch       string `json:"branch"`
	State        string `json:"state"`
	Commit       string `json:"commit,omitempty"`
	Error        string `json:"error,omitempty"`
	Managed      bool   `json:"managed"`
}

func myaiProjectsRoot(c Config) string { return filepath.Join(c.Data, "myai-projects", "Projects") }
func myaiProjectsMeta(c Config, uid string) string {
	return filepath.Join(c.Data, "myai-projects", "registry", uid)
}
func myaiProjectSlug(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if b.Len() > 0 && !strings.HasSuffix(b.String(), "-") {
			b.WriteByte('-')
		}
		if b.Len() >= 45 {
			break
		}
	}
	s := strings.Trim(b.String(), "-")
	if s == "" {
		s = "project"
	}
	return s
}
func myaiSafeTreePath(p string) error {
	for q := p; ; q = filepath.Dir(q) {
		if _, e := os.Lstat(q); e == nil {
			if e = platform.RejectReparse(q); e != nil {
				return errors.New("Projects contains a linked folder")
			}
		} else if !os.IsNotExist(e) {
			return e
		}
		if filepath.Dir(q) == q {
			break
		}
	}
	return nil
}
func (rt *Runtime) myaiProjectSave(uid string, p myaiProject) error {
	if !myaiProjectID.MatchString(p.ID) {
		return errors.New("invalid project")
	}
	root := myaiProjectsMeta(rt.Config, uid)
	if e := myaiSafeTreePath(root); e != nil {
		return e
	}
	return platform.SaveJSON(filepath.Join(root, p.ID+".json"), p)
}
func (rt *Runtime) myaiProjectLoad(uid, id string) (p myaiProject, e error) {
	if !myaiProjectID.MatchString(id) {
		return p, errors.New("project not found")
	}
	file := filepath.Join(myaiProjectsMeta(rt.Config, uid), id+".json")
	if e = myaiSafeTreePath(file); e != nil {
		return
	}
	b, e := os.ReadFile(file)
	if e != nil {
		return p, errors.New("project not found")
	}
	if json.Unmarshal(b, &p) != nil || p.ID != id || !p.Managed || p.Folder != myaiProjectSlug(p.Name)+"-"+id[:8] {
		return p, errors.New("project registration needs attention")
	}
	want := filepath.Join(myaiProjectsRoot(rt.Config), "users", uid, p.Folder)
	git := filepath.Join(rt.Config.Data, "myai-projects", "git", uid, id)
	if !strings.EqualFold(filepath.Clean(want), filepath.Clean(p.Path)) || !strings.EqualFold(filepath.Clean(git), filepath.Clean(p.GitDir)) {
		return p, errors.New("project location needs attention")
	}
	for _, v := range []string{p.Path, p.GitDir} {
		if e = myaiSafeTreePath(v); e != nil {
			return
		}
	}
	return p, nil
}
func (rt *Runtime) myaiProjectList(uid string) ([]myaiProject, error) {
	out := []myaiProject{}
	root := myaiProjectsMeta(rt.Config, uid)
	if e := myaiSafeTreePath(root); e != nil {
		return nil, e
	}
	files, e := os.ReadDir(root)
	if os.IsNotExist(e) {
		return out, nil
	}
	if e != nil {
		return nil, e
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		p, e := rt.myaiProjectLoad(uid, strings.TrimSuffix(f.Name(), ".json"))
		if e != nil {
			return nil, e
		}
		if p.Commit != "" && p.State != "setup-needed" {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			status, se := rt.myaiProjectStatus(ctx, p)
			cancel()
			if se == nil && status != "" {
				p.State = "changes"
			} else if se != nil {
				p.State = "needs-sync"
				p.Error = "Project status is unavailable; your files are retained"
			}
		}
		out = append(out, p)
	}
	return out, nil
}

type myaiGitClient struct {
	base, password string
	stateRoot      string
	client         *http.Client
}

func (rt *Runtime) myaiGit(ctx context.Context) (*myaiGitClient, error) {
	app, ok := rt.Config.Catalog.Find("gitea")
	if !ok {
		return nil, errors.New("Projects needs the Git app installed in Workspace")
	}
	file := filepath.Join(rt.Config.Data, "gitea", "recovery.dpapi")
	if e := myaiSafeTreePath(file); e != nil {
		return nil, e
	}
	blob, e := os.ReadFile(file)
	if e != nil {
		return nil, errors.New("Git is not ready yet; try again shortly")
	}
	password, e := platform.Unprotect(blob)
	if e != nil {
		return nil, errors.New("Git could not unlock its protected connection")
	}
	return &myaiGitClient{base: fmt.Sprintf("http://127.0.0.1:%d", app.Port), password: string(password), stateRoot: filepath.Join(rt.Config.Data, "myai-projects", "token-leases"), client: &http.Client{Timeout: 25 * time.Second, Transport: &http.Transport{Proxy: nil, IdleConnTimeout: 30 * time.Second}, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}
func (g *myaiGitClient) call(ctx context.Context, method, path string, body, out any) (int, error) {
	var reader io.Reader
	if body != nil {
		b, e := json.Marshal(body)
		if e != nil {
			return 0, e
		}
		reader = bytes.NewReader(b)
	}
	r, e := http.NewRequestWithContext(ctx, method, g.base+"/api/v1"+path, reader)
	if e != nil {
		return 0, e
	}
	r.SetBasicAuth("workspace-recovery", g.password)
	r.Header.Set("Content-Type", "application/json")
	resp, e := g.client.Do(r)
	if e != nil {
		return 0, errors.New("Git is unavailable; your project files are safe")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("Git could not complete this step (HTTP %d); retry without creating another project", resp.StatusCode)
	}
	if out != nil {
		if e = json.NewDecoder(io.LimitReader(resp.Body, 2*1024*1024)).Decode(out); e != nil {
			return resp.StatusCode, errors.New("Git returned an unreadable result")
		}
	}
	return resp.StatusCode, nil
}
func (rt *Runtime) myaiGitUser(ctx context.Context, g *myaiGitClient, subject, uid, name string) (string, error) {
	file := filepath.Join(rt.Config.Data, "gitea", "gitea.db")
	if e := myaiSafeTreePath(file); e != nil {
		return "", e
	}
	dbPath := filepath.ToSlash(file)
	if filepath.VolumeName(file) != "" {
		dbPath = "/" + dbPath
	}
	u := url.URL{Scheme: "file", Path: dbPath}
	q := u.Query()
	q.Set("mode", "ro")
	q.Set("_pragma", "busy_timeout(5000)")
	u.RawQuery = q.Encode()
	db, e := sql.Open("sqlite", u.String())
	if e != nil {
		return "", errors.New("Git identity is unavailable")
	}
	defer db.Close()
	var source int64
	var cfgRaw string
	if e = db.QueryRowContext(ctx, `SELECT id,cfg FROM login_source WHERE name='myaccounts' AND type=6 AND is_active=1`).Scan(&source, &cfgRaw); e != nil {
		return "", errors.New("Git is not connected to MyAccounts")
	}
	var cfg struct{ ClientID, ClientSecret, OpenIDConnectAutoDiscoveryURL string }
	if json.Unmarshal([]byte(cfgRaw), &cfg) != nil || cfg.ClientID != rt.Secrets.Clients["gitea"].ID || !platformSecretEqual(cfg.ClientSecret, rt.Secrets.Clients["gitea"].Secret) || cfg.OpenIDConnectAutoDiscoveryURL != strings.TrimRight(rt.Config.Issuer, "/")+"/.well-known/openid-configuration" {
		return "", errors.New("Git sign-in connection needs attention")
	}
	lookup := func() (string, error) {
		rows, e := db.QueryContext(ctx, `SELECT DISTINCT u.name,u.is_active,u.prohibit_login FROM user u LEFT JOIN external_login_user x ON x.user_id=u.id WHERE (u.login_source=? AND u.login_type=6 AND u.login_name=?) OR (x.login_source_id=? AND x.external_id=?)`, source, subject, source, subject)
		if e != nil {
			return "", e
		}
		defer rows.Close()
		found := ""
		for rows.Next() {
			var n string
			var active, blocked bool
			if e = rows.Scan(&n, &active, &blocked); e != nil {
				return "", e
			}
			if !active || blocked || n == "workspace-recovery" || found != "" && found != n {
				return "", errors.New("Git account is unavailable")
			}
			found = n
		}
		return found, rows.Err()
	}
	if n, e := lookup(); e != nil || n != "" {
		return n, e
	}
	username := "myai-" + uid[:24]
	_, e = g.call(ctx, "POST", "/admin/users", map[string]any{"source_id": source, "login_name": subject, "username": username, "full_name": name, "email": uid + "@myai.localhost.invalid", "must_change_password": false, "send_notify": false, "visibility": "private"}, nil)
	if e != nil {
		return "", e
	}
	n, e := lookup()
	if e != nil || n != username {
		return "", errors.New("Git could not confirm the new MyAccounts connection")
	}
	return n, nil
}
func platformSecretEqual(a, b string) bool {
	return a != "" && b != "" && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

type myaiRepo struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
	Owner   struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func (g *myaiGitClient) repository(ctx context.Context, p *myaiProject, create bool) error {
	var repo myaiRepo
	path := "/repos/" + url.PathEscape(p.GitUser) + "/" + url.PathEscape(p.Repository)
	status, e := g.call(ctx, "GET", path, nil, &repo)
	if status == 404 && create {
		_, e = g.call(ctx, "POST", "/admin/users/"+url.PathEscape(p.GitUser)+"/repos", map[string]any{"name": p.Repository, "description": "Created with My AI", "private": true, "auto_init": false, "default_branch": "main"}, &repo)
	}
	if e != nil {
		return e
	}
	if repo.ID == 0 || repo.Owner.Login != p.GitUser || repo.Name != p.Repository || !repo.Private || p.RepositoryID != 0 && repo.ID != p.RepositoryID {
		return errors.New("This project's private Git repository needs attention")
	}
	p.RepositoryID = repo.ID
	return nil
}

type myaiTokenLease struct {
	Name, User string
	ID         int64
	Created    time.Time
}

func (g *myaiGitClient) tokenLeaseRoot(username string) string {
	h := sha256.Sum256([]byte(username))
	return filepath.Join(g.stateRoot, hex.EncodeToString(h[:]))
}
func (g *myaiGitClient) reapTokens(ctx context.Context, username string) error {
	root := g.tokenLeaseRoot(username)
	if e := myaiSafeTreePath(root); e != nil {
		return e
	}
	files, e := os.ReadDir(root)
	if os.IsNotExist(e) {
		return nil
	}
	if e != nil {
		return e
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		path := filepath.Join(root, file.Name())
		if e = myaiSafeTreePath(path); e != nil {
			return e
		}
		b, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		var lease myaiTokenLease
		if json.Unmarshal(b, &lease) != nil || lease.User != username || file.Name() != lease.Name+".json" || !strings.HasPrefix(lease.Name, "myai-sync-") {
			return errors.New("Project connection recovery needs attention")
		}
		if time.Since(lease.Created) < 5*time.Minute {
			continue
		}
		if lease.ID == 0 {
			var tokens []struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			}
			if _, e = g.call(ctx, "GET", "/users/"+url.PathEscape(username)+"/tokens", nil, &tokens); e != nil {
				return e
			}
			for _, token := range tokens {
				if token.Name == lease.Name {
					lease.ID = token.ID
					break
				}
			}
		}
		if lease.ID != 0 {
			status, e := g.call(ctx, "DELETE", fmt.Sprintf("/users/%s/tokens/%d", url.PathEscape(username), lease.ID), nil, nil)
			if e != nil && status != 404 {
				return e
			}
		}
		if e = os.Remove(path); e != nil {
			return e
		}
	}
	return nil
}
func (g *myaiGitClient) token(ctx context.Context, username string) (string, func(), error) {
	if e := g.reapTokens(ctx, username); e != nil {
		return "", func() {}, e
	}
	lease := myaiTokenLease{Name: "myai-sync-" + platform.Token()[:16], User: username, Created: time.Now().UTC()}
	path := filepath.Join(g.tokenLeaseRoot(username), lease.Name+".json")
	if e := platform.SaveJSON(path, lease); e != nil {
		return "", func() {}, errors.New("Git could not protect a temporary project connection")
	}
	var out struct {
		ID  int64  `json:"id"`
		SHA string `json:"sha1"`
	}
	_, e := g.call(ctx, "POST", "/users/"+url.PathEscape(username)+"/tokens", map[string]any{"name": lease.Name, "scopes": []string{"write:repository"}}, &out)
	if e != nil || out.ID == 0 || out.SHA == "" {
		return "", func() {}, errors.New("Git could not open a temporary project connection")
	}
	lease.ID = out.ID
	closeToken := func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		status, e := g.call(c, "DELETE", fmt.Sprintf("/users/%s/tokens/%d", url.PathEscape(username), out.ID), nil, nil)
		if e == nil || status == 404 {
			_ = os.Remove(path)
		}
	}
	if e = platform.SaveJSON(path, lease); e != nil {
		closeToken()
		return "", func() {}, errors.New("Git could not protect connection recovery")
	}
	return out.SHA, closeToken, nil
}
func (rt *Runtime) myaiGitBinary() string {
	p := filepath.Join(rt.Config.Programs, "gitea-runtime-28.1.0", "git", "cmd", "git.exe")
	if f, e := os.Stat(p); e == nil && !f.IsDir() {
		return p
	}
	return "git"
}
func (rt *Runtime) myaiGitRun(ctx context.Context, p myaiProject, token string, args ...string) (string, error) {
	empty := filepath.Join(rt.Config.Data, "myai-projects", "no-hooks")
	if e := os.MkdirAll(empty, 0700); e != nil {
		return "", e
	}
	prefix := []string{"-c", "core.hooksPath=" + filepath.ToSlash(empty), "-c", "core.fsmonitor=false", "-c", "credential.helper=", "-c", "core.attributesFile=NUL", "-c", "commit.gpgSign=false", "-c", "tag.gpgSign=false", "-c", "protocol.allow=never", "-c", "protocol.http.allow=always", "-c", "protocol.https.allow=always", "--git-dir=" + p.GitDir, "--work-tree=" + p.Path}
	cmd := exec.CommandContext(ctx, rt.myaiGitBinary(), append(prefix, args...)...)
	cmd.Dir = p.Path
	cmd.Env = myaiGitEnvironment()
	if token != "" {
		cmd.Env = append(cmd.Env, "GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=http.extraHeader", "GIT_CONFIG_VALUE_0=Authorization: Basic "+base64.StdEncoding.EncodeToString([]byte(p.GitUser+":"+token)))
	}
	cmd.SysProcAttr = hiddenProcessAttrs()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = io.Discard
	if e := cmd.Run(); e != nil {
		return "", errors.New("Git could not complete this step; your project files and local history are retained")
	}
	if output.Len() > 4*1024*1024 {
		return "", errors.New("Project changes are too large to display")
	}
	return output.String(), nil
}
func myaiProjectForbidden(rel string) bool {
	for _, v := range strings.Split(filepath.ToSlash(rel), "/") {
		l := strings.ToLower(v)
		if l == ".git" || l == ".ssh" || l == ".codex" || l == ".agents" || strings.HasPrefix(l, ".env") && l != ".env.example" || strings.Contains(l, "credential") || l == "id_rsa" || l == "id_ed25519" || strings.HasSuffix(l, ".pem") || strings.HasSuffix(l, ".key") || strings.HasSuffix(l, ".dpapi") {
			return true
		}
	}
	return false
}
func (rt *Runtime) myaiProjectInspect(p myaiProject) error {
	count := 0
	var size int64
	return filepath.WalkDir(p.Path, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return errors.New("Project files are unavailable")
		}
		count++
		if count > 20000 {
			return errors.New("Project has too many files to save at once")
		}
		if e = myaiSafeTreePath(path); e != nil {
			return e
		}
		if d.IsDir() {
			if d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		f, e := d.Info()
		if e != nil {
			return e
		}
		size += f.Size()
		if f.Size() > 64*1024*1024 || size > 256*1024*1024 {
			return errors.New("Project is too large to save at once; keep large downloads outside Projects")
		}
		return nil
	})
}
func (rt *Runtime) myaiProjectInitialize(ctx context.Context, p myaiProject) error {
	if e := os.MkdirAll(p.Path, 0700); e != nil {
		return e
	}
	if e := os.MkdirAll(filepath.Dir(p.GitDir), 0700); e != nil {
		return e
	}
	if _, e := os.Stat(filepath.Join(p.GitDir, "HEAD")); os.IsNotExist(e) {
		cmd := exec.CommandContext(ctx, rt.myaiGitBinary(), "-c", "core.hooksPath=", "init", "--initial-branch=main", "--separate-git-dir", p.GitDir, p.Path)
		cmd.Env = myaiGitEnvironment()
		cmd.SysProcAttr = hiddenProcessAttrs()
		if cmd.Run() != nil {
			return errors.New("Git could not prepare this project")
		}
	}
	readme := filepath.Join(p.Path, "README.md")
	f, e := os.OpenFile(readme, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e == nil {
		_, e = f.WriteString("# " + p.Name + "\n\nCreated with My AI.\n")
		f.Close()
	}
	if e != nil && !os.IsExist(e) {
		return e
	}
	return nil
}
func (rt *Runtime) myaiProjectSync(ctx context.Context, g *myaiGitClient, p *myaiProject, prepare bool) error {
	if e := g.repository(ctx, p, false); e != nil {
		return e
	}
	if e := rt.myaiProjectInspect(*p); e != nil {
		return e
	}
	token, closeToken, e := g.token(ctx, p.GitUser)
	if e != nil {
		return e
	}
	defer closeToken()
	remote := g.base + "/" + url.PathEscape(p.GitUser) + "/" + url.PathEscape(p.Repository) + ".git"
	status, e := rt.myaiProjectStatus(ctx, *p)
	if e != nil {
		return e
	}
	refs, e := rt.myaiGitRun(ctx, *p, token, "ls-remote", remote, "refs/heads/main")
	if e != nil {
		return e
	}
	remoteHead := ""
	if parts := strings.Fields(refs); len(parts) > 0 {
		remoteHead = parts[0]
	}
	local, _ := rt.myaiGitRun(ctx, *p, "", "rev-parse", "--verify", "HEAD")
	local = strings.TrimSpace(local)
	if remoteHead != "" && remoteHead != local {
		if _, e = rt.myaiGitRun(ctx, *p, token, "fetch", "--no-tags", remote, "refs/heads/main"); e != nil {
			return e
		}
		if _, e = rt.myaiGitRun(ctx, *p, "", "merge-base", "--is-ancestor", "FETCH_HEAD", "HEAD"); e != nil {
			if status != "" {
				return errors.New("Git has newer changes from another window. Your files are safe; review the changes before continuing")
			}
			if _, e = rt.myaiGitRun(ctx, *p, "", "merge", "--ff-only", "FETCH_HEAD"); e != nil {
				return errors.New("This project has different changes in Git and My Drive. Review them before continuing; nothing was overwritten")
			}
		}
	}
	if prepare {
		p.State = "ready"
		p.Error = ""
		if status != "" {
			p.State = "changes"
		}
		head, _ := rt.myaiGitRun(ctx, *p, "", "rev-parse", "--verify", "HEAD")
		p.Commit = strings.TrimSpace(head)
		if status == "" && p.Commit == remoteHead {
			p.State = "synced"
		}
		return nil
	}
	excludes := []string{"add", "-A", "--", "."}
	for _, pattern := range []string{"**/.env", "**/.env.*", "**/.git/**", "**/.ssh/**", "**/.codex/**", "**/.agents/**", "**/*credential*", "**/*.pem", "**/*.key", "**/*.dpapi", "**/id_rsa", "**/id_ed25519", "**/node_modules/**"} {
		excludes = append(excludes, ":(exclude,glob,icase)"+pattern)
	}
	if _, e = rt.myaiGitRun(ctx, *p, "", excludes...); e != nil {
		return e
	}
	_ = filepath.WalkDir(p.Path, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == ".git") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.EqualFold(d.Name(), ".env.example") {
			rel, e := filepath.Rel(p.Path, path)
			if e != nil {
				return e
			}
			if myaiProjectForbidden(rel) {
				return nil
			}
			_, e = rt.myaiGitRun(ctx, *p, "", "add", "--", filepath.ToSlash(rel))
			return e
		}
		return nil
	})
	changed, _ := rt.myaiGitRun(ctx, *p, "", "diff", "--cached", "--name-only")
	if changed != "" {
		for _, file := range strings.Split(strings.TrimSpace(changed), "\n") {
			if myaiProjectForbidden(file) {
				return errors.New("A credential file was excluded from project history")
			}
		}
		if _, e = rt.myaiGitRun(ctx, *p, "", "commit", "-m", "Save changes from My AI"); e != nil {
			return e
		}
	}
	head, e := rt.myaiGitRun(ctx, *p, "", "rev-parse", "--verify", "HEAD")
	if e != nil {
		return e
	}
	head = strings.TrimSpace(head)
	if _, e = rt.myaiGitRun(ctx, *p, token, "push", remote, "HEAD:refs/heads/main"); e != nil {
		return errors.New("Your changes are saved in My Drive and local history, but Git has not confirmed them. Use Sync to retry")
	}
	refs, e = rt.myaiGitRun(ctx, *p, token, "ls-remote", remote, "refs/heads/main")
	if e != nil || !strings.HasPrefix(refs, head+"\t") {
		return errors.New("Git has not confirmed the latest save. Use Sync to retry")
	}
	p.Commit = head
	p.State = "synced"
	p.Error = ""
	return nil
}
func (p *Portal) myaiProjectsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/_integration/myai/projects", p.myaiProjectsRequest)
	mux.HandleFunc("/_integration/myai/projects/", p.myaiProjectsRequest)
}
func (p *Portal) myaiProjectsRequest(w http.ResponseWriter, r *http.Request) {
	host, _, e := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(host)
	want := myaiCapability(p.rt.Secrets.Clients["myai"].Secret)
	got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	subject := r.Header.Get("X-MyAI-Subject")
	if e != nil || ip == nil || !ip.IsLoopback() || !platformSecretEqual(want, got) || subject == "" || len(subject) > 256 || strings.ContainsAny(subject, "\r\n\x00") || r.Header.Get("X-MyAI-Issuer") != p.rt.Config.Issuer {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	fail := func(code int, e error) {
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": e.Error()})
	}
	if e = p.rt.verifyMyAIStorage(r.Context(), false); e != nil {
		fail(503, e)
		return
	}
	if e = p.rt.bindMyAIAccount(r.Context(), subject); e != nil {
		fail(503, errors.New("My Drive's private Projects area is unavailable"))
		return
	}
	uid := strings.TrimPrefix(myaiScope(subject, p.rt.Config.Issuer), "/users/")
	tail := strings.Trim(strings.TrimPrefix(r.URL.Path, "/_integration/myai/projects"), "/")
	parts := strings.Split(tail, "/")
	p.rt.myaiProjectsMu.Lock()
	defer p.rt.myaiProjectsMu.Unlock()
	if r.Method == "GET" && tail == "" {
		out, e := p.rt.myaiProjectList(uid)
		if e != nil {
			fail(503, e)
			return
		}
		_ = json.NewEncoder(w).Encode(out)
		return
	}
	if r.Method == "GET" && len(parts) == 1 {
		project, e := p.rt.myaiProjectLoad(uid, parts[0])
		if e != nil {
			fail(404, e)
			return
		}
		_ = json.NewEncoder(w).Encode(project)
		return
	}
	if r.Method != "POST" {
		fail(405, errors.New("unsupported project action"))
		return
	}
	var in struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&in) != nil {
		fail(400, errors.New("Enter a project name"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()
	g, e := p.rt.myaiGit(ctx)
	if e != nil {
		fail(503, e)
		return
	}
	if tail == "" {
		in.Name = strings.TrimSpace(in.Name)
		if !myaiProjectID.MatchString(in.ID) || in.Name == "" || len([]rune(in.Name)) > 100 || strings.ContainsFunc(in.Name, unicode.IsControl) {
			fail(400, errors.New("Give your project a name up to100 characters"))
			return
		}
		project, e := p.rt.myaiProjectLoad(uid, in.ID)
		if e != nil {
			if _, se := os.Stat(filepath.Join(myaiProjectsMeta(p.rt.Config, uid), in.ID+".json")); !os.IsNotExist(se) {
				fail(503, e)
				return
			}

			folder := myaiProjectSlug(in.Name) + "-" + in.ID[:8]
			project = myaiProject{ID: in.ID, Name: in.Name, Path: filepath.Join(myaiProjectsRoot(p.rt.Config), "users", uid, folder), GitDir: filepath.Join(p.rt.Config.Data, "myai-projects", "git", uid, in.ID), Folder: folder, Repository: folder, Branch: "main", State: "creating", Managed: true}
			if e = p.rt.myaiProjectDestinationAvailable(uid, project); e != nil {
				fail(409, e)
				return
			}
			if e = p.rt.myaiProjectSave(uid, project); e != nil {
				fail(503, e)
				return
			}
		}
		name := strings.TrimSpace(r.Header.Get("X-MyAI-Name"))
		if name == "" {
			name = "My AI projects"
		}
		if len([]rune(name)) > 100 {
			name = string([]rune(name)[:100])
		}
		project.GitUser, e = p.rt.myaiGitUser(ctx, g, subject, uid, name)
		if e == nil {
			project.URL = strings.TrimRight(p.rt.Config.Origin, "/") + "/git/" + url.PathEscape(project.GitUser) + "/" + url.PathEscape(project.Repository)
			project.DriveURL = strings.TrimRight(p.rt.Config.Origin, "/") + "/drive/files/Projects/" + url.PathEscape(project.Folder) + "/"
			e = g.repository(ctx, &project, true)
		}
		if e == nil {
			e = p.rt.myaiProjectInitialize(ctx, project)
		}
		if e == nil {
			e = p.rt.myaiProjectSync(ctx, g, &project, false)
		}
		if e != nil {
			project.State = "setup-needed"
			project.Error = e.Error()
		}
		if se := p.rt.myaiProjectSave(uid, project); se != nil {
			fail(503, se)
			return
		}
		if e != nil {
			fail(503, e)
			return
		}
		_ = json.NewEncoder(w).Encode(project)
		return
	}
	if len(parts) != 2 {
		fail(404, errors.New("project not found"))
		return
	}
	project, e := p.rt.myaiProjectLoad(uid, parts[0])
	if e != nil {
		fail(404, e)
		return
	}
	if parts[1] == "changes" {
		status, se := p.rt.myaiGitRun(ctx, project, "", "status", "--short", "--untracked-files=normal")
		diff, de := p.rt.myaiGitRun(ctx, project, "", "diff", "--no-ext-diff", "--no-textconv", "HEAD", "--", ".")
		if se != nil || de != nil {
			fail(503, errors.New("Project changes could not be read"))
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": status, "diff": diff})
		return
	}
	if parts[1] != "sync" && parts[1] != "prepare" {
		fail(404, errors.New("project action not found"))
		return
	}
	gitUser, e := p.rt.myaiGitUser(ctx, g, subject, uid, r.Header.Get("X-MyAI-Name"))
	if e == nil && gitUser != project.GitUser {
		e = errors.New("The Git account for this project changed")
	}
	if e == nil {
		e = p.rt.myaiProjectSync(ctx, g, &project, parts[1] == "prepare")
	}
	if e != nil {
		project.State = "needs-sync"
		project.Error = e.Error()
	}
	if se := p.rt.myaiProjectSave(uid, project); se != nil {
		fail(503, se)
		return
	}
	if e != nil {
		fail(409, e)
		return
	}
	_ = json.NewEncoder(w).Encode(project)
}

func (rt *Runtime) verifyMyAIProjectsStorage(repair bool) error {
	root := myaiProjectsRoot(rt.Config)
	if containedPath(rt.Config.Files, root) || containedPath(root, rt.Config.Files) {
		return errors.New("Projects overlaps shared Files")
	}
	if e := myaiSafeTreePath(root); e != nil {
		return e
	}
	marker := filepath.Join(root, ".workspace-projects-v1")
	if _, e := os.Stat(marker); os.IsNotExist(e) {
		if !repair {
			return errors.New("Projects storage needs a Workspace update")
		}
		if files, e := os.ReadDir(root); e == nil && len(files) > 0 {
			return errors.New("Projects storage is not empty")
		}
		if e = os.MkdirAll(root, 0700); e != nil {
			return e
		}
		return platform.AtomicWrite(marker, []byte("Private managed My AI projects v1\n"))
	}
	return nil
}

func (rt *Runtime) myaiProjectStatus(ctx context.Context, p myaiProject) (string, error) {
	raw, e := rt.myaiGitRun(ctx, p, "", "status", "--porcelain", "-z", "--untracked-files=all")
	if e != nil {
		return "", e
	}
	fields := strings.Split(raw, "\x00")
	out := []string{}
	for i := 0; i < len(fields); i++ {
		entry := fields[i]
		if len(entry) < 4 {
			continue
		}
		rel := entry[3:]
		renamed := entry[0] == 'R' || entry[1] == 'R' || entry[0] == 'C' || entry[1] == 'C'
		if renamed {
			i++
		}
		skip := myaiProjectForbidden(rel)
		for _, part := range strings.Split(filepath.ToSlash(rel), "/") {
			skip = skip || part == "node_modules" || strings.Contains(part, ".myai-")
		}
		if !skip {
			out = append(out, entry)
		}
	}
	return strings.Join(out, "\n"), nil
}

func myaiGitEnvironment() []string {
	out := []string{}
	for _, key := range []string{"SystemRoot", "WINDIR", "PATH", "PATHEXT", "COMSPEC", "TEMP", "TMP", "USERPROFILE", "APPDATA", "LOCALAPPDATA"} {
		if v := os.Getenv(key); v != "" {
			out = append(out, key+"="+v)
		}
	}
	return append(out, "GIT_CONFIG_GLOBAL=NUL", "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0", "GIT_AUTHOR_NAME=My AI", "GIT_AUTHOR_EMAIL=myai@localhost.invalid", "GIT_COMMITTER_NAME=My AI", "GIT_COMMITTER_EMAIL=myai@localhost.invalid")
}

func (rt *Runtime) myaiProjectDestinationAvailable(uid string, p myaiProject) error {
	if _, e := os.Lstat(p.Path); e == nil {
		return errors.New("A project folder already exists here. Choose another project name")
	} else if !os.IsNotExist(e) {
		return errors.New("Projects storage is unavailable")
	}
	files, e := os.ReadDir(myaiProjectsMeta(rt.Config, uid))
	if os.IsNotExist(e) {
		return nil
	}
	if e != nil {
		return errors.New("Projects registry is unavailable")
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		other, e := rt.myaiProjectLoad(uid, strings.TrimSuffix(f.Name(), ".json"))
		if e != nil {
			return e
		}
		if other.ID != p.ID && (other.Folder == p.Folder || other.Repository == p.Repository) {
			return errors.New("This project name is already in use. Choose another name")
		}
	}
	return nil
}
