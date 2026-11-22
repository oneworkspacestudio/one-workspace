package fieldnotes

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const Version = "0.1.0-preview.1"

//go:embed ui/dist
var web embed.FS

func UI() fs.FS {
	files, err := fs.Sub(web, "ui/dist")
	if err != nil {
		panic(err)
	}
	return files
}

const schema = `
CREATE TABLE research_jobs (
 id TEXT PRIMARY KEY, account TEXT NOT NULL, kind TEXT NOT NULL,
 query TEXT NOT NULL, provider TEXT NOT NULL, status TEXT NOT NULL,
 evidence TEXT NOT NULL DEFAULT '[]', error TEXT NOT NULL DEFAULT '',
 created INTEGER NOT NULL
);
CREATE UNIQUE INDEX fieldnotes_active ON research_jobs(account) WHERE status IN ('queued','running');
CREATE INDEX fieldnotes_recent ON research_jobs(account, created DESC);
CREATE TABLE fieldnotes_articles (
 id TEXT PRIMARY KEY, account TEXT NOT NULL, content TEXT NOT NULL, created INTEGER NOT NULL
);
CREATE INDEX fieldnotes_library ON fieldnotes_articles(account, created DESC);
CREATE TABLE fieldnotes_saved (
 account TEXT NOT NULL, id TEXT NOT NULL, PRIMARY KEY(account,id)
);`

type secrets struct {
	Account  string `json:"account"`
	DeepSeek string `json:"deepseek"`
	Browser  string `json:"browser"`
}

type job struct {
	ID       string   `json:"id"`
	Kind     string   `json:"kind"`
	Query    string   `json:"query"`
	Provider string   `json:"provider"`
	Status   string   `json:"status"`
	Evidence []source `json:"sources"`
	Error    string   `json:"error,omitempty"`
}

type service struct {
	*apphost.Host
	client   *http.Client
	browser  string
	wake     chan struct{}
	configMu sync.Mutex
}

func Register(host *apphost.Host) (apphost.Application, error) {
	if err := host.Migrate(1, schema); err != nil {
		return nil, err
	}
	if _, err := host.DB.Exec(`UPDATE research_jobs SET status='failed', evidence='[]',
		error='Workspace restarted. Run this request again.' WHERE status IN ('queued','running')`); err != nil {
		return nil, err
	}
	s := &service{
		Host: host, browser: os.Getenv("FIELDNOTES_BROWSER_URL"), wake: make(chan struct{}, 1),
		client: &http.Client{Timeout: 95 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}},
	}
	host.BackupPaths = append(host.BackupPaths, "providers.secret")
	host.Mux.HandleFunc("GET /api/fieldnotes", host.Owner(s.library))
	host.Mux.HandleFunc("PUT /api/providers", host.Owner(s.configure))
	host.Mux.HandleFunc("POST /api/research", host.Owner(s.enqueue))
	host.Mux.HandleFunc("PUT /api/saved", host.Owner(s.bookmark))
	return s, nil
}

func (s *service) credentials(account string) (secrets, error) {
	var values secrets
	err := apphost.LoadPrivateSetting(filepath.Join(s.Config.Data, "providers.secret"), &values)
	if errors.Is(err, os.ErrNotExist) || (err == nil && values.Account != account) {
		return secrets{Account: account}, nil
	}
	return values, err
}

func (s *service) configure(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var input struct{ Provider, Key string }
	if !apphost.Decode(w, r, &input) {
		return
	}
	if !validKey(input.Key) || (input.Provider != "deepseek" && input.Provider != "workspace") {
		apphost.Problem(w, 400, "Enter a valid provider credential.")
		return
	}
	// A single owner can manage providers, including from concurrent tabs.
	s.configMu.Lock()
	defer s.configMu.Unlock()
	values, err := s.credentials(apphost.Current(r).User.ID)
	if err != nil {
		apphost.Problem(w, 500, "Could not read provider settings.")
		return
	}
	if input.Provider == "deepseek" {
		values.DeepSeek = input.Key
	} else {
		values.Browser = input.Key
	}
	if err = apphost.SavePrivateSetting(filepath.Join(s.Config.Data, "providers.secret"), values); err != nil {
		apphost.Problem(w, 500, "Could not save provider settings.")
		return
	}
	apphost.JSON(w, 200, map[string]bool{"saved": true})
}

