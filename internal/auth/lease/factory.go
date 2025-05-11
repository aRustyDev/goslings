package lease

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/arustydev/goslings/internal/auth/shared"
)

// AuthFactory abstracts credential loading and retrieval to allow extending and mocking the available methods
type AuthFactory interface {
	// Provides a credential either real or macked
	GetCredential(
		ctx context.Context,
		CredentialCategory CredentialCategory,
		options *CredentialOptions,
	) error

	// Acquire acquires a new token
	GetToken(
		ctx context.Context,
		opts policy.TokenRequestOptions,
	) error

	// IsTokenExpired checks if the token has expired or are about to expire
	IsTokenExpired(ctx context.Context, gracePeriod time.Duration) bool

	// IsCredentialExpired checks if credentials have expired or are about to expire
	IsCredentialExpired(ctx context.Context, gracePeriod time.Duration) bool
}

// CredentialOptions is a struct for holding all options that any Credential returned by a AuthFactory GetCredential() could require
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

// lease(ctx, azcf) -> someLease
// 		c = azcf.getcred()
// 		return c.getToken(tgt)
//
// someLease {
// 		token
// 		credFact
//
// }
