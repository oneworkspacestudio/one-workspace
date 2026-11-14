package pet

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const journalSchema = `
CREATE TABLE IF NOT EXISTS pet_journal (
 id INTEGER PRIMARY KEY, tick INTEGER NOT NULL, kind TEXT NOT NULL, text TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS pet_journal_tick ON pet_journal(tick,id);
CREATE TABLE IF NOT EXISTS pet_journal_import (id INTEGER PRIMARY KEY CHECK(id=1));
`

// Use additive tables so startup is retryable and preview.2 can read its
// original snapshot format from a restored Workspace backup.
func (p *Habitat) openJournal() error {
	if _, err := p.host.DB.Exec(journalSchema); err != nil {
		return err
	}
	if p.world == nil {
		return nil
	}
	var marker int
	err := p.host.DB.QueryRow("SELECT id FROM pet_journal_import WHERE id=1").Scan(&marker)
	if err == nil {
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}
	w := *p.world
	w.events = w.Journal
	return p.save(context.Background(), w)
}

func (p *Habitat) history(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	view := struct {
		Day         int     `json:"day"`
		FirstDay    int     `json:"firstDay"`
		LastDay     int     `json:"lastDay"`
		StartedLate bool    `json:"startedLate"`
		Events      []Event `json:"events"`
	}{Events: []Event{}}
	if p.world == nil {
		apphost.JSON(w, 200, view)
		return
	}
	view.LastDay = p.world.Day
	view.Day = view.LastDay
	if value := r.URL.Query().Get("day"); value != "" {
		day, err := strconv.Atoi(value)
		if err != nil || day < 1 || day > view.LastDay {
			apphost.Problem(w, 400, "Choose an available habitat day.")
			return
		}
		view.Day = day
	}
	var oldest int64
	q := p.host.DB
	if err := q.QueryRowContext(r.Context(), "SELECT coalesce(min(tick),0) FROM pet_journal").Scan(&oldest); err != nil {
		apphost.Problem(w, 500, "The activity history could not be read.")
		return
	}
	view.FirstDay = int((oldest*2+510)/1440) + 1
	view.StartedLate = oldest > 0
	start, end := (int64(view.Day-1)*1440-510)/2, (int64(view.Day)*1440-510)/2
	rows, err := q.QueryContext(r.Context(), "SELECT tick,kind,text FROM pet_journal WHERE tick>=? AND tick<? ORDER BY tick,id", start, end)
	if err != nil {
		apphost.Problem(w, 500, "The activity history could not be read.")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var e Event
		if err = rows.Scan(&e.Tick, &e.Kind, &e.Text); err != nil {
			break
		}
		view.Events = append(view.Events, e)
	}
	if err != nil || rows.Err() != nil {
		apphost.Problem(w, 500, "The activity history could not be read.")
		return
	}
	apphost.JSON(w, 200, view)
}

func appendJournal(ctx context.Context, tx *sql.Tx, events []Event) error {
	if _, err := tx.ExecContext(ctx, "INSERT OR IGNORE INTO pet_journal_import VALUES(1)"); err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO pet_journal(tick,kind,text) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, e := range events {
		if _, err = stmt.ExecContext(ctx, e.Tick, e.Kind, e.Text); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) catchUp(now time.Time, elapsed, count int64) {
	if elapsed > count {
		w.note("clock", "Returned after a longer outage and caught up eight hours")
	}
	w.Saved = now.Unix() - (now.Unix()-w.Saved)%int64(step/time.Second)
}
