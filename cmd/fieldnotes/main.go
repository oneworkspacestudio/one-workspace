package main

import (
	"os"

	"github.com/oneworkspacestudio/one-workspace/internal/fieldnotes"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func main() {
	config := apphost.Config{
		ID: "fieldnotes", Name: "AI Fieldnotes", Version: fieldnotes.Version,
		Listen: "127.0.0.1:9050", RestoreFiles: []string{"providers.secret"},
	}
	if err := apphost.Main(config, fieldnotes.UI(), fieldnotes.Register); err != nil {
		apphost.Report(err)
		os.Exit(1)
	}
}
