package suite

import (
	"bytes"
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
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

type aiAppCall struct {
	App        string            `json:"app"`
	Operation  string            `json:"operation"`
	Parameters map[string]string `json:"parameters"`
	Query      map[string]string `json:"query"`
	Body       json.RawMessage   `json:"body"`
	// This field is set only by My AI after its own authenticated human review.
	Reviewed  bool   `json:"reviewed"`
	Prepare   bool   `json:"prepare"`
	RequestID string `json:"requestId,omitempty"`
}

type aiHTTPError struct {
	Status  int
	Message string
}

func (e *aiHTTPError) Error() string { return e.Message }

type aiAppConnection struct{ Cookies []*http.Cookie }
type aiConnectionRequired struct{ error }

func splitAIPath(path string) []string { return strings.Split(path, "/") }

// Remember native sessions only inside their current Workspace browser session.
// Never mint app sessions, adopt an account by email or use recovery/admin keys.
func (p *Portal) rememberAIAppSession(r *http.Request) {
	root, e := r.Cookie("workspace_session")
	if e != nil {
		return
	}
	info, ok := p.signInForPath(r.URL.Path)
	if !ok {
		return
	}
	cookies := []*http.Cookie{}
	for _, name := range info.cookies {
		if c, e := r.Cookie(name); e == nil && len(c.Value) > 0 && len(c.Value) < 8192 {
			cookies = append(cookies, &http.Cookie{Name: name, Value: c.Value})
		}
	}
	if len(cookies) == 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	s, ok := p.sessions[root.Value]
	if !ok || !time.Now().Before(s.Expires) {
		return
	}
	if s.AIApps == nil {
		s.AIApps = map[string]aiAppConnection{}
	}
	s.AIApps[info.id] = aiAppConnection{Cookies: cookies}
	p.sessions[root.Value] = s
}

func (p *Portal) aiBrowser(r *http.Request) (session, bool) {
	host, _, e := net.SplitHostPort(r.RemoteAddr)
	want := myaiCapability(p.rt.Secrets.Clients["myai"].Secret)
	if e != nil || net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() || !platformSecretEqual(want, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")) || r.Header.Get("X-MyAI-Issuer") != p.rt.Config.Issuer {
		return session{}, false
	}
	key := r.Header.Get("X-Workspace-Browser-Session")
	if len(key) > 192 {
		return session{}, false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	s, ok := p.sessions[key]
	if !ok || !time.Now().Before(s.Expires) || s.Subject == "" || s.Subject != r.Header.Get("X-MyAI-Subject") {
		return session{}, false
	}
	// Copy connection state before releasing the mutex. Session slices are immutable.
	copyApps := map[string]aiAppConnection{}
	for k, v := range s.AIApps {
		copyApps[k] = v
	}
	s.AIApps = copyApps
	return s, true
}

func (p *Portal) myaiAppsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/_integration/myai/workspace", func(w http.ResponseWriter, r *http.Request) {
		browser, ok := p.aiBrowser(r)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == "GET" {
			apps := []map[string]any{}
			states := p.rt.Status()
			for _, a := range p.rt.Config.Catalog.Apps {
				if a.ID == "monitor" && browser.Subject != p.rt.Config.Owner {
					continue
				}
				ops := aiOperations(a.ID)
				if ops == nil {
					ops = []aiOperation{}
				}
				probe := ""
				connectionID := a.ID
				if a.ID == "photos" {
					connectionID = "drive"
				}
				for _, signin := range suiteSignIns {
					if signin.id == connectionID {
						route := strings.TrimRight(a.Route, "/")
						probe = route + signin.probe
						if strings.HasPrefix(signin.probe, route+"/") {
							probe = signin.probe
						}
						break
					}
				}
				if a.ID == "photos" {
					probe = "/drive/api/users?username=self"
				}
				connected := len(browser.AIApps[connectionID].Cookies) > 0 || a.ID == "gitea"
				apps = append(apps, map[string]any{"id": a.ID, "name": a.Name, "description": a.Description, "url": p.rt.Config.Origin + "/_workspace/open/" + a.ID, "state": states[a.ID], "connected": connected, "probe": probe, "operations": ops, "nativeAccess": len(ops) > 0})
			}
			apphost.JSON(w, 200, map[string]any{"apps": apps, "home": p.rt.Config.Origin, "contract": "workspace-apps-v1", "notice": "App operations use your own native app sessions and permissions. Open an app once if it needs connecting. Apps without registered operations can be opened normally; credentials and administration are not exposed as generic tools."})
			return
		}
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		var call aiAppCall
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 300<<10))
		decoder.DisallowUnknownFields()
		if decoder.Decode(&call) != nil || decoder.Decode(new(any)) != io.EOF {
			apphost.Problem(w, 400, "Invalid app operation")
			return
		}
		app, installed := p.rt.Config.Catalog.Find(call.App)
		if !installed || call.App == "monitor" && browser.Subject != p.rt.Config.Owner {
			apphost.Problem(w, 404, "App is not available to this account")
			return
		}
		var operation aiOperation
		for _, op := range aiOperations(call.App) {
			if op.ID == call.Operation {
				operation = op
				break
			}
		}
		if operation.ID == "" {
			apphost.Problem(w, 400, "This action is not registered; discover the app's operations first")
			return
		}
		path, e := validateAIAppCall(operation, call)
		if e != nil {
			apphost.Problem(w, 400, e.Error())
			return
		}
		if operation.Method == "GET" {
			call.Body = nil // Empty JSON objects in model calls mean no HTTP body.
		}
		if operation.Review && !call.Reviewed && !call.Prepare {
			apphost.Problem(w, 409, "Review this app change in My AI before continuing")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer cancel()
		uid := apphost.Identity(p.rt.Config.Issuer, browser.Subject)
		var result any
		if call.Prepare {
			result, e = p.aiAppReview(ctx, browser, uid, app, operation, call)
		} else if strings.HasPrefix(path, "managed-") {
			result, e = p.aiManagedOperation(ctx, browser, uid, operation, call)
		} else {
			if app.ID == "photos" {
				app, installed = p.rt.Config.Catalog.Find("drive")
				if !installed {
					apphost.Problem(w, 409, "My Photos requires My Drive")
					return
				}
			}
			connection, csrf, err := p.aiNativeIdentity(ctx, browser, app)
			if err != nil {
				apphost.JSON(w, 409, map[string]any{"error": err.Error(), "connectUrl": p.rt.Config.Origin + "/_workspace/open/" + app.ID, "connectionRequired": true})
				return
			}
			var raw []byte
			var kind string
			body := call.Body
			if operation.MailMethod != "" {
				var args map[string]any
				_ = json.Unmarshal(body, &args)
				body, _ = json.Marshal(map[string]any{"method": operation.MailMethod, "args": args})
			}
			raw, kind, e = p.aiNativeRequest(ctx, app, connection, csrf, operation.Method, path, body, "application/json")
			if e == nil && call.App == "drive" && operation.ID == "read_file" {
				// Download MIME types for JavaScript/JSON/Markdown are not always
				// text/*. Return a regular UTF-8 source document as quoted text.
				if !aiTextFile(call.Query["file"]) || !utf8.Valid(raw) || bytes.ContainsRune(raw, 0) {
					e = errors.New("Choose a regular UTF-8 text document; open binary files in Drive")
				} else {
					kind = "text/plain; charset=utf-8"
				}
			}
			if e == nil {
				result, e = aiAppResult(raw, kind)
			}
		}
		if e != nil {
			var connectionError *aiConnectionRequired
			if errors.As(e, &connectionError) {
				apphost.JSON(w, 409, map[string]any{"error": e.Error(), "connectUrl": p.rt.Config.Origin + "/_workspace/open/" + call.App, "connectionRequired": true})
				return
			}
			apphost.Problem(w, 409, e.Error())
			return
		}
		apphost.JSON(w, 200, map[string]any{"app": call.App, "operation": operation.ID, "result": result, "untrusted_reference": true})
	})
}

