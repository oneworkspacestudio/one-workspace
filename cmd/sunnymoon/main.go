package main

import (
	"os"

	"github.com/oneworkspacestudio/one-workspace/internal/sunnymoon"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func main() {
	config := apphost.Config{
		ID: "sunnymoon", Name: "Sunny & Moon", Version: sunnymoon.Version, Listen: "127.0.0.1:9047",
	}
	if err := apphost.Main(config, sunnymoon.UI(), func(host *apphost.Host) (apphost.Application, error) {
		return nil, sunnymoon.Register(host)
	}); err != nil {
		apphost.Report(err)
		os.Exit(1)
	}
}
