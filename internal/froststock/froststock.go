package froststock

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const (
	Version  = "0.1.0-preview.1"
	photoMax = 3 << 20
)

//go:embed web
var web embed.FS

func UI() fs.FS {
	files, err := fs.Sub(web, "web")
	if err != nil {
		panic(err)
	}
	return files
}

type entry struct {
	ID           string `json:"id"`
	Revision     int64  `json:"revision"`
	Created      int64  `json:"created"`
	Used         int64  `json:"used,omitempty"`
	Name         string `json:"name,omitempty"`
	Category     string `json:"category,omitempty"`
	Quantity     int    `json:"quantity"`
	Expiry       string `json:"expiry,omitempty"`
	DateType     string `json:"dateType,omitempty"`
	Review       bool   `json:"review,omitempty"`
	ReviewReason string `json:"reviewReason,omitempty"`
	RawItem      string `json:"rawItem,omitempty"`
	RawDate      string `json:"rawDate,omitempty"`
	VisionItem   string `json:"visionItem,omitempty"`
	VisionDate   string `json:"visionDate,omitempty"`
	State        string `json:"state"`
	Job          string `json:"job"`
}

var identifier = regexp.MustCompile(`^[A-Za-z0-9_-]{12,80}$`)

func revisionTag(revision int64) string {
	return `"` + strconv.FormatInt(revision, 10) + `"`
}

func (e entry) valid() bool {
	if !identifier.MatchString(e.ID) || !identifier.MatchString(e.Job) ||
		e.Quantity < 1 || e.Quantity > 999 || e.Created <= 0 || e.Used < 0 ||
		(e.State != "queued" && e.State != "ready") {
		return false
	}
	if e.Expiry != "" {
		date, err := time.Parse("2006-01-02", e.Expiry)
		if err != nil || date.Year() < 1900 || date.Year() > 2199 {
			return false
		}
	}
	return len(e.Name) <= 400 && len(e.Category) <= 160 && len(e.DateType) <= 80 &&
		len(e.ReviewReason) <= 1200 && len(e.RawItem) <= 20000 && len(e.RawDate) <= 20000 &&
		len(e.VisionItem) <= 4000 && len(e.VisionDate) <= 4000
}

type drive struct {
	root string
	mu   sync.Mutex
}

func Register(host *apphost.Host, root string) error {
	if !filepath.IsAbs(root) || filepath.Base(root) != "Froststock" || host.Config.Issuer == "" || host.Config.Owner == "" {
		return errors.New("Froststock requires My Drive and a managed MyAccounts owner")
	}
	store := &drive{root: root}
	owner := func(handler http.HandlerFunc) http.HandlerFunc {
		return host.Owner(func(w http.ResponseWriter, r *http.Request) {
			if !apphost.Equal(apphost.Current(r).User.Subject, host.Config.Owner) {
				apphost.Problem(w, 403, "Open Froststock with the Workspace owner account")
				return
			}
			handler(w, r)
		})
	}
	host.Mux.HandleFunc("GET /api/items", owner(store.list))
	host.Mux.HandleFunc("POST /api/items/{id}", owner(store.create))
	host.Mux.HandleFunc("PUT /api/items/{id}", owner(store.update))
	host.Mux.HandleFunc("DELETE /api/items/{id}", owner(store.delete))
	host.Mux.HandleFunc("GET /api/items/{id}/{photo}", owner(store.photo))
	return nil
}

func (d *drive) open() (*os.Root, error) {
	parent := filepath.Dir(d.root)
	if err := platform.RejectReparse(parent); err != nil {
		return nil, err
	}
	if err := os.Mkdir(d.root, 0700); err != nil && !errors.Is(err, fs.ErrExist) {
		return nil, err
	}
	if err := platform.RejectReparse(d.root); err != nil {
		return nil, err
	}
	return os.OpenRoot(d.root)
}

func read(root *os.Root, id string) (entry, error) {
	var result entry
	directory, err := root.Lstat(id)
	if err != nil {
		return result, err
	}
	if !directory.IsDir() || directory.Mode()&os.ModeSymlink != 0 {
		return result, errors.New("invalid inventory directory")
	}
	file, err := openRegular(root, path.Join(id, "item.json"))
	if err != nil {
		return result, err
	}
	defer file.Close()
	err = json.NewDecoder(io.LimitReader(file, 65536)).Decode(&result)
	if err == nil && (result.ID != id || result.Revision < 1 || !result.valid()) {
		err = errors.New("invalid inventory file")
	}
	return result, err
}

func openRegular(root *os.Root, name string) (*os.File, error) {
	info, err := root.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("inventory file is not a regular file")
	}
	return root.Open(name)
}

func write(root *os.Root, name string, data []byte) error {
	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err == nil {
		err = f.Sync()
	}
	return errors.Join(err, f.Close())
}

func tempName() string {
	var data [16]byte
	_, _ = rand.Read(data[:])
	return ".froststock-" + hex.EncodeToString(data[:])
}

func unavailable(w http.ResponseWriter) {
	apphost.Problem(w, 503, "Froststock could not access its My Drive folder. Try again.")
}

