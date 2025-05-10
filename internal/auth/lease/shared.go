// Package lease provides interfaces and implementations for acquiring and renewing authentication tokens
package lease

import (
	"context"
	"time"

	"goslings/internal/auth/shared"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// Lease is the interface that wraps the basic token acquisition and renewal methods
type Lease interface {
	// Acquire acquires a new set of credentials
	Acquire(ctx context.Context, factory *CredentialFactory) (*shared.Credentials, error)

	// Renew attempts to renew existing credentials
	Renew(ctx context.Context, factory *CredentialFactory) (*shared.Credentials, error)

	// IsExpired checks if credentials have expired or are about to expire
	IsExpired(factory *CredentialFactory, gracePeriod time.Duration) bool
}

// CredentialCategory represents the category of credential being used
type CredentialCategory string

const (
	// DeviceCode represents a Device code credential type
	DeviceCode CredentialCategory = "devicecode"

	// ClientSecret represents a Client Secret credential type
	ClientSecret CredentialCategory = "clientsecret"

	// InteractiveBrowser represents a credential retrieved from an Interactive Browser session
	InteractiveBrowser CredentialCategory = "interactivebrowser"
)

// CredentialMethod represents the method used to get the credential
type CredentialMethod string

const (
	// DeviceCode represents a Device code credential type
	LocalFile  CredentialMethod = "local"
	Config     CredentialMethod = "config"
	VaultLease CredentialMethod = "vault"
)

// CredentialFactory abstracts credential creation for testing
type CredentialFactory interface {
	GetCredential(ctx context.Context, CredentialCategory CredentialCategory, options *CredentialOptions) (TokenCredential, error)
}

// CredentialOptions is a struct for holding all options that any Credential returned by a CredentialFactory GetCredential() could require
type CredentialOptions struct {
	DeviceCodeOptions         *azidentity.DeviceCodeCredentialOptions
	TenantID                  string
	ClientID                  string
	ClientSecret              string
	ClientSecretOptions       *azidentity.ClientSecretCredentialOptions
	InteractiveBrowserOptions *azidentity.InteractiveBrowserCredentialOptions
	Category                  *CredentialCategory
	Method                    *CredentialMethod
	AuthParams                *shared.AuthParams
	UserPrompt                func(ctx context.Context, deviceCode azidentity.DeviceCodeMessage) error
}

// TokenCredential abstracts the Azure TokenCredential interface for testing
type TokenCredential interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error)
}

// LeaseType represents the type of lease being sought
type LeaseType string

const (
	// Azure represents a Azure lease type
	Azure LeaseType = "azure"

	// M365 represents a M365 lease type
	M365 LeaseType = "m365"

	// D4iot represents a D4iot lease type
	D4iot LeaseType = "d4iot"
)
