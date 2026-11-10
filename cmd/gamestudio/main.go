package main

import (
	"os"

	"github.com/oneworkspacestudio/one-workspace/internal/gamestudio"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func main() {
	cfg := apphost.Config{ID: "gamestudio", Name: "Game Studio", Version: gamestudio.Version, Listen: "127.0.0.1:9045", RestoreFiles: []string{"studio-settings"}}
	if err := apphost.Main(cfg, gamestudio.UI(), func(h *apphost.Host) (apphost.Application, error) {
		_, err := gamestudio.New(h)
		return nil, err
	}); err != nil {
		apphost.Report(err)
		os.Exit(1)
	}
}
