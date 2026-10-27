//go:build windows

package platform

func protectSecret(_ string, b []byte) ([]byte, error)   { return Protect(b) }
func unprotectSecret(_ string, b []byte) ([]byte, error) { return Unprotect(b) }
