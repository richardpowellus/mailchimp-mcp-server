package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterLandingPages(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_landing_pages",
		Description: "List all landing pages.",
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
		items, err := client.FetchAll(ctx, "/landing-pages", "landing_pages")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_landing_page",
		Description: "Create a new landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Landing page data."},
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
		return client.PostRaw(ctx, "/landing-pages", p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_landing_page",
		Description: "Get a specific landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"page_id": {Type: "string", Description: "The landing page ID."},
			},
			Required: []string{"account", "page_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			PageID  string `json:"page_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/landing-pages/%s", p.PageID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_landing_page",
		Description: "Update a landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"page_id": {Type: "string", Description: "The landing page ID."},
				"body":    {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "page_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			PageID  string          `json:"page_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/landing-pages/%s", p.PageID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_landing_page",
		Description: "Delete a landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"page_id": {Type: "string", Description: "The landing page ID."},
			},
			Required: []string{"account", "page_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			PageID  string `json:"page_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/landing-pages/%s", p.PageID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "publish_landing_page",
		Description: "Publish a landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"page_id": {Type: "string", Description: "The landing page ID."},
			},
			Required: []string{"account", "page_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			PageID  string `json:"page_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/landing-pages/%s/actions/publish", p.PageID), nil)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "unpublish_landing_page",
		Description: "Unpublish a landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"page_id": {Type: "string", Description: "The landing page ID."},
			},
			Required: []string{"account", "page_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			PageID  string `json:"page_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/landing-pages/%s/actions/unpublish", p.PageID), nil)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_landing_page_content",
		Description: "Get the content of a landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"page_id": {Type: "string", Description: "The landing page ID."},
			},
			Required: []string{"account", "page_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			PageID  string `json:"page_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/landing-pages/%s/content", p.PageID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "set_landing_page_content",
		Description: "Set the content of a landing page.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"page_id": {Type: "string", Description: "The landing page ID."},
				"body":    {Type: "object", Description: "Landing page content."},
			},
			Required: []string{"account", "page_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			PageID  string          `json:"page_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PutRaw(ctx, fmt.Sprintf("/landing-pages/%s/content", p.PageID), p.Body)
	})
}
