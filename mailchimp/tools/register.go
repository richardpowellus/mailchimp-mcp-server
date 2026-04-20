package tools

import (
	"github.com/richardpowellus/mailchimp-mcp-server/mailchimp"
	"github.com/richardpowellus/mailchimp-mcp-server/mcp"
)

// RegisterAll registers all Mailchimp tools with the given server.
func RegisterAll(s mcp.ToolRegistrar, cfg *mailchimp.Config) {
	RegisterAccounts(s, cfg)
	RegisterCampaigns(s, cfg)
	RegisterCampaignContent(s, cfg)
	RegisterCampaignFeedback(s, cfg)
	RegisterCampaignFolders(s, cfg)
	RegisterAudiences(s, cfg)
	RegisterAudienceStats(s, cfg)
	RegisterMembers(s, cfg)
	RegisterMemberGoals(s, cfg)
	RegisterMergeFields(s, cfg)
	RegisterInterestCategories(s, cfg)
	RegisterSegments(s, cfg)
	RegisterTags(s, cfg)
	RegisterAutomations(s, cfg)
	RegisterCustomerJourneys(s, cfg)
	RegisterTemplates(s, cfg)
	RegisterTemplateFolders(s, cfg)
	RegisterReports(s, cfg)
	RegisterEcommerceStores(s, cfg)
	RegisterEcommerceProducts(s, cfg)
	RegisterEcommerceOrders(s, cfg)
	RegisterEcommerceCustomers(s, cfg)
	RegisterEcommerceCarts(s, cfg)
	RegisterEcommercePromos(s, cfg)
	RegisterLandingPages(s, cfg)
	RegisterFileManager(s, cfg)
	RegisterWebhooks(s, cfg)
	RegisterBatchOperations(s, cfg)
	RegisterBatchWebhooks(s, cfg)
	RegisterVerifiedDomains(s, cfg)
	RegisterAuthorizedApps(s, cfg)
	RegisterActivityFeed(s, cfg)
	RegisterConversations(s, cfg)
	RegisterConnectedSites(s, cfg)
	RegisterAccountExports(s, cfg)
	RegisterSurveys(s, cfg)
	RegisterFacebookAds(s, cfg)
	RegisterLandingPageReports(s, cfg)
	RegisterAudiencesBeta(s, cfg)
}
