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
	Kind     string    `json:"kind"`
	Style    int       `json:"style"`
	Level    int       `json:"level,omitempty"`
	Painting *Painting `json:"painting,omitempty"`
	Design   string    `json:"design,omitempty"`
	Material string    `json:"material,omitempty"`
}

type Event struct {
	Tick int64  `json:"tick"`
	Kind string `json:"kind"`
	Text string `json:"text"`
}

type World struct {
	Name       string    `json:"name"`
	Seed       uint64    `json:"seed,string"`
	Tick       int64     `json:"tick"`
	Saved      int64     `json:"saved"`
	Position   Point     `json:"position"`
	Target     Point     `json:"target"`
	Action     string    `json:"action"`
	Until      int64     `json:"until"`
	Focus      int64     `json:"focus"`
	FocusStyle int       `json:"focusStyle"`
	Mood       string    `json:"mood"`
	Energy     float64   `json:"energy"`
	Fullness   float64   `json:"fullness"`
	Happiness  float64   `json:"happiness"`
	Curiosity  float64   `json:"curiosity"`
	Minute     int       `json:"minute"`
	Day        int       `json:"day"`
	Weather    string    `json:"weather"`
	Source     string    `json:"source"`
	TreatAfter int64     `json:"treatAfter"`
	Food       []Item    `json:"food"`
	Creations  []Item    `json:"creations"`
	Journal    []Event   `json:"journal"`
	Place      string    `json:"place,omitempty"`
	Route      []Point   `json:"-"`
	Painting   *Painting `json:"painting,omitempty"`
	BuildXP    int       `json:"buildXP"`
	Project    *Project  `json:"project,omitempty"`
	VisitTrail []string  `json:"visitTrail,omitempty"`
	events     []Event
}

var (
	nest       = Point{9, -6.85}
	easel      = Point{7, 4}
	playground = Point{-6, 5}
	pond       = Point{-15, -9}
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

func (w *World) note(kind, text string) {
	event := Event{w.Tick, kind, text}
	w.Journal = append(w.Journal, event)
	w.events = append(w.events, event)
	if len(w.Journal) > 80 {
		w.Journal = append([]Event(nil), w.Journal[len(w.Journal)-80:]...)
	}
}

func distance(a, b Point) float64 { return math.Hypot(a.X-b.X, a.Z-b.Z) }

func (w *World) meadowPoint() Point {
	blocks := w.barriers()
	for {
		p := Point{float64(rand.IntN(33) - 16), float64(rand.IntN(29) - 14)}
		if clearSegment(p, p, blocks) && distance(p, pond) > 5.5 && distance(p, nest) > 3 && distance(p, easel) > 1.5 {
			return p
		}
	}
}

func (w *World) nearestFood() (Item, bool) {
	var best Item
	dist := math.Inf(1)
	blocks := w.barriers()
	for _, item := range w.Food {
		if !clearSegment(item.Point, item.Point, blocks) {
			continue
		}
		if d := distance(item.Point, w.Position); d < dist {
			best, dist = item, d
		}
	}
	return best, dist < math.Inf(1)
}

func (w *World) start(action string, style int, source string) {
	w.Action, w.Source, w.Focus, w.FocusStyle = action, source, 0, style
	w.Until, w.Route, w.Painting = 0, nil, nil
	switch action {
	case "sleep":
		w.Target, w.Place = nest, "the burrow"
	case "eat":
		food, _ := w.nearestFood()
		w.Focus, w.Target, w.Place = food.ID, food.Point, "a "+food.Kind
	case "play":
		w.Target, w.Place = playground, "the meadow ball"
	case "draw":
		w.Target, w.Place = Point{easel.X, 4.7}, "the easel"
	case "build":
		w.Target, w.Place = w.meadowPoint(), "the meadow"
	default:
		w.Target, w.Place = w.meadowPoint(), "the meadow"
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
	case "draw", "build", "gather":
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
	Action   string    `json:"action"`
	Style    int       `json:"style"`
	Place    string    `json:"place,omitempty"`
	Via      string    `json:"via,omitempty"`
	Design   string    `json:"design,omitempty"`
	Material string    `json:"material,omitempty"`
	Painting *Painting `json:"painting,omitempty"`
	Source   string    `json:"-"`
}

func (w *World) choose(idea Decision) {
	if need := w.urgent(); need != "" {
		w.start(need, 0, "Instincts")
		return
	}
	if w.validIdea(idea) {
		w.beginIdea(idea)
		return
	}
	actions := []string{"explore", "play", "draw", "build"}
	weights := []int{5 + int(w.Curiosity/18), 4 + int((100-w.Happiness)/12), 6, 9}
	if w.Project != nil {
		weights[3] = 18
	}
	total := 0
	for i, action := range actions {
		if !w.allows(action) || action == w.Action && action != "build" || action == "build" && w.Project == nil && len(w.buildOptions()) == 0 {
			weights[i] = 0
		}
		total += weights[i]
	}
	if total > 0 {
		pick := rand.IntN(total)
		for i, weight := range weights {
			pick -= weight
			if pick < 0 {
				w.beginIdea(w.fallback(actions[i]))
				return
			}
		}
	}
	w.start("explore", 0, "Instincts")
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
		w.note("explore", "Explored "+w.Place)
		w.VisitTrail = append(append([]string{}, w.VisitTrail...), w.Place)
		if len(w.VisitTrail) > 6 {
			w.VisitTrail = w.VisitTrail[1:]
		}
	case "draw":
		w.finishPainting()
	case "gather":
		w.gather()
	case "build":
		w.build()
	}
}

func (w *World) advance(idea *Decision) {
	w.Tick++
	w.updateEnvironment()
	if w.Tick%27 == 0 && len(w.Food) < 9 {
		kind := []string{"berry", "fruit", "mushroom"}[rand.IntN(3)]
		w.Food = append(w.Food, Item{ID: -w.Tick - 2, Point: w.meadowPoint(), Kind: kind})
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
			if w.Action == "build" || w.Action == "gather" {
				w.note("build", "Set the project aside to tend to basic needs")
			}
			w.start(need, 0, "Instincts")
		}
	}
	if (w.Action == "build" || w.Action == "gather") && (w.dark() || w.Weather == "Rain") {
		w.note("build", "Paused the project until daylight and clear skies")
		w.choose(Decision{})
	}
	if w.Until == 0 && w.move() {
		duration := map[string]int64{"explore": 5, "play": 12, "draw": 16, "build": 12, "gather": 5, "eat": 4, "sleep": 28}
		w.Until = w.Tick + duration[w.Action]
		if w.Action == "sleep" {
			w.note("sleep", "Curled up in the burrow")
		} else if w.Action == "build" {
			w.note("build", "Started work at the "+w.Place)
		} else if w.Action == "draw" {
			w.note("art", "Started a painting at the easel")
		}
	}
	if w.Until > 0 && w.Tick >= w.Until {
		w.complete()
		next := Decision{}
		if idea != nil && w.urgent() == "" {
			next = *idea
		}
		if idea != nil {
			*idea = Decision{}
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
