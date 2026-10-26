package suite

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var aiCommit = regexp.MustCompile(`^[a-f0-9]{40}$`)

func aiProjectPublic(p myaiProject) map[string]any {
	return map[string]any{"id": p.ID, "name": p.Name, "repository": p.URL, "drive": p.DriveURL, "commit": p.Commit, "state": p.State}
}

func (p *Portal) aiManagedOperation(ctx context.Context, browser session, uid string, op aiOperation, call aiAppCall) (any, error) {
	if op.Path == "managed-drive-write" {
		return p.aiDriveWrite(ctx, browser, call)
	}
	if op.Path == "managed-git-projects" {
		ps, e := p.rt.myaiProjectList(uid)
		if e != nil {
			return nil, e
		}
		out := []map[string]any{}
		for _, project := range ps {
			out = append(out, aiProjectPublic(project))
		}
		return out, nil
	}
	p.rt.myaiProjectsMu.Lock()
	defer p.rt.myaiProjectsMu.Unlock()
	project, e := p.rt.myaiProjectLoad(uid, call.Parameters["project"])
	if e != nil || !project.Managed {
		return nil, errors.New("Choose one of your own managed Git projects")
	}
	if op.Path == "managed-git-history" {
		history, e := p.rt.myaiGitRun(ctx, project, "", "log", "-10", "--format=%H%x09%aI%x09%s")
		return map[string]any{"project": project.ID, "history": history, "untrusted_reference": true}, e
	}
	g, e := p.rt.myaiGit(ctx)
	if e != nil {
		return nil, e
	}
	if op.Path == "managed-git-save" {
		if e = p.rt.myaiProjectSync(ctx, g, &project, false); e != nil {
			return nil, e
		}
		if e = p.rt.myaiProjectSave(uid, project); e != nil {
			return nil, e
		}
		return aiProjectPublic(project), nil
	}
	if e = g.repository(ctx, &project, false); e != nil {
		return nil, e
	}
	status, e := p.rt.myaiProjectStatus(ctx, project)
	if e != nil {
		return nil, e
	}
	head, e := p.rt.myaiGitRun(ctx, project, "", "rev-parse", "--verify", "HEAD")
	if e != nil {
		return nil, e
	}
	head = strings.TrimSpace(head)
	token, closeToken, e := g.token(ctx, project.GitUser)
	if e != nil {
		return nil, e
	}
	defer closeToken()
	remote := g.base + "/" + url.PathEscape(project.GitUser) + "/" + url.PathEscape(project.Repository) + ".git"
	refs, e := p.rt.myaiGitRun(ctx, project, token, "ls-remote", remote, "refs/heads/main")
	if e != nil {
		return nil, e
	}
	parts := strings.Fields(refs)
	remoteHead := ""
	if len(parts) > 0 {
		remoteHead = parts[0]
	}
	if op.Path == "managed-git-status" {
		return map[string]any{"project": project.ID, "name": project.Name, "repository": project.URL, "commit": head, "remoteCommit": remoteHead, "savedToGit": head == remoteHead, "pendingChanges": status}, nil
	}
	if op.Path != "managed-project-upload" {
		return nil, errors.New("Unsupported managed project operation")
	}
	var in struct {
		Commit      string `json:"commit"`
		Environment string `json:"environment"`
	}
	if json.Unmarshal(call.Body, &in) != nil || !aiCommit.MatchString(in.Commit) {
		return nil, errors.New("Read Git status and choose its exact saved commit before deploying")
	}
	if in.Environment == "" {
		in.Environment = "preview"
	}
	if in.Environment != "preview" && in.Environment != "production" {
		return nil, errors.New("Choose preview or production")
	}
	if status != "" || head != remoteHead || head != in.Commit {
		return nil, errors.New("The reviewed project revision has changed or is not pushed. Save it to Git, then review deployment again")
	}
	hosting, ok := p.rt.Config.Catalog.Find("hosting")
	if !ok {
		return nil, errors.New("My Hosting is not installed")
	}
	connection, csrf, e := p.aiNativeIdentity(ctx, browser, hosting)
	if e != nil {
		return nil, &aiConnectionRequired{e}
	}
	site := call.Parameters["site"]
	if !regexp.MustCompile(`^[a-f0-9]{32}$`).MatchString(site) {
		return nil, errors.New("Choose a site ID returned by My Hosting")
	}
	// Check native authorization before reading source or queuing work.
	if _, _, e = p.aiNativeRequest(ctx, hosting, connection, csrf, "GET", "/api/projects/"+site, nil, ""); e != nil {
		return nil, e
	}
	archive, e := p.rt.aiGitArchive(ctx, project, in.Commit)
	if e != nil {
		return nil, e
	}
	return p.aiSubmitDeployment(ctx, uid, hosting, connection, csrf, call, archive, in.Commit, in.Environment, project.URL)
}

