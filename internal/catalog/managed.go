package catalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed managed-apps.json
var managedRegistry []byte

type ManagedApp struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Route                   string `json:"route"`
	Repository              string `json:"repository"`
	OIDC                    bool   `json:"oidc"`
	Foundation              bool   `json:"foundation"`
	Recovery                string `json:"recovery"`
	UpdateGroup             string `json:"updateGroup"`
	IndependentDistribution bool   `json:"independentDistribution,omitempty"`
	SinceIntegration        int    `json:"sinceIntegration,omitempty"`
}

func ManagedApps() []ManagedApp {
	var r struct {
		Apps []ManagedApp `json:"apps"`
	}
	if err := json.Unmarshal(managedRegistry, &r); err != nil {
		panic(err)
	}
	return r.Apps
}

// Release inventory validation is separate from Catalog.Validate: old installed
// catalogs and intentionally reduced test/selection catalogs remain supported.
func (c Catalog) ValidateManagedInventory() error {
	for _, want := range ManagedApps() {
		a, ok := c.Find(want.ID)
		if !ok && c.InstallerIntegration < want.SinceIntegration {
			continue
		}
		if want.IndependentDistribution && c.InstallerIntegration >= IndependentAppIntegration {
			if ok {
				return fmt.Errorf("managed release bundles an independent app: %s", want.ID)
			}
			for id := range c.Packages {
				if id == want.ID || strings.HasPrefix(id, want.ID+"-") {
					return fmt.Errorf("managed release bundles an independent package: %s", id)
				}
			}
			continue
		}
		if !ok {
			return fmt.Errorf("managed release omits %s (%s)", want.Name, want.ID)
		}
		name := want.Name
		if want.ID == "mymusic" && c.Version == "0.3.0-preview.142" {
			name = "My Music" // The immutable preview.142 catalog used the original display name.
		}
		if a.Name != name || a.Route != want.Route || a.OIDC != want.OIDC {
			return fmt.Errorf("managed contract drift for %s", want.ID)
		}
	}
	return nil
}
