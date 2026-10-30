package main

import (
	"fmt"
	"os"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: catalogcheck catalog.json")
		os.Exit(1)
	}
	b, err := os.ReadFile(os.Args[1])
	if err == nil {
		var c catalog.Catalog
		c, err = catalog.Read(b)
		if err == nil {
			if err = c.ValidateManagedInventory(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			for _, app := range c.Apps {
				if !catalog.DefaultInstall(app.ID) {
					continue
				}
				for _, dep := range app.Depends {
					if !catalog.DefaultInstall(dep) {
						fmt.Fprintf(os.Stderr, "default app %s depends on optional app %s\n", app.ID, dep)
						os.Exit(1)
					}
				}
			}
			fmt.Printf("Validated complete catalogue: %d apps\n", len(c.Apps))
			return
		}
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
