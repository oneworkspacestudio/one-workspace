package suite

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMyAIProjectsKeepRegistryAndGitOutsideDrive(t *testing.T) {
	root := t.TempDir()
	c := Config{Data: filepath.Join(root, "Private"), Files: filepath.Join(root, "Files")}
	rt := Runtime{Config: c}
	uid := strings.Repeat("a", 64)
	id := strings.Repeat("b", 48)
	name := "Family recipes"
	folder := myaiProjectSlug(name) + "-" + id[:8]
	p := myaiProject{ID: id, Name: name, Folder: folder, Managed: true, Path: filepath.Join(myaiProjectsRoot(c), "users", uid, folder), GitDir: filepath.Join(c.Data, "myai-projects", "git", uid, id)}
	if containedPath(myaiProjectsRoot(c), p.GitDir) || containedPath(myaiProjectsRoot(c), myaiProjectsMeta(c, uid)) {
		t.Fatal("metadata exposed in Drive")
	}
	if e := rt.myaiProjectSave(uid, p); e != nil {
		t.Fatal(e)
	}
	if _, e := rt.myaiProjectLoad(uid, id); e != nil {
		t.Fatal(e)
	}
	if _, e := rt.myaiProjectLoad(strings.Repeat("c", 64), id); e == nil {
		t.Fatal("cross-account registration read")
	}
	p.Path = c.Files
	if e := rt.myaiProjectSave(uid, p); e != nil {
		t.Fatal(e)
	}
	if _, e := rt.myaiProjectLoad(uid, id); e == nil {
		t.Fatal("registry accepted arbitrary host root")
	}
}
func TestMyAIProjectsNeedsPrivateMarkerAndRefusesLinks(t *testing.T) {
	root := t.TempDir()
	rt := Runtime{Config: Config{Data: filepath.Join(root, "Private"), Files: filepath.Join(root, "Files")}}
	if rt.verifyMyAIProjectsStorage(false) == nil {
		t.Fatal("uninitialised storage accepted")
	}
	if e := rt.verifyMyAIProjectsStorage(true); e != nil {
		t.Fatal(e)
	}
	if e := rt.verifyMyAIProjectsStorage(false); e != nil {
		t.Fatal(e)
	}
	other := t.TempDir()
	link := filepath.Join(myaiProjectsRoot(rt.Config), "linked")
	if e := os.Symlink(other, link); e == nil {
		if myaiSafeTreePath(filepath.Join(link, "secret")) == nil {
			t.Fatal("linked path accepted")
		}
	}
}
func TestMyAIProjectBrokerDoesNotExposeLoopbackAuthority(t *testing.T) {
	rt := Runtime{Config: Config{Issuer: "https://accounts.example.test"}, Secrets: Secrets{Clients: map[string]Client{"myai": {Secret: "synthetic-capability"}}}}
	p := Portal{rt: &rt}
	for _, address := range []string{"203.0.113.3:4000", "127.0.0.1:4000"} {
		req := httptest.NewRequest("GET", "http://localhost/_integration/myai/projects", nil)
		req.RemoteAddr = address
		req.Header.Set("X-MyAI-Issuer", rt.Config.Issuer)
		req.Header.Set("X-MyAI-Subject", "other-person")
		w := httptest.NewRecorder()
		p.myaiProjectsRequest(w, req)
		if w.Code != http.StatusNotFound || strings.Contains(w.Body.String(), "synthetic-capability") {
			t.Fatal("project broker authority exposed")
		}
	}
}

func TestMyAIProjectsCannotShareWorkingFolder(t *testing.T) {
	root := t.TempDir()
	c := Config{Data: filepath.Join(root, "Private"), Files: filepath.Join(root, "Files")}
	rt := Runtime{Config: c}
	uid := strings.Repeat("a", 64)
	id := strings.Repeat("b", 48)
	folder := myaiProjectSlug("Recipes") + "-" + id[:8]
	p := myaiProject{ID: id, Name: "Recipes", Managed: true, Folder: folder, Repository: folder, Path: filepath.Join(myaiProjectsRoot(c), "users", uid, folder), GitDir: filepath.Join(c.Data, "myai-projects", "git", uid, id)}
	if e := rt.myaiProjectDestinationAvailable(uid, p); e != nil {
		t.Fatal(e)
	}
	if e := rt.myaiProjectSave(uid, p); e != nil {
		t.Fatal(e)
	}
	second := p
	second.ID = strings.Repeat("b", 8) + strings.Repeat("c", 40)
	if rt.myaiProjectDestinationAvailable(uid, second) == nil {
		t.Fatal("two project IDs shared a working folder")
	}
}
