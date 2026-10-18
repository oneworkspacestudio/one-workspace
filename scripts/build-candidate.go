//go:build ignore

// Local, offline acceptance bundle: verify every seed package, then replace only
// the reviewed manager, vault and vault guides. Does not publish or change pins.
package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func must(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	seed := flag.String("seed", "", "verified existing suite bundle")
	manager := flag.String("manager", "", "locally built candidate manager")
	vault := flag.String("vault", "", "locally built candidate MyVault.exe")
	docs := flag.String("vault-docs", "", "candidate vault docs folder")
	output := flag.String("output", "", "new candidate output directory")
	version := flag.String("version", "0.3.0-rc.1", "local suite candidate version")
	vaultVersion := flag.String("vault-version", "0.2.0-rc.1", "local vault candidate version")
	flag.Parse()
	for _, s := range []string{*seed, *manager, *vault, *docs, *output} {
		if s == "" {
			panic("all paths are required")
		}
	}
	if !strings.Contains(*version, "-rc.") || !strings.Contains(*vaultVersion, "-rc.") {
		panic("local acceptance bundles must be labelled release candidates")
	}
	if _, e := os.Stat(*output); !os.IsNotExist(e) {
		panic("use a new output directory")
	}
	z, e := zip.OpenReader(*seed)
	must(e)
	defer z.Close()
	files := map[string]*zip.File{}
	seen := map[string]bool{}
	for _, f := range z.File {
		key := strings.ToLower(f.Name)
		if seen[key] {
			panic("duplicate archive member")
		}
		seen[key] = true
		files[f.Name] = f
	}
	cf := files["catalog.json"]
	if cf == nil {
		panic("seed catalog missing")
	}
	r, e := cf.Open()
	must(e)
	data, e := io.ReadAll(io.LimitReader(r, 1<<20))
	must(e)
	must(r.Close())
	cat, e := catalog.Read(data)
	must(e)
	if len(files) != len(cat.Packages)+1 {
		panic("seed includes unverified archive members")
	}
	for key, p := range cat.Packages {
		f := files[p.File]
		if f == nil {
			panic("missing package: " + key)
		}
		r, e := f.Open()
		must(e)
		h := sha256.New()
		n, e := io.Copy(h, r)
		must(e)
		must(r.Close())
		if n != p.Size || hex.EncodeToString(h.Sum(nil)) != p.SHA256 {
			panic("seed package integrity mismatch: " + key)
		}
	}
	replacements := map[string][]byte{}
	replace := func(key, path, file string) {
		b, e := os.ReadFile(path)
		must(e)
		s := sha256.Sum256(b)
		cat.Packages[key] = catalog.Package{File: file, Size: int64(len(b)), SHA256: hex.EncodeToString(s[:])}
		replacements[file] = b
	}
	replace("manager", *manager, "Workspace.exe")
	replace("vault", *vault, "MyVault.exe")
	for _, name := range []string{"RECOVERY.md", "CLIENTS.md", "HOSTING.md", "ADMIN.md", "SECURITY-REVIEW.md", "SIGN-IN.md", "PRODUCTION-GATES.md", "acceptance.json"} {
		replace("vault-"+name, filepath.Join(*docs, name), "docs/MyVault/"+name)
	}
	cat.Version = *version
	for i := range cat.Apps {
		if cat.Apps[i].ID == "vault" {
			cat.Apps[i].Version = *vaultVersion
			cat.Apps[i].Source = "https://github.com/oneworkspacestudio/my-vault"
		}
	}
	must(cat.Validate())
	data, e = json.MarshalIndent(cat, "", "  ")
	must(e)
	replacements["catalog.json"] = data
	must(os.MkdirAll(*output, 0700))
	must(os.WriteFile(filepath.Join(*output, "catalog.json"), data, 0600))
	out, e := os.OpenFile(filepath.Join(*output, "workspace-bundle.zip"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	must(e)
	w := zip.NewWriter(out)
	for _, f := range z.File {
		if b, ok := replacements[f.Name]; ok {
			h := f.FileHeader
			h.Method = zip.Deflate
			dest, e := w.CreateHeader(&h)
			must(e)
			_, e = dest.Write(b)
			must(e)
			delete(replacements, f.Name)
		} else {
			must(w.Copy(f))
		}
	}
	names := make([]string, 0, len(replacements))
	for name := range replacements {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		dest, e := w.Create(name)
		must(e)
		_, e = dest.Write(replacements[name])
		must(e)
	}
	must(w.Close())
	must(out.Close())
	fmt.Println("Verified local candidate bundle; unchanged components copied verbatim. No release or organization pin changed.")
}
