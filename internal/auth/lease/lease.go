// Package lease provides interfaces and implementations for acquiring and renewing authentication tokens
package lease

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/arustydev/goslings/internal/auth/shared"
	log "github.com/sirupsen/logrus"
)

type Lease struct {
	// CloudURL    string          // cloud environment URL (commercial or government)
	AuthFactory AuthFactory        // Factory for getting credentials (allows mocking & real cred retrieval)
	Expiration  time.Time          // When the credential will expire
	Options     *CredentialOptions //
}

type LeaseInfo struct {
	Expiration time.Time
	Lessor     LeaseProvider
	Started    time.Time
	Proxy      LeaseProxy
}

// LeaseProvider represents the where the leased credential is for
type LeaseProvider string

// LeaseProxy represents the middleware for getting access to a LeaseProvider
type LeaseProxy string

const (
	Azure LeaseProvider = "azure"
	M365  LeaseProvider = "m365"
	D4iot LeaseProvider = "d4iot"

	Vault LeaseProxy = "vault"
)

func NewLease(ctx context.Context, f AuthFactory) (Lease, error) {
	return Lease{
		AuthFactory: f,
		Expiration:  time.Now().Add(time.Hour * 24),
		Options:     &CredentialOptions{},
	}, nil
}

// Acquire implements Lease.Acquire for Lease
func (l *Lease) Acquire(ctx context.Context) (*shared.Credentials, error) {
	// Load the Credential for token retrieval
	l.AuthFactory.GetCredential(ctx, DeviceCode, l.Options)

	// Create a new credentials map
	creds := &shared.Credentials{
		Tokens:        make(map[string]*shared.Token),
		LastRefreshed: time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour), // Default expiration, will be updated
	}

	// Retrieve the Token
	if err := l.AuthFactory.GetToken(ctx, policy.TokenRequestOptions{}); err != nil {
		log.Debugf("Device code authentication failed: %v", err)
	}

	return creds, nil
}

func (l *Lease) IsExpired(gracePeriod time.Duration) bool {
	return time.Now().After(l.Expiration)
}

func (l *Lease) Renew(ctx context.Context) (*shared.Credentials, error) {
	// Create a new credentials map
	creds := &shared.Credentials{
		Tokens:        make(map[string]*shared.Token),
		LastRefreshed: time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour), // Default expiration, will be updated
	}

	// Retrieve the Token
	if err := l.AuthFactory.GetToken(ctx, policy.TokenRequestOptions{}); err != nil {
		log.Debugf("Device code authentication failed: %v", err)
	}

	return creds, nil
}
