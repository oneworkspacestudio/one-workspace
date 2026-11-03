package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"path/filepath"
	"strings"
	"testing"
)

func TestBooksNativeIdentityAndStorage(t *testing.T) {
	a := catalog.App{ID: "books", Kind: "process", Executable: "MyBooks.exe", Port: 9029, Route: "/books", OIDC: true, Env: map[string]string{"MYBOOKS_LIBRARY_ROOT": "${FILES}/My Books"}}
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Data", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{a}})
	a, ok := c.Catalog.Find("books")
	if !ok || a.Port != 29039 || !foundationApp("books") || !knownUpdateStorageProfile("books") {
		t.Fatal("Books integration missing", a)
	}
	if appWorkingDir(c, "books") != filepath.Join(c.Data, "books") || len(appUpdatePrivatePaths("books")) != 1 {
		t.Fatal("Books private recovery root missing")
	}
	rt := Runtime{Config: c, Secrets: Secrets{}}
	cmd := rt.command(a)
	found := false
	for _, e := range cmd.Env {
		if e == "MYBOOKS_LIBRARY_ROOT="+c.Files+"/My Books" {
			found = true
		}
	}
	if !found {
		t.Fatal("Books originals are not routed to Files")
	}
	if !strings.Contains(cmd.Dir, "books") {
		t.Fatal("mutable work placed in program folder")
	}
	signed := false
	for _, s := range suiteSignIns {
		if s.id == "books" {
			signed = s.login == "/books/auth/login" && s.logout == "/books/auth/logout" && s.cookies[0] == "books_session"
		}
	}
	if !signed {
		t.Fatal("Books cannot share sign-in/sign-out")
	}
}
