//go:debug httpmuxgo121=0

package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

//go:embed web/*
var assets embed.FS

const (
	version       = "0.1.0-preview.4"
	defaultEngine = "http://127.0.0.1:8001"
	maxAudio      = 100 << 20
)

type track struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Format  string    `json:"format"`
	Created time.Time `json:"created"`
}

func request(ctx context.Context, endpoint, key string, body any) (*http.Response, error) {
	var input io.Reader
	method := "GET"
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		input = bytes.NewReader(data)
		method = "POST"
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, input)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	client := http.Client{Timeout: 2 * time.Minute, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("service returned HTTP %d", res.StatusCode)
	}
	return res, nil
}

func query(ctx context.Context, endpoint, key string, body, out any) error {
	res, err := request(ctx, endpoint, key, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(io.LimitReader(res.Body, 4<<20)).Decode(out)
}

type aceResponse[T any] struct {
	Data T   `json:"data"`
	Code int `json:"code"`
}

type generation struct {
	Title         string
	Prompt        string
	Lyrics        string
	Seconds       int
	BPM           int
	AutoDuration  bool
	VocalLanguage string
	Instrumental  bool
	Describe      bool
}

type application struct {
	dir    string
	config func() engineSettings
	freeze *sync.RWMutex
}

func (app application) generate(ctx context.Context, song generation) (track, error) {
	body, err := song.acePayload()
	if err != nil {
		return track{}, err
	}
	settings := app.config()
	engine, key := settings.Engine, settings.ACEKey
	var start aceResponse[struct {
		TaskID string `json:"task_id"`
	}]
	if err := query(ctx, engine+"/release_task", key, body, &start); err != nil {
		return track{}, fmt.Errorf("ACE-Step submission: %w", err)
	}
	if start.Code != 200 || start.Data.TaskID == "" {
		return track{}, errors.New("ACE-Step did not return a job ID")
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return track{}, ctx.Err()
		case <-ticker.C:
		}
		var result aceResponse[[]struct {
			Status int
			Result string
		}]
		if err := query(ctx, engine+"/query_result", key, map[string]any{"task_id_list": []string{start.Data.TaskID}}, &result); err != nil {
			return track{}, err
		}
		if result.Code != 200 {
			return track{}, errors.New("ACE-Step could not read the job")
		}
		if len(result.Data) == 0 || result.Data[0].Status == 0 {
			continue
		}
		if result.Data[0].Status != 1 {
			return track{}, errors.New("ACE-Step generation failed. See the engine log for details.")
		}
		var files []struct{ File string }
		if err := json.Unmarshal([]byte(result.Data[0].Result), &files); err != nil || len(files) == 0 {
			return track{}, errors.New("ACE-Step returned no audio")
		}
		file, err := url.Parse(files[0].File)
		if err != nil || file.IsAbs() || file.Host != "" || file.Path != "/v1/audio" || file.Query().Get("path") == "" {
			return track{}, errors.New("unexpected ACE-Step download URL")
		}
		res, err := request(ctx, engine+file.RequestURI(), key, nil)
		if err != nil {
			return track{}, err
		}
		defer res.Body.Close()
		if app.freeze != nil {
			app.freeze.RLock()
			defer app.freeze.RUnlock()
		}
		return app.save(song.Title, "wav", res.Body)
	}
}

func (app application) save(title, format string, input io.Reader) (track, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "Untitled"
	}
	if len(title) > 240 {
		return track{}, errors.New("title is too long")
	}
	buffer := bufio.NewReader(input)
	header, _ := buffer.Peek(12)
	if !audioMatches(format, header) {
		return track{}, errors.New("audio contents do not match the selected format")
	}
	id := rand.Text()
	item := track{id, title, format, time.Now()}
	audioPath := filepath.Join(app.dir, id+"."+format)
	file, err := os.OpenFile(audioPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return track{}, err
	}
	size, copyErr := io.Copy(file, io.LimitReader(buffer, maxAudio+1))
	closeErr := file.Close()
	if copyErr != nil || closeErr != nil || size == 0 || size > maxAudio {
		os.Remove(audioPath)
		return track{}, errors.New("audio could not be saved or is larger than 100 MB")
	}
	if err := platform.SaveJSON(filepath.Join(app.dir, id+".json"), item); err != nil {
		os.Remove(audioPath)
		return track{}, err
	}
	return item, nil
}

