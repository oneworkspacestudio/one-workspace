package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fahimc/one-workspace/internal/platform"
	"github.com/fahimc/one-workspace/internal/suite"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

var version = "development"

func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run() error {
	data := flag.String("data", "", "Protected Workspace data folder")
	flag.Parse()
	args := flag.Args()
	action := "run"
	if len(args) > 0 {
		action = args[0]
	}
	if action == "version" {
		fmt.Println(version)
		return nil
	}
	if *data == "" {
		exe, _ := os.Executable()
		var p struct {
			Data string `json:"data"`
		}
		if e := platform.LoadJSON(filepath.Join(filepath.Dir(exe), "workspace-location.json"), &p); e == nil {
			*data = p.Data
		} else {
			*data = filepath.Dir(exe)
		}
	}
	c, s, e := suite.Load(*data)
	if e != nil {
		return fmt.Errorf("load Workspace: %w", e)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	switch action {
	case "run":
		if len(args) == 0 {
			go func() {
				if suite.WaitHTTP(ctx, fmt.Sprintf("http://localhost:%d/healthz", c.Port), 10*time.Minute) == nil {
					platform.OpenBrowser(c.Origin)
				}
			}()
		}
		return suite.Run(ctx, c, s)
	case "open":
		platform.OpenBrowser(c.Origin)
		return nil
	case "stop":
		return suite.Stop(ctx, c)
	case "status":
		client := &http.Client{Timeout: 5 * time.Second}
		r, e := client.Get(fmt.Sprintf("http://localhost:%d/healthz", c.Port))
		if e != nil {
			return e
		}
		defer r.Body.Close()
		var v any
		if e = json.NewDecoder(r.Body).Decode(&v); e != nil {
			return e
		}
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
		return nil
	case "apps":
		for _, a := range c.Catalog.Apps {
			fmt.Printf("%s\t%s\t%s\n", a.ID, a.Version, a.Route)
		}
		return nil
	case "add-app":
		f := flag.NewFlagSet("add-app", flag.ContinueOnError)
		manifest := f.String("manifest", "", "App manifest JSON")
		binary := f.String("package", "", "Trusted Windows app EXE")
		if e = f.Parse(args[1:]); e != nil {
			return e
		}
		return suite.AddApp(ctx, *data, *manifest, *binary)
	case "install-apps":
		f := flag.NewFlagSet("install-apps", flag.ContinueOnError)
		bundle := f.String("bundle", "", "Verified organization Workspace bundle ZIP")
		if e = f.Parse(args[1:]); e != nil {
			return e
		}
		if *bundle == "" {
			return fmt.Errorf("provide --bundle with the organization suite ZIP")
		}
		return suite.InstallApps(ctx, *data, *bundle, func(message string) { fmt.Println(message) })
	case "recovery":
		fmt.Printf("MyAccounts server recovery: %s, organization built-in\nPassword: %s\nDrive recovery: workspace-recovery\nPassword: %s\n", s.AdminEmail, s.AdminPassword, s.RecoveryPassword)
		return nil
	default:
		return fmt.Errorf("unknown command %s; use run, open, stop, status, apps, add-app or recovery", action)
	}
}
