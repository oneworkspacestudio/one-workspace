// A tiny, self-contained future app. The Workspace gateway requires sign-in.
// Real apps should also implement their own authorization for sensitive actions.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.Int("port", 9004, "Loopback port")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "ok") })
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html><title>Hello Workspace</title><body style="font:20px system-ui;margin:8vw;background:#f4f6fc;color:#17233b"><h1>Hello, Workspace.</h1><p>This app was added from a manifest after installation.</p><a href="/">Back to Workspace</a>`)
	})
	server := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", *port), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(server.ListenAndServe())
}
