package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterSurveys(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "list_surveys",
		Description: "List all surveys.",
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
		return client.Get(ctx, "/reporting/surveys")
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_survey",
		Description: "Get a specific survey.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"survey_id": {Type: "string", Description: "The survey ID."},
			},
			Required: []string{"account", "survey_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			SurveyID string `json:"survey_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reporting/surveys/%s", p.SurveyID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_survey_questions",
		Description: "List all questions for a survey.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"survey_id": {Type: "string", Description: "The survey ID."},
			},
			Required: []string{"account", "survey_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			SurveyID string `json:"survey_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reporting/surveys/%s/questions", p.SurveyID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_survey_question",
		Description: "Get a specific question from a survey.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"survey_id":   {Type: "string", Description: "The survey ID."},
				"question_id": {Type: "string", Description: "The question ID."},
			},
			Required: []string{"account", "survey_id", "question_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			SurveyID   string `json:"survey_id"`
			QuestionID string `json:"question_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reporting/surveys/%s/questions/%s", p.SurveyID, p.QuestionID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_survey_question_answers",
		Description: "List all answers for a survey question.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"survey_id":   {Type: "string", Description: "The survey ID."},
				"question_id": {Type: "string", Description: "The question ID."},
			},
			Required: []string{"account", "survey_id", "question_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			SurveyID   string `json:"survey_id"`
			QuestionID string `json:"question_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reporting/surveys/%s/questions/%s/answers", p.SurveyID, p.QuestionID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "list_survey_responses",
		Description: "List all responses for a survey.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":   {Type: "string", Description: "Account name."},
				"survey_id": {Type: "string", Description: "The survey ID."},
			},
			Required: []string{"account", "survey_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account  string `json:"account"`
			SurveyID string `json:"survey_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reporting/surveys/%s/responses", p.SurveyID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "get_survey_response",
		Description: "Get a specific response from a survey.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"survey_id":   {Type: "string", Description: "The survey ID."},
				"response_id": {Type: "string", Description: "The response ID."},
			},
			Required: []string{"account", "survey_id", "response_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			SurveyID   string `json:"survey_id"`
			ResponseID string `json:"response_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/reporting/surveys/%s/responses/%s", p.SurveyID, p.ResponseID))
	})

	s.RegisterTool(mcp.Tool{
		Name:        "publish_survey",
		Description: "Publish a survey.",
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
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/surveys/%s/actions/publish", p.ListID, p.SurveyID), nil)
	})

	s.RegisterTool(mcp.Tool{
		Name:        "unpublish_survey",
		Description: "Unpublish a survey.",
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
		return client.PostRaw(ctx, fmt.Sprintf("/lists/%s/surveys/%s/actions/unpublish", p.ListID, p.SurveyID), nil)
	})
}
