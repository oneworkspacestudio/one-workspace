package pet

import (
	"math"
	"math/rand/v2"
	"time"
)

const step = 2 * time.Second

type Point struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
}

type Item struct {
	ID int64 `json:"id"`
	Point
	Kind  string `json:"kind"`
	Style int    `json:"style"`
	Level int    `json:"level,omitempty"`
}

type Event struct {
	Tick int64  `json:"tick"`
	Kind string `json:"kind"`
	Text string `json:"text"`
}

type World struct {
	Name       string  `json:"name"`
	Seed       uint64  `json:"seed,string"`
	Tick       int64   `json:"tick"`
	Saved      int64   `json:"saved"`
	Position   Point   `json:"position"`
	Target     Point   `json:"target"`
	Action     string  `json:"action"`
	Until      int64   `json:"until"`
	Focus      int64   `json:"focus"`
	FocusStyle int     `json:"focusStyle"`
	Mood       string  `json:"mood"`
	Energy     float64 `json:"energy"`
	Fullness   float64 `json:"fullness"`
	Happiness  float64 `json:"happiness"`
	Curiosity  float64 `json:"curiosity"`
	Minute     int     `json:"minute"`
	Day        int     `json:"day"`
	Weather    string  `json:"weather"`
	Source     string  `json:"source"`
	TreatAfter int64   `json:"treatAfter"`
	Food       []Item  `json:"food"`
	Creations  []Item  `json:"creations"`
	Journal    []Event `json:"journal"`
}

var (
	nest       = Point{9, -8}
	easel      = Point{7, 4}
	playground = Point{-6, 5}
	pond       = Point{-15, -9}
	buildSites = []Point{{13, 7}, {15, 2}, {13, -3}, {5, -14}, {-3, -14}, {-8, -9}}
)

func newWorld(name string, now time.Time) World {
	w := World{
		Name: name, Seed: rand.Uint64(), Saved: now.Unix(), Position: Point{0, 0},
		Energy: 86, Fullness: 77, Happiness: 74, Curiosity: 83,
		Food:      []Item{{ID: -1, Point: Point{-4, -2}, Kind: "berry"}, {ID: -2, Point: Point{3, -6}, Kind: "fruit"}},
		Creations: []Item{}, Journal: []Event{},
	}
	w.updateEnvironment()
	w.start("explore", 0, "Instincts")
	w.updateMood()
	w.note("home", "Moved into Clover Hollow")
	return w
}

func (w *World) rng() *rand.Rand {
	return rand.New(rand.NewPCG(w.Seed, uint64(w.Tick)))
}

func (w *World) note(kind, text string) {
	w.Journal = append(w.Journal, Event{w.Tick, kind, text})
	if len(w.Journal) > 80 {
		w.Journal = append([]Event(nil), w.Journal[len(w.Journal)-80:]...)
	}
}

func distance(a, b Point) float64 { return math.Hypot(a.X-b.X, a.Z-b.Z) }

func (w *World) meadowPoint(r *rand.Rand) Point {
	for {
		p := Point{float64(r.IntN(33) - 16), float64(r.IntN(29) - 14)}
		if distance(p, pond) > 5.5 && distance(p, nest) > 3 && distance(p, easel) > 1.5 {
			return p
		}
	}
}

func (w *World) nearestFood() (Item, bool) {
	var best Item
	dist := math.Inf(1)
	for _, item := range w.Food {
		if d := distance(item.Point, w.Position); d < dist {
			best, dist = item, d
		}
	}
	return best, dist < math.Inf(1)
}

func (w *World) start(action string, style int, source string) {
	w.Action, w.Source, w.Focus, w.FocusStyle = action, source, 0, style
	w.Until = 0
	switch action {
	case "sleep":
		w.Target = nest
	case "eat":
		food, _ := w.nearestFood()
		w.Focus, w.Target = food.ID, food.Point
	case "play":
		w.Target = playground
	case "draw":
		w.Target = easel
	case "build":
		w.Target = buildSites[w.rng().IntN(len(buildSites))]
	default:
		w.Target = w.meadowPoint(w.rng())
	}
}

