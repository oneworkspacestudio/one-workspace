package main

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
	"github.com/oneworkspacestudio/one-workspace/internal/payload"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
)

func repairInstalled(ctx context.Context, data string, offline bool) error {
	z, c, e := payload.Open()
	if e != nil {
		return e
	}
	var mu sync.Mutex
	enc := json.NewEncoder(os.Stdout)
	send := func(event suite.InstallerRepairEvent) { mu.Lock(); defer mu.Unlock(); _ = enc.Encode(event) }
	ctx = operationprogress.Observe(ctx, func(event operationprogress.Event) { send(suite.InstallerRepairEvent{Progress: &event}) })
	_, e = suite.ReinstallInstalled(ctx, data, z, c, offline, func(phase string) { send(suite.InstallerRepairEvent{Phase: phase}) })
	return e
}
