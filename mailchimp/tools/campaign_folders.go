package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterCampaignFolders(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_campaign_folders",
		Description: "List all campaign folders.",
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
		items, err := client.FetchAll(ctx, "/campaign-folders", "folders")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_campaign_folder",
		Description: "Create a new campaign folder.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Folder data (name)."},
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
		return client.PostRaw(ctx, "/campaign-folders", p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_folder",
		Description: "Get a specific campaign folder.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"folder_id": {Type: "string", Description: "The folder ID."},
			},
			Required: []string{"account", "folder_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			FolderID string `json:"folder_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/campaign-folders/%s", p.FolderID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_campaign_folder",
		Description: "Update a campaign folder.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"folder_id": {Type: "string", Description: "The folder ID."},
				"body":      {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "folder_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string          `json:"account"`
			FolderID string          `json:"folder_id"`
			Body     json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/campaign-folders/%s", p.FolderID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_campaign_folder",
		Description: "Delete a campaign folder.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"folder_id": {Type: "string", Description: "The folder ID."},
			},
			Required: []string{"account", "folder_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			FolderID string `json:"folder_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/campaign-folders/%s", p.FolderID))
	})
}
