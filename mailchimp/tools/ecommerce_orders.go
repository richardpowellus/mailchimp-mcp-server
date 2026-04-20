package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterEcommerceOrders(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// --- All orders (cross-store) ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_all_orders",
		Description: "List all e-commerce orders across all stores.",
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
		items, err := client.FetchAll(ctx, "/ecommerce/orders", "orders")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// --- Store-scoped orders ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_store_orders",
		Description: "List all orders for an e-commerce store.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders", p.StoreID), "orders")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_store_order",
		Description: "Add an order to an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"body":     {Type: "object", Description: "Order data (id, customer, lines, currency_code, order_total, etc.)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders", p.StoreID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_store_order",
		Description: "Get a specific order from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
			},
			Required: []string{"account", "store_id", "order_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			OrderID string `json:"order_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s", p.StoreID, p.OrderID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_store_order",
		Description: "Update an order in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
				"body":     {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "order_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			OrderID string          `json:"order_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s", p.StoreID, p.OrderID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_store_order",
		Description: "Delete an order from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
			},
			Required: []string{"account", "store_id", "order_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			OrderID string `json:"order_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s", p.StoreID, p.OrderID))
	})

	// --- Order lines ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_order_lines",
		Description: "List all line items for an order.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
			},
			Required: []string{"account", "store_id", "order_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			OrderID string `json:"order_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s/lines", p.StoreID, p.OrderID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_order_line",
		Description: "Add a line item to an order.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
				"body":     {Type: "object", Description: "Line item data (id, product_id, product_variant_id, quantity, price)."},
			},
			Required: []string{"account", "store_id", "order_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			OrderID string          `json:"order_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s/lines", p.StoreID, p.OrderID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_order_line",
		Description: "Get a specific line item from an order.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
				"line_id":  {Type: "string", Description: "The line item ID."},
			},
			Required: []string{"account", "store_id", "order_id", "line_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			OrderID string `json:"order_id"`
			LineID  string `json:"line_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s/lines/%s", p.StoreID, p.OrderID, p.LineID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_order_line",
		Description: "Update a line item in an order.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
				"line_id":  {Type: "string", Description: "The line item ID."},
				"body":     {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "order_id", "line_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			StoreID string          `json:"store_id"`
			OrderID string          `json:"order_id"`
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
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s/lines/%s", p.StoreID, p.OrderID, p.LineID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_order_line",
		Description: "Delete a line item from an order.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"order_id": {Type: "string", Description: "The order ID."},
				"line_id":  {Type: "string", Description: "The line item ID."},
			},
			Required: []string{"account", "store_id", "order_id", "line_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			StoreID string `json:"store_id"`
			OrderID string `json:"order_id"`
			LineID  string `json:"line_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/orders/%s/lines/%s", p.StoreID, p.OrderID, p.LineID))
	})
}
