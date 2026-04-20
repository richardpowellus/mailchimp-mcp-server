package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterAudienceStats(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_abuse_reports",
		Description: "List abuse reports for an audience.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/abuse-reports", p.ListID), "abuse_reports")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_abuse_report",
		Description: "Get a specific abuse report for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"list_id":   {Type: "string", Description: "The audience/list ID."},
				"report_id": {Type: "string", Description: "The abuse report ID."},
			},
			Required: []string{"account", "list_id", "report_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			ListID   string `json:"list_id"`
			ReportID string `json:"report_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/abuse-reports/%s", p.ListID, p.ReportID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_activity",
		Description: "Get recent activity for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
			},
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
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/activity", p.ListID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_signup_forms",
		Description: "List signup forms for an audience.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/signup-forms", p.ListID), "signup_forms")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "create_signup_form",
		Description: "Create a signup form for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"body":    {Type: "object", Description: "Signup form configuration."},
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
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/signup-forms", p.ListID), p.Body)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_audience_surveys",
		Description: "List surveys for an audience.",
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
		items, err := client.FetchAll(ctx, fmt.Sprintf("/lists/%s/surveys", p.ListID), "surveys")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_audience_survey",
		Description: "Get a specific survey for an audience.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"list_id":   {Type: "string", Description: "The audience/list ID."},
				"survey_id": {Type: "string", Description: "The survey ID."},
			},
			Required: []string{"account", "list_id", "survey_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			ListID   string `json:"list_id"`
			SurveyID string `json:"survey_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/lists/%s/surveys/%s", p.ListID, p.SurveyID))
	})
}
