package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterConnectedSites(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_connected_sites",
		Description: "List connected e-commerce sites.",
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
		items, err := client.FetchAll(ctx, "/connected-sites", "sites")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "add_connected_site",
		Description: "Add a connected site.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Site data (foreign_id, domain)."},
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
		return client.PostRaw(ctx, "/connected-sites", p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_connected_site",
		Description: "Get a specific connected site.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"site_id": {Type: "string", Description: "The connected site ID."},
			},
			Required: []string{"account", "site_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			SiteID  string `json:"site_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/connected-sites/%s", p.SiteID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_connected_site",
		Description: "Delete a connected site.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"site_id": {Type: "string", Description: "The connected site ID."},
			},
			Required: []string{"account", "site_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			SiteID  string `json:"site_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/connected-sites/%s", p.SiteID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "verify_connected_site",
		Description: "Verify a connected site's script installation.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"site_id": {Type: "string", Description: "The connected site ID."},
			},
			Required: []string{"account", "site_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			SiteID  string `json:"site_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/connected-sites/%s/actions/verify-script-installation", p.SiteID), nil)
	})
}
