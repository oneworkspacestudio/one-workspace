//go:build windows

package websearch

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:embed ocr.ps1
var ocrScript string

func screenText(ctx context.Context, data, image string) (string, error) {
	root, e := windows.GetWindowsDirectory()
	if e != nil {
		return "", e
	}
	child, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(child, filepath.Join(root, "System32", "WindowsPowerShell", "v1.0", "powershell.exe"), "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", ocrScript)
	hideProcess(cmd)
	cmd.Env = append(childEnvironment(), "WEBSEARCH_OCR_PROBE=0")
	if image == "" {
		cmd.Env = append(cmd.Env, "WEBSEARCH_OCR_PROBE=1")
	} else {
		png, e := base64.StdEncoding.DecodeString(image)
		if e != nil || len(png) > 16<<20 {
			return "", errors.New("invalid browser screenshot")
		}
		dir := filepath.Join(data, "ocr-temp")
		if e = os.MkdirAll(dir, 0700); e != nil {
			return "", e
		}
		file, e := os.CreateTemp(dir, "page-*.png")
		if e != nil {
			return "", e
		}
		defer os.Remove(file.Name())
		_, e = file.Write(png)
		closeErr := file.Close()
		if e != nil {
			return "", e
		}
		if closeErr != nil {
			return "", closeErr
		}
		cmd.Env = append(cmd.Env, "WEBSEARCH_OCR_IMAGE="+file.Name())
	}
	out, e := cmd.Output()
	if e != nil {
		return "", errors.New("Windows screenshot text recognition is unavailable; enable an English OCR language in Windows")
	}
	var result struct {
		Available bool   `json:"available"`
		Text      string `json:"text"`
	}
	if json.Unmarshal(out, &result) != nil || !result.Available {
		return "", errors.New("Windows screenshot text recognition did not return a valid result")
	}
	if len(result.Text) > 10000 {
		result.Text = result.Text[:10000]
	}
	return result.Text, nil
}
