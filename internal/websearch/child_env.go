package websearch

import "os"

// Browser, OCR and model children need Windows runtime variables, not the
// manager's OIDC, registry, publication or other app credentials.
func childEnvironment() []string {
	result := []string{}
	for _, key := range []string{"SystemRoot", "WINDIR", "SystemDrive", "PATH", "TEMP", "TMP", "LOCALAPPDATA", "APPDATA", "USERPROFILE", "ProgramData", "ProgramFiles", "ProgramFiles(x86)", "PROCESSOR_ARCHITECTURE", "NUMBER_OF_PROCESSORS"} {
		if value := os.Getenv(key); value != "" {
			result = append(result, key+"="+value)
		}
	}
	return result
}
