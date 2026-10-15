package catalog

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallerIntegrationGate(t *testing.T) {
	c := fixture()
	c.InstallerIntegration = 2
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = Read(data); err != nil {
		t.Fatal(err)
	}
	// Setup16's exact strict field set must reject the new integration before
	// extraction/identity configuration. Merely retaining schema1 was unsafe.
	var legacy struct {
		Schema   int                `json:"schema"`
		Version  string             `json:"version"`
		Packages map[string]Package `json:"packages"`
		Apps     []App              `json:"apps"`
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err = d.Decode(&legacy); err == nil {
		t.Fatal("older installer accepted a shared-service catalog")
	}
	c.InstallerIntegration = 3
	data, _ = json.Marshal(c)
	if _, err = Read(data); err == nil {
		t.Fatal("future integration accepted")
	}
}

func fixture() Catalog {
	b := []byte("test executable")
	return Catalog{Schema: 1, Version: "0.1.0", Packages: map[string]Package{"app": {File: "app.exe", Size: int64(len(b)), SHA256: fmt.Sprintf("%x", sha256.Sum256(b))}}, Apps: []App{{ID: "test", Name: "Test", Version: "1.0", Kind: "process", Port: 9900, Route: "/test", Executable: "app.exe", Package: "app"}}}
}
func TestCatalogRejectsInvalidGraphsAndPackages(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*Catalog)
	}{
		{"dependency cycle", func(c *Catalog) { c.Apps[0].Depends = []string{"test"} }},
		{"missing dependency", func(c *Catalog) { c.Apps[0].Depends = []string{"absent"} }},
		{"different executable", func(c *Catalog) { c.Apps[0].Executable = "unchecked.exe" }},
		{"unsafe version", func(c *Catalog) { c.Apps[0].Version = "../bad" }},
		{"overlapping route", func(c *Catalog) {
			b := c.Apps[0]
			b.ID = "other"
			b.Port++
			b.Route = "/test/child"
			c.Apps = append(c.Apps, b)
		}},
		{"unverified extra package", func(c *Catalog) { c.Packages["bad"] = Package{File: "bad.exe"} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := fixture()
			tc.change(&c)
			if c.Validate() == nil {
				t.Fatal("invalid catalog accepted")
			}
		})
	}
}
func TestVerifiedExtraction(t *testing.T) {
	for _, tc := range []struct {
		name, body, extra string
		want              bool
	}{{"valid", "test executable", "", true}, {"corrupt", "evil executable", "", false}, {"traversal", "test executable", "../escape", false}, {"case duplicate", "test executable", "APP.EXE", false}, {"ADS", "test executable", "file:evil", false}, {"device", "test executable", "NUL.txt", false}} {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			z := zip.NewWriter(&b)
			f, _ := z.Create("app.exe")
			_, _ = f.Write([]byte(tc.body))
			if tc.extra != "" {
				f, _ = z.Create(tc.extra)
				_, _ = f.Write([]byte("x"))
			}
			_ = z.Close()
			reader, _ := zip.NewReader(bytes.NewReader(b.Bytes()), int64(b.Len()))
			dest := t.TempDir()
			e := Extract(reader, fixture(), dest)
			if (e == nil) != tc.want {
				t.Fatalf("extraction error = %v", e)
			}
			if tc.want {
				data, _ := os.ReadFile(filepath.Join(dest, "app.exe"))
				if string(data) != tc.body {
					t.Fatal("wrong extracted bytes")
				}
			}
		})
	}
}
func TestReadRejectsTrailingAndUnknown(t *testing.T) {
	for _, body := range []string{`{"schema":1,"version":"x","apps":[],"packages":{},"unknown":true}`, `{"schema":1,"version":"x","apps":[],"packages":{}} {}`} {
		if _, e := Read([]byte(body)); e == nil {
			t.Fatal("accepted", strings.TrimSpace(body))
		}
	}
}
