package suite

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Curated native request shapes are validated before an action is shown for
// human review. The native app remains responsible for permissions, versions,
// quotas and business rules. No model-provided schema or executable validator runs.
func aiSchema(example any) map[string]any {
	switch v := example.(type) {
	case string:
		return map[string]any{"type": "string"}
	case bool:
		return map[string]any{"type": "boolean"}
	case int:
		return map[string]any{"type": "integer"}
	case float64:
		return map[string]any{"type": "number"}
	case []string:
		return map[string]any{"type": "array", "items": map[string]any{"type": "string"}}
	case []any:
		item := map[string]any{}
		if len(v) > 0 {
			item = aiSchema(v[0])
		}
		return map[string]any{"type": "array", "items": item}
	case map[string]any:
		props := map[string]any{}
		for k, val := range v {
			props[k] = aiSchema(val)
		}
		return map[string]any{"type": "object", "properties": props, "additionalProperties": true}
	}
	return map[string]any{}
}

func aiBodySchema(app string, op aiOperation) map[string]any {
	if !op.Review {
		return nil
	}
	example := op.Example
	if example == nil {
		example = map[string]any{}
	}
	s := aiSchema(example)
	// No-body actions do not accept arbitrary extra fields.
	s["additionalProperties"] = false
	p := s["properties"].(map[string]any)
	required := []string{}
	switch app + "/" + op.ID {
	case "tasks/create_task", "tasks/update_task":
		for _, name := range strings.Fields("id listId parentId title description notes due timezone status updated ruleId rule reminder snooze calendarLink notesLink generatedBy previousStatus") {
			p[name] = map[string]any{"type": "string"}
		}
		for _, name := range []string{"version", "priority"} {
			p[name] = map[string]any{"type": "integer"}
		}
		p["priority"] = map[string]any{"type": "integer", "minimum": 0, "maximum": 3, "description": "0 normal, up to 3 highest"}
		p["status"] = map[string]any{"type": "string", "enum": []any{"todo", "doing", "done"}}
		p["starred"] = map[string]any{"type": "boolean"}
		p["position"] = map[string]any{"type": "number"}
		for _, name := range []string{"tags", "exceptions"} {
			p[name] = map[string]any{"type": "array", "items": map[string]any{"type": "string"}}
		}
		required = []string{"title"}
		if op.ID == "update_task" {
			required = append(required, "listId", "version")
		}
	case "tasks/create_list", "notes/create_notebook", "music/create_playlist":
		required = []string{"name"}
	case "notes/create_note", "notes/update_note":
		for _, name := range strings.Fields("id owner notebook_id title body permission") {
			p[name] = map[string]any{"type": "string"}
		}
		for _, name := range strings.Fields("starred pinned archived deleted template") {
			p[name] = map[string]any{"type": "boolean"}
		}
		for _, name := range strings.Fields("reminder created updated revision") {
			p[name] = map[string]any{"type": "integer"}
		}
		p["tags"] = map[string]any{"type": "array", "items": map[string]any{"type": "string"}}
		required = []string{"title", "body"}
		if op.ID == "update_note" {
			required = append(required, "revision")
		}
	case "hosting/create_site":
		required = []string{"name", "slug", "adapter"}
		p["adapter"] = map[string]any{"type": "string", "enum": []any{"static"}}
	case "hosting/deploy_project":
		required = []string{"commit"}
		p["environment"] = map[string]any{"type": "string", "enum": []any{"preview", "production"}}
	case "drive/save_text":
		required = []string{"source", "path", "content"}
	case "calendar/save_event", "contacts/save_contact":
		required = []string{"href", "etag", "data"}
		for _, name := range []string{"name", "color"} {
			p[name] = map[string]any{"type": "string"}
		}
		p["readOnly"] = map[string]any{"type": "boolean"}
	case "social/publish_post", "clips/publish_post":
		p["location"] = map[string]any{"type": "string"}
		p["publishAt"] = map[string]any{"type": "integer"}
		p["extras"] = map[string]any{"type": "object", "additionalProperties": true}
		if app == "clips" {
			p["comments"] = map[string]any{"type": "boolean"}
			p["remix"] = map[string]any{"type": "boolean"}
		}
		required = []string{"kind", "caption", "media", "requestKey", "visibility"}
	case "maps/save_place":
		for _, name := range []string{"id", "address", "list", "notes"} {
			p[name] = map[string]any{"type": "string"}
		}
		required = []string{"name", "lat", "lon"}
	case "my-tube/create_playlist":
		required = []string{"title", "visibility"}
	case "create/save_project":
		required = []string{"name", "scene"}
		p["id"] = map[string]any{"type": "string"}
		p["source"] = map[string]any{"type": "string"}
	case "my3d/save_project":
		required = []string{"id", "name", "scene", "revision"}
	case "pdf/annotate":
		// Native PDF validates the numeric page keys and each Mark's geometry.
		return map[string]any{"type": "object", "additionalProperties": true}
	}
	if len(required) > 0 {
		s["required"] = required
	}
	if app == "notes" && op.ID == "create_notebook" {
		p["parent_id"] = map[string]any{"type": "string"}
	}
	return s
}

func aiValidateShape(schema map[string]any, value any, path string) error {
	typeName, _ := schema["type"].(string)
	valid := true
	switch typeName {
	case "object":
		_, valid = value.(map[string]any)
	case "string":
		_, valid = value.(string)
	case "boolean":
		_, valid = value.(bool)
	case "array":
		_, valid = value.([]any)
	case "number", "integer":
		n, ok := value.(float64)
		valid = ok && (typeName != "integer" || n == float64(int64(n)))
	}
	if !valid {
		return fmt.Errorf("%s must be %s; follow describe_workspace_app bodySchema", path, typeName)
	}
	if choices, ok := schema["enum"].([]any); ok {
		found := false
		for _, v := range choices {
			found = found || reflect.DeepEqual(v, value)
		}
		if !found {
			return fmt.Errorf("%s must be one of %v", path, choices)
		}
	}
	if n, ok := value.(float64); ok {
		if v, ok := schema["minimum"].(int); ok && n < float64(v) {
			return fmt.Errorf("%s must be at least %d", path, v)
		}
		if v, ok := schema["maximum"].(int); ok && n > float64(v) {
			return fmt.Errorf("%s must be at most %d", path, v)
		}
	}
	if m, ok := value.(map[string]any); ok {
		properties, _ := schema["properties"].(map[string]any)
		if required, ok := schema["required"].([]string); ok {
			for _, key := range required {
				if _, found := m[key]; !found {
					return fmt.Errorf("%s.%s is required", path, key)
				}
			}
		}
		for key, val := range m {
			field, found := properties[key]
			if !found {
				if allowed, ok := schema["additionalProperties"].(bool); ok && !allowed {
					return fmt.Errorf("Unknown field %s.%s; follow describe_workspace_app bodySchema", path, key)
				}
				continue
			}
			if e := aiValidateShape(field.(map[string]any), val, path+"."+key); e != nil {
				return e
			}
		}
	}
	if list, ok := value.([]any); ok {
		if item, ok := schema["items"].(map[string]any); ok {
			for _, val := range list {
				if e := aiValidateShape(item, val, path+"[]"); e != nil {
					return e
				}
			}
		}
	}
	return nil
}

func aiValidateBody(op aiOperation, raw json.RawMessage) error {
	if op.BodySchema == nil {
		return nil
	}
	var body any = map[string]any{}
	if len(raw) > 0 && string(raw) != "null" {
		if e := json.Unmarshal(raw, &body); e != nil {
			return e
		}
	}
	return aiValidateShape(op.BodySchema, body, "body")
}
