package apphost

import (
	"context"
	"crypto/hmac"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Equal(a, b string) bool { return a != "" && b != "" && hmac.Equal([]byte(a), []byte(b)) }
func (h *Host) cookie(w http.ResponseWriter, name, value string, age int) {
	http.SetCookie(w, &http.Cookie{Name: h.Config.ID + "_" + name, Value: value, Path: h.Base + "/", HttpOnly: true, Secure: strings.HasPrefix(h.Config.Origin, "https://"), SameSite: http.SameSiteLaxMode, MaxAge: age})
}
func (h *Host) sameOrigin(r *http.Request) bool {
	u, _ := url.Parse(h.Config.Origin)
	return r.Header.Get("Origin") == u.Scheme+"://"+u.Host
}
func (h *Host) session(r *http.Request) (s Session, ok bool) {
	c, e := r.Cookie(h.Config.ID + "_session")
	if e != nil || len(c.Value) != 96 {
		return s, false
	}
	var expiry string
	e = h.DB.QueryRow(`SELECT u.id,u.issuer,u.subject,u.name,u.email,u.role,s.csrf,s.expires FROM host_sessions s JOIN host_users u ON u.id=s.user_id WHERE s.id=?`, Hash(c.Value)).Scan(&s.User.ID, &s.User.Issuer, &s.User.Subject, &s.User.Name, &s.User.Email, &s.User.Role, &s.CSRF, &expiry)
	if e != nil || s.User.Role == "disabled" {
		return s, false
	}
	s.Expires, e = time.Parse(time.RFC3339Nano, expiry)
	if e != nil || !time.Now().Before(s.Expires) {
		return s, false
	}
	if h.Config.Issuer != "" {
		if s.User.Issuer != h.Config.Issuer {
			return s, false
		}
		s.User.Role = "member"
		if s.User.Subject == h.Config.Owner && h.Config.Owner != "" {
			s.User.Role = "owner"
		}
	} else if s.User.Issuer != "local:"+h.Config.ID {
		return s, false
	}
	return s, true
}
func (h *Host) createSession(w http.ResponseWriter, u User, expires time.Time) error {
	token, csrf := ID()+ID(), ID()+ID()
	_, e := h.DB.Exec(`INSERT INTO host_sessions(id,user_id,csrf,expires)VALUES(?,?,?,?)`, Hash(token), u.ID, csrf, expires.UTC().Format(time.RFC3339Nano))
	if e == nil {
		h.cookie(w, "session", token, int(time.Until(expires).Seconds()))
		h.Audit(u.ID, "sign-in", h.Config.ID)
	}
	return e
}
func (h *Host) discover(ctx context.Context) error {
	h.authMu.Lock()
	defer h.authMu.Unlock()
	if h.verifier != nil {
		return nil
	}
	ctx = oidc.ClientContext(ctx, &http.Client{Timeout: 20 * time.Second})
	p, e := oidc.NewProvider(ctx, h.Config.Issuer)
	if e != nil {
		return errors.New("MyAccounts is unavailable; retry when your sign-in server is running")
	}
	endpoint := p.Endpoint()
	issuer, _ := url.Parse(h.Config.Issuer)
	for _, raw := range []string{endpoint.AuthURL, endpoint.TokenURL} {
		u, e := url.Parse(raw)
		if e != nil || u.Scheme != issuer.Scheme || u.Host != issuer.Host {
			return errors.New("MyAccounts endpoints must use the configured issuer's origin")
		}
	}
	h.oauth = oauth2.Config{ClientID: h.Config.ClientID, ClientSecret: h.Config.ClientSecret, Endpoint: endpoint, RedirectURL: h.Config.Origin + "/auth/callback", Scopes: []string{oidc.ScopeOpenID, "profile", "email"}}
	h.verifier = p.Verifier(&oidc.Config{ClientID: h.Config.ClientID})
	return nil
}
func (h *Host) login(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.session(r); ok {
		http.Redirect(w, r, h.Base+"/", 303)
		return
	}
	if h.Config.Issuer == "" {
		http.Redirect(w, r, h.Base+"/#login", 303)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	if e := h.discover(ctx); e != nil {
		Problem(w, 503, e.Error())
		return
	}
	state, nonce, v := ID()+ID(), ID(), oauth2.GenerateVerifier()
	h.authMu.Lock()
	for id, f := range h.flows {
		if time.Now().After(f.Expires) {
			delete(h.flows, id)
		}
	}
	if len(h.flows) >= 500 {
		h.authMu.Unlock()
		Problem(w, 429, "Too many sign-in attempts; try shortly")
		return
	}
	h.flows[state] = flow{Verifier: v, Nonce: nonce, Return: h.Base + "/", Expires: time.Now().Add(5 * time.Minute)}
	h.authMu.Unlock()
	h.cookie(w, "login", state, 300)
	http.Redirect(w, r, h.oauth.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.S256ChallengeOption(v)), 302)
}
func (h *Host) callback(w http.ResponseWriter, r *http.Request) {
	bad := func() { Problem(w, 401, "Sign-in did not finish. Reopen this app from Workspace to reconnect.") }
	state := r.URL.Query().Get("state")
	cookie, e := r.Cookie(h.Config.ID + "_login")
	if e != nil || !Equal(cookie.Value, state) {
		bad()
		return
	}
	h.authMu.Lock()
	f, ok := h.flows[state]
	delete(h.flows, state)
	h.authMu.Unlock()
	h.cookie(w, "login", "", -1)
	if !ok || time.Now().After(f.Expires) || r.URL.Query().Get("error") != "" || h.verifier == nil {
		bad()
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	ctx = oidc.ClientContext(ctx, &http.Client{Timeout: 20 * time.Second})
	tok, e := h.oauth.Exchange(ctx, r.URL.Query().Get("code"), oauth2.VerifierOption(f.Verifier))
	if e != nil {
		bad()
		return
	}
	raw, _ := tok.Extra("id_token").(string)
	identity, e := h.verifier.Verify(ctx, raw)
	if e != nil || identity.Issuer != h.Config.Issuer || identity.Subject == "" || identity.Nonce != f.Nonce {
		bad()
		return
	}
	var claims struct {
		Name     string `json:"name"`
		Display  string `json:"displayName"`
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if identity.Claims(&claims) != nil {
		bad()
		return
	}
	if claims.Display != "" {
		claims.Name = claims.Display
	}
	if claims.Name == "" {
		claims.Name = "Workspace member"
	}
	if !claims.Verified {
		claims.Email = ""
	}
	u := User{ID: Identity(identity.Issuer, identity.Subject), Issuer: identity.Issuer, Subject: identity.Subject, Name: claims.Name, Email: claims.Email, Role: "member"}
	if u.Subject == h.Config.Owner && h.Config.Owner != "" {
		u.Role = "owner"
	}
	h.Freeze.RLock()
	defer h.Freeze.RUnlock()
	_, e = h.DB.Exec(`INSERT INTO host_users(id,issuer,subject,name,email,role)VALUES(?,?,?,?,?,?) ON CONFLICT(issuer,subject)DO UPDATE SET name=excluded.name,email=excluded.email`, u.ID, u.Issuer, u.Subject, u.Name, u.Email, u.Role)
	if e != nil {
		h.DBError(w, e)
		return
	}
	var role string
	if e = h.DB.QueryRow("SELECT role FROM host_users WHERE id=?", u.ID).Scan(&role); e != nil || role == "disabled" {
		bad()
		return
	}
	expires := time.Now().Add(8 * time.Hour)
	if identity.Expiry.Before(expires) {
		expires = identity.Expiry
	}
	if e = h.createSession(w, u, expires); e != nil {
		h.DBError(w, e)
		return
	}
	http.Redirect(w, r, f.Return, 303)
}
func (h *Host) logout(w http.ResponseWriter, r *http.Request) {
	if !h.sameOrigin(r) {
		Problem(w, 403, "Sign out from this app or Workspace")
		return
	}
	if c, e := r.Cookie(h.Config.ID + "_session"); e == nil {
		_, _ = h.DB.Exec("DELETE FROM host_sessions WHERE id=?", Hash(c.Value))
	}
	h.cookie(w, "session", "", -1)
	h.cookie(w, "login", "", -1)
	JSON(w, 200, map[string]bool{"ok": true})
}
func (h *Host) permittedLocal(r *http.Request) bool {
	host, _, e := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(host)
	return e == nil && ip != nil && ip.IsLoopback() && h.Config.Issuer == "" && h.sameOrigin(r)
}
func (h *Host) rateLimit(r *http.Request) bool {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	h.authMu.Lock()
	defer h.authMu.Unlock()
	now := time.Now()
	for key, v := range h.attempts {
		if now.Sub(v.Since) > 10*time.Minute {
			delete(h.attempts, key)
		}
	}
	a := h.attempts[host]
	if a.Since.IsZero() {
		a.Since = now
	}
	a.Count++
	h.attempts[host] = a
	return a.Count > 20 || len(h.attempts) > 5000
}
func (h *Host) Bootstrap(email, password string) error {
	if h.Config.Issuer != "" {
		return errors.New("standalone accounts are disabled in managed mode")
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if !strings.Contains(email, "@") || len(email) > 254 || len(password) < 8 || len(password) > 72 {
		return errors.New("use a valid email and an 8–72 byte password")
	}
	h.authMu.Lock()
	defer h.authMu.Unlock()
	var count int
	if e := h.DB.QueryRow("SELECT count(*) FROM host_users").Scan(&count); e != nil {
		return e
	}
	if count != 0 {
		return errors.New("the first owner is already configured")
	}
	_, e := h.addLocal(email, password, "owner")
	if e == nil {
		h.setup = ""
	}
	return e
}
func (h *Host) addLocal(email, password, role string) (User, error) {
	issuer := "local:" + h.Config.ID
	u := User{ID: Identity(issuer, email), Issuer: issuer, Subject: email, Email: email, Name: strings.Split(email, "@")[0], Role: role}
	hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return u, e
	}
	_, e = h.DB.Exec("INSERT INTO host_users(id,issuer,subject,name,email,role,password)VALUES(?,?,?,?,?,?,?)", u.ID, u.Issuer, u.Subject, u.Name, u.Email, u.Role, hash)
	return u, e
}
func (h *Host) setupUser(w http.ResponseWriter, r *http.Request) {
	if !h.permittedLocal(r) || h.rateLimit(r) {
		Problem(w, 403, "First setup requires this PC and its one-time setup link")
		return
	}
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !Decode(w, r, &in) {
		return
	}
	h.authMu.Lock()
	valid := Equal(r.Header.Get("X-Setup-Token"), h.setup)
	h.authMu.Unlock()
	if !valid {
		Problem(w, 403, "Use the one-time setup link printed by this app")
		return
	}
	if e := h.Bootstrap(in.Email, in.Password); e != nil {
		Problem(w, 400, e.Error())
		return
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))
	issuer := "local:" + h.Config.ID
	u := User{ID: Identity(issuer, email), Issuer: issuer, Subject: email, Email: email, Name: strings.Split(email, "@")[0], Role: "owner"}
	if e := h.createSession(w, u, time.Now().Add(8*time.Hour)); e != nil {
		h.DBError(w, e)
		return
	}
	JSON(w, 201, map[string]bool{"ok": true})
}
func (h *Host) localLogin(w http.ResponseWriter, r *http.Request) {
	if !h.permittedLocal(r) {
		Problem(w, 403, "Continue with MyAccounts for this managed installation")
		return
	}
	if h.rateLimit(r) {
		Problem(w, 429, "Too many attempts; wait ten minutes")
		return
	}
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !Decode(w, r, &in) {
		return
	}
	var u User
	var password []byte
	e := h.DB.QueryRow("SELECT id,issuer,subject,name,email,role,password FROM host_users WHERE issuer=? AND subject=?", "local:"+h.Config.ID, strings.ToLower(strings.TrimSpace(in.Email))).Scan(&u.ID, &u.Issuer, &u.Subject, &u.Name, &u.Email, &u.Role, &password)
	dummy := []byte("$2a$10$e5BwUaTqfYBdnxWFBCU9cO/8neA5I5.XoJiG41q2HZWVx3.C.KKZS")
	if e != nil {
		password = dummy
	}
	if bcrypt.CompareHashAndPassword(password, []byte(in.Password)) != nil || e != nil || u.Role == "disabled" {
		Problem(w, 401, "Email or password is incorrect")
		return
	}
	if e = h.createSession(w, u, time.Now().Add(8*time.Hour)); e != nil {
		h.DBError(w, e)
		return
	}
	JSON(w, 200, map[string]bool{"ok": true})
}
func (h *Host) localUser(w http.ResponseWriter, r *http.Request) {
	if h.Config.Issuer != "" {
		Problem(w, 403, "Manage accounts through MyAccounts")
		return
	}
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !Decode(w, r, &in) {
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if !strings.Contains(in.Email, "@") || len(in.Email) > 254 || len(in.Password) < 8 || len(in.Password) > 72 {
		Problem(w, 400, "Use an email and an 8–72 byte password")
		return
	}
	u, e := h.addLocal(in.Email, in.Password, "member")
	if e != nil {
		Problem(w, 409, "That account could not be created; check whether it already exists")
		return
	}
	JSON(w, 201, u)
}
