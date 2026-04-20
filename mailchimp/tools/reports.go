package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterReports(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_reports — GET /reports
	s.RegisterTool(mcp.Tool{
		Name:        "list_reports",
		Description: "List campaign reports. GET /reports. Returns performance summaries for sent campaigns including opens, clicks, bounces, and unsubscribes.",
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
		items, err := client.FetchAll(ctx, "/reports", "reports")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_campaign_report — GET /reports/{campaign_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_report",
		Description: "Get the full report for a single campaign. GET /reports/{campaign_id}. Returns opens, clicks, bounces, unsubscribes, forwards, and more.",
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
		return client.Get(ctx, fmt.Sprintf("/reports/%s", p.CampaignID))
	})

	// get_click_details — GET /reports/{campaign_id}/click-details
	s.RegisterTool(mcp.Tool{
		Name:        "get_click_details",
		Description: "Get click details for a campaign. GET /reports/{campaign_id}/click-details. Returns all tracked URLs and their click counts.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/click-details", p.CampaignID), "urls_clicked")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_click_link_details — GET /reports/{campaign_id}/click-details/{link_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_click_link_details",
		Description: "Get click details for a specific URL in a campaign. GET /reports/{campaign_id}/click-details/{link_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"link_id":     {Type: "string", Description: "The URL link ID."},
			},
			Required: []string{"account", "campaign_id", "link_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			LinkID     string `json:"link_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reports/%s/click-details/%s", p.CampaignID, p.LinkID))
	})

	// get_click_link_members — GET /reports/{campaign_id}/click-details/{link_id}/members
	s.RegisterTool(mcp.Tool{
		Name:        "get_click_link_members",
		Description: "Get members who clicked a specific URL. GET /reports/{campaign_id}/click-details/{link_id}/members.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"link_id":     {Type: "string", Description: "The URL link ID."},
			}, paging.Properties()),
			Required: []string{"account", "campaign_id", "link_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			LinkID     string `json:"link_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/click-details/%s/members", p.CampaignID, p.LinkID), "members")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_click_link_member — GET /reports/{campaign_id}/click-details/{link_id}/members/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_click_link_member",
		Description: "Get click info for a specific member on a specific URL. GET /reports/{campaign_id}/click-details/{link_id}/members/{subscriber_hash}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"link_id":     {Type: "string", Description: "The URL link ID."},
				"email":       {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "campaign_id", "link_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			LinkID     string `json:"link_id"`
			Email      string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/reports/%s/click-details/%s/members/%s", p.CampaignID, p.LinkID, hash))
	})

	// get_open_details — GET /reports/{campaign_id}/open-details
	s.RegisterTool(mcp.Tool{
		Name:        "get_open_details",
		Description: "Get open details for a campaign. GET /reports/{campaign_id}/open-details. Lists members who opened the campaign.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/open-details", p.CampaignID), "members")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_open_details_subscriber — GET /reports/{campaign_id}/open-details/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_open_details_subscriber",
		Description: "Get open details for a specific subscriber. GET /reports/{campaign_id}/open-details/{subscriber_hash}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"email":       {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "campaign_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			Email      string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/reports/%s/open-details/%s", p.CampaignID, hash))
	})

	// get_email_activity — GET /reports/{campaign_id}/email-activity
	s.RegisterTool(mcp.Tool{
		Name:        "get_email_activity",
		Description: "Get email activity for a campaign. GET /reports/{campaign_id}/email-activity. Returns per-member activity (opens, clicks, bounces).",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/email-activity", p.CampaignID), "emails")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_subscriber_email_activity — GET /reports/{campaign_id}/email-activity/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_subscriber_email_activity",
		Description: "Get email activity for a specific subscriber in a campaign. GET /reports/{campaign_id}/email-activity/{subscriber_hash}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"email":       {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "campaign_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			Email      string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/reports/%s/email-activity/%s", p.CampaignID, hash))
	})

	// get_campaign_unsubscribes — GET /reports/{campaign_id}/unsubscribed
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_unsubscribes",
		Description: "Get members who unsubscribed from a campaign. GET /reports/{campaign_id}/unsubscribed.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/unsubscribed", p.CampaignID), "unsubscribes")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_campaign_unsubscribe — GET /reports/{campaign_id}/unsubscribed/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_unsubscribe",
		Description: "Get unsubscribe info for a specific member. GET /reports/{campaign_id}/unsubscribed/{subscriber_hash}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"email":       {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "campaign_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			Email      string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/reports/%s/unsubscribed/%s", p.CampaignID, hash))
	})

	// get_domain_performance — GET /reports/{campaign_id}/domain-performance
	s.RegisterTool(mcp.Tool{
		Name:        "get_domain_performance",
		Description: "Get domain performance stats for a campaign. GET /reports/{campaign_id}/domain-performance. Shows opens/clicks broken down by email domain.",
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
		return client.Get(ctx, fmt.Sprintf("/reports/%s/domain-performance", p.CampaignID))
	})

	// get_sent_to — GET /reports/{campaign_id}/sent-to
	s.RegisterTool(mcp.Tool{
		Name:        "get_sent_to",
		Description: "Get the list of members a campaign was sent to. GET /reports/{campaign_id}/sent-to.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/sent-to", p.CampaignID), "sent_to")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_sent_to_member — GET /reports/{campaign_id}/sent-to/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_sent_to_member",
		Description: "Get sent-to info for a specific member. GET /reports/{campaign_id}/sent-to/{subscriber_hash}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"email":       {Type: "string", Description: "Member email address (hashed internally)."},
			},
			Required: []string{"account", "campaign_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			Email      string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/reports/%s/sent-to/%s", p.CampaignID, hash))
	})

	// get_sub_reports — GET /reports/{campaign_id}/sub-reports
	s.RegisterTool(mcp.Tool{
		Name:        "get_sub_reports",
		Description: "Get sub-reports for a campaign (e.g., A/B test variants). GET /reports/{campaign_id}/sub-reports.",
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
		return client.Get(ctx, fmt.Sprintf("/reports/%s/sub-reports", p.CampaignID))
	})

	// get_ecommerce_product_activity — GET /reports/{campaign_id}/ecommerce-product-activity
	s.RegisterTool(mcp.Tool{
		Name:        "get_ecommerce_product_activity",
		Description: "Get e-commerce product activity for a campaign. GET /reports/{campaign_id}/ecommerce-product-activity. Shows products purchased via this campaign.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/ecommerce-product-activity", p.CampaignID), "products")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_campaign_advice — GET /reports/{campaign_id}/advice
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_advice",
		Description: "Get advice and feedback for a campaign. GET /reports/{campaign_id}/advice. Returns actionable suggestions to improve future campaigns.",
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
		return client.Get(ctx, fmt.Sprintf("/reports/%s/advice", p.CampaignID))
	})

	// get_campaign_abuse_reports — GET /reports/{campaign_id}/abuse-reports
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_abuse_reports",
		Description: "Get abuse reports for a campaign. GET /reports/{campaign_id}/abuse-reports.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/abuse-reports", p.CampaignID), "abuse_reports")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_campaign_abuse_report — GET /reports/{campaign_id}/abuse-reports/{report_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign_abuse_report",
		Description: "Get a specific abuse report. GET /reports/{campaign_id}/abuse-reports/{report_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"report_id":   {Type: "string", Description: "The abuse report ID."},
			},
			Required: []string{"account", "campaign_id", "report_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			CampaignID string `json:"campaign_id"`
			ReportID   string `json:"report_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reports/%s/abuse-reports/%s", p.CampaignID, p.ReportID))
	})

	// get_report_locations — GET /reports/{campaign_id}/locations
	s.RegisterTool(mcp.Tool{
		Name:        "get_report_locations",
		Description: "Get top open locations for a campaign. GET /reports/{campaign_id}/locations.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/reports/%s/locations", p.CampaignID), "locations")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_eepurl_activity — GET /reports/{campaign_id}/eepurl
	s.RegisterTool(mcp.Tool{
		Name:        "get_eepurl_activity",
		Description: "Get eepurl (short URL) sharing activity for a campaign. GET /reports/{campaign_id}/eepurl.",
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
		return client.Get(ctx, fmt.Sprintf("/reports/%s/eepurl", p.CampaignID))
	})
}
