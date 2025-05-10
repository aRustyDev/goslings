package lease

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/arustydev/goslings/internal/auth/shared"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotAuthenticated   = errors.New("not authenticated")
	ErrCredentialsExpired = errors.New("credentials expired")
)

// AzureLease implements Lease for Azure authentication
type AzureLease struct {
	// cloud environment URL (commercial or government)
	CloudURL string

	// Factory for creating credentials (for testing)
	CredentialFactory CredentialFactory

	Expiration time.Time

	Options *CredentialOptions
}

// NewAzureLease creates a new Azure lease
func NewAzureLease(options CredentialOptions) *AzureLease {
	return &AzureLease{
		CloudURL:          "https://login.microsoftonline.com", // Default to commercial cloud
		CredentialFactory: &DefaultCredentialFactory{},
		Options:           &options,
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
func (l *AzureLease) Acquire(
	ctx context.Context,
	factory *CredentialFactory,
) (*shared.Credentials, error) {
	// Set the cloud URL based on parameters
	l.CloudURL = l.getCloudURL(l.Options.AuthParams)

	// Create a new credentials map
	creds := &shared.Credentials{
		Tokens:        make(map[string]*shared.Token),
		LastRefreshed: time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour), // Default expiration, will be updated
	}

	// Try authentication methods in order of preference
	var err error

	// 1. Try device code authentication first
	if err = l.acquireViaDeviceCode(ctx, l.Options.AuthParams, creds); err != nil {
		log.Debugf("Device code authentication failed: %v", err)

		// 2. Try client credentials
		if err = l.acquireViaClientCredentials(ctx, l.Options.AuthParams, creds); err != nil {
			log.Debugf("Client credentials authentication failed: %v", err)

			// 3. Try interactive browser
			if err = l.acquireViaInteractiveBrowser(ctx, l.Options.AuthParams, creds); err != nil {
				log.Debugf("Interactive browser authentication failed: %v", err)
				return nil, errors.New("all authentication methods failed")
			}
		}
	}

	// Update the expiration time to the earliest token expiration
	l.updateCredentialsExpiration(creds)

	return creds, nil
}

func (l *AzureLease) IsExpired(factory *CredentialFactory, gracePeriod time.Duration) bool {
	return time.Now().After(l.Expiration)
}

func (l *AzureLease) Renew(
	ctx context.Context,
	factory *CredentialFactory,
) (*shared.Credentials, error) {
	// Set the cloud URL based on parameters
	l.CloudURL = l.getCloudURL(l.Options.AuthParams)

	// Create a new credentials map
	creds := &shared.Credentials{
		Tokens:        make(map[string]*shared.Token),
		LastRefreshed: time.Now(),
		ExpiresAt:     time.Now().Add(time.Hour), // Default expiration, will be updated
	}

	// Try authentication methods in order of preference
	var err error

	// 1. Try device code authentication first
	if err = l.acquireViaDeviceCode(ctx, l.Options.AuthParams, creds); err != nil {
		log.Debugf("Device code authentication failed: %v", err)

		// 2. Try client credentials
		if err = l.acquireViaClientCredentials(ctx, l.Options.AuthParams, creds); err != nil {
			log.Debugf("Client credentials authentication failed: %v", err)

			// 3. Try interactive browser
			if err = l.acquireViaInteractiveBrowser(ctx, l.Options.AuthParams, creds); err != nil {
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
func (l *AzureLease) acquireViaDeviceCode(
	ctx context.Context,
	params *shared.AuthParams,
	creds *shared.Credentials,
) error {
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

	// // Create options
	// options := &azidentity.DeviceCodeCredentialOptions{}

	// Create the credential using the factory
	credential, err := l.CredentialFactory.GetCredential(ctx, DeviceCode, &CredentialOptions{
		TenantID:   params.TenantID,
		ClientID:   params.ClientID,
		UserPrompt: deviceCodeCallback,
	})
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
func (l *AzureLease) acquireViaClientCredentials(
	ctx context.Context,
	params *shared.AuthParams,
	creds *shared.Credentials,
) error {
	log.Info("Attempting to authenticate with client credentials")

	// Validate required fields
	if params.ClientID == "" || params.ClientSecret == "" || params.TenantID == "" {
		return errors.New(
			"client ID, client secret, and tenant ID are required for client credentials auth",
		)
	}

	// // Create options
	// options := &azidentity.ClientSecretCredentialOptions{}

	// Create the credential using the factory
	credential, err := l.CredentialFactory.GetCredential(ctx, ClientSecret, &CredentialOptions{
		TenantID:     params.TenantID,
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
		// options,
	})
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
func (l *AzureLease) acquireViaInteractiveBrowser(
	ctx context.Context,
	params *shared.AuthParams,
	creds *shared.Credentials,
) error {
	log.Info("Attempting to authenticate with interactive browser login")

	// Skip if no tenant or client ID
	if params.TenantID == "" || params.ClientID == "" {
		return errors.New(
			"tenant ID and client ID are required for interactive browser authentication",
		)
	}

	// // Create options
	// options := &azidentity.InteractiveBrowserCredentialOptions{
	// 	TenantID: params.TenantID,
	// 	ClientID: params.ClientID,
	// }

	// Create the credential using the factory
	credential, err := l.CredentialFactory.GetCredential(
		ctx,
		InteractiveBrowser,
		&CredentialOptions{
			TenantID: params.TenantID,
			ClientID: params.ClientID,
		},
	)
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
