package suite

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDedicatedIngressPrecedesWildcardAndPreservesOtherRoutes(t *testing.T) {
	home := map[string]any{"hostname": "workspace.example.test", "service": "http://127.0.0.1:8990"}
	wildcard := map[string]any{"hostname": "*.example.test", "service": "http://127.0.0.1:9007", "originRequest": map[string]any{"connectTimeout": "15s"}}
	last := map[string]any{"service": "http_status:404"}
	old := []any{home, wildcard, last}
	before, _ := json.Marshal(old)
	next, changed, err := dedicatedTunnelIngress(old, "chat.example.test", "http://127.0.0.1:9006")
	if err != nil || !changed || len(next) != 4 || next[0].(map[string]any)["hostname"] != "chat.example.test" || !reflect.DeepEqual(next[1:], old) {
		t.Fatal("exact Chat route did not precede the unchanged wildcard fallbacks", err)
	}
	after, _ := json.Marshal(old)
	if string(before) != string(after) {
		t.Fatal("mutated another app's original route")
	}
	stable, changed, err := dedicatedTunnelIngress(next, "chat.example.test", "http://127.0.0.1:9006")
	if err != nil || changed || !reflect.DeepEqual(stable, next) {
		t.Fatal("correct configuration was not stable", err)
	}
}

func TestDedicatedIngressRepairsOnlyItsOwnedHostOverride(t *testing.T) {
	wildcard := map[string]any{"hostname": "*.example.test", "service": "http://127.0.0.1:9007"}
	owned := map[string]any{"hostname": "chat.example.test", "service": "http://127.0.0.1:9006", "originRequest": map[string]any{"httpHostHeader": "old.example.test", "connectTimeout": "15s", "noTLSVerify": false}}
	last := map[string]any{"service": "http_status:404"}
	old := []any{wildcard, owned, last}
	before, _ := json.Marshal(old)
	next, changed, err := dedicatedTunnelIngress(old, "chat.example.test", "http://127.0.0.1:9006")
	if err != nil || !changed || !reflect.DeepEqual(next[1:], []any{wildcard, last}) {
		t.Fatal("owned route repair changed unrelated routes", err)
	}
	options := next[0].(map[string]any)["originRequest"].(map[string]any)
	if options["httpHostHeader"] != "chat.example.test" || options["connectTimeout"] != "15s" || options["noTLSVerify"] != false {
		t.Fatal("lost an owned origin option")
	}
	after, _ := json.Marshal(old)
	if string(before) != string(after) {
		t.Fatal("mutated input configuration")
	}
}

func TestDedicatedIngressRejectsConflictingOrInvalidRoutes(t *testing.T) {
	last := map[string]any{"service": "http_status:404"}
	for _, routes := range [][]any{
		{map[string]any{"hostname": "chat.example.test", "service": "http://127.0.0.1:9999"}, last},
		{map[string]any{"hostname": "CHAT.example.test", "service": "http://127.0.0.1:9999"}, last},
		{map[string]any{"hostname": "chat.example.test", "service": "http://127.0.0.1:9006", "path": "^/other"}, last},
		{map[string]any{"hostname": "chat.example.test", "service": "http://127.0.0.1:9006"}, map[string]any{"hostname": "chat.example.test", "service": "http://127.0.0.1:9006"}, last},
		{map[string]any{"hostname": "chat.example.test", "service": "http://127.0.0.1:9006"}},
		{map[string]any{"hostname": "chat.example.test", "service": "http://127.0.0.1:9006", "originRequest": []any{}}, last},
	} {
		if _, _, err := dedicatedTunnelIngress(routes, "chat.example.test", "http://127.0.0.1:9006"); err == nil {
			t.Fatal("conflicting configuration accepted")
		}
	}
}
