package main

import (
	"os"

	"github.com/oneworkspacestudio/one-workspace/internal/pet"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func main() {
	cfg := apphost.Config{ID: "pet", Name: "Pet", Version: pet.Version, Listen: "127.0.0.1:9046", RestoreFiles: []string{"decisions"}}
	if err := apphost.Main(cfg, pet.UI(), func(h *apphost.Host) (apphost.Application, error) {
		return pet.New(h)
	}); err != nil {
		apphost.Report(err)
		os.Exit(1)
	}
}
