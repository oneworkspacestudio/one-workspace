package suite

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestFroststockAndHostingCanSharePort(t *testing.T) {
	hosting, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = hosting.Close() })
	port := hosting.Addr().(*net.TCPAddr).Port

	app := catalog.App{
		ID: "froststock", Kind: "process", Port: port,
		Route: "/froststock", StripPrefix: true, Health: "/healthz",
		Args: []string{"--listen", "127.0.0.1:${PORT}"},
	}
	request := Request{Port: port - 58, Mode: "local"}
	c := NewConfig(request, catalog.Catalog{Apps: []catalog.App{app}})
	if err := CheckPorts(c); err != nil {
		t.Fatal("fresh installation with Caddy on Froststock's assigned port:", err)
	}
	// A preview.130 installation stores the same numeric port and listen flag.
	c.Origin = "https://workspace.example.test"

	rt := &Runtime{Config: c}
	address := fmt.Sprintf("[::1]:%d", port)
	if !slices.Contains(rt.command(app).Args, address) {
		t.Fatalf("saved Froststock launch didn't use %s", address)
	}

	listener, err := net.Listen("tcp6", address)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		l    net.Listener
	}{
		{"My Hosting", hosting},
		{"Froststock", listener},
	} {
		server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s %s", tc.name, r.URL.Path)
		})}
		t.Cleanup(func() { _ = server.Close() })
		go server.Serve(tc.l)
	}

	p := &Portal{rt: rt, sessions: map[string]session{
		"test": {Subject: "owner", Expires: time.Now().Add(time.Hour)},
	}}
	req := httptest.NewRequest(http.MethodGet, c.Origin+"/froststock/", nil)
	req.AddCookie(&http.Cookie{Name: "workspace_session", Value: "test"})
	response := httptest.NewRecorder()
	p.Handler("control").ServeHTTP(response, req)
	if response.Code != http.StatusOK || response.Body.String() != "Froststock /" {
		t.Fatalf("Froststock route returned %d %q", response.Code, response.Body.String())
	}

	if err := CheckPorts(c); err == nil {
		t.Fatal("another IPv6 listener must reserve Froststock's address")
	}
}
