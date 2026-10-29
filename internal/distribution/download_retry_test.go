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
	"time"
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

func TestLargeDownloadUsesBoundedRanges(t *testing.T) {
	data := make([]byte, 10<<20)
	for i := range data {
		data[i] = byte(i)
	}
	sum := sha256.Sum256(data)
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var start, end int
		if _, err := fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil || end-start+1 > 8<<20 || start < 0 || end >= len(data) {
			t.Error("invalid bounded range", r.Header.Get("Range"))
			w.WriteHeader(400)
			return
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(data)))
		w.WriteHeader(206)
		w.Write(data[start : end+1])
	}))
	defer server.Close()
	if err := downloadMailAsset(context.Background(), server.Client(), "", Asset{Name: "large", URL: server.URL, Size: int64(len(data)), Digest: "sha256:" + hex.EncodeToString(sum[:])}, filepath.Join(t.TempDir(), "asset")); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatal("expected two chunks", calls)
	}
}

func TestLargeDownloadExceedsSingleRequestTimeoutSafely(t *testing.T) {
	data := make([]byte, 10<<20)
	for i := range data {
		data[i] = byte(i)
	}
	sum := sha256.Sum256(data)
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		time.Sleep(300 * time.Millisecond)
		var start, end int
		if _, err := fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil || end-start+1 > 8<<20 || start < 0 || end >= len(data) {
			t.Error("invalid bounded range", r.Header.Get("Range"))
			w.WriteHeader(400)
			return
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(data)))
		w.WriteHeader(206)
		w.Write(data[start : end+1])
	}))
	defer server.Close()
	if err := downloadMailAsset(context.Background(), &http.Client{Timeout: 500 * time.Millisecond}, "", Asset{Name: "large", URL: server.URL, Size: int64(len(data)), Digest: "sha256:" + hex.EncodeToString(sum[:])}, filepath.Join(t.TempDir(), "asset")); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatal("expected two chunks", calls)
	}
}
