package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterCampaignFeedback(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_campaign_feedback",
		Description: "List feedback for a campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
			}, paging.Properties()),
			Required: []string{"account", "campaign_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/campaigns/%s/feedback", p.CampaignID), "feedback")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_feedback",
		Description: "Get a specific feedback message for a campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"feedback_id": {Type: "string", Description: "The feedback ID."},
			},
			Required: []string{"account", "campaign_id", "feedback_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			FeedbackID string `json:"feedback_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/campaigns/%s/feedback/%s", p.CampaignID, p.FeedbackID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "add_campaign_feedback",
		Description: "Add feedback to a campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"body":        {Type: "object", Description: "Feedback data (message, etc.)."},
			},
			Required: []string{"account", "campaign_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			CampaignID string          `json:"campaign_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/campaigns/%s/feedback", p.CampaignID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "update_campaign_feedback",
		Description: "Update a specific feedback message for a campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"feedback_id": {Type: "string", Description: "The feedback ID."},
				"body":        {Type: "object", Description: "Fields to update."},
			},
			Required: []string{"account", "campaign_id", "feedback_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			CampaignID string          `json:"campaign_id"`
			FeedbackID string          `json:"feedback_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/campaigns/%s/feedback/%s", p.CampaignID, p.FeedbackID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "delete_campaign_feedback",
		Description: "Delete a specific feedback message for a campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"feedback_id": {Type: "string", Description: "The feedback ID."},
			},
			Required: []string{"account", "campaign_id", "feedback_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			FeedbackID string `json:"feedback_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/campaigns/%s/feedback/%s", p.CampaignID, p.FeedbackID))
	})
}
