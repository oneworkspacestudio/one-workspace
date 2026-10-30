package catalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed managed-apps.json
var managedRegistry []byte

type ManagedApp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Route       string `json:"route"`
	Repository  string `json:"repository"`
	OIDC        bool   `json:"oidc"`
	Foundation  bool   `json:"foundation"`
	Recovery    string `json:"recovery"`
	UpdateGroup string `json:"updateGroup"`
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
		if !ok {
			return fmt.Errorf("managed release omits %s (%s)", want.Name, want.ID)
		}
		if a.Name != want.Name || a.Route != want.Route || a.OIDC != want.OIDC {
			return fmt.Errorf("managed contract drift for %s", want.ID)
		}
	}
	return nil
}
