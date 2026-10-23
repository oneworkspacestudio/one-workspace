package maintenance

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
)

// Recheck an idle window without silently changing the user's reviewed files.
// Nothing is downloaded, stopped or backed up until this succeeds.
func refreshReviewedPlan(ctx context.Context, current catalog.Catalog, reviewed distribution.UpdatePlan, selected []string, reinstall bool, check func(context.Context, catalog.Catalog) (distribution.UpdatePlan, error)) (distribution.UpdatePlan, bool, error) {
	if reviewed.Fresh(time.Now()) {
		return reviewed, false, nil
	}
	fresh, err := check(ctx, current)
	if err != nil {
		return reviewed, false, errors.New("Could not refresh the release check; your installation has not changed. " + err.Error())
	}
	if !fresh.Fresh(time.Now()) {
		return reviewed, false, errors.New("The refreshed release check is not current. Check your PC clock and try again")
	}
	unchanged := false
	if reinstall {
		// Validate the reviewed identity and metadata independently of its age.
		original := reviewed
		original.CheckedAt = fresh.CheckedAt
		unchanged = distribution.ValidateInstallerPlan(current, original) == nil && distribution.ValidateInstallerPlan(current, fresh) == nil && reflect.DeepEqual(reviewed.Installer, fresh.Installer)
	} else {
		before, beforeFiles, beforeErr := distribution.SelectUpdates(current, reviewed, selected)
		after, afterFiles, afterErr := distribution.SelectUpdates(current, fresh, selected)
		unchanged = beforeErr == nil && afterErr == nil && reflect.DeepEqual(before, after) && reflect.DeepEqual(beforeFiles, afterFiles)
	}
	if !unchanged {
		return fresh, true, errors.New("Available releases or your installation changed. Review the refreshed choices, then try again. Your installation has not changed")
	}
	return fresh, true, nil
}

func (s *Server) preparePlan(ctx context.Context, reviewed distribution.UpdatePlan, selected []string, reinstall bool) (distribution.UpdatePlan, error) {
	c, err := suite.LoadUpdateConfig(s.Data)
	if err != nil {
		return reviewed, err
	}
	if !reviewed.Fresh(time.Now()) {
		s.progress("Refreshing the expired release check and confirming your selected files")
	}
	plan, refreshed, err := refreshReviewedPlan(ctx, c.Catalog, reviewed, selected, reinstall, distribution.CheckUpdates)
	if refreshed {
		s.mu.Lock()
		s.Plan = &plan
		s.mu.Unlock()
	}
	return plan, err
}
