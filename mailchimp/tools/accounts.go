package tools

import (
	"context"
	"encoding/json"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterAccounts(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_accounts — no API call, returns configured accounts
	s.RegisterTool(mcp.Tool{
		Name:        "list_accounts",
		Description: "List all configured Mailchimp accounts. Returns account names and display names. Never exposes API keys.",
		InputSchema: mcp.InputSchema{Type: "object"},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		return cfg.Accounts(ctx), nil
	})

	// ping — GET /ping
	s.RegisterTool(mcp.Tool{
		Name:        "ping",
		Description: "Health check for a Mailchimp account. Calls GET /ping to verify the API key is valid and the account is reachable.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, "/ping")
	})

	// get_account_info — GET /
	s.RegisterTool(mcp.Tool{
		Name:        "get_account_info",
		Description: "Get full account details for a Mailchimp account. Calls GET / and returns account name, contact info, industry, plan, and stats.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
			},
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, "/")
	})
}
