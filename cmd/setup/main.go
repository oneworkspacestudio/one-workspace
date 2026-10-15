package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fahimc/one-workspace/internal/distribution"
	"github.com/fahimc/one-workspace/internal/payload"
	"github.com/fahimc/one-workspace/internal/setupui"
	"github.com/fahimc/one-workspace/internal/suite"
	"os"
	"os/signal"
)

func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run() error {
	offline := flag.Bool("offline", false, "Deliberately install bundled versions without checking GitHub")
	checkUpdates := flag.Bool("check-updates", false, "Download and verify the latest GitHub bundle without installing")
	noBrowser := flag.Bool("no-browser", false, "Print setup link without opening a browser")
	headless := flag.Bool("headless", false, "Read setup JSON from standard input")
	resume := flag.String("resume", "", "Show startup progress for an existing protected installation")
	flag.Parse()
	z, c, e := payload.Open()
	if e != nil {
		return e
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if *checkUpdates {
		bundle, err := distribution.Latest(ctx, func(phase string) { fmt.Println(phase) })
		if err != nil {
			return err
		}
		defer bundle.Close()
		fmt.Println("Verified latest Workspace:", bundle.Catalog.Version)
		for _, app := range bundle.Catalog.Apps {
			fmt.Println(app.ID + ": " + app.Version)
		}
		return nil
	}
	if *resume != "" {
		installed, _, e := suite.Load(*resume)
		if e != nil {
			return e
		}
		return setupui.ServeInstalled(ctx, z, c, installed, !*noBrowser)
	}
	if *headless {
		var r suite.Request
		if e = json.NewDecoder(os.Stdin).Decode(&r); e != nil {
			return e
		}
		r.Offline = *offline
		conf, e := suite.Install(ctx, r, z, c, func(phase string) { fmt.Println(phase) })
		if e == nil {
			fmt.Println("Installed:", conf.Origin, "; launch", conf.Data+"/Workspace.exe")
		}
		return e
	}
	return setupui.Serve(ctx, z, c, !*noBrowser, *offline)
}
