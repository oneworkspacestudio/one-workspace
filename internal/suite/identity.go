package suite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fahimc/one-workspace/internal/platform"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type identityAPI struct {
	base   string
	client *http.Client
}

func newIdentity(c Config) *identityAPI {
	a, _ := c.Catalog.Find("accounts")
	jar, _ := cookiejar.New(nil)
	return &identityAPI{fmt.Sprintf("http://localhost:%d", a.Port), &http.Client{Jar: jar, Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}
func (a *identityAPI) call(ctx context.Context, path string, body any, out any) error {
	var data io.Reader
	method := "GET"
	if body != nil {
		b, e := json.Marshal(body)
		if e != nil {
			return e
		}
		data = bytes.NewReader(b)
		method = "POST"
	}
	req, e := http.NewRequestWithContext(ctx, method, a.base+path, data)
	if e != nil {
		return e
	}
	req.Header.Set("Origin", a.base)
	req.Header.Set("Content-Type", "application/json")
	r, e := a.client.Do(req)
	if e != nil {
		return e
	}
	defer r.Body.Close()
	var result struct {
		Status string          `json:"status"`
		Data   json.RawMessage `json:"data"`
	}
	if r.StatusCode != 200 {
		return errors.New("identity API refused the operation")
	}
	if e = json.NewDecoder(io.LimitReader(r.Body, 8<<20)).Decode(&result); e != nil {
		return e
	}
	if result.Status != "ok" {
		return errors.New("identity configuration rejected; check component logs")
	}
	if out != nil {
		return json.Unmarshal(result.Data, out)
	}
	return nil
}
func (a *identityAPI) login(ctx context.Context, s Secrets) error {
	return a.call(ctx, "/api/login", map[string]any{"organization": "built-in", "username": s.AdminEmail, "password": s.AdminPassword, "application": "app-built-in", "type": "login", "signinMethod": "Password"}, nil)
}
func ConfigureIdentity(ctx context.Context, c *Config, s *Secrets, r Request) error {
	api := newIdentity(*c)
	if e := api.login(ctx, *s); e != nil {
		return e
	}
	// Match the fresh Workspace organization's policy to the installer. The
	// recovery administrator retains its generated long password.
	var organization map[string]any
	if e := api.call(ctx, "/api/get-organization?id=admin/myaccounts", nil, &organization); e != nil {
		return e
	}
	if organization == nil {
		return errors.New("Workspace account organization is unavailable")
	}
	organization["passwordOptions"] = []string{"AtLeast8"}
	if e := api.call(ctx, "/api/update-organization?id=admin/myaccounts", organization, nil); e != nil {
		return e
	}
	user := map[string]any{"owner": "myaccounts", "name": "owner", "displayName": r.Name, "email": r.Email, "emailVerified": false, "type": "normal-user", "isAdmin": true, "passwordType": "plain", "password": r.Password, "signupApplication": "myaccounts", "createdTime": time.Now().UTC().Format(time.RFC3339)}
	if e := api.call(ctx, "/api/add-user", user, nil); e != nil {
		return e
	}
	var saved struct {
		ID string `json:"id"`
	}
	if e := api.call(ctx, "/api/get-user?id=myaccounts/owner", nil, &saved); e != nil {
		return e
	}
	if saved.ID == "" {
		return errors.New("workspace owner was not created")
	}
	c.Owner = saved.ID
	for _, entry := range []struct{ id, callback, name string }{{"portal", c.Origin + "/auth/callback", "One Workspace"}, {"drive", c.Origin + "/drive/api/auth/oidc/callback", "Workspace Drive and Photos"}} {
		client, e := registerClient(ctx, api, entry.id, entry.name, entry.callback)
		if e != nil {
			return e
		}
		s.Clients[entry.id] = client
	}
	for _, app := range c.Catalog.Apps {
		if !app.OIDC || app.ID == "drive" || (SharedEnabled(*c) && sharedConsumer(app.ID)) {
			continue
		}
		client, e := registerClient(ctx, api, app.ID, app.Name, c.Origin+app.Callback)
		if e != nil {
			return e
		}
		s.Clients[app.ID] = client
	}
	if SharedEnabled(*c) {
		client, e := registerSharedClient(ctx, api, *c)
		if e != nil {
			return e
		}
		s.Clients["shared"] = client
		for _, app := range c.Catalog.Apps {
			if sharedConsumer(app.ID) {
				s.Clients[app.ID] = client
			}
		}
	}
	return nil
}
func registerClient(ctx context.Context, api *identityAPI, id, name, callback string) (Client, error) {
	return registerClientRedirects(ctx, api, id, name, []string{callback})
}
func registerClientRedirects(ctx context.Context, api *identityAPI, id, name string, callbacks []string) (Client, error) {
	var app map[string]any
	var result Client
	if e := api.call(ctx, "/api/get-application?id=admin/myaccounts", nil, &app); e != nil {
		return result, e
	}
	result = Client{ID: "workspace-" + id, Secret: platform.Token()}
	app["name"] = result.ID
	app["clientId"] = result.ID
	app["clientSecret"] = result.Secret
	app["displayName"] = name
	app["redirectUris"] = callbacks
	app["grantTypes"] = []string{"authorization_code"}
	app["enableSignUp"] = false
	app["disableSelfSignup"] = true
	app["enableAutoSignin"] = true
	app["enableSigninSession"] = true
	app["providers"] = []any{}
	app["scopes"] = []any{}
	return result, api.call(ctx, "/api/add-application", app, nil)
}
