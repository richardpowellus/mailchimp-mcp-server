# Mailchimp MCP Server

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green)](https://modelcontextprotocol.io)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-support-yellow?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/letri)

A Go [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server providing full read+write access to the [Mailchimp Marketing API v3](https://mailchimp.com/developer/marketing/api/) — **279 tools** covering every endpoint.

> [!NOTE]
> **Unofficial** — This project is not affiliated with, endorsed by, or sponsored by Mailchimp or Intuit.

## Features

- **279 tools** — full Mailchimp Marketing API coverage, not just analytics
- **Read + write** — create campaigns, manage subscribers, send emails, set up automations
- **Campaigns** — create, schedule, send, test, replicate, cancel, pause/resume RSS campaigns
- **Automations** — manage workflows, emails, queues, start/pause automation steps
- **Audiences** — full subscriber lifecycle: add, update, tag, archive, search, segment
- **E-commerce** — stores, products, variants, orders, carts, promo codes
- **Reports** — opens, clicks, unsubscribes, revenue, domain performance, locations
- **Templates, landing pages, surveys, file manager, webhooks** — everything else
- **Smart retry** — automatic backoff with `Retry-After` header support for rate limiting
- **Clean responses** — `_links` stripped from all responses to save tokens
- **10-connection concurrency** — respects Mailchimp's simultaneous connection limit
- **Extensible** — plug in custom credential backends via the `CredentialProvider` interface

## Quick Start

```bash
# Install
go install github.com/richardpowellus/mailchimp-mcp-server/cmd/mailchimp-mcp-server@latest

# Run
export MAILCHIMP_API_KEY="your-api-key-us2"
mailchimp-mcp-server
```

The server communicates over stdio using JSON-RPC 2.0 — connect it to any MCP-compatible client.

## Installation

### Option 1: Go Install (recommended)

```bash
go install github.com/richardpowellus/mailchimp-mcp-server/cmd/mailchimp-mcp-server@latest
```

Requires Go 1.23 or later.

### Option 2: Docker

```bash
docker run -i --rm \
  -e MAILCHIMP_API_KEY="your-api-key-us2" \
  ghcr.io/richardpowellus/mailchimp-mcp-server
```

### Option 3: Binary Download

Download prebuilt binaries for Linux, macOS, and Windows from [GitHub Releases](https://github.com/richardpowellus/mailchimp-mcp-server/releases).

## MCP Client Configuration

### Claude Desktop / VS Code / GitHub Copilot

Add to your MCP configuration file:

```json
{
  "mcpServers": {
    "mailchimp": {
      "command": "mailchimp-mcp-server",
      "env": {
        "MAILCHIMP_API_KEY": "your-api-key-us2"
      }
    }
  }
}
```

### Docker

```json
{
  "mcpServers": {
    "mailchimp": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "MAILCHIMP_API_KEY",
        "ghcr.io/richardpowellus/mailchimp-mcp-server"
      ],
      "env": {
        "MAILCHIMP_API_KEY": "your-api-key-us2"
      }
    }
  }
}
```

## Configuration

| Variable | Required | Description |
|---|---|---|
| `MAILCHIMP_API_KEY` | Yes | Mailchimp API key in `key-datacenter` format (e.g., `abc123def-us2`) |

### Getting an API Key

1. Log in to [Mailchimp](https://login.mailchimp.com/)
2. Navigate to **Account & billing** → **Extras** → **API keys**
3. Click **Create A Key**
4. Copy the full key — it includes a datacenter suffix (e.g., `-us2`)

The suffix after the last `-` identifies your Mailchimp datacenter. The server uses this to route requests to `https://<dc>.api.mailchimp.com/3.0`.

## Tools

279 tools organized across 38 categories. Click a category to see the tools.

<details>
<summary>Account (3 tools)</summary>

| Tool | Description |
|---|---|
| `list_accounts` | List configured Mailchimp accounts |
| `ping` | Ping the Mailchimp API to verify connectivity |
| `get_account_info` | Get account details (name, contact, plan, industry) |

</details>

<details>
<summary>Campaigns (16 tools)</summary>

| Tool | Description |
|---|---|
| `list_campaigns` | List campaigns with filtering and pagination |
| `create_campaign` | Create a new campaign (regular, plaintext, RSS, A/B) |
| `get_campaign` | Get a single campaign by ID |
| `update_campaign` | Update campaign settings |
| `delete_campaign` | Delete a campaign |
| `send_campaign` | Send a campaign immediately |
| `schedule_campaign` | Schedule a campaign for future delivery |
| `unschedule_campaign` | Cancel a scheduled campaign |
| `cancel_campaign` | Cancel a campaign that is sending |
| `send_test_email` | Send a test email for a campaign |
| `replicate_campaign` | Create a copy of a campaign |
| `resend_campaign` | Resend a campaign to non-openers |
| `get_send_checklist` | Get the pre-send checklist for a campaign |
| `search_campaigns` | Search campaigns by query string |
| `pause_rss_campaign` | Pause an RSS campaign |
| `resume_rss_campaign` | Resume a paused RSS campaign |

</details>

<details>
<summary>Campaign Content (2 tools)</summary>

| Tool | Description |
|---|---|
| `get_campaign_content` | Get the HTML/text content of a campaign |
| `set_campaign_content` | Set campaign content (HTML, template, URL) |

</details>

<details>
<summary>Campaign Feedback (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_campaign_feedback` | List feedback on a campaign |
| `get_campaign_feedback` | Get a specific feedback message |
| `add_campaign_feedback` | Add feedback to a campaign |
| `update_campaign_feedback` | Update a feedback message |
| `delete_campaign_feedback` | Delete a feedback message |

</details>

<details>
<summary>Campaign Folders (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_campaign_folders` | List campaign folders |
| `get_campaign_folder` | Get a campaign folder |
| `create_campaign_folder` | Create a campaign folder |
| `update_campaign_folder` | Rename a campaign folder |
| `delete_campaign_folder` | Delete a campaign folder |

</details>

<details>
<summary>Audiences (8 tools)</summary>

| Tool | Description |
|---|---|
| `list_audiences` | List audiences (subscriber lists) |
| `create_audience` | Create a new audience |
| `get_audience` | Get audience details |
| `update_audience` | Update audience settings |
| `delete_audience` | Delete an audience |
| `get_audience_growth` | Get audience growth history |
| `get_audience_locations` | Get top subscriber locations |
| `get_email_client_stats` | Get email client usage stats for an audience |

</details>

<details>
<summary>Audience Stats (7 tools)</summary>

| Tool | Description |
|---|---|
| `list_abuse_reports` | List abuse reports for an audience |
| `get_abuse_report` | Get a specific abuse report |
| `get_audience_activity` | Get recent audience activity |
| `list_signup_forms` | List signup forms for an audience |
| `create_signup_form` | Create/update a signup form |
| `list_audience_surveys` | List surveys for an audience |
| `get_audience_survey` | Get a specific audience survey |

</details>

<details>
<summary>Members (20 tools)</summary>

| Tool | Description |
|---|---|
| `list_members` | List members of an audience |
| `get_member` | Get a member by email or subscriber hash |
| `add_member` | Add a new member to an audience |
| `update_member` | Update a member's info |
| `upsert_member` | Add or update a member (PUT) |
| `archive_member` | Archive (soft delete) a member |
| `delete_member_permanent` | Permanently delete a member |
| `search_members` | Search for members across all audiences |
| `batch_subscribe` | Batch add/update members |
| `get_member_tags` | Get tags for a member |
| `manage_member_tags` | Add or remove tags on a member |
| `list_member_notes` | List notes for a member |
| `get_member_note` | Get a specific member note |
| `add_member_note` | Add a note to a member |
| `update_member_note` | Update a member note |
| `delete_member_note` | Delete a member note |
| `get_member_events` | Get events for a member |
| `add_member_event` | Create a custom event for a member |
| `get_member_activity` | Get recent activity for a member |
| `get_member_activity_feed` | Get the full activity feed for a member |

</details>

<details>
<summary>Member Goals (1 tool)</summary>

| Tool | Description |
|---|---|
| `get_member_goals` | Get goal tracking data for a member |

</details>

<details>
<summary>Merge Fields (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_merge_fields` | List merge fields for an audience |
| `get_merge_field` | Get a specific merge field |
| `create_merge_field` | Create a merge field |
| `update_merge_field` | Update a merge field |
| `delete_merge_field` | Delete a merge field |

</details>

<details>
<summary>Interest Categories & Interests (11 tools)</summary>

| Tool | Description |
|---|---|
| `list_interest_categories` | List interest categories (groups) |
| `get_interest_category` | Get an interest category |
| `create_interest_category` | Create an interest category |
| `update_interest_category` | Update an interest category |
| `delete_interest_category` | Delete an interest category |
| `list_interests` | List interests within a category |
| `get_interest` | Get a specific interest |
| `create_interest` | Create an interest |
| `update_interest` | Update an interest |
| `delete_interest` | Delete an interest |
| `get_growth_history_month` | Get growth history for a specific month |

</details>

<details>
<summary>Segments (9 tools)</summary>

| Tool | Description |
|---|---|
| `list_segments` | List segments for an audience |
| `get_segment` | Get a segment |
| `create_segment` | Create a segment (static or dynamic) |
| `update_segment` | Update a segment |
| `delete_segment` | Delete a segment |
| `list_segment_members` | List members in a segment |
| `add_segment_member` | Add a member to a static segment |
| `remove_segment_member` | Remove a member from a static segment |
| `batch_segment_members` | Batch add/remove segment members |

</details>

<details>
<summary>Tags (1 tool)</summary>

| Tool | Description |
|---|---|
| `search_tags` | Search for tags across audiences |

</details>

<details>
<summary>Automations (18 tools)</summary>

| Tool | Description |
|---|---|
| `list_automations` | List automations |
| `get_automation` | Get an automation workflow |
| `create_automation` | Create an automation |
| `start_automation` | Start an automation |
| `pause_automation` | Pause an automation |
| `archive_automation` | Archive an automation |
| `list_automation_emails` | List emails in an automation |
| `get_automation_email` | Get an automation email |
| `delete_automation_email` | Delete an automation email |
| `start_automation_email` | Start an automation email |
| `pause_automation_email` | Pause an automation email |
| `update_automation_email` | Update automation email settings |
| `get_automation_email_content` | Get automation email content |
| `set_automation_email_content` | Set automation email content |
| `list_automation_queue` | List subscribers in an automation queue |
| `get_automation_queue_subscriber` | Get a subscriber from the automation queue |
| `remove_automation_subscriber` | Remove a subscriber from an automation |
| `get_subscriber_email_activity` | Get email activity for a subscriber in an automation |

</details>

<details>
<summary>Customer Journeys (1 tool)</summary>

| Tool | Description |
|---|---|
| `trigger_journey_step` | Trigger a step in a customer journey for a contact |

</details>

<details>
<summary>Templates (6 tools)</summary>

| Tool | Description |
|---|---|
| `list_templates` | List templates |
| `get_template` | Get a template |
| `create_template` | Create a template |
| `update_template` | Update a template |
| `delete_template` | Delete a template |
| `get_template_content` | Get the HTML content of a template |

</details>

<details>
<summary>Template Folders (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_template_folders` | List template folders |
| `get_template_folder` | Get a template folder |
| `create_template_folder` | Create a template folder |
| `update_template_folder` | Rename a template folder |
| `delete_template_folder` | Delete a template folder |

</details>

<details>
<summary>Reports (22 tools)</summary>

| Tool | Description |
|---|---|
| `list_reports` | List campaign reports |
| `get_campaign_report` | Get a specific campaign report |
| `get_campaign_abuse_reports` | List abuse reports for a campaign |
| `get_campaign_abuse_report` | Get a specific campaign abuse report |
| `get_campaign_advice` | Get campaign feedback/advice |
| `get_open_details` | Get open details for a campaign |
| `get_open_details_subscriber` | Get open details for a specific subscriber |
| `get_click_details` | Get click details for a campaign |
| `get_click_link_details` | Get details for a specific tracked link |
| `get_click_link_members` | List members who clicked a specific link |
| `get_click_link_member` | Get click details for a specific member on a link |
| `get_campaign_unsubscribes` | List unsubscribes for a campaign |
| `get_campaign_unsubscribe` | Get a specific unsubscribe record |
| `get_email_activity` | Get email activity for a campaign |
| `get_sent_to` | List recipients of a campaign |
| `get_sent_to_member` | Get delivery details for a specific recipient |
| `get_domain_performance` | Get sending domain performance for a campaign |
| `get_report_locations` | Get top open locations for a campaign |
| `get_ecommerce_product_activity` | Get e-commerce product activity for a campaign |
| `get_sub_reports` | Get child campaign reports (A/B, variate) |
| `get_eepurl_activity` | Get activity for a campaign's eepurl |
| `get_removed_subscriber` | Get a removed subscriber from a campaign |

</details>

<details>
<summary>E-commerce Stores (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_stores` | List e-commerce stores |
| `get_store` | Get a store |
| `create_store` | Create a store |
| `update_store` | Update a store |
| `delete_store` | Delete a store |

</details>

<details>
<summary>E-commerce Products (15 tools)</summary>

| Tool | Description |
|---|---|
| `list_store_products` | List products in a store |
| `get_store_product` | Get a product |
| `create_store_product` | Create a product |
| `update_store_product` | Update a product |
| `delete_store_product` | Delete a product |
| `list_product_variants` | List variants of a product |
| `get_product_variant` | Get a variant |
| `create_product_variant` | Create a variant |
| `update_product_variant` | Update a variant |
| `upsert_product_variant` | Add or update a variant (PUT) |
| `delete_product_variant` | Delete a variant |
| `list_product_images` | List images for a product |
| `get_product_image` | Get a product image |
| `create_product_image` | Add an image to a product |
| `delete_product_image` | Delete a product image |

</details>

<details>
<summary>E-commerce Orders (11 tools)</summary>

| Tool | Description |
|---|---|
| `list_store_orders` | List orders in a store |
| `list_all_orders` | List orders across all stores |
| `get_store_order` | Get an order |
| `create_store_order` | Create an order |
| `update_store_order` | Update an order |
| `delete_store_order` | Delete an order |
| `list_order_lines` | List line items in an order |
| `get_order_line` | Get an order line item |
| `create_order_line` | Add a line item to an order |
| `update_order_line` | Update an order line item |
| `delete_order_line` | Delete an order line item |

</details>

<details>
<summary>E-commerce Customers (6 tools)</summary>

| Tool | Description |
|---|---|
| `list_store_customers` | List customers in a store |
| `get_store_customer` | Get a customer |
| `create_store_customer` | Create a customer |
| `update_store_customer` | Update a customer |
| `upsert_store_customer` | Add or update a customer (PUT) |
| `delete_store_customer` | Delete a customer |

</details>

<details>
<summary>E-commerce Carts (10 tools)</summary>

| Tool | Description |
|---|---|
| `list_store_carts` | List carts in a store |
| `get_store_cart` | Get a cart |
| `create_store_cart` | Create a cart |
| `update_store_cart` | Update a cart |
| `delete_store_cart` | Delete a cart |
| `list_cart_lines` | List line items in a cart |
| `get_cart_line` | Get a cart line item |
| `create_cart_line` | Add a line item to a cart |
| `update_cart_line` | Update a cart line item |
| `delete_cart_line` | Delete a cart line item |

</details>

<details>
<summary>E-commerce Promos (10 tools)</summary>

| Tool | Description |
|---|---|
| `list_promo_rules` | List promo rules in a store |
| `get_promo_rule` | Get a promo rule |
| `create_promo_rule` | Create a promo rule |
| `update_promo_rule` | Update a promo rule |
| `delete_promo_rule` | Delete a promo rule |
| `list_promo_codes` | List promo codes for a rule |
| `get_promo_code` | Get a promo code |
| `create_promo_code` | Create a promo code |
| `update_promo_code` | Update a promo code |
| `delete_promo_code` | Delete a promo code |

</details>

<details>
<summary>Landing Pages (9 tools)</summary>

| Tool | Description |
|---|---|
| `list_landing_pages` | List landing pages |
| `get_landing_page` | Get a landing page |
| `create_landing_page` | Create a landing page |
| `update_landing_page` | Update a landing page |
| `delete_landing_page` | Delete a landing page |
| `publish_landing_page` | Publish a landing page |
| `unpublish_landing_page` | Unpublish a landing page |
| `get_landing_page_content` | Get the content of a landing page |
| `set_landing_page_content` | Set the content of a landing page |

</details>

<details>
<summary>File Manager (10 tools)</summary>

| Tool | Description |
|---|---|
| `list_files` | List files in the file manager |
| `get_file` | Get a file's metadata |
| `upload_file` | Upload a file |
| `update_file` | Update a file's metadata |
| `delete_file` | Delete a file |
| `list_file_folders` | List file manager folders |
| `get_file_folder` | Get a folder |
| `create_file_folder` | Create a folder |
| `update_file_folder` | Update a folder |
| `delete_file_folder` | Delete a folder |

</details>

<details>
<summary>Webhooks (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_webhooks` | List webhooks for an audience |
| `get_webhook` | Get a webhook |
| `create_webhook` | Create a webhook |
| `update_webhook` | Update a webhook |
| `delete_webhook` | Delete a webhook |

</details>

<details>
<summary>Batch Operations (4 tools)</summary>

| Tool | Description |
|---|---|
| `list_batches` | List batch operations |
| `get_batch_status` | Get the status of a batch operation |
| `create_batch` | Start a batch operation |
| `delete_batch` | Delete a batch operation |

</details>

<details>
<summary>Batch Webhooks (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_batch_webhooks` | List batch webhooks |
| `get_batch_webhook` | Get a batch webhook |
| `create_batch_webhook` | Create a batch webhook |
| `update_batch_webhook` | Update a batch webhook |
| `delete_batch_webhook` | Delete a batch webhook |

</details>

<details>
<summary>Verified Domains (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_verified_domains` | List verified sending domains |
| `get_verified_domain` | Get a verified domain |
| `add_verified_domain` | Add a domain for verification |
| `delete_verified_domain` | Delete a verified domain |
| `verify_domain` | Submit a verification code for a domain |

</details>

<details>
<summary>Authorized Apps (2 tools)</summary>

| Tool | Description |
|---|---|
| `list_authorized_apps` | List authorized integrations |
| `get_authorized_app` | Get an authorized app |

</details>

<details>
<summary>Activity Feed (1 tool)</summary>

| Tool | Description |
|---|---|
| `get_activity_feed` | Get the Chimp Chatter activity feed |

</details>

<details>
<summary>Conversations (4 tools)</summary>

| Tool | Description |
|---|---|
| `list_conversations` | List conversations |
| `get_conversation` | Get a conversation |
| `list_conversation_messages` | List messages in a conversation |
| `get_conversation_message` | Get a conversation message |

</details>

<details>
<summary>Connected Sites (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_connected_sites` | List connected sites |
| `get_connected_site` | Get a connected site |
| `add_connected_site` | Add a connected site |
| `delete_connected_site` | Remove a connected site |
| `verify_connected_site` | Verify a connected site's script installation |

</details>

<details>
<summary>Account Exports (3 tools)</summary>

| Tool | Description |
|---|---|
| `list_account_exports` | List account data exports |
| `get_account_export` | Get an export's details |
| `create_account_export` | Create a new account export |

</details>

<details>
<summary>Surveys (9 tools)</summary>

| Tool | Description |
|---|---|
| `list_surveys` | List surveys |
| `get_survey` | Get a survey |
| `publish_survey` | Publish a survey |
| `unpublish_survey` | Unpublish a survey |
| `list_survey_questions` | List questions in a survey |
| `get_survey_question` | Get a survey question |
| `list_survey_question_answers` | List answers for a survey question |
| `list_survey_responses` | List responses to a survey |
| `get_survey_response` | Get a specific survey response |

</details>

<details>
<summary>Facebook Ads (5 tools)</summary>

| Tool | Description |
|---|---|
| `list_facebook_ads` | List Facebook ad campaigns |
| `get_facebook_ad` | Get a Facebook ad campaign |
| `list_facebook_ad_reports` | List Facebook ad reports |
| `get_facebook_ad_report` | Get a Facebook ad report |
| `get_facebook_ad_ecommerce` | Get e-commerce data for a Facebook ad |

</details>

<details>
<summary>Landing Page Reports (2 tools)</summary>

| Tool | Description |
|---|---|
| `list_landing_page_reports` | List reports for landing pages |
| `get_landing_page_report` | Get a landing page report |

</details>

<details>
<summary>Audiences BETA (8 tools)</summary>

| Tool | Description |
|---|---|
| `list_audiences_beta` | List audiences (beta endpoint) |
| `get_audience_beta` | Get an audience (beta endpoint) |
| `create_audience_contact` | Create a contact in an audience (beta) |
| `list_audience_contacts` | List contacts in an audience (beta) |
| `get_audience_contact` | Get a contact (beta) |
| `update_audience_contact` | Update a contact (beta) |
| `archive_audience_contact` | Archive a contact (beta) |
| `forget_audience_contact` | Permanently forget a contact (GDPR, beta) |

</details>

## Usage Examples

**Campaign analytics:**

> "Show me the open rate and click rate for my last 5 campaigns."

The AI will call `list_campaigns` to find recent campaigns, then `get_campaign_report` for each to retrieve performance metrics.

**Send a campaign:**

> "Create a new campaign called 'Summer Sale' for my main audience, set the HTML content, and send a test email to me first."

Uses `list_audiences` → `create_campaign` → `set_campaign_content` → `send_test_email`.

**Subscriber management:**

> "Search for john@example.com, show me their activity and tags, then add a 'VIP' tag."

Uses `search_members` → `get_member_activity` → `get_member_tags` → `manage_member_tags`.

**Abandoned cart automation:**

> "List my automations, find the abandoned cart workflow, and show me the emails and their performance."

Uses `list_automations` → `list_automation_emails` → `get_automation_email_content` and `get_campaign_report` for each email.

**E-commerce sync:**

> "Create a new product in my store with two variants and a promo code for 10% off."

Uses `create_store_product` → `create_product_variant` (×2) → `create_promo_rule` → `create_promo_code`.

## Architecture

```
┌─────────────────┐     stdio      ┌──────────────────────────┐      HTTPS       ┌─────────────────┐
│   MCP Client    │◄──────────────►│  mailchimp-mcp-server    │◄────────────────►│  Mailchimp API  │
│  (Claude, etc.) │   JSON-RPC 2.0 │  Go binary               │  Basic Auth       │  v3 (per-DC)    │
└─────────────────┘                │                          │                   └─────────────────┘
                                   │  • 279 registered tools  │
                                   │  • 10-conn semaphore     │
                                   │  • Smart retry + backoff │
                                   │  • _links stripping      │
                                   │  • Auto-pagination       │
                                   └──────────────────────────┘
```

- **Transport:** stdio JSON-RPC 2.0 (compatible with any MCP host)
- **Auth:** HTTP Basic Auth — datacenter auto-detected from API key suffix
- **Concurrency:** Channel-based semaphore limits to 10 simultaneous connections (Mailchimp's limit)
- **Retry:** GETs retry on transient errors (429, 5xx) with exponential backoff; writes retry only on 429 with `Retry-After` header support
- **Response cleaning:** `_links` arrays stripped recursively from all responses, saving significant tokens
- **Auto-pagination:** `FetchAll()` helper automatically paginates Mailchimp list endpoints using `offset`/`count`
- **Process lifecycle:** Parent-process watchdog for clean shutdown; handles SIGINT/SIGTERM

## Extensibility

The server can be used as a Go library with a custom credential provider — for example, to load API keys from a secret manager instead of environment variables:

```go
package main

import (
    "context"

    "github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
    "github.com/richardpowellus/mailchimp-mcp-server/mailchimp/tools"
    "github.com/richardpowellus/mailchimp-mcp-server/mcp"
)

// VaultProvider loads Mailchimp API keys from a secret manager.
type VaultProvider struct { /* ... */ }

func (p *VaultProvider) GetAPIKey(ctx context.Context, name string) (string, error) {
    // Fetch from Vault, AWS Secrets Manager, Azure Key Vault, etc.
    return fetchFromVault(ctx, name)
}

func (p *VaultProvider) ListAccounts(ctx context.Context) ([]mailchimp.AccountInfo, error) {
    return []mailchimp.AccountInfo{
        {Name: "production", DisplayName: "Production Account"},
        {Name: "staging", DisplayName: "Staging Account"},
    }, nil
}

func main() {
    cfg := mailchimp.NewConfig(&VaultProvider{})
    server := mcp.New("Mailchimp", "1.0.0", "Custom Mailchimp MCP server")
    tools.RegisterAll(server, cfg)
    server.Run(context.Background())
}
```

## Building from Source

```bash
git clone https://github.com/richardpowellus/mailchimp-mcp-server.git
cd mailchimp-mcp-server
go build ./cmd/mailchimp-mcp-server/
```

### Run tests

```bash
go test ./...
```

### Run linter

```bash
go vet ./...
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Building and testing
- Code style requirements
- Adding new tools
- Pull request process

## Support

If you find this project useful, consider buying me a coffee:

<a href="https://buymeacoffee.com/letri" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" width="217" height="60"></a>

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
