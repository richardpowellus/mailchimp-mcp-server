package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterCustomerJourneys(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	// trigger_journey_step — POST /customer-journeys/journeys/{journey_id}/steps/{step_id}/actions/trigger
	s.RegisterTool(mcp.Tool{
		Name:        "trigger_journey_step",
		Description: "Trigger a customer journey step for a contact. POST /customer-journeys/journeys/{journey_id}/steps/{step_id}/actions/trigger. Body must include email_address.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account":    {Type: "string", Description: "Account name."},
				"journey_id": {Type: "string", Description: "The customer journey ID."},
				"step_id":    {Type: "string", Description: "The journey step ID."},
				"body":       {Type: "object", Description: "Payload with email_address of the contact to trigger."},
			},
			Required: []string{"account", "journey_id", "step_id", "body"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account   string          `json:"account"`
			JourneyID string          `json:"journey_id"`
			StepID    string          `json:"step_id"`
			Body      json.RawMessage `json:"body"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		return client.PostRaw(ctx, fmt.Sprintf("/customer-journeys/journeys/%s/steps/%s/actions/trigger", p.JourneyID, p.StepID), p.Body)
	})
}
