package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterEcommerceProducts(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// --- Products ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_store_products",
		Description: "List all products for an e-commerce store.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/ecommerce/stores/%s/products", p.StoreID), "products")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_store_product",
		Description: "Add a product to an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":  {Type: "string", Description: "Account name."},
				"store_id": {Type: "string", Description: "The store ID."},
				"body":     {Type: "object", Description: "Product data (id, title, variants, etc.)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/products", p.StoreID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_store_product",
		Description: "Get a specific product from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
			},
			Required: []string{"account", "store_id", "product_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s", p.StoreID, p.ProductID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_store_product",
		Description: "Update a product in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"body":       {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "product_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			StoreID   string          `json:"store_id"`
			ProductID string          `json:"product_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s", p.StoreID, p.ProductID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_store_product",
		Description: "Delete a product from an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
			},
			Required: []string{"account", "store_id", "product_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s", p.StoreID, p.ProductID))
	})

	// --- Variants ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_product_variants",
		Description: "List all variants for a product in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
			},
			Required: []string{"account", "store_id", "product_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/variants", p.StoreID, p.ProductID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_product_variant",
		Description: "Add a variant to a product in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"body":       {Type: "object", Description: "Variant data (id, title, price, etc.)."},
			},
			Required: []string{"account", "store_id", "product_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			StoreID   string          `json:"store_id"`
			ProductID string          `json:"product_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/variants", p.StoreID, p.ProductID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_product_variant",
		Description: "Get a specific variant for a product.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"variant_id": {Type: "string", Description: "The variant ID."},
			},
			Required: []string{"account", "store_id", "product_id", "variant_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
			VariantID string `json:"variant_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/variants/%s", p.StoreID, p.ProductID, p.VariantID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_product_variant",
		Description: "Update a variant for a product.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"variant_id": {Type: "string", Description: "The variant ID."},
				"body":       {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "store_id", "product_id", "variant_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			StoreID   string          `json:"store_id"`
			ProductID string          `json:"product_id"`
			VariantID string          `json:"variant_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/variants/%s", p.StoreID, p.ProductID, p.VariantID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "upsert_product_variant",
		Description: "Add or update a variant for a product.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"variant_id": {Type: "string", Description: "The variant ID."},
				"body":       {Type: "object", Description: "Variant data."},
			},
			Required: []string{"account", "store_id", "product_id", "variant_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			StoreID   string          `json:"store_id"`
			ProductID string          `json:"product_id"`
			VariantID string          `json:"variant_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PutRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/variants/%s", p.StoreID, p.ProductID, p.VariantID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_product_variant",
		Description: "Delete a variant from a product.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"variant_id": {Type: "string", Description: "The variant ID."},
			},
			Required: []string{"account", "store_id", "product_id", "variant_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
			VariantID string `json:"variant_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/variants/%s", p.StoreID, p.ProductID, p.VariantID))
	})

	// --- Images ---

	s.RegisterTool(mcp.Tool{
		Name:        "list_product_images",
		Description: "List all images for a product in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
			},
			Required: []string{"account", "store_id", "product_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/images", p.StoreID, p.ProductID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_product_image",
		Description: "Add an image to a product in an e-commerce store.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"body":       {Type: "object", Description: "Image data (id, url, variant_ids)."},
			},
			Required: []string{"account", "store_id", "product_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			StoreID   string          `json:"store_id"`
			ProductID string          `json:"product_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/images", p.StoreID, p.ProductID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_product_image",
		Description: "Get a specific image for a product.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"image_id":   {Type: "string", Description: "The image ID."},
			},
			Required: []string{"account", "store_id", "product_id", "image_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
			ImageID   string `json:"image_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/images/%s", p.StoreID, p.ProductID, p.ImageID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_product_image",
		Description: "Delete an image from a product.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"store_id":   {Type: "string", Description: "The store ID."},
				"product_id": {Type: "string", Description: "The product ID."},
				"image_id":   {Type: "string", Description: "The image ID."},
			},
			Required: []string{"account", "store_id", "product_id", "image_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string `json:"account"`
			StoreID   string `json:"store_id"`
			ProductID string `json:"product_id"`
			ImageID   string `json:"image_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/ecommerce/stores/%s/products/%s/images/%s", p.StoreID, p.ProductID, p.ImageID))
	})
}
