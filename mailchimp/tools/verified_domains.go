package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterVerifiedDomains(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_verified_domains",
		Description: "List verified sending domains.",
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
		items, err := client.FetchAll(ctx, "/verified-domains", "domains")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "add_verified_domain",
		Description: "Add a domain to verify for sending.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Domain data (verification_email)."},
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
		return client.PostRaw(ctx, "/verified-domains", p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_verified_domain",
		Description: "Get a specific verified domain.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"domain_id": {Type: "string", Description: "The domain ID."},
			},
			Required: []string{"account", "domain_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			DomainID string `json:"domain_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/verified-domains/%s", p.DomainID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "verify_domain",
		Description: "Verify a domain with the verification code.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"domain_id": {Type: "string", Description: "The domain ID."},
				"body":      {Type: "object", Description: "Verification data (code)."},
			},
			Required: []string{"account", "domain_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string          `json:"account"`
			DomainID string          `json:"domain_id"`
			Body     json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/verified-domains/%s/actions/verify", p.DomainID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_verified_domain",
		Description: "Delete a verified domain.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"domain_id": {Type: "string", Description: "The domain ID."},
			},
			Required: []string{"account", "domain_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			DomainID string `json:"domain_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/verified-domains/%s", p.DomainID))
	})
}
