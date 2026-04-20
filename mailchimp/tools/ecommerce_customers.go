package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterEcommerceCustomers(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_store_customers",
		Description: "List all customers for an e-commerce store.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/ecommerce/stores/%s/customers", p.StoreID), "customers")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_store_customer",
		Description: "Add a customer to an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"body":     {Type: "object", Description: "Customer data (id, email_address, opt_in_status, etc.)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/customers", p.StoreID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_store_customer",
		Description: "Get a specific customer from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"store_id":    {Type: "string", Description: "The store ID."},
				"customer_id": {Type: "string", Description: "The customer ID."},
			},
			Required: []string{"account", "store_id", "customer_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			StoreID    string `json:"store_id"`
			CustomerID string `json:"customer_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/customers/%s", p.StoreID, p.CustomerID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "upsert_store_customer",
		Description: "Add or update a customer in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"store_id":    {Type: "string", Description: "The store ID."},
				"customer_id": {Type: "string", Description: "The customer ID."},
				"body":        {Type: "object", Description: "Customer data."},
			},
			Required: []string{"account", "store_id", "customer_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			StoreID    string          `json:"store_id"`
			CustomerID string          `json:"customer_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PutRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/customers/%s", p.StoreID, p.CustomerID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_store_customer",
		Description: "Update a customer in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"store_id":    {Type: "string", Description: "The store ID."},
				"customer_id": {Type: "string", Description: "The customer ID."},
				"body":        {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "customer_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			StoreID    string          `json:"store_id"`
			CustomerID string          `json:"customer_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/customers/%s", p.StoreID, p.CustomerID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_store_customer",
		Description: "Delete a customer from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"store_id":    {Type: "string", Description: "The store ID."},
				"customer_id": {Type: "string", Description: "The customer ID."},
			},
			Required: []string{"account", "store_id", "customer_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			StoreID    string `json:"store_id"`
			CustomerID string `json:"customer_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/customers/%s", p.StoreID, p.CustomerID))
	})
}
