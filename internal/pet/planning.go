package pet

import "math/rand/v2"

func (w World) validIdea(d Decision) bool {
	if !w.allows(d.Action) || d.Style < 0 || d.Style > 3 {
		return false
	}
	if d.Via != "" {
		if _, ok := featureByID(landmarks, d.Via); !ok {
			return false
		}
	}
	switch d.Action {
	case "explore", "play":
		if d.Place != "" {
			if _, ok := featureByID(landmarks, d.Place); !ok {
				return false
			}
		}
		return d.Painting == nil && d.Design == "" && d.Material == ""
	case "draw":
		return d.Place == "" && d.Design == "" && d.Material == "" && (d.Painting == nil || w.novelPainting(d.Painting))
	case "build":
		if d.Painting != nil {
			return false
		}
		if w.Project != nil {
			p := w.Project
			return (d.Place == "" || d.Place == p.Site) && (d.Design == "" || d.Design == p.Design) && (d.Material == "" || d.Material == p.Material)
		}
		_, ok := featureByID(w.buildOptions(), d.Place)
		return ok && contains(w.learned(), d.Design) && contains([]string{"wood", "stone", "reed"}, d.Material)
	}
	return false
}

func (w World) fallback(action string) Decision {
	d := Decision{Action: action, Source: "Instincts"}
	switch action {
	case "explore", "play":
		choices := []feature{}
		for _, f := range landmarks {
			if distance(w.Position, f.Point) > 2 && !contains(w.VisitTrail, "the "+f.ID) {
				choices = append(choices, f)
			}
		}
		if len(choices) <= 1 {
			choices = landmarks
		}
		if action == "explore" || rand.IntN(2) == 0 {
			d.Place = choices[rand.IntN(len(choices))].ID
		}
	case "draw":
		d.Painting = w.art()
	case "build":
		if w.Project == nil {
			s := w.buildOptions()
			if len(s) == 0 {
				return Decision{Action: "explore", Source: "Instincts"}
			}
			d.Place = s[rand.IntN(len(s))].ID
			d.Design = w.learned()[rand.IntN(len(w.learned()))]
			d.Material = resources[rand.IntN(len(resources))].Material
		}
	}
	return d
}

func (w *World) beginIdea(d Decision) {
	w.start(d.Action, d.Style, d.Source)
	switch d.Action {
	case "explore", "play":
		if f, ok := featureByID(landmarks, d.Place); ok {
			w.Target, w.Place = f.Point, "the "+f.ID
		}
	case "draw":
		w.Painting = d.Painting
		if w.Painting == nil {
			w.Painting = w.art()
		}
	case "build":
		if !w.planBuild(d) {
			w.start("explore", 0, "Instincts")
		}
	}
	if f, ok := featureByID(landmarks, d.Via); ok {
		blocks := w.barriers()
		first, second := route(w.Position, f.Point, blocks), route(f.Point, w.Target, blocks)
		if len(first) > 0 && len(second) > 0 {
			w.Route = append(first, second...)
			w.note("plan", "Took a route via the "+f.ID+" toward "+w.Place)
		}
	}
	if w.Source != "Instincts" {
		w.note("plan", w.Source+" planned "+d.Action+" near "+w.Place)
	}
}

func (w World) observation() map[string]any {
	blocks := w.barriers()
	actions, places, availableSites := []string{}, []feature{}, []feature{}
	for _, f := range landmarks {
		if len(route(w.Position, f.Point, blocks)) > 0 {
			places = append(places, f)
		}
	}
	for _, f := range w.buildOptions() {
		if _, ok := w.approach(f); ok {
			availableSites = append(availableSites, f)
		}
	}
	for _, a := range []string{"explore", "play", "draw", "build"} {
		if w.allows(a) && (a != "build" || w.Project != nil || len(availableSites) > 0) {
			actions = append(actions, a)
		}
	}
	type memory struct {
		Tick     int64  `json:"tick"`
		Activity string `json:"activity"`
	}
	last := []memory{}
	for _, e := range w.Journal[max(0, len(w.Journal)-6):] {
		last = append(last, memory{e.Tick, e.Kind})
	}
	buildings := []Item{}
	artTitles := []string{}
	for _, c := range w.Creations {
		if c.Kind == "build" {
			buildings = append(buildings, c)
		} else if c.Painting != nil {
			artTitles = append(artTitles, c.Painting.Title)
		}
	}
	return map[string]any{
		"actions": actions, "position": w.Position, "destination": w.Target, "currentAction": w.Action, "currentPlace": w.Place,
		"weather": w.Weather, "minuteOfDay": w.Minute, "mood": w.Mood, "energy": int(w.Energy), "fullness": int(w.Fullness),
		"happiness": int(w.Happiness), "curiosity": int(w.Curiosity), "foodLocations": w.Food,
		"reachablePlaces": places, "buildSites": availableSites, "constructionExperience": w.BuildXP, "barriers": blocks, "resourceSites": resources,
		"recentPlaces": w.VisitTrail, "recentActivities": last, "structures": buildings, "buildingDesigns": w.learned(),
		"project": w.Project, "recentArt": artTitles[max(0, len(artTitles)-5):],
	}
}