// Build from immutable Git blobs, not mutable working files or untracked secrets.
// This deliberately excludes symlinks/submodules and bounds every object before
// loading it. Git's hooks, helpers and mutable repository config stay disabled.
func (rt *Runtime) aiGitArchive(ctx context.Context, p myaiProject, commit string) ([]byte, error) {
	if !aiCommit.MatchString(commit) {
		return nil, errors.New("Invalid saved commit")
	}
	tree, e := rt.myaiGitRun(ctx, p, "", "ls-tree", "-r", "-l", "-z", commit)
	if e != nil {
		return nil, e
	}
	entries := strings.Split(strings.TrimSuffix(tree, "\x00"), "\x00")
	if len(entries) > 2000 {
		return nil, errors.New("Source has too many files for a deployment handoff")
	}
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	var total int64
	for _, entry := range entries {
		fields := strings.SplitN(entry, "\t", 2)
		if len(fields) != 2 {
			return nil, errors.New("Git source listing is invalid")
		}
		metadata := strings.Fields(fields[0])
		name := fields[1]
		if len(metadata) != 4 || metadata[1] != "blob" || (metadata[0] != "100644" && metadata[0] != "100755") || !aiCommit.MatchString(metadata[2]) || !aiArchivePath(name) {
			return nil, errors.New("Source contains a linked, credential or unsupported file; review it in Git")
		}
		size, e := strconv.ParseInt(metadata[3], 10, 64)
		total += size
		if e != nil || size < 0 || size > 4<<20 || total > 32<<20 {
			return nil, errors.New("Source exceeds the bounded32MiB handoff or4MiB per-file limit")
		}
		content, e := rt.myaiGitRun(ctx, p, "", "cat-file", "blob", metadata[2])
		if e != nil || int64(len(content)) != size {
			return nil, errors.New("Git could not read the reviewed source")
		}
		header := &zip.FileHeader{Name: name, Method: zip.Deflate}
		header.SetModTime(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC))
		header.SetMode(0644)
		writer, e := zw.CreateHeader(header)
		if e != nil {
			return nil, e
		}
		if _, e = writer.Write([]byte(content)); e != nil {
			return nil, e
		}
	}
	if e = zw.Close(); e != nil {
		return nil, e
	}
	return out.Bytes(), nil
}
func aiArchivePath(name string) bool {
	return name != "" && !strings.HasPrefix(name, "/") && !strings.ContainsAny(name, "\\:\x00\r\n") && !strings.Contains("/"+name+"/", "/../") && !strings.Contains("/"+name+"/", "/./") && !myaiProjectForbidden(name)
}

func (p *Portal) aiDriveWrite(ctx context.Context, browser session, call aiAppCall) (any, error) {
	var in struct {
		Source    string `json:"source"`
		Path      string `json:"path"`
		Content   string `json:"content"`
		Overwrite bool   `json:"overwrite"`
	}
	decoder := json.NewDecoder(bytes.NewReader(call.Body))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&in) != nil || in.Source == "" || len(in.Source) > 100 || !aiArchivePath(strings.TrimPrefix(in.Path, "/")) || len(in.Content) > 256<<10 || !utf8.ValidString(in.Content) || strings.ContainsRune(in.Content, 0) {
		return nil, errors.New("Choose a regular text file within your Drive source, up to256KiB")
	}
	app, ok := p.rt.Config.Catalog.Find("drive")
	if !ok {
		return nil, errors.New("My Drive is not installed")
	}
	connection, csrf, e := p.aiNativeIdentity(ctx, browser, app)
	if e != nil {
		return nil, &aiConnectionRequired{e}
	}
	query := url.Values{"source": {in.Source}, "path": {in.Path}, "override": {strconv.FormatBool(in.Overwrite)}, "isDir": {"false"}}
	if _, _, e = p.aiNativeRequest(ctx, app, connection, csrf, "POST", "/api/resources?"+query.Encode(), []byte(in.Content), "application/octet-stream"); e != nil {
		return nil, e
	}
	return map[string]any{"source": in.Source, "path": in.Path, "saved": true, "bytes": len(in.Content)}, nil
}
