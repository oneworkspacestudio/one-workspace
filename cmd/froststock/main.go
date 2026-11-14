package main

import (
	"os"

	"github.com/oneworkspacestudio/one-workspace/internal/froststock"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func main() {
	config := apphost.Config{
		ID: "froststock", Name: "Froststock", Version: froststock.Version,
		Listen: "127.0.0.1:9048", MediaCapture: true,
	}
	if err := apphost.Main(config, froststock.UI(), func(host *apphost.Host) (apphost.Application, error) {
		return nil, froststock.Register(host, os.Getenv("FROSTSTOCK_DRIVE_ROOT"))
	}); err != nil {
		apphost.Report(err)
		os.Exit(1)
	}
}
