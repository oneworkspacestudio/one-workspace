package suite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

var aiActionID = regexp.MustCompile(`^[a-f0-9]{64}$`)

type aiDeployIntent struct {
	Site, Commit, SourceHash, Environment, Deployment, Created string
	Rejected                                                   string `json:"rejected,omitempty"`
}
type aiDeployPending struct {
	Key string `json:"key"`
}

// Receipts distinguish a retry of one reviewed action from a later deliberate
// deployment of the same revision. A pending source receipt survives lost HTTP
// responses and process restarts; it is reconciled before any duplicate upload.
func (p *Portal) aiSubmitDeployment(ctx context.Context, uid string, hosting catalog.App, connection aiAppConnection, csrf string, call aiAppCall, archive []byte, commit, environment, repository string) (any, error) {
	if !aiActionID.MatchString(call.RequestID) || !aiActionID.MatchString(uid) {
		return nil, errors.New("Deployment needs a server-generated action receipt")
	}
	site := call.Parameters["site"]
	digest := sha256.Sum256(archive)
	sourceHash := hex.EncodeToString(digest[:])
	root := filepath.Join(p.rt.Config.Data, "myai-projects", "deployments", uid)
	if e := myaiSafeTreePath(root); e != nil {
		return nil, e
	}
	stable := call
	stable.Body, _ = json.Marshal(map[string]string{"commit": commit, "environment": environment})
	sourceKey := aiRequestDigest(uid, stable)
	requestKey := sha256.Sum256([]byte(call.RequestID + "\x00" + sourceKey))
	key := hex.EncodeToString(requestKey[:])
	path := filepath.Join(root, key+".json")
	pendingPath := filepath.Join(root, "pending-"+sourceKey+".json")
	var intent aiDeployIntent
	existing := false
	alias := ""
	if e := platform.LoadJSON(path, &intent); e == nil {
		existing = true
	} else if !os.IsNotExist(e) {
		return nil, errors.New("The protected deployment receipt is unavailable")
	} else {
		var pending aiDeployPending
		if e = platform.LoadJSON(pendingPath, &pending); e == nil {
			if !aiActionID.MatchString(pending.Key) {
				return nil, errors.New("The pending deployment receipt is invalid")
			}
			previous := filepath.Join(root, pending.Key+".json")
			var prior aiDeployIntent
			if e = platform.LoadJSON(previous, &prior); e != nil {
				return nil, errors.New("The pending deployment receipt is unavailable")
			}
			intent = prior
			existing = true
			alias = path
			path = previous
		} else if !os.IsNotExist(e) {
			return nil, errors.New("The pending deployment receipt is unreadable")
		}
	}
	if existing && (intent.Site != site || intent.Commit != commit || intent.Environment != environment || intent.SourceHash != sourceHash) {
		return nil, errors.New("The deployment receipt differs from the reviewed revision")
	}
	save := func() error {
		if e := platform.SaveJSON(path, intent); e != nil {
			return e
		}
		if alias != "" {
			if e := platform.SaveJSON(alias, intent); e != nil {
				return e
			}
		}
		return nil
	}
	clearPending := func() {
		var pending aiDeployPending
		if platform.LoadJSON(pendingPath, &pending) == nil && pending.Key+".json" == filepath.Base(path) {
			_ = os.Remove(pendingPath)
		}
	}
	result := func(raw []byte, kind, notice string) (any, error) {
		value, e := aiAppResult(raw, kind)
		if e != nil {
			return nil, e
		}
		return map[string]any{"deployment": value, "site": site, "sourceCommit": commit, "repository": repository, "sourceHash": sourceHash, "notice": notice}, nil
	}
	if existing && intent.Rejected != "" {
		if e := save(); e != nil {
			return nil, e
		}
		clearPending()
		return nil, errors.New(intent.Rejected)
	}
	if existing && intent.Deployment != "" {
		if e := save(); e != nil {
			return nil, e
		}
		clearPending()
		raw, kind, e := p.aiNativeRequest(ctx, hosting, connection, csrf, "GET", "/api/projects/"+site+"/deployments/"+intent.Deployment, nil, "")
		if e != nil {
			return nil, e
		}
		return result(raw, kind, "This reviewed action already has a deployment job. Check its final state and URL")
	}
	if existing {
		created, e := time.Parse(time.RFC3339Nano, intent.Created)
		if e != nil {
			return nil, errors.New("The pending deployment time is invalid; inspect My Hosting")
		}
		for offset := 0; offset < 300; offset += 30 {
			raw, _, e := p.aiNativeRequest(ctx, hosting, connection, csrf, "GET", fmt.Sprintf("/api/projects/%s/deployments?offset=%d", site, offset), nil, "")
			if e != nil {
				return nil, e
			}
			var page struct {
				Deployments []struct{ ID, SourceHash, Environment, Created string }
			}
			if json.Unmarshal(raw, &page) != nil {
				return nil, errors.New("Hosting history is unreadable")
			}
			for _, d := range page.Deployments {
				at, parseErr := time.Parse(time.RFC3339Nano, d.Created)
				if d.SourceHash == sourceHash && d.Environment == environment && parseErr == nil && !at.Before(created) && aiParameter.MatchString(d.ID) {
					intent.Deployment = d.ID
					if e = save(); e != nil {
						return nil, e
					}
					clearPending()
					raw, kind, e := p.aiNativeRequest(ctx, hosting, connection, csrf, "GET", "/api/projects/"+site+"/deployments/"+d.ID, nil, "")
					if e != nil {
						return nil, e
					}
					return result(raw, kind, "Recovered the existing job after an unconfirmed response; no duplicate upload was started")
				}
			}
			if len(page.Deployments) < 30 {
				break
			}
		}
		return nil, errors.New("A previous upload has an uncertain outcome. Inspect My Hosting before retrying; no duplicate deployment was started")
	}
	intent = aiDeployIntent{Site: site, Commit: commit, SourceHash: sourceHash, Environment: environment, Created: time.Now().UTC().Format(time.RFC3339Nano)}
	if e := save(); e != nil {
		return nil, e
	}
	if e := platform.SaveJSON(pendingPath, aiDeployPending{Key: key}); e != nil {
		return nil, e
	}
	raw, kind, e := p.aiNativeRequest(ctx, hosting, connection, csrf, "POST", "/api/projects/"+site+"/deployments?environment="+environment, archive, "application/zip")
	if e != nil {
		var rejected *aiHTTPError
		// Hosting's 409 can follow an ambiguous database update; preserve it as
		// uncertain. Validation/auth/body rejections cannot publish this upload.
		if errors.As(e, &rejected) && (rejected.Status == 400 || rejected.Status == 401 || rejected.Status == 403 || rejected.Status == 404 || rejected.Status == 413 || rejected.Status == 415 || rejected.Status == 422) {
			intent.Rejected = e.Error()
			if e = save(); e != nil {
				return nil, e
			}
			clearPending()
			return nil, errors.New(intent.Rejected)
		}
		return nil, errors.New("Hosting has not confirmed the upload. Retry this same revision to reconcile its existing job; no successful deployment is claimed")
	}
	var deployment struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(raw, &deployment) != nil || !aiParameter.MatchString(deployment.ID) {
		return nil, errors.New("Hosting did not confirm a deployment ID; inspect its history before retrying")
	}
	intent.Deployment = deployment.ID
	if e = save(); e != nil {
		return nil, errors.New("Hosting accepted the source but its receipt could not be saved. Inspect My Hosting before retrying")
	}
	clearPending()
	return result(raw, kind, "Source accepted by My Hosting. Read deployment status and report its real final state and URL")
}
