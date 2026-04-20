package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterCampaignContent(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// get_campaign_content — GET /campaigns/{id}/content
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_content",
		Description: "Get the HTML and plain-text content of a campaign. GET /campaigns/{campaign_id}/content.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
			},
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/campaigns/%s/content", p.CampaignID))
	})

	// set_campaign_content — PUT /campaigns/{id}/content
	s.RegisterTool(mcp.Tool{
		Name:        "set_campaign_content",
		Description: "Set campaign content. PUT /campaigns/{campaign_id}/content. Body can include html, plain_text, url, template (with id and sections).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"body":        {Type: "object", Description: "Content payload: html, plain_text, url, or template object."},
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
		return client.PutRaw(ctx, fmt.Sprintf("/campaigns/%s/content", p.CampaignID), p.Body)
	})
}
