package gamestudio

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io/fs"
)

//go:embed web/dist/* game.html arcade.js
var assets embed.FS

func UI() fs.FS {
	ui, _ := fs.Sub(assets, "web/dist")
	return ui
}

type Spec struct {
	Title    string `json:"title"`
	Template string `json:"template"`
	Hero     string `json:"hero"`
	Item     string `json:"item"`
	Theme    string `json:"theme"`
	Goal     int    `json:"goal"`
	Speed    int    `json:"speed"`
}

func (s Spec) valid() bool {
	return safeText(s.Title, 65) &&
		oneOf(s.Template, "catch", "maze", "runner", "memory") &&
		oneOf(s.Hero, "rocket", "cat", "robot", "frog", "owl", "dragon") &&
		oneOf(s.Item, "star", "gem", "apple", "moon", "cookie", "flower") &&
		oneOf(s.Theme, "space", "forest", "candy", "ocean", "sunset") &&
		s.Goal >= 4 && s.Goal <= 20 && s.Speed >= 1 && s.Speed <= 5
}

func oneOf(s string, options ...string) bool {
	for _, option := range options {
		if s == option {
			return true
		}
	}
	return false
}

var gamePage = template.Must(template.ParseFS(assets, "game.html"))
var engine = func() []byte {
	code, err := assets.ReadFile("arcade.js")
	if err != nil {
		panic(err)
	}
	// HTML parsing normalizes newlines before the browser checks inline hashes.
	code = bytes.ReplaceAll(code, []byte("\r\n"), []byte("\n"))
	return bytes.ReplaceAll(code, []byte("\r"), []byte("\n"))
}()
var engineHash = sha256.Sum256(engine)
var gamePolicy = "default-src 'none'; script-src 'sha256-" + base64.StdEncoding.EncodeToString(engineHash[:]) +
	"'; style-src 'unsafe-inline'; form-action 'none'; base-uri 'none'"

func render(s Spec) ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	err = gamePage.Execute(&b, struct {
		Title, Spec, Policy string
		Engine              template.JS
	}{s.Title, base64.StdEncoding.EncodeToString(data), gamePolicy, template.JS(engine)})
	return b.Bytes(), err
}

func archive(page []byte) ([]byte, error) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	f, err := z.Create("index.html")
	if err != nil {
		return nil, err
	}
	if _, err = f.Write(page); err != nil {
		return nil, err
	}
	err = z.Close()
	return b.Bytes(), err
}
