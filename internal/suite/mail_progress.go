package suite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

func (rt *Runtime) mailHostingProgress(ctx context.Context, app catalog.App) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	client := &http.Client{Timeout: 5 * time.Second}
	for {
		req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d/api/status", app.Port), nil)
		resp, e := client.Do(req)
		var status struct {
			MainMailURL    string `json:"mainMailURL"`
			Configured     bool   `json:"configured"`
			Busy           bool   `json:"hostingBusy"`
			Phase          string `json:"hostingPhase"`
			Error          string `json:"hostingError"`
			IdentityError  string `json:"identityError"`
			InboundState   string `json:"inboundState"`
			InboundMessage string `json:"inboundMessage"`
			Issue          struct {
				Code string `json:"code"`
			} `json:"hostingIssue"`
		}
		if e == nil {
			e = json.NewDecoder(resp.Body).Decode(&status)
			resp.Body.Close()
		}
		if e == nil && status.MainMailURL != "" {
			rt.progress.Step(app.ID, "done", "Mail uses the main Workspace. Its availability is checked when Mail opens; incoming delivery is managed there.")
			return
		}
		if e == nil && status.IdentityError != "" {
			rt.logMailSetupProblem(status.IdentityError)
			rt.progress.Step(app.ID, "warning", "Mail couldn't finish connecting your Workspace account. Open Mail and choose Reconnect Mail. Your other apps are ready to use.")
			return
		}
		if e == nil && status.Error != "" {
			rt.logMailSetupProblem(status.Error)
			message := "Mail setup couldn't finish. Open Mail and choose Retry automatic setup. Your other apps are ready to use."
			if status.Configured {
				message = "Mail is available on this computer, but Internet email needs attention. Open Mail for help finishing setup."
			}
			if status.Issue.Code == "existing_mail_dns" {
				message = "Mail found an existing email service on your domain. Open Mail to choose how you want to connect it."
			}
			rt.progress.Step(app.ID, "warning", message)
			return
		}
		if e == nil && status.Configured && !status.Busy {
			switch status.InboundState {
			case "reachable":
				rt.progress.Step(app.ID, "done", "Mail is ready, and its incoming email connection was checked. Individual message delivery is checked separately.")
				return
			case "attention", "unverified", "configured":
				rt.logMailSetupProblem(status.InboundMessage)
				rt.progress.Step(app.ID, "warning", "Mail is ready on this computer. Incoming Internet email still needs a connection check. Open Mail for help finishing this step; your other apps are ready to use.")
				return
			case "checking":
				status.Phase = "Checking your incoming email connection…"
			default:
				rt.progress.Step(app.ID, "warning", "Mail is ready on this computer. Open Mail to finish checking its incoming Internet email connection.")
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

func (rt *Runtime) logMailSetupProblem(detail string) {
	if rt.Config.Data == "" {
		return
	}
	log := &platform.Log{Path: filepath.Join(rt.Config.Data, "logs", "mail-setup.log")}
	_, _ = fmt.Fprintln(log, time.Now().UTC().Format(time.RFC3339), detail)
}
