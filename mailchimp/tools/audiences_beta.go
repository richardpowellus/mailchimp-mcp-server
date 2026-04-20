package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterAudiencesBeta(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_audiences_beta",
		Description: "List all audiences (beta API).",
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
		return client.Get(ctx, "/audiences")
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_beta",
		Description: "Get a specific audience (beta API).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The audience ID."},
			},
			Required: []string{"account", "audience_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/audiences/%s", p.AudienceID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_audience_contacts",
		Description: "List contacts in an audience (beta API).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The audience ID."},
			},
			Required: []string{"account", "audience_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/audiences/%s/contacts", p.AudienceID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_audience_contact",
		Description: "Add a contact to an audience (beta API).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The audience ID."},
				"body":        {Type: "object", Description: "Contact data (email_address, status, merge_fields, etc.)."},
			},
			Required: []string{"account", "audience_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			AudienceID string          `json:"audience_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/audiences/%s/contacts", p.AudienceID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_contact",
		Description: "Get a specific contact from an audience (beta API).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The audience ID."},
				"contact_id":  {Type: "string", Description: "The contact ID."},
			},
			Required: []string{"account", "audience_id", "contact_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
			ContactID  string `json:"contact_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/audiences/%s/contacts/%s", p.AudienceID, p.ContactID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_audience_contact",
		Description: "Update a contact in an audience (beta API).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The audience ID."},
				"contact_id":  {Type: "string", Description: "The contact ID."},
				"body":        {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "audience_id", "contact_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			AudienceID string          `json:"audience_id"`
			ContactID  string          `json:"contact_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/audiences/%s/contacts/%s", p.AudienceID, p.ContactID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "archive_audience_contact",
		Description: "Archive a contact in an audience (beta API).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The audience ID."},
				"contact_id":  {Type: "string", Description: "The contact ID."},
			},
			Required: []string{"account", "audience_id", "contact_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			AudienceID string `json:"audience_id"`
			ContactID  string `json:"contact_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/audiences/%s/contacts/%s/actions/archive", p.AudienceID, p.ContactID), nil)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "forget_audience_contact",
		Description: "Permanently delete (forget) a contact from an audience (beta API). Requires email_address in body.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"audience_id": {Type: "string", Description: "The audience ID."},
				"contact_id":  {Type: "string", Description: "The contact ID."},
				"body":        {Type: "object", Description: "Must include email_address."},
			},
			Required: []string{"account", "audience_id", "contact_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			AudienceID string          `json:"audience_id"`
			ContactID  string          `json:"contact_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/audiences/%s/contacts/%s/actions/forget", p.AudienceID, p.ContactID), p.Body)
	})
}
