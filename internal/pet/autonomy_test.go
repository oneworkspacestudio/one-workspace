package pet

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func TestWorldRouteAvoidsPondAndExistingStructures(t *testing.T) {
	a, b := Point{-19, -2}, Point{-10, -13}
	blocks := []barrier{{pond, 4.7}, {Point{-13, -4}, 1.2}}
	path := route(a, b, blocks)
	if len(path) < 2 || path[len(path)-1] != b {
		t.Fatalf("expected a path around the water, got %v", path)
	}
	for _, p := range path {
		if !clearSegment(a, p, blocks) {
			t.Fatalf("route intersects an obstacle between %v and %v", a, p)
		}
		a = p
	}
}

func TestConstructionRequiresStockAndTeachesDesigns(t *testing.T) {
	w := newWorld("Pip", time.Unix(1000, 0))
	w.Weather = "Clear"
	plan := Decision{Action: "build", Place: "south meadow", Design: "cairn", Material: "wood", Source: "AI Service"}
	unlearned := plan
	unlearned.Design = "lookout"
	if w.validIdea(unlearned) {
		t.Fatal("accepted an unlearned design")
	}
	for stage := 1; stage <= 3; stage++ {
		w.beginIdea(plan)
		if w.Action != "gather" || w.Place != "pine branches" {
			t.Fatalf("material not collected: %s %s", w.Action, w.Place)
		}
		before := len(w.Creations)
		w.build()
		if len(w.Creations) != before {
			t.Fatal("constructed without stock")
		}
		w.gather()
		w.beginIdea(Decision{Action: "build", Source: "AI Service"})
		if w.Action != "build" {
			t.Fatal("gather did not lead to construction")
		}
		w.Tick++
		w.build()
		if got := w.Creations[0].Level; got != stage {
			t.Fatalf("expected stage %d, got %d", stage, got)
		}
	}
	if w.Project != nil || !contains(w.learned(), "arch") || w.BuildXP != 3 {
		t.Fatal("completed project did not teach the next design")
	}
	w.Minute = 1300
	if w.validIdea(plan) {
		t.Fatal("construction was allowed at night")
	}
}

func TestPaintingRejectsRepeatedCompositionAndKeepsPalette(t *testing.T) {
	w := newWorld("Pip", time.Unix(1000, 0))
	w.Weather = "Clear"
	p := &Painting{Title: "Pond shade", Sky: "mint", Strokes: []Stroke{
		{"oval", "gold", []int{60, 20, 10, 10}},
		{"leaf", "fern", []int{72, 82, 36, 24}},
		{"line", "blue", []int{32, 44, 10, 5}},
	}}
	idea := Decision{Action: "draw", Painting: p, Source: "Local model"}
	if !w.validIdea(idea) {
		t.Fatal("composition rejected")
	}
	w.beginIdea(idea)
	w.finishPainting()
	if w.Creations[0].Painting == nil || w.Creations[0].Painting.Title != "Pond shade" {
		t.Fatal("composition was not saved")
	}
	if w.validIdea(idea) {
		t.Fatal("model can repeat an identical composition")
	}
}

func TestJournalImportsLegacyEntriesAndRetainsFullDays(t *testing.T) {
	cfg := apphost.Config{ID: "pet", Name: "Pet", Version: Version, Data: t.TempDir(), Listen: "127.0.0.1:9046", Origin: "http://localhost:9046"}
	h, err := apphost.New(cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if h != nil {
			h.Close()
		}
	}()
	if err := h.Migrate(1, "CREATE TABLE pet_world(id INTEGER PRIMARY KEY CHECK(id=1), snapshot BLOB NOT NULL)"); err != nil {
		t.Fatal(err)
	}
	w := newWorld("Pip", time.Unix(1000, 0))
	w.Tick = 465
	w.updateEnvironment()
	w.note("sleep", "Day two")
	data, _ := json.Marshal(w)
	if _, err = h.DB.Exec("INSERT INTO pet_world VALUES(1,?)", data); err != nil {
		t.Fatal(err)
	}
	p, err := New(h)
	if err != nil {
		t.Fatal(err)
	}
	w = *p.world
	for i := int64(466); i < 566; i++ {
		w.Tick = i
		w.note("explore", "Explored the pond bank")
	}
	if err := p.save(context.Background(), w); err != nil {
		t.Fatal(err)
	}
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
	h, err = apphost.New(cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	p, err = New(h)
	if err != nil {
		t.Fatal(err)
	}
	for _, check := range []struct{ day, count int }{{1, 1}, {2, 101}} {
		req := httptest.NewRequest("GET", "/api/history?day=1", nil)
		if check.day == 2 {
			req = httptest.NewRequest("GET", "/api/history?day=2", nil)
		}
		res := httptest.NewRecorder()
		p.history(res, req)
		var view struct {
			Day    int
			Events []Event
		}
		if err := json.Unmarshal(res.Body.Bytes(), &view); err != nil || view.Day != check.day || len(view.Events) != check.count {
			t.Fatalf("day %d: status %d, count %d, error %v", check.day, res.Code, len(view.Events), err)
		}
	}
}