func (app application) library() ([]track, error) {
	entries, err := os.ReadDir(app.dir)
	if err != nil {
		return nil, err
	}
	items := []track{}
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(app.dir, entry.Name()))
		var item track
		if err == nil && json.Unmarshal(data, &item) == nil && item.ID+".json" == entry.Name() {
			items = append(items, item)
		}
	}
	return items, nil
}

func assist(ctx context.Context, prompt, mode, key string) (string, error) {
	if key == "" {
		return "", errors.New("configure DeepSeek in Engine before using the assistant")
	}
	instruction := "Write original song lyrics, with section tags such as [Verse] and [Chorus]. Return only the lyrics."
	if mode == "prompt" {
		instruction = "Rewrite the brief as one concise music generation prompt specifying genre, instrumentation and production. Return only the prompt."
	}
	var result struct {
		Choices []struct{ Message struct{ Content string } }
	}
	err := query(ctx, "https://api.deepseek.com/chat/completions", key, map[string]any{
		"model": "deepseek-flash", "max_tokens": 1800, "thinking": map[string]string{"type": "disabled"},
		"messages": []map[string]string{{"role": "system", "content": instruction}, {"role": "user", "content": prompt}},
	}, &result)
	if err != nil {
		return "", err
	}
	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return "", errors.New("DeepSeek returned no text")
	}
	return result.Choices[0].Message.Content, nil
}

func respond(w http.ResponseWriter, value any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		value = map[string]string{"error": err.Error()}
	}
	json.NewEncoder(w).Encode(value)
}

func (app application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		html, _ := assets.ReadFile("web/index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	})
	for _, name := range []string{"app.js", "style.css"} {
		mux.HandleFunc("GET /"+name, func(w http.ResponseWriter, r *http.Request) {
			data, _ := assets.ReadFile("web/" + name)
			http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(data))
		})
	}
	app.register(mux, func(fn http.HandlerFunc) http.HandlerFunc { return fn })
	return mux
}

func (app application) register(mux *http.ServeMux, protect func(http.HandlerFunc) http.HandlerFunc) {
	handle := func(pattern string, fn http.HandlerFunc) { mux.HandleFunc(pattern, protect(fn)) }
	handle("GET /api/settings", func(w http.ResponseWriter, r *http.Request) {
		settings := app.config()
		ctx, cancel := context.WithTimeout(r.Context(), 2500*time.Millisecond)
		defer cancel()
		var health aceResponse[struct{ Status string }]
		err := query(ctx, settings.Engine+"/health", settings.ACEKey, nil, &health)
		respond(w, map[string]any{"engine": settings.Engine, "assistant": settings.DeepSeekKey != "", "ready": err == nil && health.Data.Status == "ok"}, nil)
	})
	handle("GET /api/library", func(w http.ResponseWriter, r *http.Request) {
		items, err := app.library()
		respond(w, items, err)
	})
	handle("GET /audio/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		id, format, ok := strings.Cut(name, ".")
		if !ok || len(id) != 26 || strings.Trim(id, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567") != "" || !validFormat(format) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", audioTypes[format])
		http.ServeFile(w, r, filepath.Join(app.dir, name))
	})
	handle("POST /api/import", func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("X-Audio-Name")
		title, format, ok := strings.Cut(name, "|")
		if !ok || !validFormat(format) {
			respond(w, nil, errors.New("choose WAV, MP3, OGG, FLAC or M4A audio"))
			return
		}
		decoded, err := url.PathUnescape(title)
		if err != nil {
			respond(w, nil, errors.New("invalid audio title"))
			return
		}
		item, err := app.save(decoded, format, r.Body)
		respond(w, item, err)
	})
	app.generationRoutes(handle)
	handle("POST /api/assist", func(w http.ResponseWriter, r *http.Request) {
		var input struct{ Prompt, Mode string }
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 24<<10)).Decode(&input); err != nil || strings.TrimSpace(input.Prompt) == "" {
			respond(w, nil, errors.New("enter a music brief first"))
			return
		}
		result, err := assist(r.Context(), input.Prompt, input.Mode, app.config().DeepSeekKey)
		respond(w, map[string]string{"text": result}, err)
	})
}

