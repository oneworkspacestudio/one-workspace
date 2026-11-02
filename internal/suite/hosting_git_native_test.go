package suite

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Uses a real verified native Gitea binary with a fresh synthetic database.
// No installed account, configuration, repository or credential is opened.
func TestHostingGitNativePrivateRepository(t *testing.T) {
	bin := os.Getenv("HOSTING_ACCEPTANCE_GITEA")
	if bin == "" {
		t.Skip("Set HOSTING_ACCEPTANCE_GITEA to run native My Git acceptance")
	}
	root := t.TempDir()
	gitroot := filepath.Join(root, "gitea")
	_ = os.MkdirAll(gitroot, 0700)
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	config := filepath.Join(root, "app.ini")
	ini := fmt.Sprintf("APP_NAME = Hosting Synthetic Acceptance\nRUN_MODE = prod\nWORK_PATH = %s\n[database]\nDB_TYPE = sqlite3\nPATH = %s\n[repository]\nROOT = %s\n[server]\nHTTP_ADDR = 127.0.0.1\nHTTP_PORT = %d\nROOT_URL = %s/\nAPP_DATA_PATH = %s\nDISABLE_SSH = true\nOFFLINE_MODE = true\n[security]\nINSTALL_LOCK = true\nSECRET_KEY = synthetic-gitea-acceptance-key-never-production\n[service]\nDISABLE_REGISTRATION = true\n[log]\nMODE = file\nROOT_PATH = %s\nLEVEL = Warn\n", filepath.ToSlash(root), filepath.ToSlash(filepath.Join(gitroot, "gitea.db")), filepath.ToSlash(filepath.Join(root, "repos")), port, base, filepath.ToSlash(filepath.Join(root, "runtime")), filepath.ToSlash(filepath.Join(root, "logs")))
	if e = os.WriteFile(config, []byte(ini), 0600); e != nil {
		t.Fatal(e)
	}
	command := func(args ...string) {
		t.Helper()
		args = append([]string{"--work-path", root, "--config", config}, args...)
		c := exec.Command(bin, args...)
		if b, e := c.CombinedOutput(); e != nil {
			t.Fatalf("native My Git fixture command failed: %v %.500s", e, b)
		}
	}
	command("migrate")
	password := "SyntheticFixtureOnly8!"
	command("admin", "user", "create", "--username", "workspace-recovery", "--password", password, "--email", "recovery@example.test", "--admin", "--must-change-password=false")
	log, e := os.Create(filepath.Join(root, "server.log"))
	if e != nil {
		t.Fatal(e)
	}
	defer log.Close()
	cmd := exec.Command(bin, "--work-path", root, "--config", config, "web")
	cmd.Stdout = log
	cmd.Stderr = log
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
	client := &http.Client{Timeout: 15 * time.Second}
	ready := false
	for until := time.Now().Add(30 * time.Second); time.Now().Before(until); {
		r, e := client.Get(base + "/api/healthz")
		if e == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				ready = true
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !ready {
		t.Fatal("native Git fixture not ready")
	}
	api := func(method, path, sudo string, body any, out any) int {
		t.Helper()
		var data []byte
		if body != nil {
			data, _ = json.Marshal(body)
		}
		r, _ := http.NewRequest(method, base+"/api/v1"+path, bytes.NewReader(data))
		r.SetBasicAuth("workspace-recovery", password)
		r.Header.Set("Content-Type", "application/json")
		if sudo != "" {
			r.Header.Set("Sudo", sudo)
		}
		res, e := client.Do(r)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		if out != nil {
			if e = json.NewDecoder(res.Body).Decode(out); e != nil {
				t.Fatal(e)
			}
		}
		return res.StatusCode
	}
	for _, name := range []string{"alice-fixture", "bob-fixture"} {
		status := api("POST", "/admin/users", "", map[string]any{"username": name, "password": password, "email": name + "@example.test", "must_change_password": false}, nil)
		if status != 201 {
			t.Fatal("create fixture Git user", status)
		}
	}
	var repo hostingGitRepo
	status := api("POST", "/user/repos", "alice-fixture", map[string]any{"name": "fidget-clicker-creator", "private": true, "auto_init": true, "default_branch": "trunk"}, &repo)
	if status != 201 {
		t.Fatal("create private fixture repository", status)
	}
	status = api("POST", "/repos/alice-fixture/fidget-clicker-creator/contents/index.html", "alice-fixture", map[string]any{"content": base64.StdEncoding.EncodeToString([]byte("<!doctype html><title>Native My Git Hosting acceptance</title>")), "message": "Synthetic static site", "branch": "trunk"}, nil)
	if status != 201 {
		t.Fatal("commit fixture site", status)
	}
	// Bind the synthetic account by the same persisted OAuth subject fields used
	// by MyAccounts, independently of email/user-provided Git usernames.
	db, e := sql.Open("sqlite", filepath.Join(gitroot, "gitea.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	issuer := "https://accounts.example.test"
	clientID := "synthetic-git-client"
	clientSecret := "synthetic-client-secret"
	cfg, _ := json.Marshal(map[string]string{"ClientID": clientID, "ClientSecret": clientSecret, "OpenIDConnectAutoDiscoveryURL": issuer + "/.well-known/openid-configuration"})
	result, e := db.Exec(`INSERT INTO login_source(type,name,is_active,cfg,created_unix,updated_unix) VALUES(6,'myaccounts',1,?,0,0)`, string(cfg))
	if e != nil {
		t.Fatal(e)
	}
	source, e := result.LastInsertId()
	if e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`UPDATE user SET login_source=?,login_type=6,login_name='subject-alice' WHERE name='alice-fixture'`, source); e != nil {
		t.Fatal(e)
	}
	blob, e := platform.Protect([]byte(password))
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(gitroot, "recovery.dpapi"), blob, 0600); e != nil {
		t.Fatal(e)
	}
	rt := Runtime{Config: Config{Data: root, Issuer: issuer, Origin: "https://workspace.example.test", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "gitea", Port: port}}}}, Secrets: Secrets{Clients: map[string]Client{"hosting": {Secret: "synthetic-hosting-capability"}, "gitea": {ID: clientID, Secret: clientSecret}}}}
	portal := Portal{rt: &rt}
	broker := func(op, repository, ref, file string) *httptest.ResponseRecorder {
		t.Helper()
		in := hostingGitCall{Operation: op, Issuer: issuer, Subject: "subject-alice", Repository: repository, Ref: ref, File: file}
		if op != "resolve" {
			in.RepositoryID = repo.ID
		}
		b, _ := json.Marshal(in)
		r := httptest.NewRequest("POST", "http://localhost/_integration/hosting/git", bytes.NewReader(b))
		r.RemoteAddr = "127.0.0.1:1234"
		r.Header.Set("Authorization", "Bearer "+hostingGitCapability("synthetic-hosting-capability"))
		w := httptest.NewRecorder()
		portal.hostingGitRequest(w, r)
		return w
	}
	for _, input := range []string{"fidget-clicker-creator", rt.Config.Origin + "/git/alice-fixture/fidget-clicker-creator", rt.Config.Origin + "/git/alice-fixture/fidget-clicker-creator.git"} {
		w := broker("resolve", input, "", "")
		if w.Code != 200 {
			t.Fatal(input, w.Code, w.Body.String())
		}
	}
	w := broker("commit", "alice-fixture/fidget-clicker-creator", "trunk", "")
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var commit struct {
		SHA string `json:"sha"`
	}
	if json.Unmarshal(w.Body.Bytes(), &commit) != nil || !hostingGitSHA.MatchString(commit.SHA) {
		t.Fatal("native commit shape", w.Body.String())
	}
	w = broker("contents", "alice-fixture/fidget-clicker-creator", "trunk", "index.html")
	if w.Code != 200 {
		t.Fatal("native contents", w.Code, w.Body.String())
	}
	w = broker("archive", "alice-fixture/fidget-clicker-creator", commit.SHA, "")
	if w.Code != 200 {
		t.Fatal("native archive", w.Code, w.Body.String())
	}
	gz, e := gzip.NewReader(w.Body)
	if e != nil {
		t.Fatal(e)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	found := false
	for {
		h, e := tr.Next()
		if e == io.EOF {
			break
		}
		if e != nil {
			t.Fatal(e)
		}
		if strings.HasSuffix(h.Name, "/index.html") {
			b, _ := io.ReadAll(tr)
			found = strings.Contains(string(b), "Native My Git Hosting acceptance")
		}
	}
	if !found {
		t.Fatal("native archive site missing")
	}
	if goBin := os.Getenv("HOSTING_ACCEPTANCE_GO"); goBin != "" {
		live := httptest.NewServer(http.HandlerFunc(portal.hostingGitRequest))
		defer live.Close()
		child := exec.Command(goBin, "test", "-v", "-buildvcs=false", "./internal/hosting", "-run", "^TestMyGitNativeBridgeDeployment$")
		child.Dir = os.Getenv("HOSTING_ACCEPTANCE_APP_REPO")
		child.Env = append(os.Environ(), "HOSTING_ACCEPTANCE_BRIDGE="+live.URL+"/_integration/hosting/git", "HOSTING_ACCEPTANCE_CAPABILITY="+hostingGitCapability("synthetic-hosting-capability"))
		b, err := child.CombinedOutput()
		if err != nil {
			t.Fatalf("native Hosting deployment failed: %v %s", err, b)
		}
		t.Log(string(b))
		files := map[string]string{
			"package.json":      `{"name":"synthetic-site","version":"1.0.0","scripts":{"build":"node build.cjs"}}`,
			"package-lock.json": `{"name":"synthetic-site","version":"1.0.0","lockfileVersion":3,"packages":{"":{"name":"synthetic-site","version":"1.0.0"}}}`,
			"build.cjs":         `const fs=require('node:fs');fs.mkdirSync('dist',{recursive:true});fs.writeFileSync('dist/index.html','<!doctype html><title>Native My Git Hosting acceptance</title>');`,
		}
		for file, content := range files {
			if status := api("POST", "/repos/alice-fixture/fidget-clicker-creator/contents/"+file, "alice-fixture", map[string]any{"content": base64.StdEncoding.EncodeToString([]byte(content)), "message": "Synthetic Node build fixture", "branch": "trunk"}, nil); status != 201 {
				t.Fatal("create Node fixture", file, status)
			}
		}
		child = exec.Command(goBin, "test", "-v", "-buildvcs=false", "./internal/hosting", "-run", "^TestMyGitNativeBridgeDeployment$")
		child.Dir = os.Getenv("HOSTING_ACCEPTANCE_APP_REPO")
		child.Env = append(os.Environ(), "HOSTING_ACCEPTANCE_BRIDGE="+live.URL+"/_integration/hosting/git", "HOSTING_ACCEPTANCE_CAPABILITY="+hostingGitCapability("synthetic-hosting-capability"), "HOSTING_ACCEPTANCE_NODE=1")
		b, err = child.CombinedOutput()
		if err != nil {
			t.Fatalf("native Git Node build failed: %v %s", err, b)
		}
		t.Log(string(b))

	}

	// Admin's global access must not leak through impersonation.
	var other hostingGitRepo
	if api("POST", "/user/repos", "bob-fixture", map[string]any{"name": "secret", "private": true, "auto_init": true}, &other) != 201 {
		t.Fatal("create other-user repository")
	}
	w = broker("resolve", "bob-fixture/secret", "", "")
	if w.Code != 404 {
		t.Fatal("cross-user private repository exposed", w.Code, w.Body.String())
	}
	// Removing the OAuth link makes later background fetches fail closed.
	_, _ = db.Exec(`UPDATE user SET prohibit_login=1 WHERE name='alice-fixture'`)
	w = broker("resolve", "fidget-clicker-creator", "", "")
	if w.Code != 403 {
		t.Fatal("blocked user still reads", w.Code)
	}
}
