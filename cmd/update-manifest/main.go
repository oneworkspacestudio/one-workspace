// update-manifest runs in a publisher's build pipeline, never on an end user's
// machine. The output contains no local build paths or credentials.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run() error {
	catalogPath := flag.String("catalog", "", "last compatible Workspace catalog.json")
	group := flag.String("group", "", "app update group")
	version := flag.String("version", "", "app version")
	repo := flag.String("repo", "", "oneworkspacestudio app repository")
	tag := flag.String("tag", "", "new release tag")
	fileMap := flag.String("files", "", "JSON package map: id -> {path, asset, member?}")
	output := flag.String("out", "workspace-update.json", "generated release asset")
	minimumManager := flag.String("minimum-manager", "", "minimum compatible Workspace version; selects manifest schema2")
	capabilities := flag.String("capabilities", "", "comma-separated required manager capabilities")
	flag.Parse()
	data, e := os.ReadFile(*catalogPath)
	if e != nil {
		return e
	}
	base, e := catalog.Read(data)
	if e != nil {
		return e
	}
	data, e = os.ReadFile(*fileMap)
	if e != nil {
		return e
	}
	var inputs map[string]struct {
		Path   string `json:"path"`
		Asset  string `json:"asset"`
		Member string `json:"member"`
	}
	if e = json.Unmarshal(data, &inputs); e != nil {
		return e
	}
	files := map[string]distribution.UpdatePackage{}
	for id, in := range inputs {
		old, ok := base.Packages[id]
		if !ok {
			return fmt.Errorf("unknown package %s", id)
		}
		f, e := os.Open(in.Path)
		if e != nil {
			return e
		}
		h := sha256.New()
		size, e := io.Copy(h, f)
		f.Close()
		if e != nil {
			return e
		}
		old.Size, old.SHA256 = size, hex.EncodeToString(h.Sum(nil))
		files[id] = distribution.UpdatePackage{Package: old, Asset: in.Asset, Member: in.Member}
	}
	manifest, e := distribution.MakeAppUpdateManifest(base, *group, *version, *repo, *tag, files)
	if e != nil {
		return e
	}
	var required []string
	if *capabilities != "" {
		for _, value := range strings.Split(*capabilities, ",") {
			required = append(required, strings.TrimSpace(value))
		}
	}
	if e = distribution.RequireManager(&manifest, *minimumManager, required); e != nil {
		return e
	}
	return platform.SaveJSON(*output, manifest)
}
