//go:build !windows

package websearch

func allowSandboxCode(path string) error { return nil }
