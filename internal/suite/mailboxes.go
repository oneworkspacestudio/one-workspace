package suite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fahimc/one-workspace/internal/platform"
)

type mailboxUser struct {
	Subject  string `json:"subject"`
	Username string `json:"username"`
	Name     string `json:"displayName"`
}

// MyAccounts remains authoritative. Polling reconciles users created while Mail
// was offline, without changing contact/recovery email or relying on first login.
func (rt *Runtime) maintainMailboxes(ctx context.Context) {
	lastState := ""
	api := newIdentity(rt.Config)
	for {
		attempt, cancel := context.WithTimeout(ctx, 2*time.Minute)
		count, err := rt.syncMailboxes(attempt, api)
		cancel()
		state := "ready"
		message := fmt.Sprintf("%d MyAccounts mailboxes connected. New users are checked every 15 seconds.", count)
		if err != nil {
			state = "attention"
			message = err.Error()
		}
		_ = platform.SaveJSON(filepath.Join(rt.Config.Data, "mailbox-sync.json"), map[string]any{"state": state, "message": message, "updated": time.Now().UTC()})
		if state != lastState {
			log := &platform.Log{Path: filepath.Join(rt.Config.Data, "logs", "mailbox-sync.log")}
			_, _ = fmt.Fprintln(log, time.Now().UTC().Format(time.RFC3339), message)
			lastState = state
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(15 * time.Second):
		}
	}
}

func (rt *Runtime) syncMailboxes(ctx context.Context, api *identityAPI) (int, error) {
	app, ok := rt.Config.Catalog.Find("mail")
	if !ok {
		return 0, nil
	}
	count := 0
	seen := map[string]bool{}
	for page := 1; page <= 10000; page++ {
		var users []struct {
			ID          string `json:"id"`
			Owner       string `json:"owner"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Deleted     bool   `json:"isDeleted"`
			Forbidden   bool   `json:"isForbidden"`
		}
		path := fmt.Sprintf("/api/get-users?owner=myaccounts&pageSize=100&p=%d&sortField=name&sortOrder=ascend", page)
		if err := api.call(ctx, path, nil, &users); err != nil {
			if api.login(ctx, rt.Secrets) != nil || api.call(ctx, path, nil, &users) != nil {
				return count, errors.New("Mailbox synchronization could not read the Workspace user directory. It will retry automatically.")
			}
		}
		batch := []mailboxUser{}
		for _, u := range users {
			if u.Deleted || u.Forbidden {
				continue
			}
			if u.Owner != "myaccounts" || u.ID == "" || u.Name == "" {
				return count, errors.New("MyAccounts returned an invalid mailbox identity; synchronization stopped safely")
			}
			if seen[u.ID] {
				continue
			}
			seen[u.ID] = true
			batch = append(batch, mailboxUser{u.ID, u.Name, u.DisplayName})
		}
		if len(batch) > 0 {
			body, _ := json.Marshal(map[string]any{"users": batch})
			origin := fmt.Sprintf("http://localhost:%d", app.Port)
			req, err := http.NewRequestWithContext(ctx, "POST", origin+"/api/mailboxes/provision", bytes.NewReader(body))
			if err != nil {
				return count, err
			}
			req.Header.Set("Origin", origin)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+rt.mailProvisionToken)
			client := &http.Client{Timeout: 90 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
			resp, err := client.Do(req)
			if err != nil {
				return count, errors.New("Mailboxes are waiting for Mail to start. Synchronization will retry automatically.")
			}
			var result struct {
				Mailboxes map[string]string `json:"mailboxes"`
				Error     string            `json:"error"`
			}
			err = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result)
			resp.Body.Close()
			if err != nil || resp.StatusCode != 200 || len(result.Mailboxes) != len(batch) {
				if result.Error != "" {
					return count, errors.New("Mailbox setup needs attention: " + result.Error)
				}
				return count, errors.New("Mail did not confirm all mailbox assignments. Synchronization will retry automatically.")
			}
			for _, u := range batch {
				if result.Mailboxes[u.Subject] == "" {
					return count, errors.New("Mail did not confirm the matching MyAccounts identity. Synchronization will retry automatically.")
				}
			}
			count += len(batch)
		}
		if len(users) < 100 {
			return count, nil
		}
	}
	return count, errors.New("The user directory exceeded the mailbox synchronization limit")
}
