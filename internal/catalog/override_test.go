package catalog

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReleaseOverrideVerifiedAndExtracted(t *testing.T) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	f, _ := z.Create("app.exe")
	f.Write([]byte("old bundled application"))
	z.Close()
	reader, _ := zip.NewReader(bytes.NewReader(b.Bytes()), int64(b.Len()))
	c := fixture()
	body := []byte("new release")
	c.Packages["app"] = Package{File: "app.exe", Size: int64(len(body)), SHA256: fmt.Sprintf("%x", sha256.Sum256(body))}
	path := filepath.Join(t.TempDir(), "release.exe")
	os.WriteFile(path, body, 0600)
	overrides := map[string]string{"app.exe": path}
	if err := Verify(reader, c, overrides); err != nil {
		t.Fatal(err)
	}
	dest := t.TempDir()
	if err := Extract(reader, c, dest, overrides); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(dest, "app.exe"))
	if !bytes.Equal(got, body) {
		t.Fatal("installed stale bundle")
	}
	os.WriteFile(path, []byte("bad release"), 0600)
	if Verify(reader, c, overrides) == nil {
		t.Fatal("accepted corrupt override")
	}
	os.Remove(path)
	if Verify(reader, c, overrides) == nil || Extract(reader, c, t.TempDir(), overrides) == nil {
		t.Fatal("fell back to old bundle")
	}
}
