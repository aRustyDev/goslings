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
	CredentialFactory CredentialFactory  // Factory for getting credentials (allows mocking & real cred retrieval)
	AuthFactory       AuthFactory        // Factory for external authentication; ie token retrieval (allows mocking & real cred retrieval)
	Expiration        time.Time          // When the credential will expire
	Options           *CredentialOptions //
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
	Vault LeaseProxy = "vault"
)

func NewLease(ctx context.Context, f CredentialFactory) (Lease, error) {
	return Lease{
		CredentialFactory: f,
		Expiration:        time.Now().Add(time.Hour * 24),
		Options:           &CredentialOptions{},
	}, nil
}

// Acquire implements Lease.Acquire for Lease
func (l *Lease) Acquire(ctx context.Context) (*shared.Credentials, error) {
	// Load the Credential for token retrieval
	err := l.CredentialFactory.GetCredential(ctx, DeviceCode, l.Options)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// Create a new credentials map
	creds := &shared.Credentials{
		Tokens:        make(map[string]*shared.Token),
		LastRefreshed: time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour), // Default expiration, will be updated
	}

	// Retrieve the Token
	if token, err := l.AuthFactory.AcquireToken(ctx, policy.TokenRequestOptions{}); err != nil {
		log.Debugf("Device code authentication failed: %v", err)
	} else {
		creds.Tokens["example"] = token
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
	if token, err := l.AuthFactory.AcquireToken(ctx, policy.TokenRequestOptions{}); err != nil {
		log.Debugf("Device code authentication failed: %v", err)
	} else {
		creds.Tokens["example"] = token
	}

	return creds, nil
}
