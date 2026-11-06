package websearch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io/fs"
	"os"
	"path/filepath"
)

// Chromium's Windows sandbox must be able to read its own immutable executable
// files. When Programs is exFAT, stage ONLY browser code in the app's NTFS cache.
// Profile, cookies, API keys and model data are never granted package access.
func prepareBrowser(ctx context.Context, runtimeRoot, data string) (string, error) {
	source := filepath.Join(runtimeRoot, "browser")
	if _, e := platform.CheckDataVolume(source); e == nil {
		if e = allowSandboxCode(source); e != nil {
			return "", e
		}
		return filepath.Join(source, "chrome-win64", "chrome.exe"), nil
	}
	target := filepath.Join(data, "browser-code-cache")
	if _, e := platform.CheckDataVolume(data); e != nil {
		return "", e
	}
	var manifest map[string]string
	raw, e := os.ReadFile(filepath.Join(runtimeRoot, "browser-files.json"))
	if e != nil {
		return "", errors.New("browser runtime manifest missing; repair Web Search")
	}
	if e = json.Unmarshal(raw, &manifest); e != nil {
		return "", e
	}
	for rel, want := range manifest {
		if e = ctx.Err(); e != nil {
			return "", e
		}
		from, e := platform.Inside(source, rel)
		if e != nil {
			return "", e
		}
		to, e := platform.Inside(target, rel)
		if e != nil {
			return "", e
		}
		if e = platform.RejectReparse(from); e != nil {
			return "", e
		}
		valid := func(path string) bool {
			f, e := os.Open(path)
			if e != nil {
				return false
			}
			defer f.Close()
			sum := sha256.New()
			_, e = copyBulk(sum, f)
			return e == nil && hex.EncodeToString(sum.Sum(nil)) == want
		}
		if valid(to) {
			continue
		}
		if !valid(from) {
			return "", errors.New("browser payload digest mismatch")
		}
		if e = os.MkdirAll(filepath.Dir(to), 0700); e != nil {
			return "", e
		}
		// Reject every existing target ancestor rather than following a junction.
		for p := filepath.Dir(to); p != data; p = filepath.Dir(p) {
			if e = platform.RejectReparse(p); e != nil {
				return "", e
			}
		}
		in, e := os.Open(from)
		if e != nil {
			return "", e
		}
		out, e := os.OpenFile(to, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if e != nil {
			in.Close()
			return "", e
		}
		_, e = copyBulk(out, in)
		in.Close()
		closeErr := out.Close()
		if e != nil {
			return "", e
		}
		if closeErr != nil {
			return "", closeErr
		}
	}
	if e = allowSandboxCode(target); e != nil {
		return "", e
	}
	return filepath.Join(target, "chrome-win64", "chrome.exe"), nil
}

// Used only by the checked build tool and synthetic native acceptance fixtures.
func BrowserManifest(root string) (map[string]string, error) {
	out := map[string]string{}
	e := filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		f, e := os.Open(path)
		if e != nil {
			return e
		}
		h := sha256.New()
		_, e = copyBulk(h, f)
		f.Close()
		if e != nil {
			return e
		}
		rel, _ := filepath.Rel(root, path)
		out[filepath.ToSlash(rel)] = hex.EncodeToString(h.Sum(nil))
		return nil
	})
	return out, e
}
