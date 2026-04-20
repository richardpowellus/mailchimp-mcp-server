package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterEcommerceCarts(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// --- Carts ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_store_carts",
		Description: "List all carts for an e-commerce store.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts", p.StoreID), "carts")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_store_cart",
		Description: "Add a cart to an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"body":     {Type: "object", Description: "Cart data (id, customer, lines, currency_code, order_total, etc.)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts", p.StoreID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_store_cart",
		Description: "Get a specific cart from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
			},
			Required: []string{"account", "store_id", "cart_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			CartID  string `json:"cart_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s", p.StoreID, p.CartID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_store_cart",
		Description: "Update a cart in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
				"body":     {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "cart_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			CartID  string          `json:"cart_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s", p.StoreID, p.CartID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_store_cart",
		Description: "Delete a cart from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
			},
			Required: []string{"account", "store_id", "cart_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			CartID  string `json:"cart_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s", p.StoreID, p.CartID))
	})

	// --- Cart lines ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_cart_lines",
		Description: "List all line items in a cart.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
			},
			Required: []string{"account", "store_id", "cart_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			CartID  string `json:"cart_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s/lines", p.StoreID, p.CartID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_cart_line",
		Description: "Add a line item to a cart.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
				"body":     {Type: "object", Description: "Line item data (id, product_id, product_variant_id, quantity, price)."},
			},
			Required: []string{"account", "store_id", "cart_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			CartID  string          `json:"cart_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s/lines", p.StoreID, p.CartID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_cart_line",
		Description: "Get a specific line item from a cart.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
				"line_id":  {Type: "string", Description: "The line item ID."},
			},
			Required: []string{"account", "store_id", "cart_id", "line_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			CartID  string `json:"cart_id"`
			LineID  string `json:"line_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s/lines/%s", p.StoreID, p.CartID, p.LineID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_cart_line",
		Description: "Update a line item in a cart.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
				"line_id":  {Type: "string", Description: "The line item ID."},
				"body":     {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "cart_id", "line_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			CartID  string          `json:"cart_id"`
			LineID  string          `json:"line_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s/lines/%s", p.StoreID, p.CartID, p.LineID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_cart_line",
		Description: "Delete a line item from a cart.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"cart_id":  {Type: "string", Description: "The cart ID."},
				"line_id":  {Type: "string", Description: "The line item ID."},
			},
			Required: []string{"account", "store_id", "cart_id", "line_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			CartID  string `json:"cart_id"`
			LineID  string `json:"line_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/carts/%s/lines/%s", p.StoreID, p.CartID, p.LineID))
	})
}