var aiParameter = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,127}$`)

func aiTextFile(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".txt", ".md", ".markdown", ".csv", ".tsv", ".json", ".jsonl", ".yaml", ".yml", ".toml", ".xml", ".html", ".htm", ".css", ".scss", ".js", ".mjs", ".cjs", ".jsx", ".ts", ".tsx", ".go", ".py", ".rs", ".java", ".c", ".h", ".cpp", ".sql", ".svg", ".aui":
		return true
	}
	return false
}

func validateAIAppCall(op aiOperation, c aiAppCall) (string, error) {
	if len(c.Body) > 256<<10 || len(c.Parameters) > 8 || len(c.Query) > 12 {
		return "", errors.New("App arguments exceed the limit")
	}
	if c.App == "hosting" && op.ID == "deploy_project" && !aiActionID.MatchString(c.RequestID) {
		return "", errors.New("Deployment needs a server-generated action receipt")
	}
	if e := aiValidateBody(op, c.Body); e != nil {
		return "", e
	}
	path := op.Path
	for _, key := range op.Parameters {
		v := c.Parameters[key]
		if !aiParameter.MatchString(v) || strings.Contains(v, "..") {
			return "", fmt.Errorf("Supply %s from this app's results", key)
		}
		path = strings.ReplaceAll(path, "{"+key+"}", url.PathEscape(v))
	}
	for key := range c.Parameters {
		found := false
		for _, want := range op.Parameters {
			found = found || key == want
		}
		if !found {
			return "", errors.New("Unexpected app parameter")
		}
	}
	query := url.Values{}
	for key, value := range c.Query {
		found := false
		for _, want := range op.Query {
			found = found || key == want
		}
		if !found || len(value) > 4096 || strings.ContainsRune(value, 0) {
			return "", errors.New("Unexpected app query")
		}
		query.Set(key, value)
	}
	if len(c.Body) > 0 && string(c.Body) != "null" {
		var body map[string]any
		if json.Unmarshal(c.Body, &body) != nil {
			return "", errors.New("App body must be a JSON object")
		}
		if op.Method == "GET" && len(body) != 0 {
			return "", errors.New("Read actions do not accept a body")
		}
		if c.App == "hosting" {
			if trusted, _ := body["trusted"].(bool); trusted {
				return "", errors.New("Approve executable source in My Hosting itself before enabling scripts")
			}
		}
	}
	if c.App == "drive" {
		for _, k := range []string{"path", "file"} {
			for _, part := range strings.Split(strings.TrimPrefix(query.Get(k), "/"), "/") {
				if myaiProjectForbidden(part) {
					return "", errors.New("Credential and Git internal files are excluded from AI file access")
				}
			}
		}
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return path, nil
}

func (p *Portal) aiNativeRequest(ctx context.Context, app catalog.App, connection aiAppConnection, csrf, method, path string, body []byte, kind string) ([]byte, string, error) {
	route := strings.TrimRight(app.Route, "/")
	if !app.StripPrefix && !strings.HasPrefix(path, route+"/") {
		path = route + path
	}
	req, e := http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://127.0.0.1:%d%s", app.Port, path), bytes.NewReader(body))
	if e != nil {
		return nil, "", errors.New("App request could not be prepared")
	}
	origin, _ := url.Parse(p.rt.Config.Origin)
	req.Host = origin.Host
	req.Header.Set("Origin", origin.Scheme+"://"+origin.Host)
	req.Header.Set("X-Forwarded-Host", origin.Host)
	req.Header.Set("X-Forwarded-Proto", origin.Scheme)
	req.Header.Set("Content-Type", kind)
	req.Header.Set("X-CSRF-Token", csrf)
	if app.ID == "music" {
		req.Header.Set("X-Music-CSRF", csrf)
		for _, cookie := range connection.Cookies {
			if cookie.Name == "csrf_access_token" {
				req.Header.Set("X-CSRF-Token", cookie.Value)
			}
		}
	}
	req.Header.Set("Accept", "application/json")
	for _, cookie := range connection.Cookies {
		req.AddCookie(cookie)
	}
	client := &http.Client{Timeout: 75 * time.Second, Transport: &http.Transport{Proxy: nil, ResponseHeaderTimeout: 60 * time.Second}, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	defer client.CloseIdleConnections()
	response, e := client.Do(req)
	if e != nil {
		return nil, "", errors.New("The app is unavailable; no successful action has been confirmed")
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, (1<<20)+1))
	if e != nil || len(raw) > 1<<20 {
		return nil, "", errors.New("App response is too large; use a narrower query or page")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var problem struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(raw, &problem)
		if problem.Error != "" && len(problem.Error) < 500 {
			return nil, "", &aiHTTPError{Status: response.StatusCode, Message: fmt.Sprintf("%s: %s", app.Name, problem.Error)}
		}
		return nil, "", &aiHTTPError{Status: response.StatusCode, Message: fmt.Sprintf("%s returned HTTP%d; no successful action has been confirmed", app.Name, response.StatusCode)}
	}
	return raw, response.Header.Get("Content-Type"), nil
}

func (p *Portal) aiNativeIdentity(ctx context.Context, browser session, app catalog.App) (aiAppConnection, string, error) {
	connection := browser.AIApps[app.ID]
	if len(connection.Cookies) == 0 {
		return connection, "", fmt.Errorf("Open %s once to connect your account, then retry", app.Name)
	}
	probe := "/api/me"
	if app.ID == "drive" {
		probe = "/api/users?username=self"
	}
	if app.ID == "calendar" || app.ID == "contacts" || app.ID == "mail" {
		probe = "/api/session"
	}
	if app.ID == "music" {
		probe = "/auth/user"
	}
	raw, _, e := p.aiNativeRequest(ctx, app, connection, "", "GET", probe, nil, "")
	if e != nil {
		return connection, "", e
	}
	var me struct {
		User                            struct{ ID, Issuer, Subject string }
		Issuer, Subject, Username, CSRF string
		LoginMethod                     string `json:"loginMethod"`
		Extra                           struct{ WorkspaceIssuer, WorkspaceSubject string }
	}
	if json.Unmarshal(raw, &me) != nil {
		return connection, "", errors.New("The app did not return a verifiable account")
	}
	issuer, subject := me.User.Issuer, me.User.Subject
	switch app.ID {
	case "notes":
		identity, _, err := p.aiNativeRequest(ctx, app, connection, "", "GET", "/api/identity", nil, "")
		if err != nil || json.Unmarshal(identity, &me) != nil {
			return connection, "", errors.New("Notes identity is unavailable")
		}
		issuer, subject = me.Issuer, me.Subject
	case "my-tube":
		// Tube intentionally omits OIDC subjects from public user JSON. Only its
		// protected local ID-to-issuer/subject mapping is read, never sessions or content.
		var err error
		issuer, subject, err = p.aiTubeIdentity(ctx, me.User.ID)
		if err != nil {
			return connection, "", err
		}
	case "drive":
		if me.LoginMethod == "oidc" {
			issuer, subject = p.rt.Config.Issuer, me.Username
		}
	case "calendar", "contacts", "mail":
		issuer, subject = p.rt.Config.Issuer, me.Subject
	case "music":
		issuer, subject = me.Extra.WorkspaceIssuer, me.Extra.WorkspaceSubject
	}
	if issuer != p.rt.Config.Issuer || subject != browser.Subject {
		return connection, "", errors.New("This app is signed in to another account; reconnect it from Workspace")
	}
	if app.ID == "music" {
		raw, _, e = p.aiNativeRequest(ctx, app, connection, "", "GET", "/studio/bootstrap", nil, "")
		var bootstrap struct{ CSRF string }
		if e != nil || json.Unmarshal(raw, &bootstrap) != nil || bootstrap.CSRF == "" {
			return connection, "", errors.New("Open My Music to reconnect its library")
		}
		me.CSRF = bootstrap.CSRF
	}
	return connection, me.CSRF, nil
}

func (p *Portal) aiTubeIdentity(ctx context.Context, id string) (string, string, error) {
	if !aiParameter.MatchString(id) {
		return "", "", errors.New("My Tube identity is unavailable")
	}
	file := filepath.Join(p.rt.Config.Data, "my-tube", "mytube.db")
	if e := myaiSafeTreePath(file); e != nil {
		return "", "", e
	}
	dbPath := filepath.ToSlash(file)
	if filepath.VolumeName(file) != "" {
		dbPath = "/" + dbPath
	}
	u := url.URL{Scheme: "file", Path: dbPath}
	db, e := sql.Open("sqlite", u.String()+"?mode=ro&_pragma=trusted_schema(0)")
	if e != nil {
		return "", "", errors.New("My Tube identity is unavailable")
	}
	defer db.Close()
	var issuer, subject string
	e = db.QueryRowContext(ctx, "SELECT issuer,subject FROM users WHERE id=? AND disabled=0", id).Scan(&issuer, &subject)
	if e != nil {
		return "", "", errors.New("My Tube identity is unavailable")
	}
	return issuer, subject, nil
}

func aiAppResult(raw []byte, kind string) (any, error) {
	var value any
	if strings.Contains(kind, "json") {
		if json.Unmarshal(raw, &value) != nil {
			return nil, errors.New("App returned unreadable JSON")
		}
		return aiRedact(value), nil
	}
	if !strings.Contains(kind, "text/") || !utf8.Valid(raw) || bytes.ContainsRune(raw, 0) {
		return nil, errors.New("This result is a binary file. Open the app to view or download it")
	}
	text := string(raw)
	truncated := false
	if len(text) > 24000 {
		text = text[:24000]
		truncated = true
	}
	return map[string]any{"text": text, "truncated": truncated, "untrusted_reference": true}, nil
}
func aiRedact(value any) any {
	switch x := value.(type) {
	case map[string]any:
		for k, v := range x {
			key := strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(k))
			if key == "csrf" || key == "password" || key == "token" || key == "accesstoken" || key == "refreshtoken" || key == "clientsecret" || key == "authorization" || key == "cookie" || key == "secret" {
				delete(x, k)
			} else {
				x[k] = aiRedact(v)
			}
		}
	case []any:
		for i := range x {
			x[i] = aiRedact(x[i])
		}
	}
	return value
}

// A journal ID is useful for idempotent, account-scoped deployment handoffs.
func aiRequestDigest(subject string, call aiAppCall) string {
	call.Reviewed = false
	call.Prepare = false
	call.RequestID = ""
	if call.Query == nil {
		call.Query = map[string]string{}
	}
	if call.Parameters == nil {
		call.Parameters = map[string]string{}
	}
	if len(call.Body) > 0 {
		var body any
		if json.Unmarshal(call.Body, &body) == nil {
			call.Body, _ = json.Marshal(body)
		}
	}
	raw, _ := json.Marshal(call)
	sum := sha256.Sum256(append([]byte(subject+"\x00"), raw...))
	return hex.EncodeToString(sum[:])
}

func (p *Portal) aiAppReview(ctx context.Context, browser session, uid string, app catalog.App, op aiOperation, call aiAppCall) (any, error) {
	review := map[string]any{"app": app.Name, "title": op.Title, "parameters": call.Parameters, "body": json.RawMessage(call.Body)}
	if len(call.Body) == 0 {
		delete(review, "body")
	}
	if call.App == "gitea" {
		project, e := p.rt.myaiProjectLoad(uid, call.Parameters["project"])
		if e != nil {
			return nil, e
		}
		review["project"] = project.Name
		review["repository"] = project.URL
		return review, nil
	}
	connection, csrf, e := p.aiNativeIdentity(ctx, browser, app)
	if e != nil {
		return nil, &aiConnectionRequired{e}
	}
	if app.ID == "tasks" && (op.ID == "create_task" || op.ID == "update_task") {
		var body struct {
			ListID string `json:"listId"`
		}
		_ = json.Unmarshal(call.Body, &body)
		review["list"] = "Inbox"
		if body.ListID != "" {
			raw, _, err := p.aiNativeRequest(ctx, app, connection, csrf, "GET", "/api/lists", nil, "")
			if err != nil {
				return nil, err
			}
			var lists []struct{ ID, Name, Role string }
			if json.Unmarshal(raw, &lists) != nil {
				return nil, errors.New("Task lists are unreadable")
			}
			found := false
			for _, list := range lists {
				if list.ID == body.ListID {
					found = true
					review["list"] = list.Name
					if list.Role == "viewer" {
						return nil, errors.New("This task list is read-only for your account")
					}
					break
				}
			}
			if !found {
				return nil, errors.New("Use an exact listId returned by lists, or omit listId to create in your Inbox")
			}
		}
	}
	if op.ID == "deploy_project" {
		project, e := p.rt.myaiProjectLoad(uid, call.Parameters["project"])
		if e != nil {
			return nil, e
		}
		raw, _, e := p.aiNativeRequest(ctx, app, connection, csrf, "GET", "/api/projects/"+call.Parameters["site"], nil, "")
		if e != nil {
			return nil, e
		}
		var site struct{ Name, Slug string }
		if json.Unmarshal(raw, &site) != nil {
			return nil, errors.New("My Hosting site details are unreadable")
		}
		review["project"] = project.Name
		review["repository"] = project.URL
		review["site"] = site.Name
		review["siteName"] = site.Slug
	}
	return review, nil
}
