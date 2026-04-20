package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterInterestCategories(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_interest_categories",
		Description: "List interest categories for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
			}, paging.Properties()),
			Required: []string{"account", "list_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/interest-categories", p.ListID), "categories")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_interest_category",
		Description: "Create an interest category for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"body":    {Type: "object", Description: "Category data (title, type, etc.)."},
			},
			Required: []string{"account", "list_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			ListID  string          `json:"list_id"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/interest-categories", p.ListID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_interest_category",
		Description: "Get a specific interest category.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
			},
			Required: []string{"account", "list_id", "category_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			ListID     string `json:"list_id"`
			CategoryID string `json:"category_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s", p.ListID, p.CategoryID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_interest_category",
		Description: "Update an interest category.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
				"body":        {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "list_id", "category_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			ListID     string          `json:"list_id"`
			CategoryID string          `json:"category_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s", p.ListID, p.CategoryID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_interest_category",
		Description: "Delete an interest category.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
			},
			Required: []string{"account", "list_id", "category_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			ListID     string `json:"list_id"`
			CategoryID string `json:"category_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s", p.ListID, p.CategoryID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_interests",
		Description: "List interests in an interest category.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
			}, paging.Properties()),
			Required: []string{"account", "list_id", "category_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			ListID     string `json:"list_id"`
			CategoryID string `json:"category_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s/interests", p.ListID, p.CategoryID), "interests")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_interest",
		Description: "Create an interest in a category.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
				"body":        {Type: "object", Description: "Interest data (name, etc.)."},
			},
			Required: []string{"account", "list_id", "category_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			ListID     string          `json:"list_id"`
			CategoryID string          `json:"category_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s/interests", p.ListID, p.CategoryID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_interest",
		Description: "Get a specific interest.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
				"interest_id": {Type: "string", Description: "The interest ID."},
			},
			Required: []string{"account", "list_id", "category_id", "interest_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			ListID     string `json:"list_id"`
			CategoryID string `json:"category_id"`
			InterestID string `json:"interest_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s/interests/%s", p.ListID, p.CategoryID, p.InterestID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_interest",
		Description: "Update an interest.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
				"interest_id": {Type: "string", Description: "The interest ID."},
				"body":        {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "list_id", "category_id", "interest_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			ListID     string          `json:"list_id"`
			CategoryID string          `json:"category_id"`
			InterestID string          `json:"interest_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s/interests/%s", p.ListID, p.CategoryID, p.InterestID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_interest",
		Description: "Delete an interest.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"list_id":     {Type: "string", Description: "The audience/list ID."},
				"category_id": {Type: "string", Description: "The interest category ID."},
				"interest_id": {Type: "string", Description: "The interest ID."},
			},
			Required: []string{"account", "list_id", "category_id", "interest_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			ListID     string `json:"list_id"`
			CategoryID string `json:"category_id"`
			InterestID string `json:"interest_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/lists/%s/interest-categories/%s/interests/%s", p.ListID, p.CategoryID, p.InterestID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_growth_history_month",
		Description: "Get growth history for an audience for a specific month.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"month":   {Type: "string", Description: "The month in YYYY-MM format."},
			},
			Required: []string{"account", "list_id", "month"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Month   string `json:"month"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/growth-history/%s", p.ListID, p.Month))
	})
}
