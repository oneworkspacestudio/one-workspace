package sunnymoon

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const Version = "0.1.0-preview.4"

//go:embed web
var web embed.FS

func UI() fs.FS {
	files, err := fs.Sub(web, "web")
	if err != nil {
		panic(err)
	}
	return files
}

const schema = `
CREATE TABLE routines (
  account TEXT PRIMARY KEY,
  bedtime TEXT NOT NULL,
  wake TEXT NOT NULL,
  lead INTEGER NOT NULL,
  last_visit INTEGER NOT NULL,
  night_visit INTEGER NOT NULL,
  expires INTEGER NOT NULL,
  revision INTEGER NOT NULL
);`

type settings struct {
	Bedtime string `json:"bedtime"`
	Wake    string `json:"wake"`
	Lead    int    `json:"lead"`
}

type routine struct {
	Settings   settings `json:"settings"`
	LastVisit  int64    `json:"lastVisit"`
	NightVisit int64    `json:"nightVisit"`
	Expires    int64    `json:"expires"`
	Revision   int64    `json:"revision"`
}

type response struct {
	Routine routine `json:"routine"`
	CSRF    string  `json:"csrf"`
}

func Register(host *apphost.Host) error {
	if err := host.Migrate(1, schema); err != nil {
		return err
	}
	host.Mux.HandleFunc("GET /api/routine", host.Require(func(w http.ResponseWriter, r *http.Request) {
		session := apphost.Current(r)
		value, err := load(r.Context(), host.DB, session.User.ID, time.Now().UnixMilli())
		if err != nil {
			apphost.Problem(w, 500, "Couldn’t open your routine. Try again.")
			return
		}
		apphost.JSON(w, 200, response{value, session.CSRF})
	}))
	host.Mux.HandleFunc("PUT /api/routine", host.Require(func(w http.ResponseWriter, r *http.Request) {
		session := apphost.Current(r)
		var value routine
		r.Body = http.MaxBytesReader(w, r.Body, 2048)
		if !apphost.Decode(w, r, &value) {
			return
		}
		now := time.Now().UnixMilli()
		if !value.valid(now) {
			apphost.Problem(w, 400, "Check the routine times and try again.")
			return
		}
		updated, err := store(r.Context(), host.DB, session.User.ID, value)
		if err != nil {
			apphost.Problem(w, 500, "Couldn’t save your routine. Try again.")
			return
		}
		if !updated {
			apphost.Problem(w, 409, "The routine changed. Try again.")
			return
		}
		value.Revision++
		apphost.JSON(w, 200, response{value, session.CSRF})
	}))
	return nil
}

func (s settings) valid() bool {
	bed, bedErr := time.Parse("15:04", s.Bedtime)
	wake, wakeErr := time.Parse("15:04", s.Wake)
	if bedErr != nil || wakeErr != nil || bed.Format("15:04") != s.Bedtime || wake.Format("15:04") != s.Wake {
		return false
	}
	switch s.Lead {
	case 15, 30, 45, 60, 90:
	default:
		return false
	}
	const day = 24 * 60
	daytime := (int(bed.Sub(wake)/time.Minute) + day) % day
	return daytime > s.Lead && daytime <= day-60
}

func (r routine) valid(now int64) bool {
	if !r.Settings.valid() || r.Revision < 0 || r.LastVisit < 0 || r.NightVisit < 0 || r.NightVisit > r.LastVisit {
		return false
	}
	if r.LastVisit == 0 {
		return r.Expires == 0
	}
	// Allow small clock differences and changes to the next wake time, including DST.
	return r.LastVisit <= now+time.Hour.Milliseconds() &&
		r.Expires > now-time.Hour.Milliseconds() && r.Expires > r.LastVisit &&
		r.Expires-r.LastVisit <= (27*time.Hour).Milliseconds()
}

func load(ctx context.Context, db *sql.DB, account string, now int64) (routine, error) {
	// Clearing expired check-ins preserves settings and keeps only the current day's values.
	if _, err := db.ExecContext(ctx, `UPDATE routines SET last_visit=0, night_visit=0, expires=0, revision=revision+1
		WHERE account=? AND expires>0 AND expires<=?`, account, now); err != nil {
		return routine{}, err
	}
	var value routine
	err := db.QueryRowContext(ctx, `SELECT bedtime, wake, lead, last_visit, night_visit, expires, revision
		FROM routines WHERE account=?`, account).Scan(
		&value.Settings.Bedtime, &value.Settings.Wake, &value.Settings.Lead,
		&value.LastVisit, &value.NightVisit, &value.Expires, &value.Revision)
	if errors.Is(err, sql.ErrNoRows) {
		return routine{Settings: settings{Bedtime: "21:00", Wake: "07:00", Lead: 60}}, nil
	}
	return value, err
}

func store(ctx context.Context, db *sql.DB, account string, value routine) (bool, error) {
	var result sql.Result
	var err error
	if value.Revision == 0 {
		result, err = db.ExecContext(ctx, `INSERT INTO routines VALUES (?, ?, ?, ?, ?, ?, ?, 1)
			ON CONFLICT(account) DO NOTHING`, account, value.Settings.Bedtime, value.Settings.Wake,
			value.Settings.Lead, value.LastVisit, value.NightVisit, value.Expires)
	} else {
		result, err = db.ExecContext(ctx, `UPDATE routines
			SET bedtime=?, wake=?, lead=?, last_visit=?, night_visit=?, expires=?, revision=revision+1
			WHERE account=? AND revision=?`, value.Settings.Bedtime, value.Settings.Wake, value.Settings.Lead,
			value.LastVisit, value.NightVisit, value.Expires, account, value.Revision)
	}
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count == 1, err
}
