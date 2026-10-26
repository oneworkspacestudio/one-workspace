package suite

import "strings"

// The AI connector is a catalogue of native app operations, not an HTTP proxy.
// Every call is still authenticated and authorized by its destination app.
type aiOperation struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Method      string         `json:"-"`
	Path        string         `json:"-"`
	Parameters  []string       `json:"parameters,omitempty"`
	Query       []string       `json:"query,omitempty"`
	Example     map[string]any `json:"bodyExample,omitempty"`
	Review      bool           `json:"reviewRequired"`
	Description string         `json:"description,omitempty"`
	MailMethod  string         `json:"-"`
	BodySchema  map[string]any `json:"bodySchema,omitempty"`
}

func aiRead(id, title, path string, query ...string) aiOperation {
	return aiOperation{ID: id, Title: title, Method: "GET", Path: path, Query: query}
}
func aiChange(id, title, method, path string, example map[string]any) aiOperation {
	return aiOperation{ID: id, Title: title, Method: method, Path: path, Example: example, Review: true}
}

func aiOperations(id string) []aiOperation {
	var ops []aiOperation
	read := func(id, title, path string, query ...string) { ops = append(ops, aiRead(id, title, path, query...)) }
	change := func(id, title, method, path string, body map[string]any) {
		ops = append(ops, aiChange(id, title, method, path, body))
	}
	switch id {
	case "tasks":
		read("lists", "List task lists", "/api/lists")
		read("tasks", "Find tasks", "/api/tasks", "q", "list", "tag", "view", "includeCompleted", "offset")
		read("task", "Read a task", "/api/tasks/{id}")
		change("create_list", "Create a task list", "POST", "/api/lists", map[string]any{"name": "Project tasks"})
		change("create_task", "Create a task", "POST", "/api/tasks", map[string]any{"listId": "ID from lists", "title": "Task title", "description": "Details", "timezone": "Europe/London"})
		change("update_task", "Update a task", "PUT", "/api/tasks/{id}", map[string]any{"listId": "Existing listId", "title": "New title", "description": "Existing description", "status": "done", "version": 1})
		change("delete_task", "Delete a task", "DELETE", "/api/tasks/{id}", nil)
	case "notes":
		read("state", "List notebooks and note summaries", "/api/state")
		read("notes", "Find notes", "/api/notes", "q", "notebook", "tag", "view", "sort", "offset")
		read("note", "Read a note", "/api/notes/{id}")
		read("history", "Read note history", "/api/notes/{id}/history")
		change("create_notebook", "Create a notebook", "POST", "/api/notebooks", map[string]any{"name": "Project notes"})
		change("create_note", "Create a note", "POST", "/api/notes", map[string]any{"title": "New note", "body": "Note content", "notebook_id": "ID from state"})
		change("update_note", "Save a note", "PUT", "/api/notes/{id}", map[string]any{"title": "Title", "body": "Content", "revision": 1, "notebook_id": "ID from state"})
	case "hosting":
		read("sites", "List accessible hosted sites", "/api/projects")
		read("site", "Read site settings", "/api/projects/{id}")
		read("deployments", "List deployment history", "/api/projects/{id}/deployments", "offset")
		read("deployment", "Read deployment status and actual URL", "/api/projects/{id}/deployments/{deployment}")
		change("create_site", "Create a static hosted site", "POST", "/api/projects", map[string]any{"name": "My site", "slug": "my-site", "adapter": "static", "output": ".", "previews": true})
		ops = append(ops, aiOperation{ID: "deploy_project", Title: "Deploy saved Git project to My Hosting", Method: "POST", Path: "managed-project-upload", Review: true, Parameters: []string{"project", "site"}, Example: map[string]any{"commit": "Exact commit from Git status", "environment": "preview"}, Description: "Publishes an immutable, already pushed My AI project commit through My Hosting. Use Git status first. Source comes from Git, never from an arbitrary folder or URL. The existing site must be configured in My Hosting; scripts/trust are never enabled by this connector. Returns a queued job, not a claim of successful deployment. Poll deployment to get final state and URL."})
	case "gitea":
		ops = append(ops, aiOperation{ID: "projects", Title: "List private Git projects", Method: "GET", Path: "managed-git-projects"}, aiOperation{ID: "status", Title: "Read project Git status and saved commit", Method: "GET", Path: "managed-git-status", Parameters: []string{"project"}}, aiOperation{ID: "history", Title: "Read project commit history", Method: "GET", Path: "managed-git-history", Parameters: []string{"project"}}, aiOperation{ID: "save", Title: "Commit and push managed project changes", Method: "POST", Path: "managed-git-save", Parameters: []string{"project"}, Review: true})
	case "drive":
		read("sources", "List your accessible Drive sources", "/api/settings/sources")
		read("files", "List files in an authorized Drive source", "/api/resources", "source", "path", "sort", "order")
		read("read_file", "Read a text file from Drive", "/api/resources/download", "source", "file")
		// Managed project storage and writes use the existing private Projects broker.
		ops = append(ops, aiOperation{ID: "projects", Title: "List your private project storage", Method: "GET", Path: "managed-git-projects"})
		ops = append(ops, aiOperation{ID: "save_text", Title: "Save a text document in My Drive", Method: "POST", Path: "managed-drive-write", Review: true, Example: map[string]any{"source": "My Drive", "path": "/Notes/example.txt", "content": "Document text", "overwrite": false}, Description: "Creates or explicitly replaces a UTF-8 file through Drive's own scoped API. Does not write arbitrary host paths."})
	case "photos":
		read("sources", "List accessible photo sources", "/api/photos/capabilities")
		read("library", "Find photos and videos", "/api/photos/library", "source", "path", "kind", "q", "limit", "favorites", "cursor")
	case "music":
		read("library", "Read music collections and playlists", "/studio/bootstrap")
		read("tracks", "Find music", "/studio/tracks", "q", "album", "artist", "start", "limit")
		read("track", "Read track details", "/studio/tracks/{hash}")
		read("collection", "Read an accessible music collection", "/studio/collections/{owner}/{kind}/{resource}")
		change("create_playlist", "Create a music playlist", "POST", "/studio/playlists", map[string]any{"name": "Playlist"})
		change("update_playlist", "Update your music playlist", "PUT", "/studio/playlists/{id}", map[string]any{"name": "Playlist", "description": "Description", "tracks": []string{"Exact track hash from tracks"}})
	case "mail":
		read("account", "Read your Mail account capabilities", "/api/jmap?session=1")
		for _, method := range []string{"Mailbox/get", "Mailbox/query", "Email/get", "Email/query", "Thread/get", "Identity/get", "SearchSnippet/get"} {
			ops = append(ops, aiOperation{ID: strings.ToLower(strings.ReplaceAll(method, "/", "_")), Title: "Mail " + method, Method: "POST", Path: "/api/jmap", MailMethod: method, Example: map[string]any{"accountId": "Account ID from account", "ids": []string{"Exact IDs returned by Mail"}}, Description: "Native JMAP read operation. Body is its arguments object; the server fixes the method. For queries use filter/sort/position/limit instead of ids. Private mail stays out of public web searches."})
		}
	case "calendar":
		read("account", "Read calendar account and domain", "/api/session")
		read("calendars", "List calendars", "/api/calendars")
		read("events", "Find calendar events", "/api/events", "calendar", "start", "end", "q")
		change("save_event", "Create or update an event", "PUT", "/api/events", map[string]any{"href": "Calendar href from calendars plus unique-event.ics", "etag": "*", "data": "Complete VCALENDAR with VEVENT, user-requested dates/timezone, SUMMARY and unique UID; use exact existing ETag when updating"})
	case "contacts":
		read("account", "Read address-book account and domain", "/api/session")
		read("contacts", "Find contacts", "/api/contacts", "q")
		change("save_contact", "Create or update a contact", "PUT", "/api/contacts", map[string]any{"href": "Use an existing contact href, or /dav/card/ plus URL-escaped account subject@domain plus /unique-id.vcf; read account first", "etag": "*", "data": "Complete BEGIN:VCARD ... END:VCARD with VERSION:3.0, UID, FN and the requested contact details; use exact existing ETag when updating"})
	case "my-tube":
		read("videos", "Find videos", "/api/videos", "q", "channel", "limit", "offset")
		read("video", "Read video details", "/api/videos/{id}")
		read("channels", "Find channels", "/api/channels", "q")
		read("playlists", "List playlists", "/api/playlists")
		read("playlist", "Read a playlist", "/api/playlists/{id}")
		change("create_playlist", "Create a playlist", "POST", "/api/playlists", map[string]any{"title": "Playlist", "description": "Description", "visibility": "private"})
	case "cinema":
		read("profiles", "List viewing profiles", "/api/profiles")
		read("catalog", "Browse media", "/api/catalog", "q", "type", "genre", "offset")
		read("search", "Search media", "/api/search", "q")
		read("item", "Read media details", "/api/items/{id}")
		read("my_list", "Read watch list", "/api/my-list")
		change("save_to_list", "Add to watch list", "POST", "/api/my-list/{id}", map[string]any{})
	case "social", "clips":
		read("feed", "Read your feed", "/api/feed", "cursor", "limit")
		read("post", "Read a post", "/api/posts/{id}")
		read("comments", "Read comments", "/api/posts/{id}/comments", "cursor")
		read("profiles", "Find profiles", "/api/profiles", "q")
		read("drafts", "List drafts", "/api/drafts")
		read("media", "List your saved media", "/api/media")
		change("publish_post", "Publish a post", "POST", "/api/posts", map[string]any{"kind": "post", "caption": "Post content", "alt": "Media description", "visibility": "private", "media": []string{"Exact ID from media"}, "requestKey": "A new 48-character lowercase hex identifier; retain it when reconciling a retry"})
	case "maps":
		read("search", "Find places", "/api/search", "q")
		read("places", "List saved places", "/api/places")
		change("save_place", "Save a place", "POST", "/api/places", map[string]any{"name": "Place", "lat": 0, "lon": 0})
	case "meet":
		read("meetings", "List meetings", "/api/meetings")
		read("meeting", "Read meeting details", "/api/meetings/{code}")
		change("create_meeting", "Create a meeting", "POST", "/api/meetings", map[string]any{"title": "Meeting"})
	case "create":
		read("projects", "List creative projects", "/api/projects")
		read("project", "Read a creative project", "/api/projects/{id}")
		change("save_project", "Save a creative project", "POST", "/api/projects", map[string]any{"name": "Project", "revision": 0, "scene": map[string]any{"schema": 1, "pages": []any{map[string]any{"width": 800, "height": 600, "canvas": map[string]any{"objects": []any{}}}}}})
	case "my3d":
		read("projects", "List 3D projects", "/api/projects")
		read("project", "Read a 3D project", "/api/projects/{id}")
		change("save_project", "Save a 3D project", "POST", "/api/projects", map[string]any{"id": "New unique identifier or exact existing project ID", "name": "Project", "revision": 0, "scene": map[string]any{"schema": 1, "units": "mm", "objects": []any{}}})
	case "pdf":
		read("documents", "List PDF documents", "/api/documents")
		read("annotations", "Read PDF annotations", "/api/documents/{id}/annotations")
		read("jobs", "List PDF jobs", "/api/jobs")
		change("annotate", "Save PDF annotations", "PUT", "/api/documents/{id}/annotations", map[string]any{"1": []any{map[string]any{"kind": "comment", "x": 20, "y": 20, "w": 150, "h": 30, "text": "Comment"}}})
	case "slicer":
		read("projects", "List printing projects", "/api/projects")
		read("project", "Read printing project", "/api/projects/{id}")
		read("slices", "List slices", "/api/slices")
		read("printers", "List configured printers", "/api/printers")
		read("prints", "List print jobs", "/api/prints")
	case "monitor":
		read("status", "Read Workspace monitoring status", "/api/overview")
		read("checks", "List monitoring checks", "/api/checks")
	case "search":
		read("preferences", "Read search preferences", "/api/preferences")
	}
	for i := range ops {
		ops[i].BodySchema = aiBodySchema(id, ops[i])
		if ops[i].ID == "update_task" || ops[i].ID == "update_note" {
			ops[i].Description = "Read the current record first, preserve all unchanged fields, and supply its exact current version/revision. Native conflict checks still apply."
		}
		if id == "tasks" && ops[i].ID == "create_task" {
			ops[i].Description = "Only title is required. Omit listId to use your Inbox. For another list, read lists and use its exact ID. Optional priority is 0 through 3; omit fields not requested by the person."
		}
		// Explicit special-operation parameters precede native route parameters.
		for _, part := range splitAIPath(ops[i].Path) {
			if len(part) > 2 && part[0] == '{' && part[len(part)-1] == '}' {
				ops[i].Parameters = append(ops[i].Parameters, part[1:len(part)-1])
			}
		}
	}
	return ops
}
