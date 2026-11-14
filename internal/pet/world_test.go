package pet

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func TestHabitatEvolvesWithoutVisitors(t *testing.T) {
	w := newWorld("Pip", time.Unix(1000, 0))
	w.Seed = 72
	actions := map[string]bool{}
	weather := map[string]bool{}
	for range 24 * time.Hour / step {
		w.advance(nil)
		actions[w.Action], weather[w.Weather] = true, true
		for _, value := range []float64{w.Energy, w.Fullness, w.Happiness, w.Curiosity} {
			if value < 0 || value > 100 {
				t.Fatalf("need out of range: %f", value)
			}
		}
		if distance(w.Position, pond) < 4.6 {
			t.Fatalf("pet walked into the pond: %+v", w.Position)
		}
	}
	for _, action := range []string{"explore", "sleep", "eat", "draw", "build", "play"} {
		if !actions[action] {
			t.Errorf("pet never performed %s", action)
		}
	}
	if !weather["Clear"] || !weather["Rain"] || w.Day < 50 {
		t.Fatalf("weather and day did not cycle: %+v", weather)
	}
	art, builds := 0, 0
	for _, c := range w.Creations {
		switch c.Kind {
		case "art":
			art++
		case "build":
			builds++
		}
	}
	if art == 0 || art > 24 || builds == 0 || builds > len(sites)+6 || len(w.Journal) > 80 || len(w.Food) > 9 {
		t.Fatalf("unexpected persisted collection sizes: art=%d builds=%d events=%d food=%d", art, builds, len(w.Journal), len(w.Food))
	}
}

func TestUrgentNeedsOverrideModel(t *testing.T) {
	w := newWorld("Pip", time.Unix(1000, 0))
	w.Minute, w.Weather = 510, "Clear"
	w.choose(Decision{Action: "draw", Style: 2, Source: "Local model"})
	if w.Action != "draw" || w.Source != "Local model" || w.FocusStyle != 2 {
		t.Fatal("valid creative decision was ignored")
	}
	w.Energy = 11
	w.choose(Decision{Action: "play", Style: 1, Source: "AI Service"})
	if w.Action != "sleep" || w.Target != nest || w.Source != "Instincts" {
		t.Fatal("pet did not prioritize rest")
	}
	w.Energy, w.Fullness = 90, 18
	w.choose(Decision{Action: "build", Style: 1, Source: "AI Service"})
	if w.Action != "eat" {
		t.Fatal("pet did not prioritize food")
	}
}

func TestSavedHabitatSurvivesRestartAndCatchesUp(t *testing.T) {
	cfg := apphost.Config{ID: "pet", Name: "Pet", Version: Version, Data: t.TempDir(), Listen: "127.0.0.1:9046", Origin: "http://localhost:9046"}
	host, err := apphost.New(cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	p, err := New(host)
	if err != nil {
		host.Close()
		t.Fatal(err)
	}
	ctx := context.Background()
	start := time.Unix(100000, 0)
	w := newWorld("Pip", start)
	w.treat()
	if err = p.save(ctx, w); err != nil {
		host.Close()
		t.Fatal(err)
	}
	if err = host.Close(); err != nil {
		t.Fatal(err)
	}
	host, err = apphost.New(cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer host.Close()
	p, err = New(host)
	if err != nil {
		t.Fatal(err)
	}
	w.events = nil
	if !reflect.DeepEqual(*p.world, w) {
		t.Fatal("snapshot changed after reload")
	}
	now := start.Add(24 * time.Hour)
	if err = p.advance(ctx, now); err != nil {
		t.Fatal(err)
	}
	if p.world.Tick != int64(8*time.Hour/step) || p.world.Saved != now.Unix() || p.world.TreatAfter != w.TreatAfter {
		t.Fatal("restart catch-up did not apply the bounded replay")
	}
}

func TestModelProtocolsAndActionValidation(t *testing.T) {
	shared := false
	content := `{"action":"draw","style":2}`
	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		payload, _ := json.Marshal(request)
		if strings.Contains(string(payload), "private-pet-name") || request["stream"] != false {
			t.Error("invalid simulation-only model request")
		}
		if shared {
			if r.URL.Path != "/_internal/chat" || r.Header.Get("X-AI-App") != "pet" || r.Header.Get("X-AI-User") != "test-owner" ||
				r.Header.Get("Authorization") != "Bearer "+strings.Repeat("k", 32) || request["model"] != nil {
				t.Error("incorrect AI Service request")
			}
		} else if r.URL.Path != "/api/chat" || r.Header.Get("Authorization") != "" || request["model"] != "pet-model" {
			t.Error("incorrect local Ollama request")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"message": map[string]string{"content": content}})
	}))
	defer endpoint.Close()
	b := newBrain("0")
	b.localURL, b.sharedURL = endpoint.URL+"/api/chat", endpoint.URL+"/_internal/chat"
	cfg := Settings{Model: "pet-model", Owner: "test-owner", ServiceKey: strings.Repeat("k", 32)}
	world := newWorld("private-pet-name", time.Unix(1000, 0))
	world.Weather = "Clear"
	for _, source := range []string{"Local model", "AI Service"} {
		idea, err := b.decide(context.Background(), cfg, world, shared)
		if err != nil || idea.Action != "draw" || idea.Style != 2 || idea.Source != source {
			t.Fatalf("bad decision from %s: %+v %v", source, idea, err)
		}
		shared = true
	}
	for _, bad := range []string{`{"action":"shell","style":0}`, `{"action":"draw","style":8}`, `{"action":"draw","style":0,"extra":true}`} {
		content = bad
		if _, err := b.decide(context.Background(), cfg, world, true); err == nil {
			t.Errorf("accepted invalid decision %s", bad)
		}
	}
}
