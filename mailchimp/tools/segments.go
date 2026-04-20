package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterSegments(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_segments — GET /lists/{list_id}/segments
	s.RegisterTool(mcp.Tool{
		Name:        "list_segments",
		Description: "List segments and tags for an audience. GET /lists/{list_id}/segments. Returns saved segments, static segments (tags), and auto-update segments.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
			}, paging.Properties()),
			Required: []string{"account", "list_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/segments", p.ListID), "segments")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// create_segment — POST /lists/{list_id}/segments
	s.RegisterTool(mcp.Tool{
		Name:        "create_segment",
		Description: "Create a new segment or tag. POST /lists/{list_id}/segments. Body must include name. For saved segments, include options with conditions. For static segments, include static_segment with member emails.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"body":    {Type: "object", Description: "Segment payload: name (required), options (for saved), static_segment (for static/tags)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/segments", p.ListID), p.Body)
	})

	// get_segment — GET /lists/{list_id}/segments/{segment_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_segment",
		Description: "Get a specific segment. GET /lists/{list_id}/segments/{segment_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"list_id":    {Type: "string", Description: "The audience/list ID."},
				"segment_id": {Type: "string", Description: "The segment ID."},
			},
			Required: []string{"account", "list_id", "segment_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ListID    string `json:"list_id"`
			SegmentID string `json:"segment_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/segments/%s", p.ListID, p.SegmentID))
	})

	// update_segment — PATCH /lists/{list_id}/segments/{segment_id}
	s.RegisterTool(mcp.Tool{
		Name:        "update_segment",
		Description: "Update a segment. PATCH /lists/{list_id}/segments/{segment_id}. Body can include name, options (conditions for saved segments), or static_segment.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"list_id":    {Type: "string", Description: "The audience/list ID."},
				"segment_id": {Type: "string", Description: "The segment ID."},
				"body":       {Type: "object", Description: "Fields to update: name, options, static_segment."},
			},
			Required: []string{"account", "list_id", "segment_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			ListID    string          `json:"list_id"`
			SegmentID string          `json:"segment_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/lists/%s/segments/%s", p.ListID, p.SegmentID), p.Body)
	})

	// delete_segment — DELETE /lists/{list_id}/segments/{segment_id}
	s.RegisterTool(mcp.Tool{
		Name:        "delete_segment",
		Description: "Delete a segment. DELETE /lists/{list_id}/segments/{segment_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"list_id":    {Type: "string", Description: "The audience/list ID."},
				"segment_id": {Type: "string", Description: "The segment ID."},
			},
			Required: []string{"account", "list_id", "segment_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ListID    string `json:"list_id"`
			SegmentID string `json:"segment_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/lists/%s/segments/%s", p.ListID, p.SegmentID))
	})

	// batch_segment_members — POST /lists/{list_id}/segments/{segment_id}
	s.RegisterTool(mcp.Tool{
		Name:        "batch_segment_members",
		Description: "Batch add/remove members from a static segment. POST /lists/{list_id}/segments/{segment_id}. Body: members_to_add (array of emails), members_to_remove (array of emails).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"list_id":    {Type: "string", Description: "The audience/list ID."},
				"segment_id": {Type: "string", Description: "The segment ID."},
				"body":       {Type: "object", Description: "Batch payload: members_to_add (array of emails), members_to_remove (array of emails)."},
			},
			Required: []string{"account", "list_id", "segment_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			ListID    string          `json:"list_id"`
			SegmentID string          `json:"segment_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/segments/%s", p.ListID, p.SegmentID), p.Body)
	})

	// list_segment_members — GET /lists/{list_id}/segments/{segment_id}/members
	s.RegisterTool(mcp.Tool{
		Name:        "list_segment_members",
		Description: "List members in a segment. GET /lists/{list_id}/segments/{segment_id}/members.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"list_id":    {Type: "string", Description: "The audience/list ID."},
				"segment_id": {Type: "string", Description: "The segment ID."},
			}, paging.Properties()),
			Required: []string{"account", "list_id", "segment_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ListID    string `json:"list_id"`
			SegmentID string `json:"segment_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/segments/%s/members", p.ListID, p.SegmentID), "members")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// add_segment_member — POST /lists/{list_id}/segments/{segment_id}/members
	s.RegisterTool(mcp.Tool{
		Name:        "add_segment_member",
		Description: "Add a member to a static segment. POST /lists/{list_id}/segments/{segment_id}/members. Body: email_address.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"list_id":    {Type: "string", Description: "The audience/list ID."},
				"segment_id": {Type: "string", Description: "The segment ID."},
				"body":       {Type: "object", Description: "Payload with email_address of the member to add."},
			},
			Required: []string{"account", "list_id", "segment_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			ListID    string          `json:"list_id"`
			SegmentID string          `json:"segment_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/segments/%s/members", p.ListID, p.SegmentID), p.Body)
	})

	// remove_segment_member — DELETE /lists/{list_id}/segments/{segment_id}/members/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "remove_segment_member",
		Description: "Remove a member from a static segment. DELETE /lists/{list_id}/segments/{segment_id}/members/{subscriber_hash}. Accepts email and hashes internally.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"list_id":    {Type: "string", Description: "The audience/list ID."},
				"segment_id": {Type: "string", Description: "The segment ID."},
				"email":      {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "list_id", "segment_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			ListID    string `json:"list_id"`
			SegmentID string `json:"segment_id"`
			Email     string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Delete(ctx, fmt.Sprintf("/lists/%s/segments/%s/members/%s", p.ListID, p.SegmentID, hash))
	})
}
