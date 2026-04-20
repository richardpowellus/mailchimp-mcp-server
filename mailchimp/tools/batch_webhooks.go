package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterBatchWebhooks(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_batch_webhooks",
		Description: "List all batch webhooks.",
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
		items, err := client.FetchAll(ctx, "/batch-webhooks", "webhooks")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_batch_webhook",
		Description: "Create a new batch webhook.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Webhook data (url)."},
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
		return client.PostRaw(ctx, "/batch-webhooks", p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_batch_webhook",
		Description: "Get a specific batch webhook.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":          {Type: "string", Description: "Account name."},
				"batch_webhook_id": {Type: "string", Description: "The batch webhook ID."},
			},
			Required: []string{"account", "batch_webhook_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string `json:"account"`
			BatchWebhookID string `json:"batch_webhook_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/batch-webhooks/%s", p.BatchWebhookID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_batch_webhook",
		Description: "Update a batch webhook.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":          {Type: "string", Description: "Account name."},
				"batch_webhook_id": {Type: "string", Description: "The batch webhook ID."},
				"body":             {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "batch_webhook_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string          `json:"account"`
			BatchWebhookID string          `json:"batch_webhook_id"`
			Body           json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/batch-webhooks/%s", p.BatchWebhookID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_batch_webhook",
		Description: "Delete a batch webhook.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":          {Type: "string", Description: "Account name."},
				"batch_webhook_id": {Type: "string", Description: "The batch webhook ID."},
			},
			Required: []string{"account", "batch_webhook_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string `json:"account"`
			BatchWebhookID string `json:"batch_webhook_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/batch-webhooks/%s", p.BatchWebhookID))
	})
}
