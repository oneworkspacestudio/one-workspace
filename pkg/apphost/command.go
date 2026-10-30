package apphost

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

type Application interface{ Run(context.Context) error }

func Main(config Config, ui fs.FS, build func(*Host) (Application, error)) error {
	f := flag.NewFlagSet(config.ID, flag.ContinueOnError)
	data := f.String("data", "", "Dedicated NTFS app data folder")
	listen := f.String("listen", config.Listen, "Loopback listener")
	origin := f.String("origin", "", "Public app URL, including its proxy subpath")
	noBrowser := f.Bool("no-browser", false, "Do not open the browser")
	bootstrap := f.String("bootstrap-email", "", "Standalone first owner; password comes from WORKSPACE_BOOTSTRAP_PASSWORD")
	backup := f.String("backup", "", "Export encrypted backup to this new file")
	restore := f.String("restore", "", "Restore encrypted backup into a new data folder")
	version := f.Bool("version", false, "Print version")
	if e := f.Parse(os.Args[1:]); e != nil {
		return e
	}
	if len(f.Args()) != 0 {
		return errors.New("unexpected command argument; check quoted paths and flags")
	}
	if *version {
		fmt.Println(config.Version)
		return nil
	}
	if *data == "" {
		exe, e := os.Executable()
		if e != nil {
			return e
		}
		*data = filepath.Join(filepath.Dir(exe), "data")
	}
	abs, e := filepath.Abs(*data)
	if e != nil {
		return e
	}
	config.Data = abs
	config.Listen = *listen
	if *origin == "" {
		*origin = "http://" + strings.Replace(*listen, "127.0.0.1", "localhost", 1)
	}
	config.Origin = *origin
	config = FromEnv(config)
	if *restore != "" {
		if *backup != "" {
			return errors.New("choose backup or restore")
		}
		return RestoreWithDatabaseFiles(config.Data, config.ID, *restore, os.Getenv("WORKSPACE_BACKUP_PASSWORD"), config.RestoreFiles, config.RestoreDatabases)
	}
	h, e := New(config, ui)
	if e != nil {
		return e
	}
	defer h.Close()
	if *bootstrap != "" {
		if e = h.Bootstrap(*bootstrap, os.Getenv("WORKSPACE_BOOTSTRAP_PASSWORD")); e != nil {
			return e
		}
		os.Unsetenv("WORKSPACE_BOOTSTRAP_PASSWORD")
	}
	app, e := build(h)
	if e != nil {
		return e
	}
	h.BackupRoute()
	if *backup != "" {
		out, e := os.OpenFile(*backup, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		e = h.Snapshot(out, os.Getenv("WORKSPACE_BACKUP_PASSWORD"))
		closeErr := out.Close()
		if e != nil {
			os.Remove(*backup)
			return e
		}
		return closeErr
	}
	link := h.Config.Origin + "/"
	if h.SetupURL() != "" {
		link = h.SetupURL()
	}
	fmt.Println(h.Config.Name + ": " + link)
	parent := config.ParentContext
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := signal.NotifyContext(parent, os.Interrupt)
	defer cancel()
	if !*noBrowser {
		platform.OpenBrowser(link)
	}
	workerDone := make(chan error, 1)
	if app != nil {
		go func() {
			e := app.Run(ctx)
			if e != nil {
				h.Log.Error("application worker stopped", "kind", fmt.Sprintf("%T", e))
				cancel()
			}
			workerDone <- e
		}()
	} else {
		workerDone <- nil
	}
	serverErr := h.Run(ctx)
	cancel()
	workerErr := <-workerDone
	if serverErr != nil {
		return serverErr
	}
	return workerErr
}
