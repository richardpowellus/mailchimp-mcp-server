package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
)

func RegisterMemberGoals(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	s.RegisterTool(mcp.Tool{
		Name:        "get_member_goals",
		Description: "Get goal events for a list member by email address.",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"account": {Type: "string", Description: "Account name."},
				"list_id": {Type: "string", Description: "The audience/list ID."},
				"email":   {Type: "string", Description: "The member's email address."},
			},
			Required: []string{"account", "list_id", "email"},
		},
	}, func(ctx context.Context, params json.RawMessage) (any, error) {
		var p struct {
			Account string `json:"account"`
			ListID  string `json:"list_id"`
			Email   string `json:"email"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		client, err := cfg.GetClient(ctx, p.Account)
		if err != nil {
			return nil, err
		}
		hash := mailchimp.SubscriberHash(p.Email)
		return client.Get(ctx, fmt.Sprintf("/lists/%s/members/%s/goals", p.ListID, hash))
	})
}