var audioTypes = map[string]string{
	"wav": "audio/wav", "mp3": "audio/mpeg", "ogg": "audio/ogg", "flac": "audio/flac", "m4a": "audio/mp4",
}

func validFormat(format string) bool { return audioTypes[format] != "" }

func audioMatches(format string, data []byte) bool {
	if len(data) < 12 {
		return false
	}
	switch format {
	case "wav":
		return string(data[0:4]) == "RIFF" && string(data[8:12]) == "WAVE"
	case "mp3":
		return string(data[0:3]) == "ID3" || data[0] == 255 && data[1]&224 == 224
	case "ogg":
		return string(data[0:4]) == "OggS"
	case "flac":
		return string(data[0:4]) == "fLaC"
	case "m4a":
		return string(data[4:8]) == "ftyp"
	}
	return false
}

func validEngine(engine string) bool {
	base, err := url.Parse(engine)
	return err == nil && base.Scheme == "http" && (base.Hostname() == "127.0.0.1" || base.Hostname() == "::1") && base.User == nil && base.Path == "" && base.RawQuery == "" && base.Fragment == "" && !base.ForceQuery
}

func main() {
	if len(os.Args) > 1 {
		if err := runManaged(); err != nil {
			log.Fatal(err)
		}
		return
	}
	settings := engineSettings{Engine: os.Getenv("MYMUSIC_ACE_URL"), ACEKey: os.Getenv("ACESTEP_API_KEY"), DeepSeekKey: os.Getenv("DEEPSEEK_API_KEY")}
	if settings.Engine == "" {
		settings.Engine = defaultEngine
	}
	if !validEngine(settings.Engine) {
		log.Fatal("MYMUSIC_ACE_URL must be a loopback HTTP origin")
	}
	config, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := os.Getenv("MYMUSIC_DATA")
	if dir == "" {
		dir = filepath.Join(config, "One Workspace", "My Music", "Audio")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	host := listener.Addr().String()
	origin := "http://" + host
	prefix := "/" + rand.Text()
	mux := http.StripPrefix(prefix, application{dir: dir, config: func() engineSettings { return settings }}.routes())
	server := http.Server{ReadHeaderTimeout: 10 * time.Second, MaxHeaderBytes: 16384, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != host || !strings.HasPrefix(r.URL.Path, prefix+"/") || (r.Method == "POST" && r.Header.Get("Origin") != origin) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; media-src 'self' blob:; style-src 'self' 'unsafe-inline'; frame-ancestors 'none'; base-uri 'none'")
		mux.ServeHTTP(w, r)
	})}
	launch := origin + prefix + "/"
	fmt.Println("Vuno • " + launch)
	if os.Getenv("MYMUSIC_NO_BROWSER") != "1" {
		edge := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "Edge", "Application", "msedge.exe")
		if _, err := os.Stat(edge); err == nil {
			if err := exec.Command(edge, "--app="+launch).Start(); err != nil {
				log.Print(err)
			}
		} else {
			if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", launch).Start(); err != nil {
				log.Print(err)
			}
		}
	}
	log.Fatal(server.Serve(listener))
}
