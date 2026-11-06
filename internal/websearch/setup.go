package websearch

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"os"
	"path/filepath"
)

// InstallModelCLI is the supported local owner setup path, including an offline
// import of the exact pinned model. It changes no accounts, sessions or grants.
func InstallModelCLI(args []string) error {
	f := flag.NewFlagSet("install-model", flag.ContinueOnError)
	data := f.String("data", "", "Existing Web Search private data directory")
	source := f.String("source", "", "Optional directory with the exact pinned GGUF files")
	if e := f.Parse(args); e != nil {
		return e
	}
	if *data == "" || len(f.Args()) != 0 {
		return errors.New("specify --data for the Web Search installation")
	}
	root, e := filepath.Abs(*data)
	if e != nil {
		return e
	}
	if _, e = os.Stat(filepath.Join(root, "app.json")); e != nil {
		return errors.New("start Web Search once before installing its model")
	}
	if e = apphost.PrepareData(root, "websearch"); e != nil {
		return e
	}
	target := filepath.Join(root, "models")
	if *source != "" {
		input, e := filepath.Abs(*source)
		if e != nil {
			return e
		}
		if e = verifyModels(input); e != nil {
			return e
		}
		if e = os.MkdirAll(target, 0700); e != nil {
			return e
		}
		for _, pin := range modelPins {
			from, e := os.Open(filepath.Join(input, pin.name))
			if e != nil {
				return e
			}
			tmp, e := os.CreateTemp(target, ".import-")
			if e != nil {
				from.Close()
				return e
			}
			_, e = copyBulk(tmp, from)
			from.Close()
			closeErr := tmp.Close()
			if e != nil || closeErr != nil {
				os.Remove(tmp.Name())
				return errors.New("model import failed")
			}
			dest := filepath.Join(target, pin.name)
			if e = os.Remove(dest); e != nil && !os.IsNotExist(e) {
				os.Remove(tmp.Name())
				return e
			}
			if e = os.Rename(tmp.Name(), dest); e != nil {
				return e
			}
		}
		e = verifyModels(target)
	} else {
		e = installModels(context.Background(), target, func(int64) {})
	}
	if e != nil {
		return e
	}
	_, e = io.WriteString(os.Stdout, fmt.Sprintf("Verified Web Search model installed (%d files).\n", len(modelPins)))
	return e
}
