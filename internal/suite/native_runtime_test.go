package suite

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"os"
	"path/filepath"
	"testing"
)

func TestSharedRuntimeVerifiedLocalDeploymentAndTamperRetention(t *testing.T) {
	root := t.TempDir()
	programs := filepath.Join(root, "Apps")
	engine := filepath.Join(root, "Private", "shared", "bin", "stalwart.exe")
	if err := os.MkdirAll(filepath.Dir(engine), 0700); err != nil {
		t.Fatal(err)
	}
	c := Config{Programs: programs, Catalog: catalog.Catalog{Packages: map[string]catalog.Package{}}}
	for _, name := range nativeRuntimeDLLs {
		rel := "runtime/msvc/" + name
		b := []byte("verified test bytes " + name)
		sum := sha256.Sum256(b)
		path := filepath.Join(programs, filepath.FromSlash(rel))
		os.MkdirAll(filepath.Dir(path), 0700)
		if err := os.WriteFile(path, b, 0600); err != nil {
			t.Fatal(err)
		}
		c.Catalog.Packages["msvc-"+name[:len(name)-4]] = catalog.Package{File: rel, SHA256: hex.EncodeToString(sum[:]), Size: int64(len(b))}
	}
	if err := prepareSharedRuntime(c, engine); err != nil {
		t.Fatal(err)
	}
	if err := prepareSharedRuntime(c, engine); err != nil {
		t.Fatal("restart must reuse identical runtime", err)
	}
	dest := filepath.Join(filepath.Dir(engine), "msvcp140.dll")
	os.WriteFile(dest, []byte("unfamiliar"), 0600)
	if err := prepareSharedRuntime(c, engine); err == nil {
		t.Fatal("overwrote an unverified DLL")
	}
	b, _ := os.ReadFile(dest)
	if string(b) != "unfamiliar" {
		t.Fatal("did not retain unfamiliar bytes")
	}
}
