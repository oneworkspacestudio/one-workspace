package gamestudio

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const Version = "0.1.0-preview.1"

const schema = `
CREATE TABLE games (
 id TEXT PRIMARY KEY, account TEXT NOT NULL, latest TEXT NOT NULL, updated INTEGER NOT NULL
);
CREATE TABLE builds (
 id TEXT PRIMARY KEY, game TEXT NOT NULL REFERENCES games(id) ON DELETE CASCADE, spec TEXT NOT NULL, page BLOB NOT NULL
);
CREATE TABLE usage (day TEXT PRIMARY KEY, turns INTEGER NOT NULL);
`

var (
	idPattern      = regexp.MustCompile(`^[a-f0-9]{48}$`)
	accountPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)
	repoPattern    = regexp.MustCompile(`^[a-zA-Z0-9_-][a-zA-Z0-9_.-]{0,79}/[a-zA-Z0-9_-][a-zA-Z0-9_.-]{0,79}$`)
	modelPattern   = regexp.MustCompile(`(?i)^deepseek-r1(?::[0-9]+(?:\.[0-9]+)?b)?$`)
	privateText    = regexp.MustCompile(`(?i)https?://|www\.|@|\+?[0-9][0-9(). -]{6,}[0-9]`)
)

type settings struct {
	Provider   string `json:"provider"`
	LocalModel string `json:"localModel"`
	Child      string `json:"child"`
	GitRepo    string `json:"gitRepo"`
	CloudKey   string `json:"cloudKey,omitempty"`
	GitKey     string `json:"gitKey,omitempty"`
}

func (s settings) allows(u apphost.User) bool {
	return u.Role == "owner" || u.ID == s.Child
}

type Studio struct {
	host                       *apphost.Host
	mu                         sync.Mutex
	settings                   settings
	turn                       chan struct{}
	client                     *http.Client
	cloudURL, localURL, gitURL string
}

func New(h *apphost.Host) (*Studio, error) {
	if err := h.Migrate(1, schema); err != nil {
		return nil, err
	}
	origin, err := url.Parse(h.Config.Origin)
	if err != nil {
		return nil, err
	}
	s := &Studio{
		host: h, turn: make(chan struct{}, 1),
		client:   &http.Client{Timeout: 50 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }},
		cloudURL: "https://api.deepseek.com/chat/completions",
		localURL: "http://127.0.0.1:11434/v1/chat/completions",
		gitURL:   origin.Scheme + "://" + origin.Host + "/git/api/v1",
	}
	if err := platform.LoadSecret(s.settingsPath(), &s.settings); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("Game Studio settings could not be opened")
	}
	h.BackupPaths = append(h.BackupPaths, "studio-settings")
	h.Mux.HandleFunc("GET /api/studio", h.Require(s.status))
	h.Mux.HandleFunc("PUT /api/settings", h.Owner(s.configure))
	h.Mux.HandleFunc("GET /api/games", h.Require(s.list))
	h.Mux.HandleFunc("POST /api/turn", h.Require(s.generate))
	h.Mux.HandleFunc("GET /play/{build}/", h.Require(s.preview))
	h.Mux.HandleFunc("GET /api/builds/{build}/archive", h.Owner(s.download))
	h.Mux.HandleFunc("POST /api/builds/{build}/git", h.Owner(s.commit))
	return s, nil
}

func (s *Studio) settingsPath() string { return filepath.Join(s.host.Config.Data, "studio-settings") }

func decode(r io.Reader, dst any) error {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	if err := d.Decode(dst); err != nil {
		return err
	}
	if d.Decode(new(any)) != io.EOF {
		return errors.New("expected a single JSON object")
	}
	return nil
}

func read(w http.ResponseWriter, r *http.Request, dst any) bool {
	if decode(http.MaxBytesReader(w, r.Body, 16384), dst) != nil {
		apphost.Problem(w, 400, "Check the fields and try again.")
		return false
	}
	return true
}

func safeText(text string, limit int) bool {
	return strings.TrimSpace(text) != "" && utf8.RuneCountInString(text) <= limit &&
		!strings.ContainsFunc(text, func(r rune) bool { return unicode.IsControl(r) && r != '\n' && r != '\t' }) &&
		!privateText.MatchString(text)
}

func (s *Studio) snapshot() settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.settings
}

func (s *Studio) status(w http.ResponseWriter, r *http.Request) {
	u := apphost.Current(r).User
	v := s.snapshot()
	reply := map[string]any{"allowed": v.allows(u), "provider": v.Provider}
	if u.Role == "owner" {
		reply["child"], reply["localModel"], reply["gitRepo"] = v.Child, v.LocalModel, v.GitRepo
		reply["hasCloudKey"], reply["hasGitKey"] = v.CloudKey != "", v.GitKey != ""
	}
	apphost.JSON(w, 200, reply)
}

