package mock

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/arustydev/goslings/internal/auth/lease"
)

// CredentialFactory abstracts credential creation for testing
type MockM365TokenCredentialFactory struct{}

// MockM365TokenCredential is a simple mock for TokenCredential
type MockM365TokenCredential struct {
	// Token is the token to be returned by GetToken.
	Credential string
	// ExpiresOn is the expiration time for the token.
	ExpiresOn time.Time
	// RefreshOn is a suggested time to refresh the token.
	// Clients should ignore this value when it's zero.
	RefreshOn time.Time
	// Error is an optional error to be returned by GetToken.
	Error error
}

func (cf *MockM365TokenCredentialFactory) GetCredential(
	ctx context.Context,
	CredentialCategory lease.CredentialCategory,
	options *lease.CredentialOptions,
) (*MockM365TokenCredential, error) {
	return &MockM365TokenCredential{}, nil
}

// GetToken implements the TokenCredential interface
func (mc *MockM365TokenCredential) GetToken(
	ctx context.Context,
	opts policy.TokenRequestOptions,
) (*azcore.AccessToken, error) {
	// Case: Error not nil
	if mc.Error != nil {
		return nil, mc.Error
	}

	// Case: Token not given
	token := mc.Credential
	if token == "" {
		token = "mock-access-token"
	}

	// Case: Token is Expired
	expiresOn := mc.ExpiresOn
	if expiresOn.IsZero() {
		expiresOn = time.Now().Add(1 * time.Hour)
	}

	// Case: Token needs to be refreshed
	RefreshOn := mc.RefreshOn
	if RefreshOn.IsZero() {
		RefreshOn = time.Now().Add(1 * time.Hour)
	}
	return &azcore.AccessToken{
		Token:     token,
		ExpiresOn: expiresOn.UTC(),
		RefreshOn: RefreshOn.UTC(),
	}, nil
}
