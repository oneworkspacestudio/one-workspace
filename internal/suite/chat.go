package suite

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Clearing a browser cookie alone cannot close an already upgraded WebSocket.
// End the same Chat session on its loopback gateway before Workspace sign-out.
// Forward only Chat's own cookie, never the provider or portal credential.
func (p *Portal) closeChatSession(r *http.Request) {
	app, installed := p.rt.Config.Catalog.Find("chat")
	cookie, err := r.Cookie("chat_session")
	if !installed || err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/auth/logout", app.Port), nil)
	if err != nil {
		return
	}
	req.AddCookie(cookie)
	req.Header.Set("Origin", p.rt.Config.Origin)
	client := &http.Client{Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(req)
	if err == nil {
		response.Body.Close()
	}
}