func (d *drive) list(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	root, err := d.open()
	if err != nil {
		unavailable(w)
		return
	}
	defer root.Close()
	names, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		unavailable(w)
		return
	}
	items := make([]entry, 0, len(names))
	incomplete := 0
	for _, name := range names {
		if name.IsDir() && identifier.MatchString(name.Name()) {
			if value, err := read(root, name.Name()); err == nil {
				items = append(items, value)
			} else if !errors.Is(err, fs.ErrNotExist) {
				incomplete++
			}
		}
	}
	apphost.JSON(w, 200, map[string]any{"items": items, "incomplete": incomplete})
}

func (d *drive) create(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !identifier.MatchString(id) {
		http.NotFound(w, r)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2*photoMax+65536)
	parts, err := r.MultipartReader()
	if err != nil {
		apphost.Problem(w, 400, "Send both photos and the item record")
		return
	}
	data := make(map[string][]byte, 3)
	for {
		part, err := parts.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			apphost.Problem(w, 400, "The photo pair could not be read")
			return
		}
		name := part.FormName()
		max := photoMax
		switch name {
		case "item", "expiry":
		case "record":
			max = 65536
		default:
			apphost.Problem(w, 400, "Unexpected upload field")
			return
		}
		bytes, err := io.ReadAll(io.LimitReader(part, int64(max+1)))
		if err != nil || len(bytes) > max || data[name] != nil {
			apphost.Problem(w, 400, "Check the photo sizes and try again")
			return
		}
		data[name] = bytes
	}
	var value entry
	if json.Unmarshal(data["record"], &value) != nil || value.ID != id || value.Revision != 0 || !value.valid() ||
		http.DetectContentType(data["item"]) != "image/jpeg" || http.DetectContentType(data["expiry"]) != "image/jpeg" {
		apphost.Problem(w, 400, "Both JPEG photos and valid item details are required")
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	root, err := d.open()
	if err != nil {
		unavailable(w)
		return
	}
	defer root.Close()
	if old, err := read(root, id); err == nil {
		apphost.JSON(w, 200, old)
		return
	}
	if _, err := root.Lstat(id); !errors.Is(err, fs.ErrNotExist) {
		apphost.Problem(w, 409, "A folder with this item ID already exists in My Drive")
		return
	}
	stage := tempName()
	if err := root.Mkdir(stage, 0700); err != nil {
		unavailable(w)
		return
	}
	defer root.RemoveAll(stage)
	value.Revision = 1
	record, _ := json.MarshalIndent(value, "", "  ")
	for name, b := range map[string][]byte{"item.json": record, "item.jpg": data["item"], "expiry.jpg": data["expiry"]} {
		if err := write(root, path.Join(stage, name), b); err != nil {
			unavailable(w)
			return
		}
	}
	if err := root.Rename(stage, id); err != nil {
		unavailable(w)
		return
	}
	apphost.JSON(w, 201, value)
}

func (d *drive) update(w http.ResponseWriter, r *http.Request) {
	var value entry
	r.Body = http.MaxBytesReader(w, r.Body, 65536)
	if !apphost.Decode(w, r, &value) {
		return
	}
	id := r.PathValue("id")
	if id != value.ID || !value.valid() {
		apphost.Problem(w, 400, "Invalid item details")
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	root, err := d.open()
	if err != nil {
		unavailable(w)
		return
	}
	defer root.Close()
	old, err := read(root, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if old.Revision != value.Revision {
		apphost.Problem(w, 409, "This item changed on another device. Refresh and try again.")
		return
	}
	value.Created = old.Created
	value.Revision++
	temp := path.Join(id, tempName())
	defer root.Remove(temp)
	record, _ := json.MarshalIndent(value, "", "  ")
	if err := write(root, temp, record); err != nil {
		unavailable(w)
		return
	}
	if err := root.Rename(temp, path.Join(id, "item.json")); err != nil {
		unavailable(w)
		return
	}
	apphost.JSON(w, 200, value)
}

func (d *drive) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !identifier.MatchString(id) {
		http.NotFound(w, r)
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	root, err := d.open()
	if err != nil {
		unavailable(w)
		return
	}
	defer root.Close()
	value, err := read(root, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if r.Header.Get("If-Match") != revisionTag(value.Revision) {
		apphost.Problem(w, 409, "This item changed on another device. Refresh and try again.")
		return
	}
	for _, name := range []string{"item.jpg", "expiry.jpg", "item.json"} {
		if err := root.Remove(path.Join(id, name)); err != nil && !errors.Is(err, fs.ErrNotExist) {
			unavailable(w)
			return
		}
	}
	// Preserve any extra files the user put in the item folder through Drive.
	_ = root.Remove(id)
	w.WriteHeader(204)
}

func (d *drive) photo(w http.ResponseWriter, r *http.Request) {
	id, name := r.PathValue("id"), r.PathValue("photo")
	if !identifier.MatchString(id) || (name != "item.jpg" && name != "expiry.jpg") {
		http.NotFound(w, r)
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	root, err := d.open()
	if err != nil {
		unavailable(w)
		return
	}
	defer root.Close()
	if _, err := read(root, id); err != nil {
		http.NotFound(w, r)
		return
	}
	file, err := openRegular(root, path.Join(id, name))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() > photoMax {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename="+strings.TrimSpace(name))
	http.ServeContent(w, r, name, time.Time{}, file)
}
