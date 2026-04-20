package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterAudiences(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_audiences — GET /lists
	s.RegisterTool(mcp.Tool{
		Name:        "list_audiences",
		Description: "List all audiences (lists). GET /lists. Returns audience summaries with member counts, stats, and settings.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, "/lists", "lists")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// create_audience — POST /lists
	s.RegisterTool(mcp.Tool{
		Name:        "create_audience",
		Description: "Create a new audience (list). POST /lists. Body must include name, contact, permission_reminder, campaign_defaults, and email_type_option.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Audience creation payload: name, contact, permission_reminder, campaign_defaults, email_type_option."},
			},
			Required: []string{"account", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, "/lists", p.Body)
	})

	// get_audience — GET /lists/{list_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_audience",
		Description: "Get details of a single audience. GET /lists/{list_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
			},
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s", p.ListID))
	})

	// update_audience — PATCH /lists/{list_id}
	s.RegisterTool(mcp.Tool{
		Name:        "update_audience",
		Description: "Update an audience. PATCH /lists/{list_id}. Body can include name, contact, permission_reminder, campaign_defaults, etc.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"body":    {Type: "object", Description: "Fields to update on the audience."},
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
		return client.PatchRaw(ctx, fmt.Sprintf("/lists/%s", p.ListID), p.Body)
	})

	// delete_audience — DELETE /lists/{list_id}
	s.RegisterTool(mcp.Tool{
		Name:        "delete_audience",
		Description: "Delete an audience. DELETE /lists/{list_id}. This permanently deletes the list and all its members.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
			},
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/lists/%s", p.ListID))
	})

	// get_audience_growth — GET /lists/{list_id}/growth-history
	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_growth",
		Description: "Get growth history for an audience. GET /lists/{list_id}/growth-history. Returns monthly subscriber counts.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/growth-history", p.ListID), "history")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_audience_locations — GET /lists/{list_id}/locations
	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_locations",
		Description: "Get top locations for an audience by opens. GET /lists/{list_id}/locations.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/locations", p.ListID), "locations")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_email_client_stats — GET /lists/{list_id}/clients
	s.RegisterTool(mcp.Tool{
		Name:        "get_email_client_stats",
		Description: "Get email client usage stats for an audience. GET /lists/{list_id}/clients. Shows which email clients subscribers use.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/clients", p.ListID), "clients")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})
}
