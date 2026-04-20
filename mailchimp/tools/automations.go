package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/paging"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterAutomations(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// list_automations — GET /automations
	s.RegisterTool(mcp.Tool{
		Name:        "list_automations",
		Description: "List all automations (classic automation workflows). GET /automations. Returns automation summaries with status, trigger, and email count.",
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
		items, err := client.FetchAll(ctx, "/automations", "automations")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// create_automation — POST /automations
	s.RegisterTool(mcp.Tool{
		Name:        "create_automation",
		Description: "Create a new classic automation. POST /automations. Body must include recipients (list_id) and trigger_settings.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"body":    {Type: "object", Description: "Automation creation payload: recipients, trigger_settings, settings."},
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
		return client.PostRaw(ctx, "/automations", p.Body)
	})

	// get_automation — GET /automations/{workflow_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_automation",
		Description: "Get details of a single automation workflow. GET /automations/{workflow_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
			},
			Required: []string{"account", "workflow_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/automations/%s", p.WorkflowID))
	})

	// start_automation — POST /automations/{workflow_id}/actions/start-all-emails
	s.RegisterTool(mcp.Tool{
		Name:        "start_automation",
		Description: "Start all emails in an automation workflow. POST /automations/{workflow_id}/actions/start-all-emails.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
			},
			Required: []string{"account", "workflow_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/automations/%s/actions/start-all-emails", p.WorkflowID), nil)
	})

	// pause_automation — POST /automations/{workflow_id}/actions/pause-all-emails
	s.RegisterTool(mcp.Tool{
		Name:        "pause_automation",
		Description: "Pause all emails in an automation workflow. POST /automations/{workflow_id}/actions/pause-all-emails.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
			},
			Required: []string{"account", "workflow_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/automations/%s/actions/pause-all-emails", p.WorkflowID), nil)
	})

	// archive_automation — POST /automations/{workflow_id}/actions/archive
	s.RegisterTool(mcp.Tool{
		Name:        "archive_automation",
		Description: "Archive an automation workflow. POST /automations/{workflow_id}/actions/archive. Archived automations cannot be restarted.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
			},
			Required: []string{"account", "workflow_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/automations/%s/actions/archive", p.WorkflowID), nil)
	})

	// list_automation_emails — GET /automations/{workflow_id}/emails
	s.RegisterTool(mcp.Tool{
		Name:        "list_automation_emails",
		Description: "List the emails in an automation workflow. GET /automations/{workflow_id}/emails. Returns individual email steps with delay, status, and position.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
			}, paging.Properties()),
			Required: []string{"account", "workflow_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/automations/%s/emails", p.WorkflowID), "emails")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_automation_email — GET /automations/{workflow_id}/emails/{email_id}
	s.RegisterTool(mcp.Tool{
		Name:        "get_automation_email",
		Description: "Get a single email from an automation workflow. GET /automations/{workflow_id}/emails/{email_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
			},
			Required: []string{"account", "workflow_id", "email_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
			EmailID    string `json:"email_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/automations/%s/emails/%s", p.WorkflowID, p.EmailID))
	})

	// update_automation_email — PATCH /automations/{workflow_id}/emails/{email_id}
	s.RegisterTool(mcp.Tool{
		Name:        "update_automation_email",
		Description: "Update an automation email. PATCH /automations/{workflow_id}/emails/{email_id}. Body can include settings, delay, etc.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
				"body":        {Type: "object", Description: "Fields to update: settings, delay, etc."},
			},
			Required: []string{"account", "workflow_id", "email_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			WorkflowID string          `json:"workflow_id"`
			EmailID    string          `json:"email_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PatchRaw(ctx, fmt.Sprintf("/automations/%s/emails/%s", p.WorkflowID, p.EmailID), p.Body)
	})

	// delete_automation_email — DELETE /automations/{workflow_id}/emails/{email_id}
	s.RegisterTool(mcp.Tool{
		Name:        "delete_automation_email",
		Description: "Delete an automation email. DELETE /automations/{workflow_id}/emails/{email_id}.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
			},
			Required: []string{"account", "workflow_id", "email_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
			EmailID    string `json:"email_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Delete(ctx, fmt.Sprintf("/automations/%s/emails/%s", p.WorkflowID, p.EmailID))
	})

	// get_automation_email_content — GET /automations/{workflow_id}/emails/{email_id}/content
	s.RegisterTool(mcp.Tool{
		Name:        "get_automation_email_content",
		Description: "Get the HTML content of an automation email. GET /automations/{workflow_id}/emails/{email_id}/content.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
			},
			Required: []string{"account", "workflow_id", "email_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
			EmailID    string `json:"email_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Get(ctx, fmt.Sprintf("/automations/%s/emails/%s/content", p.WorkflowID, p.EmailID))
	})

	// set_automation_email_content — PUT /automations/{workflow_id}/emails/{email_id}/content
	s.RegisterTool(mcp.Tool{
		Name:        "set_automation_email_content",
		Description: "Set the content of an automation email. PUT /automations/{workflow_id}/emails/{email_id}/content. Body: html, plain_text.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
				"body":        {Type: "object", Description: "Content payload: html, plain_text."},
			},
			Required: []string{"account", "workflow_id", "email_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			WorkflowID string          `json:"workflow_id"`
			EmailID    string          `json:"email_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PutRaw(ctx, fmt.Sprintf("/automations/%s/emails/%s/content", p.WorkflowID, p.EmailID), p.Body)
	})

	// start_automation_email — POST /automations/{workflow_id}/emails/{email_id}/actions/start
	s.RegisterTool(mcp.Tool{
		Name:        "start_automation_email",
		Description: "Start a single automation email. POST /automations/{workflow_id}/emails/{email_id}/actions/start.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
			},
			Required: []string{"account", "workflow_id", "email_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
			EmailID    string `json:"email_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/automations/%s/emails/%s/actions/start", p.WorkflowID, p.EmailID), nil)
	})

	// pause_automation_email — POST /automations/{workflow_id}/emails/{email_id}/actions/pause
	s.RegisterTool(mcp.Tool{
		Name:        "pause_automation_email",
		Description: "Pause a single automation email. POST /automations/{workflow_id}/emails/{email_id}/actions/pause.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
			},
			Required: []string{"account", "workflow_id", "email_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
			EmailID    string `json:"email_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.Post(ctx, fmt.Sprintf("/automations/%s/emails/%s/actions/pause", p.WorkflowID, p.EmailID), nil)
	})

	// list_automation_queue — GET /automations/{workflow_id}/emails/{email_id}/queue
	s.RegisterTool(mcp.Tool{
		Name:        "list_automation_queue",
		Description: "List subscribers in the queue for an automation email. GET /automations/{workflow_id}/emails/{email_id}/queue.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: mcp.MergeProps(map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
			}, paging.Properties()),
			Required: []string{"account", "workflow_id", "email_id"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
			EmailID    string `json:"email_id"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		pp := paging.ParseParams(params)
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		items, err := client.FetchAll(ctx, fmt.Sprintf("/automations/%s/emails/%s/queue", p.WorkflowID, p.EmailID), "queue")
		if err != nil {
			return nil, err
		}
		return paging.EmitAny(items, pp), nil
	})

	// get_automation_queue_subscriber — GET /automations/{workflow_id}/emails/{email_id}/queue/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_automation_queue_subscriber",
		Description: "Get a specific subscriber in the automation email queue. GET /automations/{workflow_id}/emails/{email_id}/queue/{subscriber_hash}. Accepts email and hashes internally.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email_id":    {Type: "string", Description: "The automation email ID."},
				"email":       {Type: "string", Description: "Subscriber email address (hashed to subscriber_hash automatically)."},
			},
			Required: []string{"account", "workflow_id", "email_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
			EmailID    string `json:"email_id"`
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
		return client.Get(ctx, fmt.Sprintf("/automations/%s/emails/%s/queue/%s", p.WorkflowID, p.EmailID, hash))
	})

	// remove_automation_subscriber — POST /automations/{workflow_id}/removed-subscribers
	s.RegisterTool(mcp.Tool{
		Name:        "remove_automation_subscriber",
		Description: "Remove a subscriber from an automation workflow so they will no longer receive emails from it. POST /automations/{workflow_id}/removed-subscribers. Body: email_address.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"body":        {Type: "object", Description: "Payload with email_address to remove."},
			},
			Required: []string{"account", "workflow_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string          `json:"account"`
			WorkflowID string          `json:"workflow_id"`
			Body       json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/automations/%s/removed-subscribers", p.WorkflowID), p.Body)
	})

	// get_removed_subscriber — GET /automations/{workflow_id}/removed-subscribers/{subscriber_hash}
	s.RegisterTool(mcp.Tool{
		Name:        "get_removed_subscriber",
		Description: "Get info about a subscriber removed from an automation. GET /automations/{workflow_id}/removed-subscribers/{subscriber_hash}. Accepts email and hashes internally.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":     {Type: "string", Description: "Account name."},
				"workflow_id": {Type: "string", Description: "The automation workflow ID."},
				"email":       {Type: "string", Description: "Subscriber email address (hashed to subscriber_hash automatically)."},
			},
			Required: []string{"account", "workflow_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account    string `json:"account"`
			WorkflowID string `json:"workflow_id"`
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
		return client.Get(ctx, fmt.Sprintf("/automations/%s/removed-subscribers/%s", p.WorkflowID, hash))
	})
}
