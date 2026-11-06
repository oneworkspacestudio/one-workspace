package websearch

import (
	"embed"
	"io/fs"
)

//go:embed web/dist/*
var ui embed.FS

func UI() fs.FS { root, _ := fs.Sub(ui, "web/dist"); return root }
