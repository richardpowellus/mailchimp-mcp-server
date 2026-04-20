package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterMergeFields(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_merge_fields",
		Description: "List merge fields for an audience.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/merge-fields", p.ListID), "merge_fields")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_merge_field",
		Description: "Create a new merge field for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"body":    {Type: "object", Description: "Merge field data (name, type, etc.)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/merge-fields", p.ListID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_merge_field",
		Description: "Get a specific merge field for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"list_id":        {Type: "string", Description: "The audience/list ID."},
				"merge_field_id": {Type: "string", Description: "The merge field ID."},
			},
			Required: []string{"account", "list_id", "merge_field_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			ListID       string `json:"list_id"`
			MergeFieldID string `json:"merge_field_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/merge-fields/%s", p.ListID, p.MergeFieldID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_merge_field",
		Description: "Update a merge field for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"list_id":        {Type: "string", Description: "The audience/list ID."},
				"merge_field_id": {Type: "string", Description: "The merge field ID."},
				"body":           {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "list_id", "merge_field_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string          `json:"account"`
			ListID       string          `json:"list_id"`
			MergeFieldID string          `json:"merge_field_id"`
			Body         json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/lists/%s/merge-fields/%s", p.ListID, p.MergeFieldID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_merge_field",
		Description: "Delete a merge field from an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":        {Type: "string", Description: "Account name."},
				"list_id":        {Type: "string", Description: "The audience/list ID."},
				"merge_field_id": {Type: "string", Description: "The merge field ID."},
			},
			Required: []string{"account", "list_id", "merge_field_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account      string `json:"account"`
			ListID       string `json:"list_id"`
			MergeFieldID string `json:"merge_field_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/lists/%s/merge-fields/%s", p.ListID, p.MergeFieldID))
	})
}
