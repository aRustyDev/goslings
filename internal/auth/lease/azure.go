// Package lease provides interfaces and implementations for acquiring and renewing authentication tokens
package lease

import (
	"context"
	"errors"
	"fmt"
	"time"

	"goslings/internal/auth/shared"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	log "github.com/sirupsen/logrus"
)

// Common errors
var (
	ErrCredentialsExpired = errors.New("credentials have expired")
	ErrNotAuthenticated   = errors.New("not authenticated")
	ErrRenewalFailed      = errors.New("failed to renew credentials")
)

// Lease is the interface that wraps the basic token acquisition and renewal methods
type Lease interface {
	// Acquire acquires a new set of credentials
	Acquire(ctx context.Context, params *shared.AuthParams) (*shared.Credentials, error)

	// Renew attempts to renew existing credentials
	Renew(ctx context.Context, creds *shared.Credentials, params *shared.AuthParams) (*shared.Credentials, error)

	// IsExpired checks if credentials have expired or are about to expire
	IsExpired(creds *shared.Credentials, gracePeriod time.Duration) bool
}

// AzureLease implements Lease for Azure authentication
type AzureLease struct {
	// cloud environment URL (commercial or government)
	CloudURL string
}

// NewAzureLease creates a new Azure lease
func NewAzureLease() *AzureLease {
	return &AzureLease{
		CloudURL: "https://login.microsoftonline.com", // Default to commercial cloud
	}
}

// getCloudURL returns the appropriate cloud URL based on params
func (l *AzureLease) getCloudURL(params *shared.AuthParams) string {
	if params.UsGovernment {
		return "https://login.microsoftonline.us"
	}
	return "https://login.microsoftonline.com"
}

// Acquire implements Lease.Acquire for AzureLease
func (l *AzureLease) Acquire(ctx context.Context, params *shared.AuthParams) (*shared.Credentials, error) {
	// Set the cloud URL based on parameters
	l.CloudURL = l.getCloudURL(params)

	// Create a new credentials map
	creds := &shared.Credentials{
		Tokens:        make(map[string]*shared.Token),
		LastRefreshed: time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour), // Default expiration, will be updated
	}

	// Try authentication methods in order of preference
	var err error

	// 1. Try device code authentication first
	if err = l.acquireViaDeviceCode(ctx, params, creds); err != nil {
		log.Debugf("Device code authentication failed: %v", err)

		// 2. Try client credentials
		if err = l.acquireViaClientCredentials(ctx, params, creds); err != nil {
			log.Debugf("Client credentials authentication failed: %v", err)

			// 3. Try interactive browser
			if err = l.acquireViaInteractiveBrowser(ctx, params, creds); err != nil {
				log.Debugf("Interactive browser authentication failed: %v", err)
				return nil, errors.New("all authentication methods failed")
			}
		}
	}

	// Update the expiration time to the earliest token expiration
	l.updateCredentialsExpiration(creds)

	return creds, nil
}

// acquireViaDeviceCode attempts to authenticate using device code flow
func (l *AzureLease) acquireViaDeviceCode(ctx context.Context, params *shared.AuthParams, creds *shared.Credentials) error {
	log.Info("Attempting to authenticate via device code. You may have to accept MFA prompts.")

	// Skip if no tenant or client ID
	if params.TenantID == "" || params.ClientID == "" {
		return errors.New("tenant ID and client ID are required for device code authentication")
	}

	// Define the device code callback
	deviceCodeCallback := func(ctx context.Context, deviceCode azidentity.DeviceCodeMessage) error {
		log.Infof("Device code authentication - Your MFA code is: %s", deviceCode.UserCode)
		log.Infof("Please authenticate at: %s", deviceCode.VerificationURL)
		return nil
	}

	// Create options
	options := &azidentity.DeviceCodeCredentialOptions{
		TenantID:   params.TenantID,
		ClientID:   params.ClientID,
		UserPrompt: deviceCodeCallback,
	}

	// Create the credential
	credential, err := azidentity.NewDeviceCodeCredential(options)
	if err != nil {
		return fmt.Errorf("failed to create device code credential: %w", err)
	}

	// Define scopes
	scopes := []string{"https://graph.microsoft.com/.default"}

	// Get token
	azToken, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: scopes})
	if err != nil {
		return fmt.Errorf("failed to get token with device code: %w", err)
	}

	// Convert to our token format
	token := &shared.Token{
		Value:     azToken.Token,
		Type:      "Bearer",
		ExpiresAt: azToken.ExpiresOn,
		Scopes:    scopes,
		Resource:  "https://graph.microsoft.com",
	}

	// Add to credentials
	creds.Tokens["graph"] = token
	creds.AuthType = shared.DeviceCodeAuth

	log.Debug("Successfully authenticated with device code")

	return nil
}

