package pet

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const Version = "0.1.0-preview.3"

//go:embed web/dist
var web embed.FS

func UI() fs.FS {
	root, err := fs.Sub(web, "web/dist")
	if err != nil {
		panic(err)
	}
	return root
}

type Habitat struct {
	host        *apphost.Host
	mu          sync.Mutex
	world       *World
	settings    Settings
	brain       brain
	nextTurn    time.Time
	sharedNext  bool
	cancel      context.CancelFunc
	generation  uint64
	pending     Decision
	modelStatus string
}

func New(h *apphost.Host) (*Habitat, error) {
	if err := h.Migrate(1, `CREATE TABLE pet_world(id INTEGER PRIMARY KEY CHECK(id=1), snapshot BLOB NOT NULL);`); err != nil {
		return nil, err
	}
	p := &Habitat{host: h, brain: newBrain(os.Getenv("PET_AI_SERVICE_PORT")), settings: Settings{Mode: "instincts"}}
	var raw []byte
	switch err := h.DB.QueryRow(`SELECT snapshot FROM pet_world WHERE id=1`).Scan(&raw); {
	case err == nil:
		p.world = &World{}
		if err = json.Unmarshal(raw, p.world); err != nil {
			return nil, err
		}
	case !errors.Is(err, sql.ErrNoRows):
		return nil, err
	}
	if err := platform.LoadSecret(p.settingsPath(), &p.settings); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("could not open Pet decision settings")
	}
	if err := p.openJournal(); err != nil {
		return nil, err
	}
	if p.world != nil && p.world.Place == "" {
		w := p.world
		remaining := w.Until
		if w.Action == "build" {
			w.beginIdea(w.fallback("build"))
		} else {
			w.start(w.Action, w.FocusStyle, w.Source)
			if w.Action == "sleep" || w.Action == "draw" {
				w.Until = remaining
			}
		}
		if err := p.save(context.Background(), *w); err != nil {
			return nil, err
		}
	}
	h.BackupPaths = append(h.BackupPaths, "decisions")
	h.Mux.HandleFunc("GET /api/world", h.Require(p.state))
	h.Mux.HandleFunc("GET /api/history", h.Require(p.history))
	h.Mux.HandleFunc("POST /api/adopt", h.Owner(p.adopt))
	h.Mux.HandleFunc("POST /api/treat", h.Require(p.treat))
	h.Mux.HandleFunc("PUT /api/decisions", h.Owner(p.configure))
	return p, nil
}

func (p *Habitat) settingsPath() string { return filepath.Join(p.host.Config.Data, "decisions") }

func (p *Habitat) save(ctx context.Context, w World) error {
	raw, err := json.Marshal(w)
	if err != nil {
		return err
	}
	tx, err := p.host.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `INSERT INTO pet_world VALUES(1,?) ON CONFLICT(id) DO UPDATE SET snapshot=excluded.snapshot`, raw); err != nil {
		return err
	}
	if err = appendJournal(ctx, tx, w.events); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	w.events = nil
	p.world = &w
	return nil
}

func (p *Habitat) advance(ctx context.Context, now time.Time) error {
	if p.world == nil {
		return nil
	}
	seconds := int64(step / time.Second)
	elapsed := (now.Unix() - p.world.Saved) / seconds
	if elapsed == 0 {
		return nil
	}
	next := *p.world
	next.Food = append([]Item{}, next.Food...)
	next.Creations = append([]Item{}, next.Creations...)
	next.Journal = append([]Event{}, next.Journal...)
	if elapsed < 0 {
		next.Saved = now.Unix()
		next.TreatAfter = min(next.TreatAfter, next.Saved+45)
		return p.save(ctx, next)
	}
	pending := p.pending
	if elapsed > 2 {
		pending = Decision{}
	}
	count := min(elapsed, int64(8*time.Hour/step))
	for i := int64(0); i < count; i++ {
		next.advance(&pending)
	}
	next.catchUp(now, elapsed, count)
	if err := p.save(ctx, next); err != nil {
		return err
	}
	p.pending = pending
	return nil
}

func (p *Habitat) Run(ctx context.Context) error {
	var workers sync.WaitGroup
	brainCtx, stop := context.WithCancel(ctx)
	defer func() { stop(); workers.Wait() }()
	ticker := time.NewTicker(step)
	defer ticker.Stop()
	for {
		p.host.Freeze.RLock()
		p.mu.Lock()
		now := time.Now()
		err := p.advance(ctx, now)
		if err == nil {
			p.schedule(brainCtx, &workers, now)
		}
		p.mu.Unlock()
		p.host.Freeze.RUnlock()
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func (p *Habitat) state(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := map[string]any{"world": p.world, "mode": p.settings.Mode, "modelStatus": p.modelStatus}
	if apphost.Current(r).User.Role == "owner" {
		result["localModel"] = p.settings.Model
		result["paired"] = p.settings.ServiceKey != ""
		result["hasService"] = p.brain.sharedURL != ""
	}
	apphost.JSON(w, 200, result)
}

func (p *Habitat) adopt(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 22 || strings.ContainsFunc(name, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsMark(c) && c != ' ' && c != '-' && c != '\''
	}) {
		apphost.Problem(w, 400, "Choose a name of up to 22 letters and numbers.")
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.world != nil {
		apphost.Problem(w, 409, "A pet already lives in this hollow.")
		return
	}
	if err := p.save(r.Context(), newWorld(name, time.Now())); err != nil {
		apphost.Problem(w, 500, "The habitat could not be saved.")
		return
	}
	apphost.JSON(w, 201, p.world)
}

func (p *Habitat) treat(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.world == nil {
		apphost.Problem(w, 409, "The hollow is still waiting for its pet.")
		return
	}
	if p.world.TreatAfter > p.world.Saved {
		apphost.Problem(w, 429, "Give your pet some time to enjoy the last treat.")
		return
	}
	next := *p.world
	next.Journal = append([]Event{}, next.Journal...)
	next.treat()
	if err := p.save(r.Context(), next); err != nil {
		apphost.Problem(w, 500, "The treat could not be saved.")
		return
	}
	apphost.JSON(w, 200, p.world)
}
