package payload

import (
	"archive/zip"
	"embed"
	"errors"
	"github.com/fahimc/one-workspace/internal/catalog"
	"io"
)

//go:embed bundle/*
var files embed.FS

func Open() (*zip.Reader, catalog.Catalog, error) {
	var c catalog.Catalog
	f, e := files.Open("bundle/payload.zip")
	if e != nil {
		return nil, c, errors.New("build this installer using scripts/build.ps1 to embed its offline payload")
	}
	// Read the embedded archive in place. Copying a 300 MB payload at startup
	// delays the wizard and doubles its working set before any app is installed.
	info, e := f.Stat()
	if e != nil {
		return nil, c, e
	}
	reader, ok := f.(io.ReaderAt)
	if !ok {
		return nil, c, errors.New("embedded bundle does not support random access")
	}
	z, e := zip.NewReader(reader, info.Size())
	if e != nil {
		return nil, c, e
	}
	for _, f := range z.File {
		if f.Name == "catalog.json" {
			if f.UncompressedSize64 > 1<<20 {
				return nil, c, errors.New("catalog too large")
			}
			r, e := f.Open()
			if e != nil {
				return nil, c, e
			}
			data, e := io.ReadAll(io.LimitReader(r, 1<<20))
			r.Close()
			if e != nil {
				return nil, c, e
			}
			c, e = catalog.Read(data)
			return z, c, e
		}
	}
	return nil, c, errors.New("bundle has no app catalog")
}
