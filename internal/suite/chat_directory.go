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

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type chatDirectoryUser struct {
	Subject  string `json:"subject"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

func (rt *Runtime) maintainChatDirectory(ctx context.Context) {
	api := newIdentity(rt.Config)
	for {
		attempt, cancel := context.WithTimeout(ctx, 90*time.Second)
		count, err := rt.syncChatDirectory(attempt, api)
		cancel()
		state, message := "ready", fmt.Sprintf("%d MyAccounts members connected. New accounts are checked every 15 seconds.", count)
		if err != nil {
			state, message = "attention", err.Error()
		}
		_ = platform.SaveJSON(filepath.Join(rt.Config.Data, "chat-directory-sync.json"), map[string]any{"state": state, "message": message, "updated": time.Now().UTC()})
		select {
		case <-ctx.Done():
			return
		case <-time.After(15 * time.Second):
		}
	}
}

func (rt *Runtime) syncChatDirectory(ctx context.Context, api *identityAPI) (int, error) {
	app, ok := rt.Config.Catalog.Find("chat")
	if !ok {
		return 0, nil
	}
	members := []chatDirectoryUser{}
	seen := map[string]bool{}
	// Send a complete snapshot only after every directory page succeeds.
	for page := 1; page <= 101; page++ {
		var users []struct {
			ID          string `json:"id"`
			Owner       string `json:"owner"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Email       string `json:"email"`
			Verified    bool   `json:"emailVerified"`
			Deleted     bool   `json:"isDeleted"`
			Forbidden   bool   `json:"isForbidden"`
		}
		path := fmt.Sprintf("/api/get-users?owner=myaccounts&pageSize=100&p=%d&sortField=name&sortOrder=ascend", page)
		if err := api.call(ctx, path, nil, &users); err != nil {
			if api.login(ctx, rt.Secrets) != nil || api.call(ctx, path, nil, &users) != nil {
				return 0, errors.New("Chat is waiting for the MyAccounts member directory. It will retry automatically.")
			}
		}
		for _, u := range users {
			if u.Owner != "myaccounts" || u.ID == "" || u.Name == "" || seen[u.ID] {
				return 0, errors.New("MyAccounts returned an invalid Chat identity; synchronization stopped safely")
			}
			seen[u.ID] = true
			email := ""
			if u.Verified {
				email = u.Email
			}
			members = append(members, chatDirectoryUser{u.ID, u.Name, u.DisplayName, email, u.Deleted || u.Forbidden})
		}
		if len(members) > 10000 {
			return 0, errors.New("The Chat directory exceeds the current 10,000 member limit")
		}
		if len(users) < 100 {
			break
		}
	}
	body, _ := json.Marshal(map[string]any{"users": members})
	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://127.0.0.1:%d/api/directory/sync", app.Port), bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+rt.chatDirectoryToken)
	client := &http.Client{Timeout: 65 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(request)
	if err != nil {
		return 0, errors.New("Chat is reconnecting; member synchronization will retry automatically")
	}
	defer response.Body.Close()
	var result struct {
		Members int `json:"members"`
	}
	if response.StatusCode != 200 || json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result) != nil {
		return 0, errors.New("Chat could not synchronize MyAccounts members. It will retry automatically.")
	}
	return result.Members, nil
}
