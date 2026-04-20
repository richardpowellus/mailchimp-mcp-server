package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterCampaigns(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_campaigns — GET /campaigns
	s.RegisterTool(mcp.Tool{
		Name:        "list_campaigns",
		Description: "List campaigns with optional filters. GET /campaigns. Supports filtering by status, type, send time range, folder, and sorting.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":          {Type: "string", Description: "Account name."},
				"status":           {Type: "string", Description: "Filter by status: save, paused, schedule, sending, sent."},
				"type":             {Type: "string", Description: "Filter by type: regular, plaintext, absplit, rss, variate."},
				"before_send_time": {Type: "string", Description: "Filter campaigns sent before this ISO 8601 date."},
				"since_send_time":  {Type: "string", Description: "Filter campaigns sent after this ISO 8601 date."},
				"folder_id":        {Type: "string", Description: "Filter by campaign folder ID."},
				"sort_field":       {Type: "string", Description: "Sort field: create_time, send_time."},
				"sort_dir":         {Type: "string", Description: "Sort direction: ASC, DESC."},
			}, paging.Properties()),
			Required: []string{"account"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account        string `json:"account"`
			Status         string `json:"status"`
			Type           string `json:"type"`
			BeforeSendTime string `json:"before_send_time"`
			SinceSendTime  string `json:"since_send_time"`
			FolderID       string `json:"folder_id"`
			SortField      string `json:"sort_field"`
			SortDir        string `json:"sort_dir"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		q := url.Values{}
		if p.Status != "" {
			q.Set("status", p.Status)
		}
		if p.Type != "" {
			q.Set("type", p.Type)
		}
		if p.BeforeSendTime != "" {
			q.Set("before_send_time", p.BeforeSendTime)
		}
		if p.SinceSendTime != "" {
			q.Set("since_send_time", p.SinceSendTime)
		}
		if p.FolderID != "" {
			q.Set("folder_id", p.FolderID)
		}
		if p.SortField != "" {
			q.Set("sort_field", p.SortField)
		}
		if p.SortDir != "" {
			q.Set("sort_dir", p.SortDir)
		}
		endpoint := "/campaigns"
		if len(q) > 0 {
			endpoint += "?" + q.Encode()
		}
		items, err := client.FetchAll(ctx, endpoint, "campaigns")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// create_campaign — POST /campaigns
	s.RegisterTool(mcp.Tool{
		Name:        "create_campaign",
		Description: "Create a new campaign. POST /campaigns. Body must include type (regular, plaintext, absplit, rss, variate) and typically recipients and settings.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Campaign creation payload: type, recipients, settings, etc."},
			},
			Required: []string{"account", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string          `json:"account"`
			Body    json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, "/campaigns", p.Body)
	})

	// get_campaign — GET /campaigns/{id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_campaign",
		Description: "Get details of a single campaign by ID. GET /campaigns/{campaign_id}.",
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
		return client.Get(ctx, fmt.Sprintf("/campaigns/%s", p.CampaignID))
	})

	// update_campaign — PATCH /campaigns/{id}
	s.RegisterTool(mcp.Tool{
		Name:        "update_campaign",
		Description: "Update a campaign. PATCH /campaigns/{campaign_id}. Body can include recipients, settings, variate_settings, tracking, rss_opts, social_card.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"body":        {Type: "object", Description: "Fields to update on the campaign."},
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
		return client.PatchRaw(ctx, fmt.Sprintf("/campaigns/%s", p.CampaignID), p.Body)
	})

	// delete_campaign — DELETE /campaigns/{id}
	s.RegisterTool(mcp.Tool{
		Name:        "delete_campaign",
		Description: "Delete a campaign. DELETE /campaigns/{campaign_id}. Only works on campaigns that have not been sent.",
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
		return client.Delete(ctx, fmt.Sprintf("/campaigns/%s", p.CampaignID))
	})

	// send_campaign — POST /campaigns/{id}/actions/send
	s.RegisterTool(mcp.Tool{
		Name:        "send_campaign",
		Description: "Send a campaign immediately. POST /campaigns/{campaign_id}/actions/send. Campaign must be in 'save' status with valid content.",
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
		return client.Post(ctx, fmt.Sprintf("/campaigns/%s/actions/send", p.CampaignID), nil)
	})

	// schedule_campaign — POST /campaigns/{id}/actions/schedule
	s.RegisterTool(mcp.Tool{
		Name:        "schedule_campaign",
		Description: "Schedule a campaign for future delivery. POST /campaigns/{campaign_id}/actions/schedule. Body must include schedule_time (ISO 8601). Optional: timewarp, batch_delivery.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"body":        {Type: "object", Description: "Schedule payload: schedule_time (required), timewarp (bool), batch_delivery (object)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/campaigns/%s/actions/schedule", p.CampaignID), p.Body)
	})

	// unschedule_campaign — POST /campaigns/{id}/actions/unschedule
	s.RegisterTool(mcp.Tool{
		Name:        "unschedule_campaign",
		Description: "Unschedule a scheduled campaign, returning it to 'save' status. POST /campaigns/{campaign_id}/actions/unschedule.",
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
		return client.Post(ctx, fmt.Sprintf("/campaigns/%s/actions/unschedule", p.CampaignID), nil)
	})

	// cancel_campaign — POST /campaigns/{id}/actions/cancel-send
	s.RegisterTool(mcp.Tool{
		Name:        "cancel_campaign",
		Description: "Cancel a campaign that is currently sending. POST /campaigns/{campaign_id}/actions/cancel-send.",
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
		return client.Post(ctx, fmt.Sprintf("/campaigns/%s/actions/cancel-send", p.CampaignID), nil)
	})

	// send_test_email — POST /campaigns/{id}/actions/test
	s.RegisterTool(mcp.Tool{
		Name:        "send_test_email",
		Description: "Send a test email for a campaign. POST /campaigns/{campaign_id}/actions/test. Body: test_emails (array of addresses), send_type (html or plaintext).",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID."},
				"body":        {Type: "object", Description: "Test email payload: test_emails (array), send_type (html/plaintext)."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/campaigns/%s/actions/test", p.CampaignID), p.Body)
	})

	// replicate_campaign — POST /campaigns/{id}/actions/replicate
	s.RegisterTool(mcp.Tool{
		Name:        "replicate_campaign",
		Description: "Replicate (copy) a campaign. POST /campaigns/{campaign_id}/actions/replicate. Returns the new campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The campaign ID to replicate."},
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
		return client.Post(ctx, fmt.Sprintf("/campaigns/%s/actions/replicate", p.CampaignID), nil)
	})

	// resend_campaign — POST /campaigns/{id}/actions/create-resend
	s.RegisterTool(mcp.Tool{
		Name:        "resend_campaign",
		Description: "Create a resend of a sent campaign to non-openers. POST /campaigns/{campaign_id}/actions/create-resend. Returns the new campaign.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The sent campaign ID to resend."},
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
		return client.Post(ctx, fmt.Sprintf("/campaigns/%s/actions/create-resend", p.CampaignID), nil)
	})

	// get_send_checklist — GET /campaigns/{id}/send-checklist
	s.RegisterTool(mcp.Tool{
		Name:        "get_send_checklist",
		Description: "Get the send checklist for a campaign. GET /campaigns/{campaign_id}/send-checklist. Returns items that must be resolved before sending.",
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
		return client.Get(ctx, fmt.Sprintf("/campaigns/%s/send-checklist", p.CampaignID))
	})

	// search_campaigns — GET /search-campaigns?query={query}
	s.RegisterTool(mcp.Tool{
		Name:        "search_campaigns",
		Description: "Search campaigns by query string. GET /search-campaigns?query={query}. Searches campaign titles and content.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"query":   {Type: "string", Description: "Search query string."},
			},
			Required: []string{"account", "query"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			Query   string `json:"query"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/search-campaigns?query=%s", url.QueryEscape(p.Query)))
	})

	// pause_rss_campaign — POST /campaigns/{id}/actions/pause
	s.RegisterTool(mcp.Tool{
		Name:        "pause_rss_campaign",
		Description: "Pause an RSS-driven campaign. POST /campaigns/{campaign_id}/actions/pause.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The RSS campaign ID."},
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
		return client.Post(ctx, fmt.Sprintf("/campaigns/%s/actions/pause", p.CampaignID), nil)
	})

	// resume_rss_campaign — POST /campaigns/{id}/actions/resume
	s.RegisterTool(mcp.Tool{
		Name:        "resume_rss_campaign",
		Description: "Resume a paused RSS-driven campaign. POST /campaigns/{campaign_id}/actions/resume.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"campaign_id": {Type: "string", Description: "The RSS campaign ID."},
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
		return client.Post(ctx, fmt.Sprintf("/campaigns/%s/actions/resume", p.CampaignID), nil)
	})
}
