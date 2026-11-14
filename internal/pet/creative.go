package pet

import (
	"math/rand/v2"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Stroke struct {
	Shape string `json:"shape"`
	Color string `json:"color"`
	// x, y, width, height, in percentages of the canvas.
	Box []int `json:"box"`
}

type Painting struct {
	Title   string   `json:"title"`
	Sky     string   `json:"sky"`
	Strokes []Stroke `json:"strokes"`
}

type Project struct {
	Site     string `json:"site"`
	Design   string `json:"design"`
	Material string `json:"material"`
	Stock    int    `json:"stock"`
	Stage    int    `json:"stage"`
}

type materialSource struct {
	feature
	Material string `json:"material"`
}

var resources = []materialSource{
	{feature{"pond stones", Point{-9, -7}}, "stone"},
	{feature{"pine branches", Point{12, -10}}, "wood"},
	{feature{"meadow reeds", Point{-5, 10}}, "reed"},
}
var sites = []feature{
	{"south meadow", Point{11, 10}}, {"pine clearing", Point{16, -5}},
	{"north trail", Point{2, -13}}, {"garden edge", Point{5, 11}},
	{"western clearing", Point{-12, 6}}, {"east clearing", Point{16, 1}},
	{"clover clearing", Point{-2, 3}}, {"foothill", Point{-6, -6}},
}

var designs = []string{"cairn", "arch", "shelter", "lookout"}
var colors = []string{"peach", "gold", "rose", "mint", "fern", "blue", "violet", "cream", "navy", "clay"}

func titleOK(title string) bool {
	return len(title) > 0 && utf8.RuneCountInString(title) <= 52 && !strings.ContainsFunc(title, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsMark(r) && !strings.ContainsRune(" '-’", r)
	})
}

