package froststock

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestStorePhotoPairAndRevision(t *testing.T) {
	root := filepath.Join(t.TempDir(), "Froststock")
	store := drive{root: root}
	const id = "food-label-0001"
	item := entry{ID: id, Created: 1, Job: "read-label-0001", Quantity: 1, State: "queued"}
	var photo bytes.Buffer
	if err := jpeg.Encode(&photo, image.NewGray(image.Rect(0, 0, 80, 60)), nil); err != nil {
		t.Fatal(err)
	}
	call := func(method string, body []byte, contentType, photoName, ifMatch string, handler http.HandlerFunc) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "/", bytes.NewReader(body))
		r.SetPathValue("id", id)
		r.SetPathValue("photo", photoName)
		r.Header.Set("Content-Type", contentType)
		r.Header.Set("If-Match", ifMatch)
		w := httptest.NewRecorder()
		handler(w, r)
		return w
	}
	upload := func(second bool) *httptest.ResponseRecorder {
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		record, _ := json.Marshal(item)
		_ = form.WriteField("record", string(record))
		_ = form.WriteField("item", photo.String())
		if second {
			_ = form.WriteField("expiry", photo.String())
		}
		_ = form.Close()
		return call("POST", body.Bytes(), form.FormDataContentType(), "", "", store.create)
	}
	if w := upload(false); w.Code != 400 {
		t.Fatalf("single photo accepted: %d", w.Code)
	}
	if w := upload(true); w.Code != 201 {
		t.Fatalf("pair save: %d %s", w.Code, w.Body.String())
	}
	if w := upload(true); w.Code != 200 {
		t.Fatalf("repeat upload: %d %s", w.Code, w.Body.String())
	}
	for _, name := range []string{"item.jpg", "expiry.jpg"} {
		b, err := os.ReadFile(filepath.Join(root, id, name))
		if err != nil || !bytes.Equal(b, photo.Bytes()) {
			t.Fatalf("Drive %s differs: %v", name, err)
		}
		w := call("GET", nil, "", name, "", store.photo)
		if w.Code != 200 || !bytes.Equal(w.Body.Bytes(), b) {
			t.Fatalf("served %s differs: %d", name, w.Code)
		}
	}
	item.Expiry, item.State, item.RawDate = "2027-11-22", "ready", "BEST BEFORE 22 NOV 2027"
	item.Revision = 1
	body, _ := json.Marshal(item)
	if w := call("PUT", body, "application/json", "", "", store.update); w.Code != 200 {
		t.Fatalf("update: %d %s", w.Code, w.Body.String())
	}
	if w := call("PUT", body, "application/json", "", "", store.update); w.Code != 409 {
		t.Fatalf("stale revision accepted: %d", w.Code)
	}
	if w := call("GET", nil, "", "", "", store.list); w.Code != 200 || !bytes.Contains(w.Body.Bytes(), []byte(item.RawDate)) {
		t.Fatalf("list: %d %s", w.Code, w.Body.String())
	}
	if w := call("DELETE", nil, "", "", revisionTag(1), store.delete); w.Code != 409 {
		t.Fatalf("stale delete: %d", w.Code)
	}
	userFile := filepath.Join(root, id, "notes.txt")
	if err := os.WriteFile(userFile, []byte("Keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if w := call("DELETE", nil, "", "", revisionTag(2), store.delete); w.Code != 204 {
		t.Fatalf("delete: %d %s", w.Code, w.Body.String())
	}
	if b, err := os.ReadFile(userFile); err != nil || string(b) != "Keep" {
		t.Fatalf("extra Drive file was changed: %s %v", b, err)
	}
}
