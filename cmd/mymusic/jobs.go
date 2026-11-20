package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

type generationJob struct {
	ID    string `json:"id"`
	Done  bool   `json:"done"`
	Track *track `json:"track,omitempty"`
	Error string `json:"error,omitempty"`
}

func (app application) generationRoutes(handle func(string, http.HandlerFunc)) {
	var mu sync.Mutex
	var active string
	jobs := make(map[string]generationJob)
	var order []string
	handle("POST /api/generate", func(w http.ResponseWriter, r *http.Request) {
		var song generation
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10)).Decode(&song); err != nil {
			respond(w, nil, errors.New("invalid generation request"))
			return
		}
		mu.Lock()
		if active != "" {
			job := jobs[active]
			mu.Unlock()
			respond(w, job, nil)
			return
		}
		job := generationJob{ID: rand.Text()}
		active = job.ID
		jobs[job.ID] = job
		order = append(order, job.ID)
		if len(order) > 8 {
			delete(jobs, order[0])
			order = order[1:]
		}
		mu.Unlock()
		go func(job generationJob) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()
			item, err := app.generate(ctx, song)
			job.Done = true
			if err != nil {
				job.Error = err.Error()
			} else {
				job.Track = &item
			}
			mu.Lock()
			defer mu.Unlock()
			jobs[job.ID] = job
			active = ""
		}(job)
		respond(w, job, nil)
	})
	handle("GET /api/generate/{id}", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		job, ok := jobs[r.PathValue("id")]
		mu.Unlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		respond(w, job, nil)
	})
}
