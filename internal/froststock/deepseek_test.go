package froststock

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestDeepSeekSettingsAndReading(t *testing.T) {
	data := t.TempDir()
	d, err := newDeepSeek(data)
	if err != nil {
		t.Fatal(err)
	}
	status := httptest.NewRecorder()
	d.status(status, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(status.Body.String(), `"configured":false`) {
		t.Fatalf("unexpected initial status: %s", status.Body)
	}

	const key = "test-key-froststock"
	configure := func(body string, want int) {
		t.Helper()
		w := httptest.NewRecorder()
		d.configure(w, httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body)))
		if w.Code != want || strings.Contains(w.Body.String(), key) {
			t.Fatalf("unexpected configure result: %d %s", w.Code, w.Body)
		}
	}
	configure(`{"key":"test-key-froststock","remove":true}`, http.StatusBadRequest)
	configure(`{"key":"test-key-froststock"}`, http.StatusOK)
	d, err = newDeepSeek(data)
	if err != nil {
		t.Fatal(err)
	}
	if d.key() != key {
		t.Fatal("key was not restored from private settings")
	}

	var upstreamCalls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := upstreamCalls.Add(1)
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer "+key {
			t.Error("unexpected upstream method or authorization")
		}
		var request struct {
			Model    string `json:"model"`
			Messages []struct {
				Content []struct {
					Type  string `json:"type"`
					Text  string `json:"text"`
					Image struct {
						URL    string `json:"url"`
						Detail string `json:"detail"`
					} `json:"image_url"`
				} `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Model != deepSeekModel || len(request.Messages) != 1 {
			t.Error("expected the vision model with one user message")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		content := request.Messages[0].Content
		wantParts := 5
		if call == 2 {
			wantParts = 3
		}
		if len(content) != wantParts || content[len(content)-2].Text != "The next image is the expiry photo." ||
			content[len(content)-1].Type != "image_url" || content[len(content)-1].Image.Detail != "original" ||
			!strings.HasPrefix(content[len(content)-1].Image.URL, "data:image/jpeg;base64,") {
			t.Error("expected original JPEG photos with an identified expiry photo")
		}
		if call == 2 {
			http.Error(w, key, http.StatusUnauthorized)
			return
		}
		reading, _ := json.Marshal(deepSeekReading{
			Name: "Chicken wings", RawDate: "PACKED 24.11.2026\nUSE BY 12.01.2027", ExpiryText: "12.01.2027",
		})
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{
			"finish_reason": "stop", "message": map[string]string{"content": string(reading)},
		}}})
	}))
	defer upstream.Close()
	d.endpoint = upstream.URL

	scan := func(photos []string, want int) string {
		t.Helper()
		body := new(bytes.Buffer)
		form := multipart.NewWriter(body)
		for _, photo := range photos {
			part, err := form.CreateFormFile(photo, photo+".jpg")
			if err != nil {
				t.Fatal(err)
			}
			_, _ = part.Write([]byte{0xff, 0xd8, 0xff, 0xd9})
		}
		if err := form.Close(); err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", form.FormDataContentType())
		w := httptest.NewRecorder()
		d.read(w, req)
		if w.Code != want || strings.Contains(w.Body.String(), key) {
			t.Fatalf("unexpected reading: %d %s", w.Code, w.Body)
		}
		return w.Body.String()
	}
	if body := scan([]string{"item", "expiry"}, http.StatusOK); !strings.Contains(body, "12.01.2027") {
		t.Fatalf("missing expiry transcription: %s", body)
	}
	scan([]string{"item"}, http.StatusBadRequest)
	scan([]string{"expiry"}, http.StatusBadGateway)
	configure(`{"remove":true}`, http.StatusOK)
	scan([]string{"item", "expiry"}, http.StatusNoContent)
	if calls := upstreamCalls.Load(); calls != 2 {
		t.Fatalf("unexpected number of upstream calls: %d", calls)
	}
}
