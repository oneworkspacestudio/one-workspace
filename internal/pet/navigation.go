package pet

import (
	"math"
	"math/rand/v2"
)

type feature struct {
	ID string `json:"id"`
	Point
}

type barrier struct {
	Point
	Radius float64 `json:"radius"`
}

// These match the fixed landmarks in the viewer. Trees at the outer edge are
// beyond the pet's walking area.
var landmarks = []feature{
	{"pond bank", Point{-10, -6}}, {"wildflower meadow", Point{-6, 9}},
	{"high clearing", Point{3, -12}}, {"pine edge", Point{13, -7}},
	{"clover trail", Point{2, 9}}, {"east meadow", Point{14, 10}},
}
var obstacles = []barrier{
	{pond, 4.7}, {Point{9, -9}, 2.05}, {Point{7, 3.55}, .79},
	{Point{-10, -1}, .7}, {Point{-10, -13}, .7}, {Point{13, -12}, .7},
	{Point{15, 6}, .7}, {Point{-4, -12}, .7}, {Point{5, -12}, .7},
	{Point{-9, 11}, .7}, {Point{-19, -1}, .7},
}

func clearSegment(a, b Point, blocks []barrier) bool {
	dx, dz := b.X-a.X, b.Z-a.Z
	d2 := dx*dx + dz*dz
	for _, o := range blocks {
		t := 0.0
		if d2 > 0 {
			t = max(0, min(1, ((o.X-a.X)*dx+(o.Z-a.Z)*dz)/d2))
		}
		if distance(Point{a.X + t*dx, a.Z + t*dz}, o.Point) < o.Radius {
			return false
		}
	}
	return true
}

func (w World) barriers() []barrier {
	out := append([]barrier(nil), obstacles...)
	for _, item := range w.Creations {
		if item.Kind == "build" {
			out = append(out, barrier{item.Point, 1.1})
		}
	}
	return out
}

// Route searches walkable grid points, then smooths corners where there is a
// direct line. A visit to the same destination can use a different approach.
func route(a, b Point, blocks []barrier) []Point {
	if !clearSegment(b, b, blocks) {
		return nil
	}
	if clearSegment(a, b, blocks) {
		return []Point{b}
	}
	neighbors := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
	type node struct {
		p      Point
		parent int
	}
	origin := Point{math.Round(a.X), math.Round(a.Z)}
	queue := []node{}
	seen := map[Point]bool{}
	for _, d := range append([]Point{{0, 0}}, neighbors...) {
		p := Point{origin.X + d.X, origin.Z + d.Z}
		if clearSegment(a, p, blocks) {
			queue = append(queue, node{p, -1})
			seen[p] = true
		}
	}
	for i := 0; i < len(queue); i++ {
		p := queue[i].p
		if distance(p, b) < 1.6 && clearSegment(p, b, blocks) {
			path := []Point{b}
			for k := i; k >= 0; k = queue[k].parent {
				path = append(path, queue[k].p)
			}
			var smooth []Point
			from := a
			for k := len(path) - 2; k >= 0; k-- {
				if !clearSegment(from, path[k], blocks) {
					smooth = append(smooth, path[k+1])
					from = path[k+1]
				}
			}
			return append(smooth, b)
		}
		for _, d := range neighbors {
			q := Point{p.X + d.X, p.Z + d.Z}
			if q.X < -20 || q.X > 19 || q.Z < -16 || q.Z > 16 || seen[q] || !clearSegment(p, q, blocks) {
				continue
			}
			seen[q] = true
			queue = append(queue, node{q, i})
		}
	}
	return nil
}

func (w *World) leaveBarrier() {
	blocks := w.barriers()
	if clearSegment(w.Position, w.Position, blocks) {
		return
	}
	// Older saves can put the pet inside a tree, building, easel or burrow.
	// Move to the closest clear grid point the first time it walks again.
	best := math.Inf(1)
	p := w.Position
	for x := math.Floor(p.X) - 6; x <= p.X+6; x++ {
		for z := math.Floor(p.Z) - 6; z <= p.Z+6; z++ {
			q := Point{x, z}
			if d := distance(p, q); d < best && clearSegment(q, q, blocks) {
				w.Position, best = q, d
			}
		}
	}
}

func (w *World) move() bool {
	w.leaveBarrier()
	blocks := w.barriers()
	if len(w.Route) == 0 {
		// Opportunistic scenic routes when there is time and energy.
		if w.Action == "explore" && w.Energy > 55 && rand.IntN(3) == 0 {
			i := rand.IntN(len(landmarks))
			via := landmarks[i].Point
			toVia, onward := route(w.Position, via, blocks), route(via, w.Target, blocks)
			if len(toVia) > 0 && len(onward) > 0 && distance(via, w.Target) < 15 {
				w.Route = append(toVia, onward...)
			}
		}
		if len(w.Route) == 0 {
			w.Route = route(w.Position, w.Target, blocks)
		}
		if len(w.Route) == 0 {
			w.note("plan", "Could not reach "+w.Place+" and looked for another route")
			w.Target = w.meadowPoint()
			w.Route = route(w.Position, w.Target, blocks)
			w.Action, w.Place = "explore", "the meadow"
			return false
		}
	}
	dest := w.Route[0]
	d := distance(w.Position, dest)
	pace := 1.45
	if w.Action == "play" {
		pace = 1.85
	}
	if d <= pace {
		w.Position = dest
		w.Route = w.Route[1:]
		return len(w.Route) == 0
	}
	w.Position = Point{w.Position.X + (dest.X-w.Position.X)*pace/d, w.Position.Z + (dest.Z-w.Position.Z)*pace/d}
	return false
}
