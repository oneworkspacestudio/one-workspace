package sunnymoon

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestRoutineIsPrivateAndRetainsOnlyCurrentCheckins(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err = db.Exec(schema); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	r, err := load(ctx, db, "child", now.UnixMilli())
	if err != nil || r.Settings != (settings{"21:00", "07:00", 60}) {
		t.Fatalf("default routine: %+v, %v", r, err)
	}
	r.Settings = settings{"20:30", "07:30", 45}
	r.LastVisit = now.UnixMilli()
	r.NightVisit = now.UnixMilli()
	r.Expires = now.Add(time.Hour).UnixMilli()
	if !r.valid(now.UnixMilli()) {
		t.Fatal("current check-in rejected")
	}
	if changed, err := store(ctx, db, "child", r); !changed || err != nil {
		t.Fatalf("store: %v, %v", changed, err)
	}
	if changed, err := store(ctx, db, "child", r); changed || err != nil {
		t.Fatalf("stale writer: %v, %v", changed, err)
	}
	other, err := load(ctx, db, "other", now.UnixMilli())
	if err != nil || other.LastVisit != 0 || other.Settings == r.Settings {
		t.Fatalf("isolated account: %+v, %v", other, err)
	}
	saved, err := load(ctx, db, "child", now.UnixMilli())
	if err != nil || saved.LastVisit != r.LastVisit || saved.NightVisit != r.NightVisit || saved.Revision != 1 {
		t.Fatalf("saved routine: %+v, %v", saved, err)
	}
	morning, err := load(ctx, db, "child", r.Expires)
	if err != nil || morning.LastVisit != 0 || morning.NightVisit != 0 || morning.Settings != r.Settings {
		t.Fatalf("morning reset: %+v, %v", morning, err)
	}
}
