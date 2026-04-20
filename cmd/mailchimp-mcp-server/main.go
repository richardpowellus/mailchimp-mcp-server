package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/richardpowellus/mailchimp-mcp-server/internal/tz"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/version"
	"github.com/richardpowellus/mailchimp-mcp-server/internal/watchdog"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp/tools"
	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
)

var instructions = `Mailchimp MCP Server — Email marketing and automation platform.

Unofficial — not affiliated with or endorsed by Mailchimp or Intuit.

This server provides full read+write access to the Mailchimp Marketing API v3
with 280+ tools covering every endpoint.

Key capabilities:
- Campaigns: Create, edit, schedule, send, and analyze email campaigns
- Automations: Manage automated workflows (abandoned cart, welcome series, etc.)
- Audiences: Manage subscriber lists, segments, tags, and merge fields
- Members: Add, update, tag, and track individual subscribers
- Templates: Create and manage reusable email templates
- Reports: Campaign/automation analytics (opens, clicks, revenue)
- E-commerce: Stores, products, orders, carts, promo codes
- Landing Pages: Create and publish landing pages
- Customer Journeys: Trigger automation flow steps

Configuration:
Set MAILCHIMP_API_KEY environment variable with your API key (format: key-datacenter, e.g. abc123-us2).
` + tz.Suffix()

func main() {
	log.SetOutput(os.Stderr)

	ctx := context.Background()
	ctx, cancel := watchdog.Start(ctx)
	defer cancel()

	cfg := mailchimp.NewEnvConfig()

	server := mcp.New("Mailchimp", version.Version, instructions)

	tools.RegisterAll(server, cfg)

	if err := server.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
