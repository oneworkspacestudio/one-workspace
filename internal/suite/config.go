package suite

import (
	"errors"
	"fmt"
	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/platform"
	"net"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	Schema    int             `json:"schema"`
	Version   string          `json:"version"`
	Programs  string          `json:"programs"`
	Data      string          `json:"data"`
	Files     string          `json:"files"`
	Origin    string          `json:"origin"`
	Issuer    string          `json:"issuer"`
	Port      int             `json:"port"`
	Owner     string          `json:"owner"`
	Autostart bool            `json:"autostart"`
	Tunnel    bool            `json:"tunnel"`
	Catalog   catalog.Catalog `json:"catalog"`
}
type Client struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}
type Secrets struct {
	AdminEmail       string            `json:"adminEmail"`
	AdminPassword    string            `json:"adminPassword"`
	RecoveryPassword string            `json:"recoveryPassword"`
	Clients          map[string]Client `json:"clients"`
	TunnelToken      string            `json:"tunnelToken,omitempty"`
}
type Request struct {
	Programs  string `json:"programs"`
	Data      string `json:"data"`
	Files     string `json:"files"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Confirm   string `json:"confirm"`
	Port      int    `json:"port"`
	Mode      string `json:"mode"`
	Origin    string `json:"origin"`
	Issuer    string `json:"issuer"`
	Photos    bool   `json:"photos"`
	Autostart bool   `json:"autostart"`
	CFToken   string `json:"cfToken"`
	AccountID string `json:"accountId"`
	ZoneID    string `json:"zoneId"`
	Authorize bool   `json:"authorize"`
}

func (r Request) Validate() error {
	if r.Name == "" || len(r.Name) > 100 {
		return errors.New("enter your name")
	}
	email, e := mail.ParseAddress(r.Email)
	if e != nil || email.Address != r.Email {
		return errors.New("enter a valid email address")
	}
	if len(r.Password) < 12 || len(r.Password) > 128 || r.Password != r.Confirm {
		return errors.New("use matching passwords of 12 to 128 characters")
	}
	if r.Port < 1024 || r.Port > 65530 {
		return errors.New("choose a base port between 1024 and 65530")
	}
	if r.Mode != "local" && r.Mode != "manual" && r.Mode != "cloudflare" {
		return errors.New("choose local, manual HTTPS or Cloudflare")
	}
	for _, p := range []string{r.Programs, r.Data, r.Files} {
		if !filepath.IsAbs(p) || filepath.Dir(filepath.Clean(p)) == filepath.Clean(p) {
			return errors.New("choose dedicated absolute folders rather than drive roots")
		}
	}
	a, b, c := strings.ToLower(filepath.Clean(r.Programs)), strings.ToLower(filepath.Clean(r.Data)), strings.ToLower(filepath.Clean(r.Files))
	if a == b || a == c || b == c || strings.HasPrefix(a, b+string(os.PathSeparator)) || strings.HasPrefix(b, a+string(os.PathSeparator)) {
		return errors.New("program, private data and files folders must be separate")
	}
	for _, pair := range [][2]string{{a, c}, {b, c}, {c, a}, {c, b}} {
		if strings.HasPrefix(pair[0], pair[1]+string(os.PathSeparator)) {
			return errors.New("program, private data and files folders must not contain each other")
		}
	}
	if r.Mode != "local" {
		for _, raw := range []string{r.Origin, r.Issuer} {
			u, e := url.Parse(raw)
			if e != nil || u.Scheme != "https" || u.Hostname() == "" || u.Port() != "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
				return errors.New("public addresses must be HTTPS origins without paths or ports")
			}
		}
		if strings.EqualFold(strings.TrimRight(r.Origin, "/"), strings.TrimRight(r.Issuer, "/")) {
			return errors.New("Workspace and Accounts need distinct hostnames")
		}
	}
	if r.Mode == "cloudflare" && (!r.Authorize || r.CFToken == "" || !regexp.MustCompile(`^[a-fA-F0-9]{32}$`).MatchString(r.AccountID) || !regexp.MustCompile(`^[a-fA-F0-9]{32}$`).MatchString(r.ZoneID)) {
		return errors.New("authorize a dedicated tunnel and supply its scoped Cloudflare account, zone and API token")
	}
	return nil
}
func NewConfig(r Request, c catalog.Catalog) Config {
	conf := Config{Schema: 1, Version: c.Version, Programs: filepath.Clean(r.Programs), Data: filepath.Clean(r.Data), Files: filepath.Clean(r.Files), Port: r.Port, Autostart: r.Autostart, Tunnel: r.Mode == "cloudflare", Catalog: c}
	conf.Origin = fmt.Sprintf("http://localhost:%d", r.Port)
	conf.Issuer = fmt.Sprintf("http://localhost:%d", r.Port+1)
	if r.Mode != "local" {
		conf.Origin = strings.ToLower(strings.TrimRight(r.Origin, "/"))
		conf.Issuer = strings.ToLower(strings.TrimRight(r.Issuer, "/"))
	}
	apps := []catalog.App{}
	for _, app := range conf.Catalog.Apps {
		switch app.ID {
		case "accounts":
			app.Port = r.Port + 1
		case "drive":
			app.Port = r.Port + 2
		case "photos":
			app.Port = r.Port + 3
			if !r.Photos {
				continue
			}
		case "gitea":
			app.Port = r.Port + 4
		case "music":
			app.Port = r.Port + 5
		case "mail":
			app.Port = r.Port + 6
		}
		apps = append(apps, app)
	}
	conf.Catalog.Apps = apps
	return conf
}
func Load(data string) (Config, Secrets, error) {
	var c Config
	var s Secrets
	e := platform.LoadJSON(filepath.Join(data, "workspace.json"), &c)
	if e != nil {
		return c, s, e
	}
	if c.Schema != 1 {
		return c, s, errors.New("unsupported workspace state")
	}
	if e = c.Catalog.Validate(); e != nil {
		return c, s, e
	}
	e = platform.LoadSecret(filepath.Join(data, "secrets.dpapi"), &s)
	return c, s, e
}
func Save(c Config, s Secrets) error {
	if e := platform.SaveSecret(filepath.Join(c.Data, "secrets.dpapi"), s); e != nil {
		return e
	}
	return platform.SaveJSON(filepath.Join(c.Data, "workspace.json"), c)
}
func CheckPorts(c Config) error {
	ports := []int{c.Port}
	for _, a := range c.Catalog.Apps {
		if a.Kind == "process" {
			ports = append(ports, a.Port)
			if a.ID == "mail" {
				ports = append(ports, a.Port+1, a.Port+2)
			}
		}
	}
	listeners := []net.Listener{}
	defer func() {
		for _, l := range listeners {
			l.Close()
		}
	}()
	for _, p := range ports {
		if p > 65535 {
			return errors.New("Workspace port range exceeds 65535; choose a lower base port")
		}
		l, e := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", p))
		if e != nil {
			return fmt.Errorf("port %d is already in use", p)
		}
		listeners = append(listeners, l)
	}
	return nil
}
