package suite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Hosting receives a read-only broker capability, never Git's administrator
// password. Git itself authorizes every request as the linked MyAccounts user.
func hostingGitCapability(secret string) string {
	if secret == "" {
		return ""
	}
	s := sha256.Sum256([]byte("workspace-hosting-git-v1\x00" + secret))
	return hex.EncodeToString(s[:])
}

type hostingGitCall struct {
	Operation    string `json:"operation"`
	Issuer       string `json:"issuer"`
	Subject      string `json:"subject"`
	Repository   string `json:"repository"`
	RepositoryID int64  `json:"repositoryId,omitempty"`
	Ref          string `json:"ref,omitempty"`
	File         string `json:"file,omitempty"`
	Page         int    `json:"page,omitempty"`
}

type hostingGitRepo struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
	Empty         bool   `json:"empty"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}

var hostingGitSegment = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,100}$`)
var hostingGitRef = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_./-]{0,150}$`)
var hostingGitSHA = regexp.MustCompile(`^[a-f0-9]{40}$`)

func hostingRepoName(raw, base string) (string, error) {
	s := strings.TrimSpace(raw)
	if strings.Contains(s, "://") {
		u, e := url.Parse(s)
		b, be := url.Parse(base)
		if e != nil || be != nil || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Scheme != b.Scheme || !strings.EqualFold(u.Host, b.Host) || !strings.HasPrefix(u.Path, strings.TrimRight(b.Path, "/")+"/") || u.RawPath != "" {
			return "", errors.New("Use a repository from this Workspace's My Git")
		}
		s = strings.TrimPrefix(u.Path, strings.TrimRight(b.Path, "/")+"/")
	}
	s = strings.TrimSuffix(strings.TrimRight(s, "/"), ".git")
	parts := strings.Split(s, "/")
	if len(parts) < 1 || len(parts) > 2 {
		return "", errors.New("Enter a repository name, owner/repository, or its My Git URL")
	}
	for _, p := range parts {
		if !hostingGitSegment.MatchString(p) || p == "." || p == ".." {
			return "", errors.New("Invalid repository name")
		}
	}
	return strings.Join(parts, "/"), nil
}

// Read-only identity lookup; neither emails nor user-supplied Git usernames are
// used to link accounts. An absent/disabled/ambiguous connection is denied.
func (rt *Runtime) hostingGitUsername(ctx context.Context, subject string) (string, error) {
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
		return "", errors.New("My Git identity is unavailable")
	}
	defer db.Close()
	var source int64
	var raw string
	if e = db.QueryRowContext(ctx, `SELECT id,cfg FROM login_source WHERE name='myaccounts' AND type=6 AND is_active=1`).Scan(&source, &raw); e != nil {
		return "", errors.New("My Git is not connected to MyAccounts")
	}
	var cfg struct{ ClientID, ClientSecret, OpenIDConnectAutoDiscoveryURL string }
	client := rt.Secrets.Clients["gitea"]
	if json.Unmarshal([]byte(raw), &cfg) != nil || cfg.ClientID != client.ID || !platformSecretEqual(cfg.ClientSecret, client.Secret) || cfg.OpenIDConnectAutoDiscoveryURL != strings.TrimRight(rt.Config.Issuer, "/")+"/.well-known/openid-configuration" {
		return "", errors.New("My Git sign-in connection needs attention")
	}
	rows, e := db.QueryContext(ctx, `SELECT DISTINCT u.name,u.is_active,u.prohibit_login FROM user u LEFT JOIN external_login_user x ON x.user_id=u.id WHERE (u.login_source=? AND u.login_type=6 AND u.login_name=?) OR (x.login_source_id=? AND x.external_id=?)`, source, subject, source, subject)
	if e != nil {
		return "", errors.New("My Git identity is unavailable")
	}
	defer rows.Close()
	name := ""
	for rows.Next() {
		var n string
		var active, blocked bool
		if e = rows.Scan(&n, &active, &blocked); e != nil {
			return "", e
		}
		if !active || blocked || n == "workspace-recovery" || !hostingGitSegment.MatchString(n) || name != "" && name != n {
			return "", errors.New("My Git account is unavailable")
		}
		name = n
	}
	if e = rows.Err(); e != nil {
		return "", e
	}
	if name == "" {
		return "", errors.New("Open My Git once to link your MyAccounts sign-in, then retry")
	}
	return name, nil
}

func hostingGitRead(ctx context.Context, g *myaiGitClient, user, endpoint string) (*http.Response, error) {
	r, e := http.NewRequestWithContext(ctx, http.MethodGet, g.base+"/api/v1"+endpoint, nil)
	if e != nil {
		return nil, e
	}
	r.SetBasicAuth("workspace-recovery", g.password)
	r.Header.Set("Sudo", user)
	r.Header.Set("Accept", "application/json")
	response, e := g.client.Do(r)
	if e != nil {
		return nil, errors.New("My Git is unavailable; retry shortly")
	}
	return response, nil
}

func hostingGitJSON(ctx context.Context, g *myaiGitClient, user, endpoint string, out any) (int, error) {
	r, e := hostingGitRead(ctx, g, user, endpoint)
	if e != nil {
		return 0, e
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return r.StatusCode, errors.New("Repository not found or read access is not granted in My Git")
	}
	if e = json.NewDecoder(io.LimitReader(r.Body, 8<<20)).Decode(out); e != nil {
		return 502, errors.New("My Git returned an unreadable response")
	}
	return 200, nil
}

func hostingResolveRepo(ctx context.Context, g *myaiGitClient, user, name string) (hostingGitRepo, int, error) {
	var repo hostingGitRepo
	if strings.Contains(name, "/") {
		status, e := hostingGitJSON(ctx, g, user, "/repos/"+name, &repo)
		return repo, status, e
	}
	matches := []hostingGitRepo{}
	for page := 1; page <= 50; page++ {
		var rs []hostingGitRepo
		status, e := hostingGitJSON(ctx, g, user, "/user/repos?limit=100&page="+strconv.Itoa(page), &rs)
		if e != nil {
			return repo, status, e
		}
		for _, r := range rs {
			if strings.EqualFold(r.Name, name) {
				if strings.EqualFold(r.Owner.Login, user) {
					return r, 200, nil
				}
				matches = append(matches, r)
			}
		}
		if len(rs) < 100 {
			if len(matches) == 1 {
				return matches[0], 200, nil
			}
			if len(matches) > 1 {
				return repo, 409, errors.New("Several accessible repositories have that name; enter owner/repository to choose one")
			}
			return repo, 404, errors.New("Repository not found in your My Git account; check its name or sharing access")
		}
	}
	return repo, 400, errors.New("Enter owner/repository to search a large My Git account")
}

func hostingGitEndpoint(in hostingGitCall, fullName string) (string, error) {
	base := "/repos/" + fullName
	ref := in.Ref
	if ref != "" && (!hostingGitRef.MatchString(ref) || strings.Contains(ref, "..")) {
		return "", errors.New("Invalid branch or commit")
	}
	switch in.Operation {
	case "metadata":
		return base, nil
	case "branches":
		return base + "/branches?limit=100", nil
	case "commit":
		if ref != "" {
			return base + "/commits/" + url.PathEscape(ref), nil
		}
	case "commits":
		if in.Page < 0 || in.Page > 100 {
			return "", errors.New("Invalid history page")
		}
		page := in.Page
		if page == 0 {
			page = 1
		}
		return base + "/commits?limit=100&page=" + strconv.Itoa(page) + "&sha=" + url.QueryEscape(ref), nil
	case "contents":
		if in.File == "" || path.Clean(in.File) != in.File || strings.HasPrefix(in.File, "/") || strings.ContainsAny(in.File, "\\\x00") || strings.Contains("/"+in.File+"/", "/../") {
			return "", errors.New("Invalid repository file")
		}
		parts := strings.Split(in.File, "/")
		for i := range parts {
			parts[i] = url.PathEscape(parts[i])
		}
		return base + "/contents/" + strings.Join(parts, "/") + "?ref=" + url.QueryEscape(ref), nil
	case "archive":
		if hostingGitSHA.MatchString(ref) {
			return base + "/archive/" + ref + ".tar.gz", nil
		}
	}
	return "", errors.New("Unsupported read-only repository operation")
}

func (p *Portal) hostingGitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/_integration/hosting/git", p.hostingGitRequest)
}

func (p *Portal) hostingGitRequest(w http.ResponseWriter, r *http.Request) {
	host, _, e := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(host)
	want := hostingGitCapability(p.rt.Secrets.Clients["hosting"].Secret)
	if e != nil || ip == nil || !ip.IsLoopback() || want == "" || !platformSecretEqual(want, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST", 405)
		return
	}
	var in hostingGitCall
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16384))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&in) != nil || decoder.Decode(new(any)) != io.EOF || in.Issuer != p.rt.Config.Issuer || in.Subject == "" || len(in.Subject) > 512 {
		http.Error(w, "Invalid account or repository request", 400)
		return
	}
	name, e := hostingRepoName(in.Repository, strings.TrimRight(p.rt.Config.Origin, "/")+"/git")
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	user, e := p.rt.hostingGitUsername(r.Context(), in.Subject)
	if e != nil {
		http.Error(w, e.Error(), 403)
		return
	}
	g, e := p.rt.myaiGit(r.Context())
	if e != nil {
		http.Error(w, e.Error(), 503)
		return
	}
	repo, status, e := hostingResolveRepo(r.Context(), g, user, name)
	if e != nil {
		if status < 400 {
			status = 502
		}
		http.Error(w, e.Error(), status)
		return
	}
	if repo.FullName == "" {
		repo.FullName = repo.Owner.Login + "/" + repo.Name
	}
	canonical, e := hostingRepoName(repo.FullName, "")
	if e != nil || !strings.Contains(canonical, "/") || repo.ID <= 0 || in.RepositoryID != 0 && in.RepositoryID != repo.ID {
		http.Error(w, "Repository identity changed; reconnect the source", 409)
		return
	}
	if in.Operation == "resolve" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"provider": "mygit", "repository": strings.TrimRight(p.rt.Config.Origin, "/") + "/git/" + canonical, "fullName": canonical, "id": repo.ID, "defaultBranch": repo.DefaultBranch, "private": repo.Private, "empty": repo.Empty})
		return
	}
	endpoint, e := hostingGitEndpoint(in, canonical)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	resp, e := hostingGitRead(r.Context(), g, user, endpoint)
	if e != nil {
		http.Error(w, e.Error(), 502)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		status = resp.StatusCode
		if status < 400 {
			status = 502
		}
		http.Error(w, fmt.Sprintf("My Git source request failed (HTTP %d)", resp.StatusCode), status)
		return
	}
	max := int64(8 << 20)
	w.Header().Set("Content-Type", "application/json")
	if in.Operation == "archive" {
		max = 64 << 20
		w.Header().Set("Content-Type", "application/gzip")
	}
	if resp.ContentLength > max {
		http.Error(w, "Repository response exceeds its size limit", 413)
		return
	}
	// The receiver enforces the same streaming bound and hashes the completed
	// archive before running a build. Upstream credentials/headers never escape.
	_, _ = io.Copy(w, io.LimitReader(resp.Body, max+1))
}
