package suite

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// Complete additions inside an older controller's existing backup transaction.
// Ordinary launches never install software or change an app selection. The
// runtime lock is held, and a failed launch makes the controller roll back.
func startupDefaults(ctx context.Context, c Config, s Secrets, report func(string)) (Config, Secrets, error) {
	if c.Catalog.Selection != nil {
		return c, s, nil
	}
	var journal UpdateJournal
	token := os.Getenv("WORKSPACE_UPDATE_TRANSACTION")
	if token == "" || platform.LoadJSON(updateJournalPath(c), &journal) != nil || journal.Phase != "testing" || journal.Token != token {
		return c, s, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()
	report("Checking newly included apps for this Workspace update")
	plan, err := distribution.StartupDefaultPlan(ctx, c.Catalog)
	if err != nil {
		return c, s, fmt.Errorf("finish new app installation: %w", err)
	}
	if len(plan.Choices) == 0 {
		c.Catalog.Selection = c.Catalog.EffectiveSelection()
		return c, s, Save(c, s)
	}
	var selected []string
	for _, choice := range plan.Choices {
		selected = append(selected, choice.ID)
	}
	target, downloads, err := distribution.SelectUpdates(c.Catalog, plan, selected)
	if err != nil {
		return c, s, err
	}
	// The old controller's launch plan must remain valid. It already selected
	// existing app upgrades; this bridge may only add apps and provenance.
	for _, old := range c.Catalog.Apps {
		next, exists := target.Find(old.ID)
		if !exists || !reflect.DeepEqual(old, next) {
			return c, s, fmt.Errorf("new app changes %s; finish with the newest installer", old.Name)
		}
	}
	for id, old := range c.Catalog.Packages {
		if target.Packages[id] != old {
			return c, s, errors.New("new app changes existing packages; finish with the newest installer")
		}
	}
	next := c
	next.Catalog = target
	if err = CheckPorts(next); err != nil {
		return c, s, err
	}
	stage, err := os.MkdirTemp(filepath.Join(c.Data, "updates"), "new-apps-")
	if err != nil {
		return c, s, err
	}
	defer os.RemoveAll(stage)
	if err = distribution.StageUpdateFiles(ctx, target, downloads, stage, report); err != nil {
		return c, s, err
	}
	for id, p := range target.Packages {
		if _, exists := c.Catalog.Packages[id]; exists {
			continue
		}
		src, e := platform.Inside(filepath.Join(stage, "packages"), p.File)
		if e != nil {
			return c, s, e
		}
		dst, e := platform.Inside(c.Programs, p.File)
		if e != nil {
			return c, s, e
		}
		if e = writeReplace(src, dst); e != nil {
			return c, s, e
		}
	}
	report("Connecting newly included apps to your existing sign-in")
	if err = connectAddedApps(ctx, next, &s); err != nil {
		return c, s, err
	}
	if err = Save(next, s); err != nil {
		return c, s, err
	}
	return next, s, nil
}
