package distribution

import (
	"errors"
	"fmt"
	"regexp"
)

var managerRequirementVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9]+(?:[.-][A-Za-z0-9]+)*)?$`)

// Add a capability only after its manager integration exists and is tested.
// The origin capability is exercised with actual native SSO, catalog and logout
// fixtures. It does not claim application E2EE, calling or native mobile support.
var managerCapabilities = map[string]bool{"self-hosted-push-v1": true, "app-update-requirements-v1": true, "chat-state-v2": true, "isolated-chat-origin-v1": true, "guarded-social-storage-v1": true, "scoped-social-storage-v2": true, "scoped-clips-storage-v1": true, "scoped-forms-storage-v1": true}

func checkRequirements(current, minimum string, capabilities []string) error {
	if len(minimum) > 80 || !managerRequirementVersion.MatchString(minimum) {
		return errors.New("invalid minimum Workspace version")
	}
	if current != minimum && !newerWorkspace(current, minimum) {
		return fmt.Errorf("Update One Workspace to %s first, then check app updates again. Your installed apps are unchanged", minimum)
	}
	if len(capabilities) > 24 {
		return errors.New("too many required manager capabilities")
	}
	seen := map[string]bool{}
	for _, capability := range capabilities {
		if !managerCapabilities[capability] || seen[capability] {
			return fmt.Errorf("This app needs a Workspace capability not supported by this manager: %s. Update Workspace first", capability)
		}
		seen[capability] = true
	}
	return nil
}
func checkManifestRequirements(m AppUpdateManifest, current string) error {
	if m.Schema == 1 && m.Integration == 2 {
		if m.MinimumManager != "" || len(m.RequiredCapabilities) > 0 {
			return errors.New("manager requirements need the version 2 update contract")
		}
		return nil
	}
	if m.Schema != 2 || m.Integration != 3 {
		return errors.New("unsupported app update contract")
	}
	return checkRequirements(current, m.MinimumManager, m.RequiredCapabilities)
}

// Schema2 is intentionally rejected by old managers before downloading or
// applying packages. Requirements are never silently added to a schema1 file.
func RequireManager(m *AppUpdateManifest, minimum string, capabilities []string) error {
	if minimum == "" && len(capabilities) == 0 {
		return nil
	}
	if err := checkRequirements(minimum, minimum, capabilities); err != nil {
		return err
	}
	m.Schema = 2
	m.Integration = 3
	m.MinimumManager = minimum
	m.RequiredCapabilities = append([]string(nil), capabilities...)
	return nil
}
