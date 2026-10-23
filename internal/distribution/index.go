package distribution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const officialRepository = "oneworkspacestudio/one-workspace"
const publicIndexURL = "https://raw.githubusercontent.com/oneworkspacestudio/one-workspace/main/updates/releases-v1.json"
const maxIndexSize = 8 << 20

// The publisher promotes verified public releases into one update channel.
// Clients do not spend the shared anonymous GitHub API allowance per app/check.
type releaseIndex struct {
	Schema       int                  `json:"schema"`
	PublishedAt  time.Time            `json:"publishedAt"`
	Repositories map[string][]Release `json:"repositories"`
}
type indexContextKey struct{}

func withReleaseIndex(ctx context.Context, token string) (context.Context, error) {
	if token != "" || Repository != officialRepository || ctx.Value(indexContextKey{}) != nil {
		return ctx, nil
	}
	meta, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := request(meta, publicIndexURL, "")
	if err != nil {
		return ctx, err
	}
	h := &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(r *http.Request, via []*http.Request) error {
		if len(via) > 3 || r.URL.Scheme != "https" || r.URL.Host != "raw.githubusercontent.com" || r.URL.User != nil {
			return errors.New("update list redirect rejected")
		}
		r.Header.Del("Authorization")
		return nil
	}}
	r, err := h.Do(req)
	if err != nil {
		return ctx, errors.New("Could not load the published update list. Check your connection and retry; your installation has not changed")
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("Published update list is unavailable (HTTP %d). Retry later; your installation has not changed", r.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(r.Body, maxIndexSize+1))
	if err != nil || len(b) > maxIndexSize {
		return ctx, errors.New("Published update list could not be read safely; retry later")
	}
	var index releaseIndex
	if err = json.Unmarshal(b, &index); err != nil || index.Schema != 1 || index.PublishedAt.IsZero() || index.PublishedAt.After(time.Now().Add(5*time.Minute)) || len(index.Repositories) == 0 || len(index.Repositories) > 100 {
		return ctx, errors.New("Published update list is invalid; contact the Workspace publisher")
	}
	for repo, releases := range index.Repositories {
		if !strings.HasPrefix(repo, "oneworkspacestudio/") || strings.Count(repo, "/") != 1 || strings.ContainsAny(repo, "?#\\ ") || len(releases) > 100 {
			return ctx, errors.New("Published update list has an invalid repository")
		}
		tags := map[string]bool{}
		for _, release := range releases {
			if release.Tag == "" || release.Draft || release.Published.IsZero() || tags[release.Tag] || release.Published.After(index.PublishedAt.Add(5*time.Minute)) {
				return ctx, errors.New("Published update list has invalid release information")
			}
			tags[release.Tag] = true
		}
	}
	return context.WithValue(ctx, indexContextKey{}, &index), nil
}

func indexPublishedAt(ctx context.Context) time.Time {
	if index, ok := ctx.Value(indexContextKey{}).(*releaseIndex); ok {
		return index.PublishedAt
	}
	return time.Time{}
}

func assetEndpoint(a Asset, token string) string {
	if token == "" {
		return a.URL
	}
	return a.APIURL
}

func releaseHTTPError(repo string, r *http.Response) error {
	if r.StatusCode == 429 || (r.StatusCode == 403 && (r.Header.Get("X-RateLimit-Remaining") == "0" || r.Header.Get("Retry-After") != "")) {
		retry := "Retry later"
		if seconds, err := strconv.ParseInt(r.Header.Get("Retry-After"), 10, 64); err == nil && seconds >= 0 {
			retry = fmt.Sprintf("Retry in %d seconds", seconds)
		} else if reset, err := strconv.ParseInt(r.Header.Get("X-RateLimit-Reset"), 10, 64); err == nil {
			retry = "Retry after " + time.Unix(reset, 0).Local().Format("15:04 MST")
		}
		return fmt.Errorf("GitHub's API request limit was reached while checking %s. %s, or use the latest public Workspace installer. Your installation has not changed", repo, retry)
	}
	return fmt.Errorf("release check for %s failed (HTTP %d); check repository access and retry", repo, r.StatusCode)
}