// acquireViaClientCredentials attempts to authenticate using client credentials
func (l *AzureLease) acquireViaClientCredentials(ctx context.Context, params *shared.AuthParams, creds *shared.Credentials) error {
	log.Info("Attempting to authenticate with client credentials")

	// Validate required fields
	if params.ClientID == "" || params.ClientSecret == "" || params.TenantID == "" {
		return errors.New("client ID, client secret, and tenant ID are required for client credentials auth")
	}

	// Create options
	options := &azidentity.ClientSecretCredentialOptions{}

	// Create the credential
	credential, err := azidentity.NewClientSecretCredential(
		params.TenantID,
		params.ClientID,
		params.ClientSecret,
		options,
	)
	if err != nil {
		return fmt.Errorf("failed to create client secret credential: %w", err)
	}

	// Define scopes
	scopes := []string{"https://graph.microsoft.com/.default"}

	// Get token
	azToken, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: scopes})
	if err != nil {
		return fmt.Errorf("failed to get token with client credentials: %w", err)
	}

	// Convert to our token format
	token := &shared.Token{
		Value:     azToken.Token,
		Type:      "Bearer",
		ExpiresAt: azToken.ExpiresOn,
		Scopes:    scopes,
		Resource:  "https://graph.microsoft.com",
	}

	// Add to credentials
	creds.Tokens["graph"] = token
	creds.AuthType = shared.ClientCredentialsAuth

	log.Debug("Successfully authenticated with client credentials")

	return nil
}

// acquireViaInteractiveBrowser attempts to authenticate using interactive browser
func (l *AzureLease) acquireViaInteractiveBrowser(ctx context.Context, params *shared.AuthParams, creds *shared.Credentials) error {
	log.Info("Attempting to authenticate with interactive browser login")

	// Skip if no tenant or client ID
	if params.TenantID == "" || params.ClientID == "" {
		return errors.New("tenant ID and client ID are required for interactive browser authentication")
	}

	// Create options
	options := &azidentity.InteractiveBrowserCredentialOptions{
		TenantID: params.TenantID,
		ClientID: params.ClientID,
	}

	// Create the credential
	credential, err := azidentity.NewInteractiveBrowserCredential(options)
	if err != nil {
		return fmt.Errorf("failed to create interactive browser credential: %w", err)
	}

	// Define scopes
	scopes := []string{"https://graph.microsoft.com/.default"}

	// Get token
	azToken, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: scopes})
	if err != nil {
		return fmt.Errorf("failed to get token with interactive browser login: %w", err)
	}

	// Convert to our token format
	token := &shared.Token{
		Value:     azToken.Token,
		Type:      "Bearer",
		ExpiresAt: azToken.ExpiresOn,
		Scopes:    scopes,
		Resource:  "https://graph.microsoft.com",
	}

	// Add to credentials
	creds.Tokens["graph"] = token
	creds.AuthType = shared.InteractiveAuth

	log.Debug("Successfully authenticated with interactive browser login")

	return nil
}

