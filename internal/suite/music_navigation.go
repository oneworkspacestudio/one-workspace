package suite

import (
	"bytes"
	_ "embed"
	"errors"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/fahimc/one-workspace/internal/platform"
)

//go:embed assets/music-worker-navigation.js
var musicWorkerNavigation []byte

// Update the approved Music client's worker at the mounted integration boundary.
// Its original cache rules remain after our navigation handler. No component
// executable, account, OIDC callback or private app data is replaced here.
func patchMusicWorker(resp *http.Response, prefix string, strip bool) error {
	if !strip || prefix != "/music" || resp.Request.URL.Path != "/sw.js" || resp.StatusCode != http.StatusOK || !strings.Contains(resp.Header.Get("Content-Type"), "javascript") || (resp.Header.Get("Content-Encoding") != "" && resp.Header.Get("Content-Encoding") != "identity") {
		return nil
	}
	const limit = 4 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	resp.Body.Close()
	if err != nil {
		return err
	}
	if len(body) > limit {
		return errors.New("Music service worker exceeds the integration size limit")
	}
	patched := append(append([]byte{}, musicWorkerNavigation...), '\n')
	patched = append(patched, body...)
	resp.Body = io.NopCloser(bytes.NewReader(patched))
	resp.ContentLength = int64(len(patched))
	resp.Header.Set("Content-Length", strconv.Itoa(len(patched)))
	resp.Header.Set("Service-Worker-Allowed", "/music/")
	resp.Header.Set("Cache-Control", "no-cache")
	resp.Header.Del("Content-Encoding")
	resp.Header.Del("ETag")
	resp.Header.Del("Content-MD5")
	return nil
}

func musicOpening(w http.ResponseWriter, target string) {
	nonce := platform.Token()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'nonce-"+nonce+"'; style-src 'unsafe-inline'; frame-ancestors 'none'; base-uri 'none'")
	_ = musicOpeningPage.Execute(w, map[string]string{"Nonce": nonce, "Target": target})
}

// This page is outside /music/, so an existing Music worker cannot intercept it.
// Only its exact registration is retired. Cached media/assets and every other
// app's worker/storage/cookies remain. The next authenticated page registers the
// patched worker, whose navigation requests always reach native authentication.
var musicOpeningPage = template.Must(template.New("opening-music").Parse(`<!doctype html><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Opening Music · One Workspace</title><style>body{font:17px/1.6 system-ui;color:#202124;background:#f8fafd;max-width:540px;margin:18vh auto;padding:28px}a{color:#0b57d0}</style><h1>Opening Music…</h1><p id="status">Connecting your account.</p><p id="retry" hidden><a href="/_workspace/open/music">Try again</a> · <a href="/">Back to Workspace</a></p><script nonce="{{.Nonce}}">(async()=>{try{if('serviceWorker' in navigator){const scope=new URL('/music/',location.origin).href;const registrations=await navigator.serviceWorker.getRegistrations();await Promise.all(registrations.filter(registration=>registration.scope===scope).map(registration=>registration.unregister()))}location.replace({{.Target}})}catch(error){document.getElementById('status').textContent='Music could not open. Reload this page to try again.';document.getElementById('retry').hidden=false}})();</script>`))
