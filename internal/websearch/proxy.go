package websearch

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

// Browser egress is pinned to a checked public IP, including every redirect,
// subresource and websocket. No private-network bypass or DNS rebinding.
func publicIP(ctx context.Context, host string) (string, error) {
	ips, e := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	if e != nil {
		return "", e
	}
	if len(ips) == 0 {
		return "", errors.New("no address")
	}
	for _, ip := range ips {
		ip = ip.Unmap()
		if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() || inReserved(ip) {
			return "", errors.New("private network blocked")
		}
	}
	return ips[0].String(), nil
}
func inReserved(ip netip.Addr) bool {
	for _, s := range []string{"100.64.0.0/10", "192.0.0.0/24", "192.0.2.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "240.0.0.0/4", "2001:db8::/32", "64:ff9b::/96"} {
		if netip.MustParsePrefix(s).Contains(ip) {
			return true
		}
	}
	return false
}
func startProxy(ctx context.Context) (string, *http.Server, error) {
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		return "", nil, e
	}
	srv := &http.Server{ReadHeaderTimeout: 5 * time.Second, MaxHeaderBytes: 16384}
	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT" {
			http.Error(w, "HTTPS required", 403)
			return
		}
		host, port, e := net.SplitHostPort(r.Host)
		if e != nil || port != "443" || strings.ContainsAny(host, "%\\") {
			http.Error(w, "destination blocked", 403)
			return
		}
		check, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		ip, e := publicIP(check, host)
		if e != nil {
			http.Error(w, "destination blocked", 403)
			return
		}
		up, e := (&net.Dialer{Timeout: 15 * time.Second}).DialContext(ctx, "tcp", net.JoinHostPort(ip, "443"))
		if e != nil {
			http.Error(w, "unavailable", 502)
			return
		}
		down, rw, e := w.(http.Hijacker).Hijack()
		if e != nil {
			up.Close()
			return
		}
		defer down.Close()
		defer up.Close()
		_, _ = rw.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n")
		if rw.Flush() != nil {
			return
		}
		done := make(chan struct{})
		go func() { _, _ = io.Copy(up, rw); close(done) }()
		go func() {
			select {
			case <-ctx.Done():
				down.Close()
				up.Close()
			case <-done:
			}
		}()
		_, _ = io.Copy(down, up)
	})
	go func() { _ = srv.Serve(l) }()
	go func() { <-ctx.Done(); _ = srv.Close() }()
	return "http://" + l.Addr().String(), srv, nil
}
