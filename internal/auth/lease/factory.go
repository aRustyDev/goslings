package lease

//go:generate mockgen -source=factory.go -destination=factory_mocks_test.go -package=lease

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/arustydev/goslings/internal/auth/shared"
)

// "github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
// "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
// "github.com/AzureAD/microsoft-authentication-library-for-go/apps/managedidentity"

// CredentialFactoryabstracts credential loading and retrieval to allow extending and mocking the available methods
type CredentialFactory interface {
	// Provides a credential either real or macked
	GetCredential(
		ctx context.Context,
		AcquisitionMethod AcquisitionMethod,
		options *CredentialOptions,
	) error

	// IsTokenExpired checks if the token has expired or are about to expire
	IsTokenExpired(ctx context.Context, gracePeriod time.Duration) bool

	// IsCredentialExpired checks if credentials have expired or are about to expire
	IsCredentialExpired(ctx context.Context, gracePeriod time.Duration) bool
}

// CredentialOptions is a struct for holding all options that any Credential returned by a CredentialFactoryGetCredential() could require
type CredentialOptions struct {
	DeviceCodeOptions         *azidentity.DeviceCodeCredentialOptions
	TenantID                  string
	ClientID                  string
	ClientSecret              string
	ClientSecretOptions       *azidentity.ClientSecretCredentialOptions
	InteractiveBrowserOptions *azidentity.InteractiveBrowserCredentialOptions
	Category                  *AcquisitionMethod
	Method                    *CredentialMethod
	AuthParams                *shared.AuthParams
	UserPrompt                func(ctx context.Context, deviceCode azidentity.DeviceCodeMessage) error
}

// AuthFactory is the abstraction that enables Lease to be a simple interface enabling agnosticism with regard to external auth methods and services
// This is NOT mocked, but the super-interface that includes it is (<Target>Factory)
type AuthFactory interface {
	AcquireToken(ctx context.Context, options policy.TokenRequestOptions) (*shared.Token, error) // Use AuthFactory to get a Token
	SetRequestMethod(ctx context.Context, method *CredentialMethod) error                        // Set the method the AuthFactory should use to Acquire Tokens
}

type DefaultCredentialFactory struct {
	Options    CredentialOptions
	Expiration time.Time
	Token      *azcore.AccessToken
	Credential string
}

func (f DefaultCredentialFactory) GetCredential(
	ctx context.Context,
	CredentialMethod CredentialMethod,
	options *CredentialOptions,
) error {
	switch CredentialMethod {
	case LocalFile:
		// Read from local file
		return nil
	case VaultRead:
		// Read from vault conn
		// https://pkg.go.dev/github.com/hashicorp/vault/api
		return nil
	default: // EnvVars
		// f.Credential = viper.GetString("GOSLING_CRED")
		return nil
	}
}

func (f DefaultCredentialFactory) IsTokenExpired(ctx context.Context, gracePeriod time.Duration) bool {
	return time.Now().After(f.Expiration.Add(-gracePeriod))
}

func (f DefaultCredentialFactory) IsCredentialExpired(ctx context.Context, gracePeriod time.Duration) bool {
	return time.Now().After(f.Expiration.Add(-gracePeriod))
}

type DefaultAuthFactory struct {
	Options    CredentialOptions
	Expiration time.Time
	Token      *azcore.AccessToken
	Client     *http.Client
}

func (f DefaultAuthFactory) AcquireToken(ctx context.Context, options policy.TokenRequestOptions) (*shared.Token, error) {
	return &shared.Token{}, nil
}

func (f DefaultAuthFactory) SetRequestMethod(ctx context.Context, method *CredentialMethod) error {
	return nil
}
