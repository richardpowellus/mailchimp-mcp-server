package mailchimp

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// CredentialProvider abstracts how Mailchimp API keys are obtained.
// The default EnvCredentialProvider reads from environment variables.
// Custom providers (e.g., Azure Key Vault) can be injected for enterprise use.
type CredentialProvider interface {
	GetAPIKey(ctx context.Context, accountName string) (string, error)
	ListAccounts(ctx context.Context) ([]AccountInfo, error)
}

// AccountInfo holds non-secret account metadata for display.
type AccountInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// EnvCredentialProvider reads the Mailchimp API key from the
// MAILCHIMP_API_KEY environment variable. Single-account mode.
type EnvCredentialProvider struct{}

func (p *EnvCredentialProvider) GetAPIKey(_ context.Context, _ string) (string, error) {
	key := os.Getenv("MAILCHIMP_API_KEY")
	if key == "" {
		return "", fmt.Errorf("MAILCHIMP_API_KEY environment variable is not set")
	}
	return key, nil
}

func (p *EnvCredentialProvider) ListAccounts(_ context.Context) ([]AccountInfo, error) {
	key := os.Getenv("MAILCHIMP_API_KEY")
	if key == "" {
		return nil, nil
	}
	return []AccountInfo{{Name: "default", DisplayName: "Default"}}, nil
}

// Config manages Mailchimp accounts and their API clients.
type Config struct {
	provider CredentialProvider
	clients  map[string]*Client
	mu       sync.Mutex
}

// NewConfig creates a Config with the given credential provider.
func NewConfig(provider CredentialProvider) *Config {
	return &Config{
		provider: provider,
		clients:  make(map[string]*Client),
	}
}

// NewEnvConfig creates a Config that reads credentials from environment variables.
// This is the default for open-source users.
func NewEnvConfig() *Config {
	return NewConfig(&EnvCredentialProvider{})
}

// GetClient returns the API client for the named account, creating it on first access.
func (c *Config) GetClient(ctx context.Context, name string) (*Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Default account name for single-account mode
	if name == "" {
		name = "default"
	}

	if cl, ok := c.clients[name]; ok {
		return cl, nil
	}

	apiKey, err := c.provider.GetAPIKey(ctx, name)
	if err != nil {
		names, _ := c.AccountNames(ctx)
		return nil, fmt.Errorf("account %q: %w (available: %s)", name, err, strings.Join(names, ", "))
	}

	client := NewClient(apiKey)
	c.clients[name] = client
	return client, nil
}

// Accounts returns non-secret metadata for all configured accounts.
func (c *Config) Accounts(ctx context.Context) []AccountInfo {
	accounts, _ := c.provider.ListAccounts(ctx)
	return accounts
}

// AccountNames returns all account names.
func (c *Config) AccountNames(ctx context.Context) ([]string, error) {
	accounts, err := c.provider.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(accounts))
	for i, a := range accounts {
		names[i] = a.Name
	}
	return names, nil
}
