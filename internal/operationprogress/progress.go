// Package operationprogress carries optional live maintenance detail without
// changing command-line reporters or making progress part of the transaction.
package operationprogress

import (
	"context"
	"time"
)

type Event struct {
	Stage  string `json:"stage"`
	Detail string `json:"detail"`
	Done   int64  `json:"done"`
	Total  int64  `json:"total"`
	Unit   string `json:"unit,omitempty"`
}

type observerKey struct{}

func Observe(ctx context.Context, fn func(Event)) context.Context {
	return context.WithValue(ctx, observerKey{}, fn)
}

func Report(ctx context.Context, stage, detail string, done, total int64, unit string) {
	if fn, ok := ctx.Value(observerKey{}).(func(Event)); ok {
		fn(Event{Stage: stage, Detail: detail, Done: done, Total: total, Unit: unit})
	}
}

// Meter observes streaming I/O, including hashing and disk copies. Completion
// of the meter means bytes were processed, not that the operation succeeded.
type Meter struct {
	ctx           context.Context
	stage, detail string
	done, total   int64
	last          time.Time
}

func Bytes(ctx context.Context, stage, detail string, total int64) *Meter {
	m := &Meter{ctx: ctx, stage: stage, detail: detail, total: total}
	m.emit()
	return m
}

func (m *Meter) Write(b []byte) (int, error) {
	m.done += int64(len(b))
	if time.Since(m.last) >= 250*time.Millisecond {
		m.emit()
	}
	return len(b), nil
}

func (m *Meter) Finish() { m.emit() }
func (m *Meter) emit() {
	m.last = time.Now()
	Report(m.ctx, m.stage, m.detail, m.done, m.total, "bytes")
}