// Renew implements Lease.Renew for AzureLease
func (l *AzureLease) Renew(ctx context.Context, creds *shared.Credentials, params *shared.AuthParams) (*shared.Credentials, error) {
	log.Debug("Attempting to renew credentials")

	// Set the cloud URL based on parameters
	l.CloudURL = l.getCloudURL(params)

	// If the credentials don't have any tokens, we can't renew
	if len(creds.Tokens) == 0 {
		return nil, ErrNotAuthenticated
	}

	// Check if credentials are expired
	if l.IsExpired(creds, 0) {
		// If far past expiration, get new credentials instead
		if time.Since(creds.ExpiresAt) > time.Hour {
			log.Info("Credentials expired too long ago, acquiring new ones")
			return l.Acquire(ctx, params)
		}
	}

	// Create a new credentials object to hold renewed tokens
	newCreds := &shared.Credentials{
		Tokens:        make(map[string]*shared.Token),
		AuthType:      creds.AuthType,
		LastRefreshed: time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour), // Default expiration, will be updated
	}

	// Renew tokens based on auth type
	var err error

	switch creds.AuthType {
	case shared.DeviceCodeAuth:
		err = l.renewViaDeviceCode(ctx, params, creds, newCreds)
	case shared.ClientCredentialsAuth:
		err = l.renewViaClientCredentials(ctx, params, creds, newCreds)
	case shared.InteractiveAuth:
		err = l.renewViaInteractiveBrowser(ctx, params, creds, newCreds)
	default:
		err = fmt.Errorf("unsupported auth type: %s", creds.AuthType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to renew credentials: %w", err)
	}

	// Update the expiration time to the earliest token expiration
	l.updateCredentialsExpiration(newCreds)

	return newCreds, nil
}

// renewViaDeviceCode renews credentials using device code flow
func (l *AzureLease) renewViaDeviceCode(ctx context.Context, params *shared.AuthParams, oldCreds *shared.Credentials, newCreds *shared.Credentials) error {
	// For device code, we generally need to re-authenticate
	// Most tokens obtained this way don't support refresh tokens
	return l.acquireViaDeviceCode(ctx, params, newCreds)
}

// renewViaClientCredentials renews credentials using client credentials
func (l *AzureLease) renewViaClientCredentials(ctx context.Context, params *shared.AuthParams, oldCreds *shared.Credentials, newCreds *shared.Credentials) error {
	// For client credentials, just get a new token
	return l.acquireViaClientCredentials(ctx, params, newCreds)
}

// renewViaInteractiveBrowser renews credentials using interactive browser
func (l *AzureLease) renewViaInteractiveBrowser(ctx context.Context, params *shared.AuthParams, oldCreds *shared.Credentials, newCreds *shared.Credentials) error {
	// For interactive browser, we generally need to re-authenticate
	// since refresh tokens are typically not accessible
	return l.acquireViaInteractiveBrowser(ctx, params, newCreds)
}

// IsExpired implements Lease.IsExpired for AzureLease
func (l *AzureLease) IsExpired(creds *shared.Credentials, gracePeriod time.Duration) bool {
	// If no credentials, consider them expired
	if creds == nil || len(creds.Tokens) == 0 {
		return true
	}

	// Check if the credentials as a whole are expired
	if time.Now().Add(gracePeriod).After(creds.ExpiresAt) {
		return true
	}

	return false
}

// updateCredentialsExpiration updates the overall expiration time of the credentials
// to be the earliest expiration time of any token
func (l *AzureLease) updateCredentialsExpiration(creds *shared.Credentials) {
	var earliestExpiration time.Time
	first := true

	// Find the earliest expiration time
	for _, token := range creds.Tokens {
		if first || token.ExpiresAt.Before(earliestExpiration) {
			earliestExpiration = token.ExpiresAt
			first = false
		}
	}

	// If we found an expiration time, update the credentials
	if !first {
		creds.ExpiresAt = earliestExpiration
	}
}
