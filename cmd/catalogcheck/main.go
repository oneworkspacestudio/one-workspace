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
			fmt.Printf("Validated complete catalogue: %d apps\n", len(c.Apps))
			return
		}
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
