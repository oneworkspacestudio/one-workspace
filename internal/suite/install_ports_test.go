package suite

import (
	"fmt"
	"net"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func installPortFixture(t *testing.T) (Request, catalog.Catalog) {
	t.Helper()
	cat := catalog.Catalog{Apps: []catalog.App{
		{ID: "accounts", Kind: "process"}, {ID: "monitor", Kind: "process"},
		{ID: "mail", Kind: "process"}, {ID: "vault", Kind: "process"},
		{ID: "meet", Kind: "process"}, {ID: "create", Kind: "process"},
		{ID: "remote", Kind: "process"},
	}}
	for port := 32000; port < 40000; port += 32 {
		r := Request{Port: port, Mode: "local"}
		if CheckPorts(NewConfig(r, cat)) == nil {
			return r, cat
		}
	}
	t.Fatal("no free fixture port group")
	return Request{}, cat
}

func TestFreshInstallRetainsAvailablePorts(t *testing.T) {
	r, cat := installPortFixture(t)
	c, err := availableInstallConfig(r, cat)
	if err != nil || c.Port != r.Port {
		t.Fatalf("available requested port changed: %d, %v", c.Port, err)
	}
}

func TestFreshInstallMovesWholeGroupOnCollision(t *testing.T) {
	for _, tc := range []struct {
		name   string
		offset int
		udp    bool
	}{{"portal", 0, false}, {"accounts", 1, false}, {"mail engine", 7, false},
		{"vault engine", 15, false}, {"Remote broker", 136, false}, {"Remote host", 137, false}, {"Meet media TCP", 1026, false}, {"Meet media UDP", 1027, true}} {
		t.Run(tc.name, func(t *testing.T) {
			r, cat := installPortFixture(t)
			address := fmt.Sprintf("127.0.0.1:%d", r.Port+tc.offset)
			if tc.udp {
				address = fmt.Sprintf("0.0.0.0:%d", r.Port+tc.offset)
				l, err := net.ListenPacket("udp4", address)
				if err != nil {
					t.Fatal(err)
				}
				defer l.Close()
			} else {
				l, err := net.Listen("tcp4", address)
				if err != nil {
					t.Fatal(err)
				}
				defer l.Close()
			}
			c, err := availableInstallConfig(r, cat)
			if err != nil {
				t.Fatal(err)
			}
			if c.Port == r.Port {
				t.Fatal("kept the occupied group")
			}
			if err = CheckPorts(c); err != nil {
				t.Fatal("replacement group unavailable", err)
			}
			if c.Origin != fmt.Sprintf("http://localhost:%d", c.Port) || c.Issuer != fmt.Sprintf("http://localhost:%d", c.Port+1) {
				t.Fatal("local identity addresses use the old port")
			}
			create, _ := c.Catalog.Find("create")
			accounts, _ := c.Catalog.Find("accounts")
			if create.Port != c.Port+27 || accounts.Port != c.Port+1 {
				t.Fatal("app ports did not move together")
			}
		})
	}
}

func TestFreshPortFallbackPreservesPublicIdentityOrigins(t *testing.T) {
	r, cat := installPortFixture(t)
	r.Mode, r.Origin, r.Issuer = "manual", "https://workspace.example.test", "https://accounts.example.test"
	l, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", r.Port))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	c, err := availableInstallConfig(r, cat)
	if err != nil {
		t.Fatal(err)
	}
	if c.Port == r.Port || c.Origin != r.Origin || c.Issuer != r.Issuer {
		t.Fatal("public identity changed during fallback")
	}
}
