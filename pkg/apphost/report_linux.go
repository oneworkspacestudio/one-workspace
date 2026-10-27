//go:build linux

package apphost

import (
	"fmt"
	"os"
)

func Report(e error) { fmt.Fprintln(os.Stderr, e) }
