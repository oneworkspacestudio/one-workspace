package distribution

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/operationprogress"
)

const mailRepository = "oneworkspacestudio/mail"
const mailManifestName = "workspace-mail.json"

type mailManifest struct {
	Schema      int    `json:"schema"`
	Integration int    `json:"integration"`
	Version     string `json:"version"`
	Asset       string `json:"asset"`
}

func latestMail(ctx context.Context, client *http.Client, token, dir string, cat *catalog.Catalog, progress func(string)) (map[string]string, error) {
	app, ok := cat.Find("mail")
	if !ok {
		return nil, nil
	}
	if app.Kind != "process" || app.Executable != "Mail.exe" || cat.Packages[app.Package].File != "Mail.exe" {
		return nil, errors.New("unsupported Mail catalog integration")
	}
	progress("Checking oneworkspacestudio/mail for the latest Mail release")
	metaCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := request(metaCtx, "https://api.github.com/repos/"+mailRepository+"/releases?per_page=100", token)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Mail release check failed (HTTP %d)", resp.StatusCode)
	}
	var releases []Release
	if err = json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&releases); err != nil {
		return nil, err
	}
	release, exe, marker, err := chooseMail(releases)
	if err != nil {
		return nil, err
	}
	if !validNamedAsset(exe, mailRepository, release.Tag, "Mail.exe", 256<<20) || !validNamedAsset(marker, mailRepository, release.Tag, mailManifestName, 32<<10) {
		return nil, errors.New("invalid organization Mail asset URL, size or checksum")
	}
	markerPath := filepath.Join(dir, mailManifestName)
	if err = downloadMailAsset(ctx, client, token, marker, markerPath); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(markerPath)
	if err != nil {
		return nil, err
	}
	manifest, err := readMailManifest(data, release.Tag)
	if err != nil {
		return nil, err
	}
	overrides := map[string]string{}
	current := cat.Packages[app.Package]
	if current.SHA256 != strings.TrimPrefix(exe.Digest, "sha256:") || current.Size != exe.Size {
		progress("Downloading Mail " + manifest.Version + " from oneworkspacestudio/mail")
		path := filepath.Join(dir, "Mail.exe")
		if err = downloadMailAsset(ctx, client, token, exe, path); err != nil {
			return nil, err
		}
		overrides[current.File] = path
	}
	cat.Packages[app.Package] = catalog.Package{File: current.File, Size: exe.Size, SHA256: strings.TrimPrefix(exe.Digest, "sha256:")}
	for i := range cat.Apps {
		if cat.Apps[i].ID == "mail" {
			cat.Apps[i].Version = manifest.Version
			cat.Apps[i].Source = "https://github.com/" + mailRepository + "/releases/tag/" + release.Tag
		}
	}
	if err = cat.Validate(); err != nil {
		return nil, err
	}
	progress("Verified Mail " + manifest.Version + " from oneworkspacestudio/mail")
	return overrides, nil
}

func chooseMail(releases []Release) (Release, Asset, Asset, error) {
	sort.SliceStable(releases, func(i, j int) bool { return releases[i].Published.After(releases[j].Published) })
	for _, r := range releases {
		if r.Draft || r.Published.IsZero() || !strings.HasPrefix(r.Tag, "workspace-mail-v") {
			continue
		}
		var exe, marker Asset
		for _, a := range r.Assets {
			if a.Name == "Mail.exe" {
				exe = a
			}
			if a.Name == mailManifestName {
				marker = a
			}
		}
		if exe.Name != "" {
			if marker.Name == "" {
				return Release{}, Asset{}, Asset{}, errors.New("newest Mail release is missing its Workspace compatibility manifest")
			}
			return r, exe, marker, nil
		}
	}
	return Release{}, Asset{}, Asset{}, errors.New("no published organization Mail release is available")
}

func readMailManifest(data []byte, tag string) (mailManifest, error) {
	var m mailManifest
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&m); err != nil {
		return m, err
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF || m.Schema != 1 || m.Integration != 1 || m.Asset != "Mail.exe" || tag != "workspace-mail-v"+m.Version || m.Version == "" {
		return m, errors.New("unsupported or mismatched Mail compatibility manifest")
	}
	return m, nil
}

func downloadMailAsset(ctx context.Context, client *http.Client, token string, asset Asset, path string) error {
	meter := operationprogress.Bytes(ctx, "download", "Downloading "+asset.Name, asset.Size)
	req, err := request(ctx, asset.APIURL, token)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Release download failed (HTTP %d)", resp.StatusCode)
	}
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	hash := sha256.New()
	n, copyErr := io.Copy(io.MultiWriter(out, hash, meter), io.LimitReader(resp.Body, asset.Size+1))
	meter.Finish()
	closeErr := out.Close()
	if copyErr != nil || closeErr != nil || n != asset.Size || "sha256:"+hex.EncodeToString(hash.Sum(nil)) != asset.Digest {
		return errors.New("Release download incomplete or SHA256 verification failed")
	}
	return nil
}
