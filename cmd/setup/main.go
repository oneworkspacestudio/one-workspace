package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
	noBrowser := flag.Bool("no-browser", false, "Print setup link without opening a browser")
	headless := flag.Bool("headless", false, "Read setup JSON from standard input")
	flag.Parse()
	z, c, e := payload.Open()
	if e != nil {
		return e
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if *headless {
		var r suite.Request
		if e = json.NewDecoder(os.Stdin).Decode(&r); e != nil {
			return e
		}
		conf, e := suite.Install(ctx, r, z, c, func(phase string) { fmt.Println(phase) })
		if e == nil {
			fmt.Println("Installed:", conf.Origin, "; launch", conf.Data+"/Workspace.exe")
		}
		return e
	}
	return setupui.Serve(ctx, z, c, !*noBrowser)
}
