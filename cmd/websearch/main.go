package main

import (
	"github.com/oneworkspacestudio/one-workspace/internal/websearch"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "install-model" {
		if e := websearch.InstallModelCLI(os.Args[2:]); e != nil {
			apphost.Report(e)
			os.Exit(1)
		}
		return
	}
	if e := websearch.OwnChildren(); e != nil {
		apphost.Report(e)
		os.Exit(1)
	}
	if e := apphost.Main(apphost.Config{ID: "websearch", Name: "Web Search", Version: websearch.Version, Listen: "127.0.0.1:9043", RestoreFiles: []string{"browser-profile/"}}, websearch.UI(), func(h *apphost.Host) (apphost.Application, error) { return websearch.New(h) }); e != nil {
		apphost.Report(e)
		os.Exit(1)
	}
}
