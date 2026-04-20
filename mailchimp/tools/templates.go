package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterTemplates(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_templates — GET /templates
	s.RegisterTool(mcp.Tool{
		Name:        "list_templates",
		Description: "List all templates. GET /templates. Returns template summaries with name, type, and date info.",
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
		items, err := client.FetchAll(ctx, "/templates", "templates")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// create_template — POST /templates
	s.RegisterTool(mcp.Tool{
		Name:        "create_template",
		Description: "Create a new template. POST /templates. Body must include name and html.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Template payload: name (required), html (required), folder_id (optional)."},
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
		return client.PostRaw(ctx, "/templates", p.Body)
	})

	// get_template — GET /templates/{template_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_template",
		Description: "Get a single template by ID. GET /templates/{template_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"template_id": {Type: "string", Description: "The template ID."},
			},
			Required: []string{"account", "template_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			TemplateID string `json:"template_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/templates/%s", p.TemplateID))
	})

	// update_template — PATCH /templates/{template_id}
	s.RegisterTool(mcp.Tool{
		Name:        "update_template",
		Description: "Update a template. PATCH /templates/{template_id}. Body can include name, html, folder_id.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"template_id": {Type: "string", Description: "The template ID."},
				"body":        {Type: "object", Description: "Fields to update: name, html, folder_id."},
			},
			Required: []string{"account", "template_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			TemplateID string          `json:"template_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/templates/%s", p.TemplateID), p.Body)
	})

	// delete_template — DELETE /templates/{template_id}
	s.RegisterTool(mcp.Tool{
		Name:        "delete_template",
		Description: "Delete a template. DELETE /templates/{template_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"template_id": {Type: "string", Description: "The template ID."},
			},
			Required: []string{"account", "template_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			TemplateID string `json:"template_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/templates/%s", p.TemplateID))
	})

	// get_template_content — GET /templates/{template_id}/default-content
	s.RegisterTool(mcp.Tool{
		Name:        "get_template_content",
		Description: "Get the default content (HTML sections) of a template. GET /templates/{template_id}/default-content.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"template_id": {Type: "string", Description: "The template ID."},
			},
			Required: []string{"account", "template_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			TemplateID string `json:"template_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/templates/%s/default-content", p.TemplateID))
	})
}