func (s *service) library(w http.ResponseWriter, r *http.Request) {
	ctx, account := r.Context(), apphost.Current(r).User.ID
	keys, err := s.credentials(account)
	if err != nil {
		apphost.Problem(w, 500, "Could not load provider settings.")
		return
	}
	jobs, articles, saved, err := s.readLibrary(ctx, account)
	if err != nil {
		apphost.Problem(w, 500, "Could not load your reading room.")
		return
	}
	apphost.JSON(w, 200, map[string]any{
		"jobs": jobs, "articles": articles, "saved": saved, "deepseek": keys.DeepSeek != "",
		"browser": s.browser != "", "browserReady": s.browser != "" && keys.Browser != "",
	})
}

func (s *service) readLibrary(ctx context.Context, account string) ([]job, []article, []string, error) {
	jobs, articles, saved := []job{}, []article{}, []string{}
	rows, err := s.DB.QueryContext(ctx, `SELECT id,kind,query,provider,status,evidence,error
		FROM research_jobs WHERE account=? ORDER BY created DESC,id DESC LIMIT 10`, account)
	if err != nil {
		return nil, nil, nil, err
	}
	for rows.Next() {
		var row job
		var raw string
		err = rows.Scan(&row.ID, &row.Kind, &row.Query, &row.Provider, &row.Status, &raw, &row.Error)
		if err == nil {
			err = json.Unmarshal([]byte(raw), &row.Evidence)
		}
		if err != nil {
			rows.Close()
			return nil, nil, nil, err
		}
		for i := range row.Evidence {
			row.Evidence[i].Text = clip(row.Evidence[i].Text, 420)
		}
		jobs = append(jobs, row)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, nil, err
	}
	rows, err = s.DB.QueryContext(ctx, `SELECT content FROM fieldnotes_articles
		WHERE account=? ORDER BY created DESC,id DESC`, account)
	if err != nil {
		return nil, nil, nil, err
	}
	for rows.Next() {
		var raw string
		var row article
		err = rows.Scan(&raw)
		if err == nil {
			err = json.Unmarshal([]byte(raw), &row)
		}
		if err != nil {
			rows.Close()
			return nil, nil, nil, err
		}
		articles = append(articles, row)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, nil, err
	}
	rows, err = s.DB.QueryContext(ctx, "SELECT id FROM fieldnotes_saved WHERE account=?", account)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, nil, nil, err
		}
		saved = append(saved, id)
	}
	return jobs, articles, saved, rows.Err()
}

func (s *service) enqueue(w http.ResponseWriter, r *http.Request) {
	ctx, account := r.Context(), apphost.Current(r).User.ID
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var input struct {
		Query, Provider, Search string
		Indices                 []int
	}
	if !apphost.Decode(w, r, &input) {
		return
	}
	keys, err := s.credentials(account)
	if err != nil {
		apphost.Problem(w, 500, "Could not open provider settings.")
		return
	}
	kind, query, provider := "search", strings.TrimSpace(input.Query), input.Provider
	evidence := []source{}
	if input.Search == "" {
		if query == "" || utf8.RuneCountInString(query) > 128 || (provider != "deepseek" && provider != "workspace") {
			apphost.Problem(w, 400, "Enter a research query of 1 to 128 characters and choose a provider.")
			return
		}
		if (provider == "deepseek" && keys.DeepSeek == "") || (provider == "workspace" && (keys.Browser == "" || s.browser == "")) {
			apphost.Problem(w, 409, "Connect the chosen search provider in Settings.")
			return
		}
	} else {
		if keys.DeepSeek == "" {
			apphost.Problem(w, 409, "Connect DeepSeek in Settings to generate a briefing.")
			return
		}
		var raw string
		err := s.DB.QueryRowContext(ctx, `SELECT query,provider,evidence FROM research_jobs
			WHERE id=? AND account=? AND kind='search' AND status='done'`,
			input.Search, account).Scan(&query, &provider, &raw)
		var candidates []source
		if err != nil || json.Unmarshal([]byte(raw), &candidates) != nil || len(input.Indices) == 0 || len(input.Indices) > 5 {
			apphost.Problem(w, 400, "Choose up to five sources from a completed search.")
			return
		}
		seen := map[int]bool{}
		for _, i := range input.Indices {
			if i < 0 || i >= len(candidates) || seen[i] {
				apphost.Problem(w, 400, "The selected sources are invalid.")
				return
			}
			seen[i] = true
			evidence = append(evidence, candidates[i])
		}
		kind = "brief"
	}
	raw, _ := json.Marshal(evidence)
	id := hex.EncodeToString(randomBytes())
	result, err := s.DB.ExecContext(ctx, `INSERT INTO research_jobs(id,account,kind,query,provider,status,evidence,created)
		VALUES(?,?,?,?,?,'queued',?,?) ON CONFLICT DO NOTHING`, id, account, kind, query, provider, string(raw), time.Now().UnixMilli())
	if err != nil {
		apphost.Problem(w, 500, "Could not queue the request.")
		return
	}
	if count, _ := result.RowsAffected(); count != 1 {
		apphost.Problem(w, 409, "Wait for your current research request to finish.")
		return
	}
	select {
	case s.wake <- struct{}{}:
	default:
	}
	apphost.JSON(w, 202, map[string]string{"id": id})
}

