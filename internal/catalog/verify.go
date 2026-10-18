package catalog

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io"
	"os"
	"strings"
)

// Verify checks every executable and license before an online install writes files.
func Verify(z *zip.Reader, c Catalog, overrides ...map[string]string) error {
	if err := c.Validate(); err != nil {
		return err
	}
	entries := map[string]*zip.File{}
	for _, f := range z.File {
		if f.Mode()&os.ModeSymlink != 0 {
			return errors.New("package links are not accepted")
		}
		if _, err := platform.Inside("root", f.Name); err != nil {
			return err
		}
		key := strings.ToLower(f.Name)
		if entries[key] != nil {
			return errors.New("duplicate archive path")
		}
		entries[key] = f
	}
	for _, p := range c.Packages {
		in, err := openPackage(entries, p, overrides)
		if err != nil {
			return err
		}
		hash := sha256.New()
		n, err := io.Copy(hash, io.LimitReader(in, p.Size+1))
		in.Close()
		if err != nil || n != p.Size || hex.EncodeToString(hash.Sum(nil)) != p.SHA256 {
			return errors.New("payload checksum verification failed")
		}
	}
	return nil
}

// Overrides are verified release files owned by the downloader, never user paths.
// Both verification and extraction hash them against the updated catalog.
func openPackage(entries map[string]*zip.File, p Package, overrides []map[string]string) (io.ReadCloser, error) {
	if len(overrides) > 0 {
		if path, ok := overrides[0][p.File]; ok {
			in, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			info, err := in.Stat()
			if err != nil || !info.Mode().IsRegular() || info.Size() != p.Size {
				in.Close()
				return nil, errors.New("release override missing or size mismatch")
			}
			return in, nil
		}
	}
	f := entries[strings.ToLower(p.File)]
	if f == nil || int64(f.UncompressedSize64) != p.Size {
		return nil, errors.New("payload missing or size mismatch")
	}
	return f.Open()
}