func validPainting(p *Painting) bool {
	if p == nil || !titleOK(p.Title) || !contains(colors, p.Sky) || len(p.Strokes) < 3 || len(p.Strokes) > 12 {
		return false
	}
	for _, stroke := range p.Strokes {
		if !contains(colors, stroke.Color) || !contains([]string{"oval", "leaf", "line", "triangle"}, stroke.Shape) || len(stroke.Box) != 4 {
			return false
		}
		for i, n := range stroke.Box {
			if n > 100 || n < 0 || i >= 2 && n < 2 {
				return false
			}
		}
	}
	return true
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func featureByID(features []feature, id string) (feature, bool) {
	for _, f := range features {
		if f.ID == id {
			return f, true
		}
	}
	return feature{}, false
}

func (w World) novelPainting(p *Painting) bool {
	if !validPainting(p) {
		return false
	}
	for _, c := range w.Creations {
		if other := c.Painting; other != nil && p.Sky == other.Sky && reflect.DeepEqual(p.Strokes, other.Strokes) {
			return false
		}
	}
	return true
}

func (w World) freeSites() []feature {
	out := []feature{}
	for _, site := range sites {
		clear := clearSegment(site.Point, site.Point, w.barriers())
		for _, c := range w.Creations {
			if c.Kind == "build" && distance(c.Point, site.Point) < 3 {
				clear = false
			}
		}
		if clear {
			out = append(out, site)
		}
	}
	return out
}

func (w World) buildOptions() []feature {
	if free := w.freeSites(); len(free) > 0 {
		return free
	}
	out := []feature{}
	for _, site := range sites {
		for _, c := range w.Creations {
			if c.Kind == "build" && c.Point == site.Point {
				out = append(out, site)
				break
			}
		}
	}
	return out
}

func (w World) learned() []string {
	return designs[:min(len(designs), 1+w.BuildXP/3)]
}

func (w World) art() *Painting {
	skies := []string{"peach", "blue", "mint", "rose"}
	if w.Weather == "Rain" {
		skies = []string{"blue", "violet"}
	}
	names := []string{"Light at", "A morning near", "Colours of", "Clouds over", "Quiet at"}
	places := []string{"the pond", "Clover Hollow", "the pine edge", "the meadow", "the clover trail"}
	near := landmarks[0]
	for _, place := range landmarks[1:] {
		if distance(w.Position, place.Point) < distance(w.Position, near.Point) {
			near = place
		}
	}
	places = append(places, near.ID, near.ID)
	p := &Painting{Title: names[rand.IntN(len(names))] + " " + places[rand.IntN(len(places))], Sky: skies[rand.IntN(len(skies))]}
	p.Strokes = []Stroke{
		{"oval", "cream", []int{15 + rand.IntN(70), 12 + rand.IntN(20), 9 + rand.IntN(9), 9 + rand.IntN(10)}},
		{"oval", "mint", []int{25 + rand.IntN(35), 76 + rand.IntN(15), 75, 45}},
		{"oval", "fern", []int{55 + rand.IntN(40), 84 + rand.IntN(16), 82, 44}},
	}
	count := 4 + rand.IntN(6)
	for range count {
		color := colors[rand.IntN(len(colors))]
		p.Strokes = append(p.Strokes, Stroke{[]string{"leaf", "oval", "line", "triangle"}[rand.IntN(4)], color,
			[]int{rand.IntN(100), 50 + rand.IntN(46), 3 + rand.IntN(12), 3 + rand.IntN(12)}})
	}
	return p
}

func (w *World) finishPainting() {
	p := w.Painting
	if p == nil {
		p = w.art()
	}
	for i, c := range w.Creations {
		if c.Kind == "art" {
			count := 0
			for _, other := range w.Creations[i:] {
				if other.Kind == "art" {
					count++
				}
			}
			if count >= 24 {
				w.Creations = append(w.Creations[:i], w.Creations[i+1:]...)
			}
			break
		}
	}
	w.Creations = append(w.Creations, Item{ID: w.Tick, Point: easel, Kind: "art", Style: w.FocusStyle, Painting: p})
	w.Happiness += 12
	w.note("art", "Painted “"+p.Title+"”")
}

func (w World) approach(site feature) (Point, bool) {
	blocks := w.barriers()
	for _, offset := range []Point{{0, 1.55}, {1.55, 0}, {0, -1.55}, {-1.55, 0}} {
		p := Point{site.X + offset.X, site.Z + offset.Z}
		if len(route(w.Position, p, blocks)) > 0 {
			return p, true
		}
	}
	return Point{}, false
}

func (w *World) planBuild(d Decision) bool {
	if w.Project == nil {
		site, ok := featureByID(w.buildOptions(), d.Place)
		if !ok {
			return false
		}
		w.Project = &Project{Site: site.ID, Design: d.Design, Material: d.Material}
		w.note("build", "Planned a "+d.Material+" "+d.Design+" at the "+site.ID)
	}
	project := *w.Project
	if project.Stock == 0 {
		for _, source := range resources {
			if source.Material == project.Material {
				w.Action, w.Target, w.Place = "gather", source.Point, source.ID
				return true
			}
		}
	}
	f, ok := featureByID(sites, project.Site)
	if !ok {
		return false
	}
	w.Action, w.Place = "build", f.ID
	// Approach the edge so existing foundations stay clear of the footpath.
	target, ok := w.approach(f)
	w.Target = target
	return ok
}

func (w *World) gather() {
	if w.Project == nil {
		return
	}
	p := *w.Project
	p.Stock = 1
	if w.BuildXP >= 3 {
		p.Stock = 2
	}
	w.Project = &p
	w.note("gather", "Collected "+p.Material+" for the "+p.Design)
}

func (w *World) build() {
	if w.Project == nil {
		return
	}
	p := *w.Project
	site, ok := featureByID(sites, p.Site)
	if !ok || p.Stock < 1 {
		return
	}
	p.Stock--
	p.Stage++
	found := false
	for i, c := range w.Creations {
		if c.Kind == "build" && c.Point == site.Point {
			w.Creations[i].Level, w.Creations[i].Design, w.Creations[i].Material = p.Stage, p.Design, p.Material
			found = true
			break
		}
	}
	if !found {
		w.Creations = append(w.Creations, Item{ID: w.Tick, Point: site.Point, Kind: "build", Level: p.Stage, Design: p.Design, Material: p.Material})
	}
	if p.Stage < 3 {
		w.Project = &p
		w.note("build", "Added a stage to the "+p.Design+" at the "+p.Site)
		return
	}
	w.Project = nil
	w.Happiness += 12
	w.note("build", "Finished the "+p.Material+" "+p.Design+" at the "+p.Site)
	old := len(w.learned())
	w.BuildXP = min(w.BuildXP+3, 9999)
	if learned := w.learned(); len(learned) > old {
		w.note("learn", "Learned how to build a "+learned[len(learned)-1])
	}
}
