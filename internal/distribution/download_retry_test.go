package distribution

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadResumesInterruptedResponse(t *testing.T) {
	data := []byte("verified download content")
	sum := sha256.Sum256(data)
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.Header().Set("Content-Length", fmt.Sprint(len(data)))
			w.Write(data[:4])
			return
		}
		if r.Header.Get("Range") != "bytes=4-" {
			t.Errorf("missing resume range: %s", r.Header.Get("Range"))
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes 4-%d/%d", len(data)-1, len(data)))
		w.WriteHeader(206)
		w.Write(data[4:])
	}))
	defer server.Close()
	path := filepath.Join(t.TempDir(), "download")
	err := downloadMailAsset(context.Background(), server.Client(), "", Asset{Name: "fixture", URL: server.URL, Size: int64(len(data)), Digest: "sha256:" + hex.EncodeToString(sum[:])}, path)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(data) || calls != 2 {
		t.Fatalf("incorrect content or attempts: %d", calls)
	}
}
func TestDownloadRetriesCorruptResponse(t *testing.T) {
	data := []byte("good")
	sum := sha256.Sum256(data)
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.Write([]byte("evil"))
		} else {
			w.Write(data)
		}
	}))
	defer server.Close()
	path := filepath.Join(t.TempDir(), "download")
	if err := downloadMailAsset(context.Background(), server.Client(), "", Asset{Name: "fixture", URL: server.URL, Size: 4, Digest: "sha256:" + hex.EncodeToString(sum[:])}, path); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatal(calls)
	}
}