func randomBytes() []byte {
	id := make([]byte, 16)
	rand.Read(id)
	return id
}

func (s *service) bookmark(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var input struct {
		ID    string `json:"id"`
		Saved bool   `json:"saved"`
	}
	if !apphost.Decode(w, r, &input) {
		return
	}
	ctx, account := r.Context(), apphost.Current(r).User.ID
	switch input.ID {
	case "loop-engineering", "five-stages", "inclusive-interface":
	default:
		var exists int
		if err := s.DB.QueryRowContext(ctx, "SELECT 1 FROM fieldnotes_articles WHERE id=? AND account=?", input.ID, account).Scan(&exists); err != nil {
			apphost.Problem(w, 404, "This briefing is unavailable.")
			return
		}
	}
	var err error
	if input.Saved {
		_, err = s.DB.ExecContext(ctx, "INSERT OR IGNORE INTO fieldnotes_saved VALUES (?,?)", account, input.ID)
	} else {
		_, err = s.DB.ExecContext(ctx, "DELETE FROM fieldnotes_saved WHERE account=? AND id=?", account, input.ID)
	}
	if err != nil {
		apphost.Problem(w, 500, "Could not update your reading list.")
		return
	}
	apphost.JSON(w, 200, input)
}

func (s *service) Run(ctx context.Context) error {
	timer := time.NewTicker(3 * time.Second)
	defer timer.Stop()
	for {
		if err := s.work(ctx); err != nil {
			s.Log.Error("Fieldnotes could not advance a research request")
		}
		select {
		case <-ctx.Done():
			return nil
		case <-s.wake:
		case <-timer.C:
		}
	}
}

func (s *service) work(ctx context.Context) error {
	var next job
	var account, raw string
	err := s.DB.QueryRowContext(ctx, `SELECT id,account,kind,query,provider,evidence
		FROM research_jobs WHERE status='queued' ORDER BY created LIMIT 1`).
		Scan(&next.ID, &account, &next.Kind, &next.Query, &next.Provider, &raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	s.Freeze.RLock()
	_, err = s.DB.ExecContext(ctx, "UPDATE research_jobs SET status='running' WHERE id=?", next.ID)
	s.Freeze.RUnlock()
	if err != nil {
		return err
	}

	keys, err := s.credentials(account)
	var briefing article
	var found []source
	work, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()
	switch {
	case err != nil:
		err = errors.New("Could not open provider settings.")
	case next.Kind == "brief":
		if keys.DeepSeek == "" || json.Unmarshal([]byte(raw), &next.Evidence) != nil {
			err = errors.New("Connect DeepSeek and run a new search.")
		} else {
			briefing, err = s.generate(work, next.Query, keys.DeepSeek, next.Evidence)
		}
	case next.Provider == "deepseek" && keys.DeepSeek == "", next.Provider == "workspace" && (keys.Browser == "" || s.browser == ""):
		err = errors.New("Reconnect the search provider and run the request again.")
	default:
		found, err = s.search(work, next.Query, next.Provider, keys)
	}
	return s.finish(next, account, briefing, found, err)
}

func (s *service) finish(next job, account string, briefing article, found []source, workErr error) error {
	s.Freeze.RLock()
	defer s.Freeze.RUnlock()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	status, issue := "done", ""
	if workErr != nil {
		status, issue = "failed", workErr.Error()
	} else if next.Kind == "brief" {
		briefing.ID, briefing.Topic, briefing.Date, briefing.Saved = next.ID, next.Query, time.Now().UTC().Format("2006-01-02"), false
		content, _ := json.Marshal(briefing)
		if _, err = tx.ExecContext(ctx, "INSERT INTO fieldnotes_articles VALUES (?,?,?,?)", next.ID, account, string(content), time.Now().UnixMilli()); err != nil {
			return err
		}
	}
	raw, _ := json.Marshal(found)
	if _, err = tx.ExecContext(ctx, "UPDATE research_jobs SET status=?, evidence=?, error=? WHERE id=?",
		status, string(raw), issue, next.ID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM research_jobs WHERE account=? AND status IN ('done','failed')
		AND id NOT IN (SELECT id FROM research_jobs WHERE account=? ORDER BY created DESC,id DESC LIMIT 20)`, account, account); err != nil {
		return err
	}
	return tx.Commit()
}
