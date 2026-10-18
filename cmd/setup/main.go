package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/maintenance"
	"github.com/oneworkspacestudio/one-workspace/internal/payload"
	"github.com/oneworkspacestudio/one-workspace/internal/setupui"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
	"os"
	"os/signal"
	"strings"
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
	update := flag.String("update", "", "Open updates for an existing protected private data folder")
	apply := flag.Bool("apply", false, "Apply selected updates headlessly; requires --update")
	apps := flag.String("apps", "", "Comma-separated update groups; empty selects all available updates")
	bundlePath := flag.String("bundle", "", "Explicit verified offline update bundle ZIP")
	backupRoot := flag.String("backup-root", "", "Local NTFS backup folder (default: Private/backups)")
	rollback := flag.String("rollback", "", "Restore a recorded backup ID; requires --update")
	recoverUpdate := flag.Bool("recover", false, "Recover an interrupted update; requires --update")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if *update != "" {
		report := func(s string) { fmt.Println(s) }
		if *recoverUpdate {
			return suite.RecoverUpdate(ctx, *update, report)
		}
		if *rollback != "" {
			return suite.RollbackUpdate(ctx, *update, *rollback, report)
		}
		if !*apply && !*checkUpdates {
			return maintenance.Serve(ctx, *update, !*noBrowser)
		}
		var selected []string
		if *apps != "" {
			for _, id := range strings.Split(*apps, ",") {
				selected = append(selected, strings.TrimSpace(id))
			}
		}
		if *apply && *bundlePath != "" {
			return suite.UpdateFromBundle(ctx, *update, *bundlePath, *backupRoot, selected, report)
		}
		installed, e := suite.LoadUpdateConfig(*update)
		if e != nil {
			return e
		}
		plan, e := distribution.CheckUpdates(ctx, installed.Catalog)
		if e != nil {
			return e
		}
		if !*apply {
			return json.NewEncoder(os.Stdout).Encode(plan)
		}
		if len(selected) == 0 {
			for _, ch := range plan.Choices {
				if ch.Available {
					selected = append(selected, ch.ID)
				}
			}
		}
		if len(selected) == 0 {
			fmt.Println("Workspace and compatible app releases are up to date")
			return nil
		}
		return suite.ApplyUpdates(ctx, *update, plan, suite.UpdateRequest{Selected: selected, BackupRoot: *backupRoot}, report)
	}
	if *apply || *recoverUpdate || *rollback != "" || *bundlePath != "" || *backupRoot != "" || *apps != "" {
		return fmt.Errorf("update options require --update PRIVATE_DATA_FOLDER")
	}
	z, c, e := payload.Open()
	if e != nil {
		return e
	}
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
