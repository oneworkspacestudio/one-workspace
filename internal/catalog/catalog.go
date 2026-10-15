package catalog

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fahimc/one-workspace/internal/platform"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Package struct {
	File   string `json:"file"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}
type App struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Kind        string            `json:"kind"`
	Depends     []string          `json:"depends,omitempty"`
	Package     string            `json:"package,omitempty"`
	Executable  string            `json:"executable,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Port        int               `json:"port,omitempty"`
	Route       string            `json:"route"`
	StripPrefix bool              `json:"stripPrefix,omitempty"`
	Health      string            `json:"health,omitempty"`
	Callback    string            `json:"callback,omitempty"`
	OIDC        bool              `json:"oidc,omitempty"`
	Source      string            `json:"source,omitempty"`
}
type Catalog struct {
	Schema   int                `json:"schema"`
	Version  string             `json:"version"`
	Packages map[string]Package `json:"packages"`
	Apps     []App              `json:"apps"`
}

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{0,39}$`)
var versionPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.+_-]{0,79}$`)

func (c Catalog) Validate() error {
	if c.Schema != 1 || c.Version == "" {
		return errors.New("unsupported or unversioned app catalog")
	}
	seen := map[string]bool{}
	routes := map[string]bool{}
	ports := map[int]bool{}
	files := map[string]bool{}
	for _, p := range c.Packages {
		if _, e := platform.Inside("root", p.File); e != nil {
			return e
		}
		if p.Size <= 0 || p.Size > 2<<30 || len(p.SHA256) != 64 {
			return errors.New("invalid package size or checksum")
		}
		if _, e := hex.DecodeString(p.SHA256); e != nil {
			return e
		}
		key := strings.ToLower(p.File)
		if files[key] {
			return errors.New("duplicate package file")
		}
		files[key] = true
	}
	for _, a := range c.Apps {
		if !idPattern.MatchString(a.ID) || seen[a.ID] || a.Name == "" || !versionPattern.MatchString(a.Version) {
			return fmt.Errorf("invalid or duplicate app: %s", a.ID)
		}
		seen[a.ID] = true
		if a.Kind != "process" && a.Kind != "builtin" {
			return errors.New("app kind must be process or builtin")
		}
		if a.Kind == "process" {
			if a.Port < 1024 || a.Port > 65535 || ports[a.Port] {
				return errors.New("process ports must be distinct and unprivileged")
			}
			ports[a.Port] = true
			p, ok := c.Packages[a.Package]
			if !ok {
				return fmt.Errorf("package missing for %s", a.ID)
			}
			if _, e := platform.Inside("root", a.Executable); e != nil {
				return e
			}
			if p.Size <= 0 || p.Size > 2<<30 || len(p.SHA256) != 64 {
				return errors.New("invalid package size or checksum")
			}
			if b, e := hex.DecodeString(p.SHA256); e != nil || len(b) != 32 {
				return errors.New("invalid package SHA256")
			}
			if _, e := platform.Inside("root", p.File); e != nil {
				return e
			}
		}
		if a.Kind == "process" && a.Executable != c.Packages[a.Package].File {
			return errors.New("executable must match the verified package file")
		}
		if a.ID != "accounts" {
			if !strings.HasPrefix(a.Route, "/") || strings.ContainsAny(a.Route, "?#\\") || strings.Contains(a.Route, "..") || a.Route == "/" {
				return errors.New("app route must be a non-root absolute path")
			}
			if a.Kind == "process" {
				route := strings.TrimRight(a.Route, "/")
				for previous := range routes {
					if route == previous || strings.HasPrefix(route, previous+"/") || strings.HasPrefix(previous, route+"/") {
						return errors.New("overlapping app routes")
					}
				}
				routes[route] = true
			}
		}
	}
	_, e := c.Order()
	return e
}
func (c Catalog) Order() ([]App, error) {
	by := map[string]App{}
	for _, a := range c.Apps {
		by[a.ID] = a
	}
	mark := map[string]int{}
	out := []App{}
	var visit func(string) error
	visit = func(id string) error {
		if mark[id] == 1 {
			return errors.New("dependency cycle")
		}
		if mark[id] == 2 {
			return nil
		}
		a, ok := by[id]
		if !ok {
			return fmt.Errorf("missing dependency %s", id)
		}
		mark[id] = 1
		for _, d := range a.Depends {
			if e := visit(d); e != nil {
				return e
			}
		}
		mark[id] = 2
		out = append(out, a)
		return nil
	}
	for _, a := range c.Apps {
		if e := visit(a.ID); e != nil {
			return nil, e
		}
	}
	return out, nil
}
func (c Catalog) Find(id string) (App, bool) {
	for _, a := range c.Apps {
		if a.ID == id {
			return a, true
		}
	}
	return App{}, false
}
func Read(data []byte) (Catalog, error) {
	var c Catalog
	d := json.NewDecoder(strings.NewReader(string(data)))
	d.DisallowUnknownFields()
	e := d.Decode(&c)
	if e == nil {
		var extra any
		if d.Decode(&extra) != io.EOF {
			e = errors.New("trailing catalog data")
		}
	}
	if e == nil {
		e = c.Validate()
	}
	return c, e
}
func Extract(z *zip.Reader, c Catalog, dest string, overrides ...map[string]string) error {
	if e := c.Validate(); e != nil {
		return e
	}
	entries := map[string]*zip.File{}
	for _, f := range z.File {
		if f.Mode()&os.ModeSymlink != 0 {
			return errors.New("package links are not accepted")
		}
		if _, e := platform.Inside(dest, f.Name); e != nil {
			return e
		}
		key := strings.ToLower(f.Name)
		if entries[key] != nil {
			return errors.New("duplicate archive path")
		}
		entries[key] = f
	}
	for _, p := range c.Packages {
		path, e := platform.Inside(dest, p.File)
		if e != nil {
			return e
		}
		if e = os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			return e
		}
		out, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		in, e := openPackage(entries, p, overrides)
		if e != nil {
			out.Close()
			return e
		}
		hash := sha256.New()
		n, e := io.Copy(io.MultiWriter(out, hash), io.LimitReader(in, p.Size+1))
		in.Close()
		if e == nil {
			e = out.Sync()
		}
		out.Close()
		if e != nil || n != p.Size || hex.EncodeToString(hash.Sum(nil)) != p.SHA256 {
			return errors.New("payload checksum verification failed")
		}
	}
	return nil
}