func (w *World) updateEnvironment() {
	minute := w.Tick*2 + 510
	w.Day, w.Minute = int(minute/1440)+1, int(minute%1440)
	w.Weather = "Clear"
	if (w.Seed+uint64(minute/180)*17)%7 == 1 {
		w.Weather = "Rain"
	}
}

func (w World) dark() bool { return w.Minute >= 1200 || w.Minute < 360 }

func (w *World) updateMood() {
	switch {
	case w.Action == "sleep" && w.Until != 0:
		w.Mood = "Dreaming"
	case w.Energy < 34:
		w.Mood = "Sleepy"
	case w.Fullness < 33:
		w.Mood = "Peckish"
	case w.Happiness < 40:
		w.Mood = "Blue"
	case w.Action == "draw" || w.Action == "build":
		w.Mood = "Inspired"
	case w.Happiness >= 82:
		w.Mood = "Delighted"
	case w.Curiosity >= 69:
		w.Mood = "Curious"
	default:
		w.Mood = "Content"
	}
}

func (w World) urgent() string {
	if w.Energy < 29 || w.Energy < 64 && w.dark() {
		return "sleep"
	}
	if _, ok := w.nearestFood(); ok && w.Fullness < 53 {
		return "eat"
	}
	return ""
}

func (w World) allows(action string) bool {
	switch action {
	case "draw", "build":
		return !w.dark() && w.Weather != "Rain" && w.Energy >= 45 && w.Fullness >= 50
	case "play":
		return w.Energy >= 38 && w.Fullness >= 45
	case "explore":
		return true
	default:
		return false
	}
}

type Decision struct {
	Action string `json:"action"`
	Style  int    `json:"style"`
	Source string `json:"-"`
}

func (w *World) choose(idea Decision) {
	if need := w.urgent(); need != "" {
		w.start(need, 0, "Instincts")
		return
	}
	if w.allows(idea.Action) && idea.Style >= 0 && idea.Style < 4 {
		w.start(idea.Action, idea.Style, idea.Source)
		return
	}
	r := w.rng()
	actions := []string{"explore", "play", "draw", "build"}
	weights := []int{5 + int(w.Curiosity/18), 6 + int((100-w.Happiness)/12), 5, 4}
	total := 0
	for i, action := range actions {
		if !w.allows(action) || action == w.Action {
			weights[i] = 0
		}
		total += weights[i]
	}
	if total > 0 {
		pick := r.IntN(total)
		for i, weight := range weights {
			pick -= weight
			if pick < 0 {
				w.start(actions[i], r.IntN(4), "Instincts")
				return
			}
		}
	}
	w.start("explore", 0, "Instincts")
}

func (w *World) move() bool {
	dest := w.Target
	p := w.Position
	d := distance(p, dest)
	if d < 1.5 {
		w.Position = dest
		return true
	}
	pace := 1.45
	if w.Action == "play" {
		pace = 1.85
	}
	next := Point{p.X + (dest.X-p.X)*pace/d, p.Z + (dest.Z-p.Z)*pace/d}
	if distance(next, pond) < 4.65 {
		// Move around the shore instead of crossing the pond.
		px, pz := p.X-pond.X, p.Z-pond.Z
		sign := 1.0
		if px*(dest.Z-pond.Z)-pz*(dest.X-pond.X) < 0 {
			sign = -1
		}
		radius := math.Atan2(pz, px) + sign*0.29
		next = Point{pond.X + 4.85*math.Cos(radius), pond.Z + 4.85*math.Sin(radius)}
	}
	w.Position = next
	return false
}

