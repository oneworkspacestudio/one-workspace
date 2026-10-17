package suite

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/fahimc/one-workspace/internal/platform"
)

// Configure only upstream OIDC authentication. The vault has its own origin,
// master password and independent login; no portal cookie or key is forwarded.
func (rt *Runtime) prepareVaultIdentity(ctx context.Context) error {
	app, ok := rt.Config.Catalog.Find("vault")
	if !ok {
		return nil
	}
	if _, ok = rt.Config.Catalog.Find("accounts"); !ok {
		return nil
	}
	root := filepath.Join(rt.Config.Data, "my-vault")
	if e := os.MkdirAll(root, 0700); e != nil {
		return e
	}
	if e := platform.ProtectDir(root); e != nil {
		return e
	}
	path := filepath.Join(root, "myvault-sso.json")
	var existing struct {
		Enabled bool `json:"enabled"`
	}
	if e := platform.LoadJSON(path, &existing); e == nil && !existing.Enabled {
		return nil // Local owner's explicit opt-out persists across suite starts.
	} else if e != nil && !os.IsNotExist(e) {
		return errors.New("vault identity configuration needs local repair")
	}
	var settings struct {
		Origin string `json:"origin"`
	}
	e := platform.LoadJSON(filepath.Join(root, "myvault.json"), &settings)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if e = platform.LoadJSON(filepath.Join(root, "myvault-pending.json"), &settings); e != nil && !os.IsNotExist(e) {
		return e
	}
	if settings.Origin == "" {
		settings.Origin = fmt.Sprintf("https://localhost:%d", app.Port+1)
	}
	u, e := url.Parse(settings.Origin)
	workspace, _ := url.Parse(rt.Config.Origin)
	if e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || strings.ContainsAny(settings.Origin, "\r\n\\") || sameVaultOrigin(u, workspace) {
		return errors.New("vault needs its own HTTPS address before connecting MyAccounts")
	}
	issuer, e := url.Parse(rt.Config.Issuer)
	if e != nil || issuer.Host == "" || issuer.User != nil || issuer.RawQuery != "" || issuer.Fragment != "" || (issuer.Scheme != "https" && !(issuer.Scheme == "http" && localVaultHostname(issuer.Hostname()))) {
		return errors.New("MyAccounts needs HTTPS or a local loopback address")
	}
	api := newIdentity(rt.Config)
	if e = api.login(ctx, rt.Secrets); e != nil {
		return e
	}
	const id = "admin/workspace-vault"
	var identityApp map[string]any
	if e = api.call(ctx, "/api/get-application?id="+id, nil, &identityApp); e != nil {
		return e
	}
	callback := strings.TrimRight(settings.Origin, "/") + "/identity/connect/oidc-signin"
	if identityApp == nil {
		if e = api.call(ctx, "/api/get-application?id=admin/myaccounts", nil, &identityApp); e != nil {
			return e
		}
		if identityApp == nil {
			return errors.New("MyAccounts application template is unavailable")
		}
		identityApp["name"] = "workspace-vault"
		identityApp["clientId"] = "workspace-vault"
		identityApp["clientSecret"] = platform.Token()
		identityApp["displayName"] = "My Vault"
		identityApp["redirectUris"] = []string{callback}
		identityApp["grantTypes"] = []string{"authorization_code"}
		identityApp["tokenFormat"] = "JWT-Standard"
		identityApp["enableSignUp"] = false
		identityApp["disableSelfSignup"] = true
		identityApp["enableAutoSignin"] = true
		identityApp["enableSigninSession"] = true
		identityApp["providers"] = []any{}
		identityApp["scopes"] = []any{}
		if e = api.call(ctx, "/api/add-application", identityApp, nil); e != nil {
			return e
		}
	} else {
		if identityApp["clientId"] != "workspace-vault" || identityApp["name"] != "workspace-vault" {
			return errors.New("My Vault identity client did not match")
		}
		identityApp["redirectUris"] = []string{callback}
		identityApp["tokenFormat"] = "JWT-Standard"
		identityApp["enableAutoSignin"] = true
		identityApp["enableSigninSession"] = true
		if e = api.call(ctx, "/api/update-application?id="+id, identityApp, nil); e != nil {
			return e
		}
	}
	secret, ok := identityApp["clientSecret"].(string)
	if !ok || secret == "" {
		return errors.New("My Vault identity client has no secret")
	}
	return platform.SaveJSON(path, map[string]any{"enabled": true, "issuer": rt.Config.Issuer, "clientId": "workspace-vault", "clientSecret": secret, "provisionAccounts": true})
}
