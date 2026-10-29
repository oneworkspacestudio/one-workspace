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

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
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
	releases, err := releasesFor(metaCtx, client, token, mailRepository)
	if err != nil {
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
	out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if attempt > 0 {
			operationprogress.Report(ctx, "download", fmt.Sprintf("Download needs another attempt; automatically retrying %s (attempt %d of 3)", asset.Name, attempt+1), 0, asset.Size, "bytes")
			timer := time.NewTimer(time.Duration(attempt) * time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
		info, err := out.Stat()
		if err != nil {
			return err
		}
		offset := info.Size()
		if offset >= asset.Size {
			if err := out.Truncate(0); err != nil {
				return err
			}
			offset = 0
		}
		req, err := request(ctx, assetEndpoint(asset, token), token)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/octet-stream")
		req.Header.Set("Accept-Encoding", "identity")
		end := asset.Size - 1
		if end-offset >= 8<<20 {
			end = offset + (8 << 20) - 1
		}
		expected := asset.Size - offset
		if asset.Size > 8<<20 {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, end))
		} else if offset > 0 {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
		}
		resp, err := client.Do(req)
		if err != nil {
			last = fmt.Errorf("download connection interrupted: %w", err)
			continue
		}
		if resp.StatusCode == http.StatusPartialContent && (offset > 0 || asset.Size > 8<<20) {
			wanted := fmt.Sprintf("bytes %d-%d/%d", offset, end, asset.Size)
			if resp.Header.Get("Content-Range") != wanted {
				resp.Body.Close()
				return errors.New("Release server returned an invalid download range")
			}
			expected = end - offset + 1
		} else if resp.StatusCode == http.StatusOK {
			if err := out.Truncate(0); err != nil {
				resp.Body.Close()
				return err
			}
			offset = 0
			expected = asset.Size
		} else {
			resp.Body.Close()
			last = fmt.Errorf("Release download failed (HTTP %d)", resp.StatusCode)
			if resp.StatusCode != 408 && resp.StatusCode != 429 && resp.StatusCode < 500 {
				return last
			}
			continue
		}
		if _, err := out.Seek(offset, io.SeekStart); err != nil {
			resp.Body.Close()
			return err
		}
		meterCtx := operationprogress.Observe(ctx, func(e operationprogress.Event) {
			operationprogress.Report(ctx, e.Stage, e.Detail, offset+e.Done, asset.Size, e.Unit)
		})
		meter := operationprogress.Bytes(meterCtx, "download", "Downloading "+asset.Name, expected)
		n, copyErr := io.Copy(io.MultiWriter(out, meter), io.LimitReader(resp.Body, expected+1))
		resp.Body.Close()
		meter.Finish()
		if copyErr != nil || n != expected {
			last = errors.New("Release download interrupted or incomplete; automatic retries exhausted")
			continue
		}
		if offset+n < asset.Size {
			attempt = -1
			continue
		}
		if _, err := out.Seek(0, io.SeekStart); err != nil {
			return err
		}
		hash := sha256.New()
		if _, err := io.Copy(hash, out); err != nil {
			return err
		}
		if "sha256:"+hex.EncodeToString(hash.Sum(nil)) != asset.Digest {
			last = errors.New("Release SHA256 verification failed after automatic retries; no unverified update was installed")
			if err := out.Truncate(0); err != nil {
				return err
			}
			continue
		}
		return out.Sync()
	}
	return last
}