func (s *Studio) configure(w http.ResponseWriter, r *http.Request) {
	var in struct {
		settings
		CloudConsent bool `json:"cloudConsent"`
		ClearGit     bool `json:"clearGit"`
		ClearCloud   bool `json:"clearCloud"`
	}
	if !read(w, r, &in) {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	next := in.settings
	next.Child = strings.TrimSpace(next.Child)
	next.GitRepo = strings.TrimSpace(next.GitRepo)
	next.LocalModel = strings.TrimSpace(next.LocalModel)
	if next.CloudKey == "" && !in.ClearCloud {
		next.CloudKey = s.settings.CloudKey
	}
	if next.GitKey == "" && !in.ClearGit {
		next.GitKey = s.settings.GitKey
	}
	if next.Child != "" && !accountPattern.MatchString(next.Child) ||
		next.GitRepo != "" && !repoPattern.MatchString(next.GitRepo) ||
		next.LocalModel != "" && !modelPattern.MatchString(next.LocalModel) ||
		!oneOf(next.Provider, "", "local", "cloud") ||
		!validToken(next.CloudKey) || !validToken(next.GitKey) ||
		next.Provider == "local" && next.LocalModel == "" ||
		next.Provider == "cloud" && (next.CloudKey == "" || !in.CloudConsent) {
		apphost.Problem(w, 400, "Check the account code, model, repository and keys. Cloud mode requires guardian consent.")
		return
	}
	if err := platform.SaveSecret(s.settingsPath(), next); err != nil {
		apphost.Problem(w, 500, "Settings could not be saved.")
		return
	}
	s.settings = next
	s.host.Audit(apphost.Current(r).User.ID, "settings-updated", "gamestudio")
	apphost.JSON(w, 200, map[string]bool{"saved": true})
}

func validToken(t string) bool {
	return len(t) <= 250 && !strings.ContainsFunc(t, func(r rune) bool { return r < 33 || r > 126 })
}

type game struct {
	ID       string `json:"id"`
	Build    string `json:"build"`
	Spec     Spec   `json:"spec"`
	Editable bool   `json:"editable"`
}

func (s *Studio) list(w http.ResponseWriter, r *http.Request) {
	u := apphost.Current(r).User
	if !s.snapshot().allows(u) {
		apphost.Problem(w, 403, "Ask the owner to approve your account code.")
		return
	}
	rows, err := s.host.DB.QueryContext(r.Context(), `SELECT g.id,b.id,b.spec,g.account FROM games g JOIN builds b ON b.id=g.latest WHERE g.account=? OR ? ORDER BY g.updated DESC LIMIT 50`, u.ID, u.Role == "owner")
	if err != nil {
		apphost.Problem(w, 500, "Games could not be loaded.")
		return
	}
	defer rows.Close()
	result := []game{}
	for rows.Next() {
		var g game
		var data []byte
		var account string
		if rows.Scan(&g.ID, &g.Build, &data, &account) != nil || json.Unmarshal(data, &g.Spec) != nil {
			apphost.Problem(w, 500, "A saved game could not be read.")
			return
		}
		g.Editable = account == u.ID
		result = append(result, g)
	}
	if rows.Err() != nil {
		apphost.Problem(w, 500, "Games could not be loaded.")
		return
	}
	apphost.JSON(w, 200, result)
}

type turnInput struct {
	Game    string `json:"game"`
	Build   string `json:"build"`
	Prompt  string `json:"prompt"`
	Context string `json:"context"`
}

func (s *Studio) generate(w http.ResponseWriter, r *http.Request) {
	u := apphost.Current(r).User
	cfg := s.snapshot()
	if !cfg.allows(u) {
		apphost.Problem(w, 403, "Ask the owner to approve your account code.")
		return
	}
	if cfg.Provider == "" {
		apphost.Problem(w, 409, "Ask the owner to turn on a DeepSeek model.")
		return
	}
	var in turnInput
	if !read(w, r, &in) {
		return
	}
	if !safeText(in.Prompt, 600) || in.Context != "" && !safeText(in.Context, 1600) {
		apphost.Problem(w, 400, "Describe your game in up to 600 characters and leave out personal details and links.")
		return
	}
	select {
	case s.turn <- struct{}{}:
		defer func() { <-s.turn }()
	default:
		apphost.Problem(w, 429, "A game idea is already being worked on. Try again in a moment.")
		return
	}
	ctx := r.Context()
	var spec Spec
	if in.Game != "" {
		if !idPattern.MatchString(in.Game) || !idPattern.MatchString(in.Build) {
			http.NotFound(w, r)
			return
		}
		var last string
		var raw []byte
		if err := s.host.DB.QueryRowContext(ctx, `SELECT g.latest,b.spec FROM games g JOIN builds b ON b.id=g.latest WHERE g.id=? AND g.account=?`, in.Game, u.ID).Scan(&last, &raw); err != nil || json.Unmarshal(raw, &spec) != nil {
			http.NotFound(w, r)
			return
		}
		if in.Build != last {
			apphost.Problem(w, 409, "This game changed. Reload its latest version before editing.")
			return
		}
	}
	today := time.Now().UTC().Format("2006-01-02")
	_, _ = s.host.DB.ExecContext(ctx, `DELETE FROM usage WHERE day < ?`, today)
	result, err := s.host.DB.ExecContext(ctx, `INSERT INTO usage VALUES(?,1) ON CONFLICT(day) DO UPDATE SET turns=turns+1 WHERE turns<30`, today)
	if err != nil {
		apphost.Problem(w, 500, "Could not check the daily build allowance.")
		return
	}
	if n, _ := result.RowsAffected(); n != 1 {
		apphost.Problem(w, 429, "Today’s 30 idea turns are used. Try again tomorrow.")
		return
	}
	out, err := s.plan(ctx, cfg, in, spec)
	if err != nil {
		apphost.Problem(w, 502, err.Error())
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.settings != cfg || !s.host.SessionActive(r, u.ID) {
		apphost.Problem(w, 409, "Your session or the owner’s settings changed. Reload before continuing.")
		return
	}
	if out.Game == nil {
		apphost.JSON(w, 200, map[string]any{"message": out.Message})
		return
	}
	g, err := s.saveBuild(ctx, u.ID, in.Game, *out.Game)
	if err != nil {
		apphost.Problem(w, 409, "Could not save this game. The studio keeps up to 50 games and 30 builds per game.")
		return
	}
	apphost.JSON(w, 200, map[string]any{"message": out.Message, "game": g})
}

func (s *Studio) saveBuild(ctx context.Context, account, gameID string, spec Spec) (game, error) {
	var g game
	page, err := render(spec)
	if err != nil {
		return g, err
	}
	raw, err := json.Marshal(spec)
	if err != nil {
		return g, err
	}
	tx, err := s.host.DB.BeginTx(ctx, nil)
	if err != nil {
		return g, err
	}
	defer tx.Rollback()
	build := apphost.ID()
	if gameID == "" {
		gameID = apphost.ID()
		res, e := tx.ExecContext(ctx, `INSERT INTO games SELECT ?,?,?,? WHERE (SELECT count(*) FROM games)<50`, gameID, account, build, time.Now().UnixMilli())
		if err = oneRow(res, e); err != nil {
			return g, err
		}
	} else if _, err = tx.ExecContext(ctx, `UPDATE games SET latest=?,updated=? WHERE id=?`, build, time.Now().UnixMilli(), gameID); err != nil {
		return g, err
	}
	res, err := tx.ExecContext(ctx, `INSERT INTO builds SELECT ?,?,?,? WHERE (SELECT count(*) FROM builds WHERE game=?)<30`, build, gameID, raw, page, gameID)
	if err = oneRow(res, err); err != nil {
		return g, err
	}
	return game{ID: gameID, Build: build, Spec: spec, Editable: true}, tx.Commit()
}

func oneRow(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	if n, e := result.RowsAffected(); e != nil || n != 1 {
		return errors.New("storage limit reached")
	}
	return nil
}

func (s *Studio) getBuild(w http.ResponseWriter, r *http.Request) ([]byte, string, bool) {
	u := apphost.Current(r).User
	build := r.PathValue("build")
	if !idPattern.MatchString(build) || !s.snapshot().allows(u) {
		http.NotFound(w, r)
		return nil, "", false
	}
	var body []byte
	var id string
	err := s.host.DB.QueryRowContext(r.Context(), `SELECT b.page,g.id FROM builds b JOIN games g ON b.game=g.id WHERE b.id=? AND (g.account=? OR ?)`, build, u.ID, u.Role == "owner").Scan(&body, &id)
	if err != nil {
		http.NotFound(w, r)
		return nil, "", false
	}
	return body, id, true
}

func (s *Studio) preview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/play/"+r.PathValue("build")+"/" {
		http.NotFound(w, r)
		return
	}
	if page, _, ok := s.getBuild(w, r); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", gamePolicy+"; sandbox allow-scripts; frame-ancestors 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		_, _ = w.Write(page)
	}
}

func (s *Studio) download(w http.ResponseWriter, r *http.Request) {
	page, _, ok := s.getBuild(w, r)
	if !ok {
		return
	}
	data, err := archive(page)
	if err != nil {
		apphost.Problem(w, 500, "Could not package the game.")
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="game.zip"`)
	_, _ = w.Write(data)
}
