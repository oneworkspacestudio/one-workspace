package websearch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Models are private app-owned libraries (like My AI models), not executables.
// Only an explicit owner setup request installs these pinned immutable weights.
// Normal startup and activation never fetch executable code or model weights.
var modelPins = []struct {
	name, url, sha string
	size           int64
}{
	{"model.gguf", "https://huggingface.co/runanywhere/Fara1.5-4B-GGUF/resolve/a21bd238d82de0e39e263b39a0bdabd13e5e51f4/Fara1.5-4B-Q4_K_M.gguf", "8ea7715d804a15d9341783cec88b9eaec5eca39d7b53a2a6af8a9a293356f04c", 2707514176},
	{"mmproj.gguf", "https://huggingface.co/runanywhere/Fara1.5-4B-GGUF/resolve/a21bd238d82de0e39e263b39a0bdabd13e5e51f4/mmproj-Fara1.5-4B-f16.gguf", "90867de044cd67fd5d4cdc6fb6c68efe80e199e507ca029dbb07314c4e619603", 675569248},
}

func verifyModels(root string) error {
	for _, p := range modelPins {
		if e := verifyModelFile(filepath.Join(root, p.name), p.sha, p.size); e != nil {
			return e
		}
	}
	return nil
}
func verifyModelFile(path, want string, size int64) error {
	f, e := os.Open(path)
	if e != nil {
		return errors.New("install the local model first")
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil || info.Size() != size {
		return errors.New("model is incomplete; reinstall it")
	}
	h := sha256.New()
	_, e = copyBulk(h, f)
	if e != nil || hex.EncodeToString(h.Sum(nil)) != want {
		return errors.New("model verification failed; reinstall it")
	}
	return nil
}
func modelsPresent(root string) bool {
	for _, pin := range modelPins {
		info, e := os.Stat(filepath.Join(root, pin.name))
		if e != nil || info.Size() != pin.size {
			return false
		}
	}
	return true
}
func installModels(ctx context.Context, root string, progress func(int64)) error {
	if e := os.MkdirAll(root, 0700); e != nil {
		return e
	}
	client := &http.Client{Timeout: 30 * time.Minute}
	for _, pin := range modelPins {
		target := filepath.Join(root, pin.name)
		if verifyModelFile(target, pin.sha, pin.size) == nil {
			continue
		}
		req, _ := http.NewRequestWithContext(ctx, "GET", pin.url, nil)
		res, e := client.Do(req)
		if e != nil {
			return e
		}
		if res.StatusCode != 200 {
			res.Body.Close()
			return errors.New("model download unavailable")
		}
		tmp, e := os.CreateTemp(root, ".model-*")
		if e != nil {
			res.Body.Close()
			return e
		}
		name := tmp.Name()
		h := sha256.New()
		n, e := copyBulk(io.MultiWriter(tmp, h, &progressWriter{fn: progress}), io.LimitReader(res.Body, pin.size+1))
		res.Body.Close()
		closeErr := tmp.Close()
		if e != nil || closeErr != nil || n != pin.size || hex.EncodeToString(h.Sum(nil)) != pin.sha {
			os.Remove(name)
			return errors.New("model download failed its size or digest check")
		}
		if e = os.Remove(target); e != nil && !os.IsNotExist(e) {
			os.Remove(name)
			return e
		}
		if e = os.Rename(name, target); e != nil {
			os.Remove(name)
			return e
		}
	}
	return verifyModels(root)
}

type progressWriter struct{ fn func(int64) }

func (p *progressWriter) Write(b []byte) (int, error) { p.fn(int64(len(b))); return len(b), nil }
