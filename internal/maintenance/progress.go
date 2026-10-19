package maintenance

import (
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
)

type operationStep struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	State string `json:"state"`
}
type operationLog struct {
	At   time.Time `json:"at"`
	Text string    `json:"text"`
}
type operationState struct {
	Kind             string          `json:"kind"`
	Status           string          `json:"status"`
	Stage            string          `json:"stage"`
	Detail           string          `json:"detail"`
	Steps            []operationStep `json:"steps"`
	History          []operationLog  `json:"history"`
	StartedAt        time.Time       `json:"startedAt"`
	StageStartedAt   time.Time       `json:"stageStartedAt"`
	LastActivityAt   time.Time       `json:"lastActivityAt"`
	FinishedAt       time.Time       `json:"finishedAt,omitempty"`
	Done             int64           `json:"done"`
	Total            int64           `json:"total"`
	Unit             string          `json:"unit,omitempty"`
	ElapsedSeconds   int64           `json:"elapsedSeconds"`
	BytesPerSecond   float64         `json:"bytesPerSecond"`
	RemainingSeconds int64           `json:"remainingSeconds"`
	Recovering       bool            `json:"recovering"`
}

func (s *Server) beginLocked(kind string) {
	now := time.Now().UTC()
	steps := []operationStep{{ID: "prepare", Label: "Prepare update"}, {ID: "download", Label: "Download releases"}, {ID: "stage", Label: "Check downloaded packages"}, {ID: "preflight", Label: "Check storage and restart"}, {ID: "backup", Label: "Back up your workspace"}, {ID: "verify", Label: "Verify backup"}, {ID: "install", Label: "Install updates"}, {ID: "restart", Label: "Restart and check apps"}, {ID: "cleanup", Label: "Clean up old update files"}}
	if kind == "check" {
		steps = []operationStep{{ID: "check", Label: "Check published app releases"}}
	}
	if kind == "recover" || kind == "rollback" {
		steps = []operationStep{{ID: "recovery", Label: "Preserve current data and restore verified backup"}, {ID: "restart", Label: "Restart and check apps"}}
	}
	for i := range steps {
		steps[i].State = "pending"
	}
	steps[0].State = "running"
	s.Operation = &operationState{Kind: kind, Status: "running", Stage: steps[0].ID, Detail: "Preparing your request…", StartedAt: now, StageStartedAt: now, LastActivityAt: now, Steps: steps, RemainingSeconds: -1}
	s.Phase = "Preparing your request…"
}

func (s *Server) detail(e operationprogress.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.Operation
	if p == nil || p.Status != "running" {
		return
	}
	stage := e.Stage
	enteringRecovery := stage == "recovery" && !p.Recovering
	if stage == "recovery" {
		p.Recovering = true
	}
	if p.Recovering || p.Kind == "rollback" || p.Kind == "recover" {
		if stage != "restart" {
			stage = "recovery"
		}
	}
	index := -1
	for i := range p.Steps {
		if p.Steps[i].ID == stage {
			index = i
			break
		}
	}
	if index < 0 {
		p.Steps = append(p.Steps, operationStep{ID: stage, Label: "Restore previous installation", State: "pending"})
		index = len(p.Steps) - 1
	}
	if p.Stage != stage {
		for i := range p.Steps {
			if p.Steps[i].State == "running" {
				p.Steps[i].State = "done"
				if enteringRecovery {
					p.Steps[i].State = "failed"
				}
			}
			if i < index && p.Steps[i].State == "pending" {
				p.Steps[i].State = "skipped"
			}
		}
		p.Steps[index].State = "running"
	}
	now := time.Now().UTC()
	if p.Stage != stage || p.Detail != e.Detail || e.Done < p.Done || p.Unit != e.Unit || p.Total != e.Total {
		p.StageStartedAt = now
	}
	if p.Stage != stage || p.Detail != e.Detail || p.Done != e.Done || p.Total != e.Total {
		p.LastActivityAt = now
	}
	p.Stage, p.Detail, p.Done, p.Total, p.Unit = stage, e.Detail, e.Done, e.Total, e.Unit
	s.Phase = e.Detail
}

func (s *Server) snapshotLocked() *operationState {
	if s.Operation == nil {
		return nil
	}
	p := *s.Operation
	p.Steps = append([]operationStep(nil), p.Steps...)
	p.History = append([]operationLog(nil), p.History...)
	now := time.Now()
	if !p.FinishedAt.IsZero() {
		now = p.FinishedAt
	}
	p.ElapsedSeconds = int64(now.Sub(p.StartedAt).Seconds())
	p.RemainingSeconds = -1
	seconds := now.Sub(p.StageStartedAt).Seconds()
	if p.Status == "running" && p.Unit == "bytes" && p.Done > 0 && p.Total > p.Done && seconds >= 2 && time.Since(p.LastActivityAt) < 15*time.Second {
		p.BytesPerSecond = float64(p.Done) / seconds
		p.RemainingSeconds = int64(float64(p.Total-p.Done) / p.BytesPerSecond)
	}
	return &p
}
