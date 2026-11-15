package suite

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func workspaceUnavailable(w http.ResponseWriter, message string) {
	w.Header().Set("Retry-After", "10")
	w.Header().Set("Cache-Control", "no-store")
	http.Error(w, message, http.StatusServiceUnavailable)
}

func workspaceOrigin(port int, handler http.Handler) (*http.Server, <-chan error, error) {
	listener, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, nil, err
	}
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	return server, done, nil
}

func updateMaintenanceOrigin(port int) (*http.Server, error) {
	server, _, err := workspaceOrigin(port, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		workspaceUnavailable(w, "Workspace is being updated. Try again shortly.")
	}))
	return server, err
}