func (w *World) complete() {
	switch w.Action {
	case "sleep":
		w.note("sleep", "Woke up in the burrow")
	case "eat":
		for i, food := range w.Food {
			if food.ID == w.Focus {
				w.Food = append(w.Food[:i], w.Food[i+1:]...)
				w.Fullness += 33
				w.Happiness += 5
				w.note("eat", "Found a snack in the meadow")
				break
			}
		}
	case "play":
		w.Happiness += 15
		w.note("play", "Played with the meadow ball")
	case "explore":
		w.Happiness += 4
		w.Curiosity -= 17
		w.note("explore", "Explored a new corner of the hollow")
	case "draw":
		names := []string{"Sun over Clover Hill", "Moonlit Pond", "Wildflower Study", "After the Rain"}
		count := 0
		for _, item := range w.Creations {
			if item.Kind == "art" {
				count++
			}
		}
		w.Creations = append(w.Creations, Item{ID: w.Tick, Point: easel, Kind: "art", Style: w.FocusStyle, Level: count + 1})
		if count >= 23 {
			for i, item := range w.Creations {
				if item.Kind == "art" {
					w.Creations = append(w.Creations[:i], w.Creations[i+1:]...)
					break
				}
			}
		}
		w.Happiness += 12
		w.note("art", "Painted “"+names[w.FocusStyle]+"”")
	case "build":
		for i := range w.Creations {
			item := &w.Creations[i]
			if item.Kind == "build" && item.Point == w.Target {
				item.Level = min(3, item.Level+1)
				item.Style = w.FocusStyle
				w.Happiness += 8
				w.note("build", "Reworked a meadow sculpture")
				return
			}
		}
		w.Creations = append(w.Creations, Item{ID: w.Tick, Point: w.Target, Kind: "build", Style: w.FocusStyle, Level: 1})
		w.Happiness += 10
		w.note("build", "Built a new meadow sculpture")
	}
}

func (w *World) advance(idea *Decision) {
	w.Tick++
	w.updateEnvironment()
	if w.Tick%27 == 0 && len(w.Food) < 9 {
		r := w.rng()
		kind := []string{"berry", "fruit", "mushroom"}[r.IntN(3)]
		w.Food = append(w.Food, Item{ID: -w.Tick - 2, Point: w.meadowPoint(r), Kind: kind})
	}
	w.Fullness -= 0.36
	w.Curiosity += 0.43
	if w.Weather == "Rain" && w.Action != "sleep" {
		w.Happiness -= 0.34
	}
	if w.Fullness < 27 || w.Energy < 24 {
		w.Happiness -= 0.42
	}
	if w.Action == "sleep" && w.Until > 0 {
		w.Energy += 2.8
		w.Happiness += 0.12
	} else {
		w.Energy -= 0.47
		if w.Action == "play" || w.Action == "build" {
			w.Energy -= 0.26
		}
	}
	// Instincts interrupt recreational plans before care becomes critical.
	if w.Action != "sleep" && w.Action != "eat" {
		if need := w.urgent(); need != "" {
			w.start(need, 0, "Instincts")
		}
	}
	if w.Until == 0 && w.move() {
		duration := map[string]int64{"explore": 5, "play": 12, "draw": 16, "build": 19, "eat": 4, "sleep": 28}
		w.Until = w.Tick + duration[w.Action]
		if w.Action == "sleep" {
			w.note("sleep", "Curled up in the burrow")
		}
	}
	if w.Until > 0 && w.Tick >= w.Until {
		w.complete()
		next := Decision{}
		if idea != nil && w.urgent() == "" {
			next, *idea = *idea, Decision{}
		}
		w.choose(next)
	}
	w.normalize()
}

func (w *World) normalize() {
	w.Energy = min(100, max(0, w.Energy))
	w.Fullness = min(100, max(0, w.Fullness))
	w.Happiness = min(100, max(0, w.Happiness))
	w.Curiosity = min(100, max(0, w.Curiosity))
	w.updateMood()
}

func (w *World) treat() {
	w.Happiness += 13
	w.Fullness += 14
	w.TreatAfter = w.Saved + 45
	w.note("treat", "A visitor brought a clover treat")
	w.normalize()
}
