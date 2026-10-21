package suite

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type ProgressStep struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	State  string `json:"state"`
	Detail string `json:"detail,omitempty"`
}
type Progress struct {
	State     string         `json:"state"`
	Phase     string         `json:"phase"`
	Error     string         `json:"error,omitempty"`
	Steps     []ProgressStep `json:"steps"`
	Completed int            `json:"completed"`
	Total     int            `json:"total"`
	Percent   int            `json:"percent"`
	StartedAt time.Time      `json:"startedAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	PID       int            `json:"pid"`
}
type ProgressTracker struct {
	mu    sync.Mutex
	value Progress
	path  string
}

func NewProgress(state string, steps []ProgressStep, path string) *ProgressTracker {
	for i := range steps {
		steps[i].State = "pending"
	}
	now := time.Now().UTC()
	return &ProgressTracker{value: Progress{State: state, Steps: steps, Total: len(steps), StartedAt: now, UpdatedAt: now, PID: os.Getpid()}, path: path}
}
func (p *ProgressTracker) save() {
	p.value.UpdatedAt = time.Now().UTC()
	p.value.Completed = 0
	for _, step := range p.value.Steps {
		if step.State == "done" || step.State == "warning" || step.State == "failed" || step.State == "blocked" {
			p.value.Completed++
		}
	}
	if p.value.Total > 0 {
		p.value.Percent = 100 * p.value.Completed / p.value.Total
	}
	if p.path != "" {
		_ = platform.SaveJSON(p.path, p.value)
	}
}
func (p *ProgressTracker) Snapshot() Progress {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := p.value
	result.Steps = append([]ProgressStep(nil), result.Steps...)
	return result
}
func (p *ProgressTracker) Step(id, state, detail string) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.value.Steps {
		if p.value.Steps[i].ID == id {
			p.value.Steps[i].State = state
			p.value.Steps[i].Detail = detail
			if state == "running" {
				p.value.Phase = p.value.Steps[i].Label
			}
			break
		}
	}
	p.save()
}
func (p *ProgressTracker) Finish(state, phase string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.value.State = state
	p.value.Phase = phase
	p.value.Error = ""
	if err != nil {
		p.value.Error = err.Error()
		for i := range p.value.Steps {
			if p.value.Steps[i].State == "running" {
				p.value.Steps[i].State = "failed"
				p.value.Steps[i].Detail = err.Error()
			}
		}
	}
	p.save()
}
func InstallProgress(tunnel bool) *ProgressTracker {
	steps := []ProgressStep{{ID: "prepare", Label: "Prepare storage and ports"}, {ID: "extract", Label: "Download and verify latest apps"}, {ID: "identity", Label: "Create your account and app sign-ins"}, {ID: "apps", Label: "Connect your storage and apps"}}
	if tunnel {
		steps = append(steps, ProgressStep{ID: "domain", Label: "Create secure domain connection"})
	}
	steps = append(steps, ProgressStep{ID: "finish", Label: "Save settings and shortcuts"})
	return NewProgress("installing", steps, "")
}
func (p *ProgressTracker) InstallPhase(phase string) {
	id := "prepare"
	switch {
	case strings.HasPrefix(phase, "Checking GitHub"), strings.HasPrefix(phase, "Checking oneworkspacestudio/mail"), strings.HasPrefix(phase, "Verified Mail"), strings.HasPrefix(phase, "Downloading"), strings.HasPrefix(phase, "Verifying"), strings.HasPrefix(phase, "Using this installer"):
		id = "extract"
	case strings.HasPrefix(phase, "Creating Accounts"):
		id = "identity"
	case strings.HasPrefix(phase, "Connecting Drive"):
		id = "apps"
	case strings.HasPrefix(phase, "Creating the authorized"):
		id = "domain"
	case strings.HasPrefix(phase, "Ready to launch"):
		id = "finish"
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.value.Steps {
		if p.value.Steps[i].ID == id {
			p.value.Steps[i].State = "running"
			p.value.Steps[i].Detail = phase
			break
		}
		p.value.Steps[i].State = "done"
	}
	p.value.Phase = phase
	p.save()
}

func StartupSteps(c Config) []ProgressStep {
	steps := []ProgressStep{{ID: "checks", Label: "Check storage and app ports"}}
	if c.Tunnel {
		steps = append(steps, ProgressStep{ID: "tunnel", Label: "Start secure connection"})
	}
	apps, _ := c.Catalog.Order()
	for _, app := range apps {
		if app.Kind != "process" {
			continue
		}
		steps = append(steps, ProgressStep{ID: app.ID, Label: "Start " + app.Name})
		if app.ID == "accounts" && c.Tunnel {
			steps = append(steps, ProgressStep{ID: "public-accounts", Label: "Make sign-in reachable on your domain"})
		}
	}
	return append(steps, ProgressStep{ID: "portal", Label: "Open your Workspace"})
}
