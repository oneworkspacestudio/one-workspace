package catalog

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestContinueReusesVerifiedFilesAndRetainsInterruptedBytes(t *testing.T) {
	c := fixture()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	file, _ := writer.Create("app.exe")
	file.Write([]byte("test executable"))
	writer.Close()
	reader, e := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	if e != nil {
		t.Fatal(e)
	}
	root := t.TempDir()
	p := filepath.Join(root, "app.exe")
	os.WriteFile(p, []byte("test executable"), 0600)
	before, _ := os.Stat(p)
	if e = ExtractContinue(reader, c, root, true, nil); e != nil {
		t.Fatal(e)
	}
	after, _ := os.Stat(p)
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatal("verified file was replaced")
	}
	os.WriteFile(p, []byte("interrupted bytes"), 0600)
	if e = ExtractContinue(reader, c, root, true, nil); e != nil {
		t.Fatal(e)
	}
	raw, _ := os.ReadFile(p)
	if string(raw) != "test executable" {
		t.Fatal("package not completed")
	}
	saved, _ := filepath.Glob(p + ".incomplete-*")
	if len(saved) != 1 {
		t.Fatal("interrupted file was not retained")
	}
	raw, _ = os.ReadFile(saved[0])
	if string(raw) != "interrupted bytes" {
		t.Fatal("original bytes lost")
	}
}
