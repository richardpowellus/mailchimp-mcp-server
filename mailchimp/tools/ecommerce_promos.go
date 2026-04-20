package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterEcommercePromos(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// --- Promo rules ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_promo_rules",
		Description: "List all promo rules for an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
			}, paging.Properties()),
			Required: []string{"account", "store_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules", p.StoreID), "promo_rules")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_promo_rule",
		Description: "Create a promo rule for an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"body":     {Type: "object", Description: "Promo rule data (id, title, description, amount, type, target, etc.)."},
			},
			Required: []string{"account", "store_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules", p.StoreID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_promo_rule",
		Description: "Get a specific promo rule from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
			},
			Required: []string{"account", "store_id", "rule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			RuleID  string `json:"rule_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s", p.StoreID, p.RuleID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_promo_rule",
		Description: "Update a promo rule in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
				"body":     {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "rule_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			RuleID  string          `json:"rule_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s", p.StoreID, p.RuleID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_promo_rule",
		Description: "Delete a promo rule from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
			},
			Required: []string{"account", "store_id", "rule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			RuleID  string `json:"rule_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s", p.StoreID, p.RuleID))
	})

	// --- Promo codes ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_promo_codes",
		Description: "List all promo codes for a promo rule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
			},
			Required: []string{"account", "store_id", "rule_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			RuleID  string `json:"rule_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s/promo-codes", p.StoreID, p.RuleID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_promo_code",
		Description: "Create a promo code for a promo rule.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
				"body":     {Type: "object", Description: "Promo code data (id, code, redemption_url, etc.)."},
			},
			Required: []string{"account", "store_id", "rule_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			RuleID  string          `json:"rule_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s/promo-codes", p.StoreID, p.RuleID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_promo_code",
		Description: "Get a specific promo code.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
				"code_id":  {Type: "string", Description: "The promo code ID."},
			},
			Required: []string{"account", "store_id", "rule_id", "code_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			RuleID  string `json:"rule_id"`
			CodeID  string `json:"code_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s/promo-codes/%s", p.StoreID, p.RuleID, p.CodeID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_promo_code",
		Description: "Update a promo code.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
				"code_id":  {Type: "string", Description: "The promo code ID."},
				"body":     {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "rule_id", "code_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			RuleID  string          `json:"rule_id"`
			CodeID  string          `json:"code_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s/promo-codes/%s", p.StoreID, p.RuleID, p.CodeID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_promo_code",
		Description: "Delete a promo code.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"rule_id":  {Type: "string", Description: "The promo rule ID."},
				"code_id":  {Type: "string", Description: "The promo code ID."},
			},
			Required: []string{"account", "store_id", "rule_id", "code_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			RuleID  string `json:"rule_id"`
			CodeID  string `json:"code_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/promo-rules/%s/promo-codes/%s", p.StoreID, p.RuleID, p.CodeID))
	})
}
