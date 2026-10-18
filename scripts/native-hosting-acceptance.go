//go:build ignore

// Run explicitly with `go run scripts/native-hosting-acceptance.go ...`.
// This creates only a fresh disposable installation; never point it at live data.
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/suite"
)

func main() {
	root := flag.String("root", "", "New disposable NTFS test root")
	manager := flag.String("manager", "", "Locally built Workspace.exe")
	accounts := flag.String("accounts", "", "Verified organization MyAccounts executable")
	hosting := flag.String("hosting", "", "Locally built MyHosting.exe")
	bundle := flag.String("bundle", "", "Verified full organization bundle ZIP (instead of local binaries)")
	port := flag.Int("port", 19300, "Isolated base port")
	flag.Parse()
	if e := install(*root, *manager, *accounts, *hosting, *bundle, *port); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func install(root, manager, accounts, hosting, bundle string, port int) error {
	password := os.Getenv("WORKSPACE_TEST_PASSWORD")
	if password == "" || root == "" || !filepath.IsAbs(root) {
		return fmt.Errorf("Supply a new absolute --root and WORKSPACE_TEST_PASSWORD for this disposable owner")
	}
	if _, e := os.Lstat(root); !os.IsNotExist(e) {
		return fmt.Errorf("Acceptance requires a completely new test root")
	}
	var archive bytes.Buffer
	z := zip.NewWriter(&archive)
	packages := map[string]catalog.Package{}
	items := []struct{ id, path, name string }{{"manager", manager, "Workspace.exe"}, {"accounts", accounts, "MyAccounts.exe"}, {"hosting", hosting, "MyHosting.exe"}}
	if bundle != "" {
		items = nil
	}
	for _, item := range items {
		f, e := os.Open(item.path)
		if e != nil {
			return e
		}
		entry, e := z.Create(item.name)
		if e != nil {
			f.Close()
			return e
		}
		hash := sha256.New()
		n, e := io.Copy(io.MultiWriter(entry, hash), f)
		f.Close()
		if e != nil {
			return e
		}
		packages[item.id] = catalog.Package{File: item.name, SHA256: hex.EncodeToString(hash.Sum(nil)), Size: n}
	}
	if e := z.Close(); e != nil {
		return e
	}
	reader, e := zip.NewReader(bytes.NewReader(archive.Bytes()), int64(archive.Len()))
	if e != nil {
		return e
	}
	c := catalog.Catalog{Schema: 1, InstallerIntegration: 2, Version: "0.3.0-preview.36-dev", Packages: packages, Apps: []catalog.App{
		{ID: "accounts", Name: "MyAccounts", Version: "preview-20261015.1", Kind: "process", Package: "accounts", Executable: "MyAccounts.exe", Port: 8991, Route: "/", Health: "/healthz", Args: []string{"-data", "${DATA}/accounts", "-no-browser"}, Source: "https://github.com/oneworkspacestudio/myaccounts"},
		{ID: "hosting", Name: "My Hosting", Version: "0.1.0-preview.1", Kind: "process", Package: "hosting", Executable: "MyHosting.exe", Port: 9007, Route: "/hosting", StripPrefix: true, Health: "/healthz", Depends: []string{"accounts"}, OIDC: true, Callback: "/hosting/auth/callback", Args: []string{"--data", "${DATA}/my-hosting", "--listen", "127.0.0.1:${PORT}", "--origin", "${ORIGIN}/hosting", "--no-browser"}, Source: "https://github.com/oneworkspacestudio/my-hosting"},
	}}
	if bundle != "" {
		full, openErr := zip.OpenReader(bundle)
		if openErr != nil {
			return openErr
		}
		defer full.Close()
		reader = &full.Reader
		manifest, openErr := reader.Open("catalog.json")
		if openErr != nil {
			return openErr
		}
		raw, readErr := io.ReadAll(io.LimitReader(manifest, 4<<20))
		manifest.Close()
		if readErr != nil {
			return readErr
		}
		c, e = catalog.Read(raw)
		if e != nil {
			return e
		}
		if _, present := c.Find("hosting"); !present {
			return fmt.Errorf("Bundle does not contain My Hosting")
		}
	}
	// Isolate the installer's Start Menu shortcut from the user's real profile.
	roaming := filepath.Join(root, "roaming")
	if e = os.MkdirAll(filepath.Join(roaming, "Microsoft", "Windows", "Start Menu", "Programs"), 0700); e != nil {
		return e
	}
	if e = os.Setenv("APPDATA", roaming); e != nil {
		return e
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	r := suite.Request{Programs: filepath.Join(root, "Apps"), Data: filepath.Join(root, "Private"), Files: filepath.Join(root, "Files"), Name: "Hosting acceptance owner", Email: "hosting-owner@example.test", Password: password, Confirm: password, Port: port, Mode: "local", Offline: true, SelectedApps: []string{"accounts", "hosting"}}
	installed, e := suite.Install(ctx, r, reader, c, func(s string) { fmt.Println(s) })
	if e != nil {
		return e
	}
	fmt.Println("Fresh test installation ready:", installed.Origin, "; manager:", filepath.Join(installed.Data, "Workspace.exe"))
	return nil
}
