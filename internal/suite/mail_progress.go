package suite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func (rt *Runtime) mailHostingProgress(ctx context.Context, app catalog.App) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	client := &http.Client{Timeout: 5 * time.Second}
	for {
		req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d/api/status", app.Port), nil)
		resp, e := client.Do(req)
		var status struct {
			Configured     bool   `json:"configured"`
			Busy           bool   `json:"hostingBusy"`
			Phase          string `json:"hostingPhase"`
			Error          string `json:"hostingError"`
			IdentityError  string `json:"identityError"`
			InboundState   string `json:"inboundState"`
			InboundMessage string `json:"inboundMessage"`
		}
		if e == nil {
			e = json.NewDecoder(resp.Body).Decode(&status)
			resp.Body.Close()
		}
		if e == nil && status.IdentityError != "" {
			rt.progress.Step(app.ID, "warning", status.IdentityError+" Open Mail settings to repair mailbox access.")
			return
		}
		if e == nil && status.Error != "" {
			rt.progress.Step(app.ID, "warning", status.Error)
			return
		}
		if e == nil && status.Configured && !status.Busy {
			switch status.InboundState {
			case "reachable":
				rt.progress.Step(app.ID, "done", "Mail is available and its incoming TCP25 connection was verified. Message delivery and spam placement are separate checks.")
				return
			case "attention", "unverified", "configured":
				rt.progress.Step(app.ID, "warning", status.InboundMessage+" Open Mail setup for network guidance. Other apps are ready to use.")
				return
			case "checking":
				status.Phase = status.InboundMessage
			default:
				rt.progress.Step(app.ID, "warning", "Mail is available locally, but incoming Internet email has not been verified. Open Mail setup to check router TCP25, Windows Security and your ISP.")
				return
			}
		}
		detail := status.Phase
		if detail == "" {
			detail = "Preparing your mailbox service"
		}
		rt.progress.Step(app.ID, "running", detail)
		select {
		case <-ctx.Done():
			rt.progress.Step(app.ID, "warning", "Mail setup is still incomplete. Open Mail for details; the other Workspace apps can be used.")
			return
		case <-time.After(2 * time.Second):
		}
	}
}
