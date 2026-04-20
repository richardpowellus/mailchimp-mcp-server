package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterMembers(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_members — GET /lists/{list_id}/members
	s.RegisterTool(mcp.Tool{
		Name:        "list_members",
		Description: "List members (subscribers) of an audience. GET /lists/{list_id}/members. Supports filtering by status and opt-in date.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":             {Type: "string", Description: "Account name."},
				"list_id":             {Type: "string", Description: "The audience/list ID."},
				"status":              {Type: "string", Description: "Filter by status: subscribed, unsubscribed, cleaned, pending, transactional, archived."},
				"since_timestamp_opt": {Type: "string", Description: "Filter members who opted in after this ISO 8601 date."},
			}, paging.Properties()),
			Required: []string{"account", "list_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account          string `json:"account"`
			ListID           string `json:"list_id"`
			Status           string `json:"status"`
			SinceTimestampOpt string `json:"since_timestamp_opt"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		q := url.Values{}
		if p.Status != "" {
			q.Set("status", p.Status)
		}
		if p.SinceTimestampOpt != "" {
			q.Set("since_timestamp_opt", p.SinceTimestampOpt)
		}
		endpoint := fmt.Sprintf("/lists/%s/members", p.ListID)
		if len(q) > 0 {
			endpoint += "?" + q.Encode()
		}
		items, err := client.FetchAll(ctx, endpoint, "members")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// add_member — POST /lists/{list_id}/members
	s.RegisterTool(mcp.Tool{
		Name:        "add_member",
		Description: "Add a new member to an audience. POST /lists/{list_id}/members. Body must include email_address and status (subscribed, unsubscribed, cleaned, pending, transactional). Optional: merge_fields, tags, language, vip, location, interests.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"body":    {Type: "object", Description: "Member payload: email_address (required), status (required), merge_fields, tags, etc."},
			},
			Required: []string{"account", "list_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/members", p.ListID), p.Body)
	})

	// get_member — GET /lists/{list_id}/members/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_member",
		Description: "Get a member by email address. GET /lists/{list_id}/members/{subscriber_hash}. Email is hashed to subscriber_hash automatically.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/lists/%s/members/%s", p.ListID, hash))
	})

	// update_member — PATCH /lists/{list_id}/members/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "update_member",
		Description: "Update a member. PATCH /lists/{list_id}/members/{subscriber_hash}. Body can include status, merge_fields, interests, language, vip, etc.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
				"body":    {Type: "object", Description: "Fields to update: status, merge_fields, interests, language, vip, etc."},
			},
			Required: []string{"account", "list_id", "email", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Email   string          `json:"email"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.PatchRaw(ctx, fmt.Sprintf("/lists/%s/members/%s", p.ListID, hash), p.Body)
	})

	// upsert_member — PUT /lists/{list_id}/members/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "upsert_member",
		Description: "Add or update a member (upsert). PUT /lists/{list_id}/members/{subscriber_hash}. Creates if not exists, updates if exists. Body must include email_address and status_if_new.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally for URL)."},
				"body":    {Type: "object", Description: "Upsert payload: email_address, status_if_new (required for new), merge_fields, etc."},
			},
			Required: []string{"account", "list_id", "email", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Email   string          `json:"email"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.PutRaw(ctx, fmt.Sprintf("/lists/%s/members/%s", p.ListID, hash), p.Body)
	})

	// archive_member — DELETE /lists/{list_id}/members/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "archive_member",
		Description: "Archive (soft delete) a member. DELETE /lists/{list_id}/members/{subscriber_hash}. The member is archived but not permanently removed.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Delete(ctx, fmt.Sprintf("/lists/%s/members/%s", p.ListID, hash))
	})

	// delete_member_permanent — POST /lists/{list_id}/members/{subscriber_hash}/actions/delete-permanent
	s.RegisterTool(mcp.Tool{
		Name:        "delete_member_permanent",
		Description: "Permanently delete a member. POST /lists/{list_id}/members/{subscriber_hash}/actions/delete-permanent. Cannot be undone — removes all data for this contact.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Post(ctx, fmt.Sprintf("/lists/%s/members/%s/actions/delete-permanent", p.ListID, hash), nil)
	})

	// get_member_activity — GET /lists/{list_id}/members/{subscriber_hash}/activity
	s.RegisterTool(mcp.Tool{
		Name:        "get_member_activity",
		Description: "Get the last 50 activity events for a member. GET /lists/{list_id}/members/{subscriber_hash}/activity.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/lists/%s/members/%s/activity", p.ListID, hash))
	})

	// get_member_activity_feed — GET /lists/{list_id}/members/{subscriber_hash}/activity-feed
	s.RegisterTool(mcp.Tool{
		Name:        "get_member_activity_feed",
		Description: "Get the activity feed for a member. GET /lists/{list_id}/members/{subscriber_hash}/activity-feed. More detailed than activity endpoint.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/lists/%s/members/%s/activity-feed", p.ListID, hash))
	})

	// get_member_tags — GET /lists/{list_id}/members/{subscriber_hash}/tags
	s.RegisterTool(mcp.Tool{
		Name:        "get_member_tags",
		Description: "Get all tags for a member. GET /lists/{list_id}/members/{subscriber_hash}/tags.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			}, paging.Properties()),
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/members/%s/tags", p.ListID, hash), "tags")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// manage_member_tags — POST /lists/{list_id}/members/{subscriber_hash}/tags
	s.RegisterTool(mcp.Tool{
		Name:        "manage_member_tags",
		Description: "Add or remove tags from a member. POST /lists/{list_id}/members/{subscriber_hash}/tags. Body: tags array of objects with name and status (active or inactive).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
				"body":    {Type: "object", Description: "Tags payload: tags (array of {name, status} where status is 'active' or 'inactive')."},
			},
			Required: []string{"account", "list_id", "email", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Email   string          `json:"email"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/members/%s/tags", p.ListID, hash), p.Body)
	})

	// list_member_notes — GET /lists/{list_id}/members/{subscriber_hash}/notes
	s.RegisterTool(mcp.Tool{
		Name:        "list_member_notes",
		Description: "List notes for a member. GET /lists/{list_id}/members/{subscriber_hash}/notes.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			}, paging.Properties()),
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/members/%s/notes", p.ListID, hash), "notes")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// add_member_note — POST /lists/{list_id}/members/{subscriber_hash}/notes
	s.RegisterTool(mcp.Tool{
		Name:        "add_member_note",
		Description: "Add a note to a member. POST /lists/{list_id}/members/{subscriber_hash}/notes. Body: note (text string).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
				"body":    {Type: "object", Description: "Note payload: note (text string)."},
			},
			Required: []string{"account", "list_id", "email", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Email   string          `json:"email"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/members/%s/notes", p.ListID, hash), p.Body)
	})

	// get_member_note — GET /lists/{list_id}/members/{subscriber_hash}/notes/{note_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_member_note",
		Description: "Get a specific note for a member. GET /lists/{list_id}/members/{subscriber_hash}/notes/{note_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
				"note_id": {Type: "string", Description: "The note ID."},
			},
			Required: []string{"account", "list_id", "email", "note_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
			NoteID  string `json:"note_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/lists/%s/members/%s/notes/%s", p.ListID, hash, p.NoteID))
	})

	// update_member_note — PATCH /lists/{list_id}/members/{subscriber_hash}/notes/{note_id}
	s.RegisterTool(mcp.Tool{
		Name:        "update_member_note",
		Description: "Update a member note. PATCH /lists/{list_id}/members/{subscriber_hash}/notes/{note_id}. Body: note (text string).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
				"note_id": {Type: "string", Description: "The note ID."},
				"body":    {Type: "object", Description: "Updated note payload: note (text string)."},
			},
			Required: []string{"account", "list_id", "email", "note_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Email   string          `json:"email"`
			NoteID  string          `json:"note_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.PatchRaw(ctx, fmt.Sprintf("/lists/%s/members/%s/notes/%s", p.ListID, hash, p.NoteID), p.Body)
	})

	// delete_member_note — DELETE /lists/{list_id}/members/{subscriber_hash}/notes/{note_id}
	s.RegisterTool(mcp.Tool{
		Name:        "delete_member_note",
		Description: "Delete a member note. DELETE /lists/{list_id}/members/{subscriber_hash}/notes/{note_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
				"note_id": {Type: "string", Description: "The note ID."},
			},
			Required: []string{"account", "list_id", "email", "note_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
			NoteID  string `json:"note_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Delete(ctx, fmt.Sprintf("/lists/%s/members/%s/notes/%s", p.ListID, hash, p.NoteID))
	})

	// get_member_events — GET /lists/{list_id}/members/{subscriber_hash}/events
	s.RegisterTool(mcp.Tool{
		Name:        "get_member_events",
		Description: "Get events for a member. GET /lists/{list_id}/members/{subscriber_hash}/events. Returns custom events tracked for this contact.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
			}, paging.Properties()),
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/members/%s/events", p.ListID, hash), "events")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// add_member_event — POST /lists/{list_id}/members/{subscriber_hash}/events
	s.RegisterTool(mcp.Tool{
		Name:        "add_member_event",
		Description: "Add a custom event for a member. POST /lists/{list_id}/members/{subscriber_hash}/events. Body: name (required), properties (optional object), is_syncing (optional bool), occurred_at (optional ISO 8601).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "Member email address (hashed internally)."},
				"body":    {Type: "object", Description: "Event payload: name (required), properties (optional object)."},
			},
			Required: []string{"account", "list_id", "email", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Email   string          `json:"email"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/members/%s/events", p.ListID, hash), p.Body)
	})

	// batch_subscribe — POST /lists/{list_id}
	s.RegisterTool(mcp.Tool{
		Name:        "batch_subscribe",
		Description: "Batch add or update members. POST /lists/{list_id}. Body: members (array of member objects), update_existing (bool). Each member needs email_address and status.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"body":    {Type: "object", Description: "Batch payload: members (array), update_existing (bool)."},
			},
			Required: []string{"account", "list_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s", p.ListID), p.Body)
	})

	// search_members — GET /search-members?query={query}&list_id={list_id}
	s.RegisterTool(mcp.Tool{
		Name:        "search_members",
		Description: "Search for members across audiences. GET /search-members?query={query}. Optionally filter by list_id. Searches by email, name, and merge fields.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"query":   {Type: "string", Description: "Search query (email, name, or merge field value)."},
				"list_id": {Type: "string", Description: "Optional audience/list ID to restrict search."},
			},
			Required: []string{"account", "query"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Query   string `json:"query"`
			ListID  string `json:"list_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		q := url.Values{}
		q.Set("query", p.Query)
		if p.ListID != "" {
			q.Set("list_id", p.ListID)
		}
		return client.Get(ctx, "/search-members?"+q.Encode())
	})
}
