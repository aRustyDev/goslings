package lease

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/arustydev/goslings/internal/auth/shared"
)

// AzureAuthFactory implements AuthFactory for Azure authentication
type AzureAuthFactory struct {
	Options    CredentialOptions
	CloudURL   string
	Expiration time.Time
	Token      *azcore.AccessToken
}

func (f *AzureAuthFactory) GetCredential(
	ctx context.Context,
	CredentialCategory CredentialCategory,
	options *CredentialOptions,
) error {
	switch CredentialCategory {
	case DeviceCode:
		azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
			TenantID:   options.TenantID,
			ClientID:   options.ClientID,
			UserPrompt: options.UserPrompt,
		})
	case ClientSecret:
		azidentity.NewClientSecretCredential(
			options.TenantID,
			options.ClientID,
			options.ClientSecret,
			&azidentity.ClientSecretCredentialOptions{})
	case InteractiveBrowser:
		azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
			TenantID: options.TenantID,
			ClientID: options.ClientID,
		})
	}
	return nil
}

func (f *AzureAuthFactory) GetToken(
	ctx context.Context,
	opts policy.TokenRequestOptions,
) error {
	// Set the cloud URL based on parameters
	f.CloudURL = f.getCloudURL(f.Options.AuthParams)
	f.Token = &azcore.AccessToken{}
	return nil
}

// getCloudURL returns the appropriate cloud URL based on params
func (l *AzureAuthFactory) getCloudURL(params *shared.AuthParams) string {
	if params.UsGovernment {
		return "https://login.microsoftonline.us"
	}
	return "https://login.microsoftonline.com"
}

func (l *AzureAuthFactory) IsTokenExpired(ctx context.Context, gracePeriod time.Duration) bool {
	return time.Now().After(l.Expiration.Add(-gracePeriod))
}

func (l *AzureAuthFactory) IsCredentialExpired(ctx context.Context, gracePeriod time.Duration) bool {
	return time.Now().After(l.Expiration.Add(-gracePeriod))
}

// // acquireViaDeviceCode attempts to authenticate using device code flow
// func (l *AzureAuthFactory) acquireViaDeviceCode(
// 	ctx context.Context,
// 	params *shared.AuthParams,
// 	creds *shared.Credentials,
// ) error {
// 	log.Info("Attempting to authenticate via device code. You may have to accept MFA prompts.")

// 	// Skip if no tenant or client ID
// 	if params.TenantID == "" || params.ClientID == "" {
// 		return errors.New("tenant ID and client ID are required for device code authentication")
// 	}

// 	// Define the device code callback
// 	deviceCodeCallback := func(ctx context.Context, deviceCode azidentity.DeviceCodeMessage) error {
// 		log.Infof("Device code authentication - Your MFA code is: %s", deviceCode.UserCode)
// 		log.Infof("Please authenticate at: %s", deviceCode.VerificationURL)
// 		return nil
// 	}

// 	// // Create options
// 	// options := &azidentity.DeviceCodeCredentialOptions{}

// 	// Create the credential using the factory
// 	credential, err := l.AuthFactory.GetCredential(ctx, DeviceCode, &CredentialOptions{
// 		TenantID:   params.TenantID,
// 		ClientID:   params.ClientID,
// 		UserPrompt: deviceCodeCallback,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create device code credential: %w", err)
// 	}

// 	// Define scopes
// 	scopes := []string{"https://graph.microsoft.com/.default"}

// 	// Get token
// 	azToken, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: scopes})
// 	if err != nil {
// 		return fmt.Errorf("failed to get token with device code: %w", err)
// 	}

// 	// Convert to our token format
// 	token := &shared.Token{
// 		Value:     azToken.Token,
// 		Type:      "Bearer",
// 		ExpiresAt: azToken.ExpiresOn,
// 		Scopes:    scopes,
// 		Resource:  "https://graph.microsoft.com",
// 	}

// 	// Add to credentials
// 	creds.Tokens["graph"] = token
// 	creds.AuthType = shared.DeviceCodeAuth

// 	log.Debug("Successfully authenticated with device code")

// 	return nil
// }

// // acquireViaClientCredentials attempts to authenticate using client credentials
// func (l *AzureAuthFactory) acquireViaClientCredentials(
// 	ctx context.Context,
// 	params *shared.AuthParams,
// 	creds *shared.Credentials,
// ) error {
// 	log.Info("Attempting to authenticate with client credentials")

// 	// Validate required fields
// 	if params.ClientID == "" || params.ClientSecret == "" || params.TenantID == "" {
// 		return errors.New(
// 			"client ID, client secret, and tenant ID are required for client credentials auth",
// 		)
// 	}

// 	// // Create options
// 	// options := &azidentity.ClientSecretCredentialOptions{}

// 	// Create the credential using the factory
// 	credential, err := l.AuthFactory.GetCredential(ctx, ClientSecret, &CredentialOptions{
// 		TenantID:     params.TenantID,
// 		ClientID:     params.ClientID,
// 		ClientSecret: params.ClientSecret,
// 		// options,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create client secret credential: %w", err)
// 	}

// 	// Define scopes
// 	scopes := []string{"https://graph.microsoft.com/.default"}

// 	// Get token
// 	azToken, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: scopes})
// 	if err != nil {
// 		return fmt.Errorf("failed to get token with client credentials: %w", err)
// 	}

// 	// Convert to our token format
// 	token := &shared.Token{
// 		Value:     azToken.Token,
// 		Type:      "Bearer",
// 		ExpiresAt: azToken.ExpiresOn,
// 		Scopes:    scopes,
// 		Resource:  "https://graph.microsoft.com",
// 	}

// 	// Add to credentials
// 	creds.Tokens["graph"] = token
// 	creds.AuthType = shared.ClientCredentialsAuth

// 	log.Debug("Successfully authenticated with client credentials")

// 	return nil
// }

// // acquireViaInteractiveBrowser attempts to authenticate using interactive browser
// func (l *AzureAuthFactory) acquireViaInteractiveBrowser(
// 	ctx context.Context,
// 	params *shared.AuthParams,
// 	creds *shared.Credentials,
// ) error {
// 	log.Info("Attempting to authenticate with interactive browser login")

// 	// Skip if no tenant or client ID
// 	if params.TenantID == "" || params.ClientID == "" {
// 		return errors.New(
// 			"tenant ID and client ID are required for interactive browser authentication",
// 		)
// 	}

// 	// // Create options
// 	// options := &azidentity.InteractiveBrowserCredentialOptions{
// 	// 	TenantID: params.TenantID,
// 	// 	ClientID: params.ClientID,
// 	// }

// 	// Create the credential using the factory
// 	credential, err := l.AuthFactory.GetCredential(
// 		ctx,
// 		InteractiveBrowser,
// 		&CredentialOptions{
// 			TenantID: params.TenantID,
// 			ClientID: params.ClientID,
// 		},
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to create interactive browser credential: %w", err)
// 	}

// 	// Define scopes
// 	scopes := []string{"https://graph.microsoft.com/.default"}

// 	// Get token
// 	err := l.GetToken(ctx, policy.TokenRequestOptions{Scopes: scopes})
// 	if err != nil {
// 		return fmt.Errorf("failed to get token with interactive browser login: %w", err)
// 	}

// 	// Convert to our token format
// 	token := &shared.Token{
// 		Value:     azToken.Token,
// 		Type:      "Bearer",
// 		ExpiresAt: azToken.ExpiresOn,
// 		Scopes:    scopes,
// 		Resource:  "https://graph.microsoft.com",
// 	}

// 	// Add to credentials
// 	creds.Tokens["graph"] = token
// 	creds.AuthType = shared.InteractiveAuth

// 	log.Debug("Successfully authenticated with interactive browser login")

// 	return nil
// }

// // updateCredentialsExpiration updates the overall expiration time of the credentials
// // to be the earliest expiration time of any token
// func (l *AzureAuthFactory) updateCredentialsExpiration(creds *shared.Credentials) {
// 	var earliestExpiration time.Time
// 	first := true

// 	// Find the earliest expiration time
// 	for _, token := range creds.Tokens {
// 		if first || token.ExpiresAt.Before(earliestExpiration) {
// 			earliestExpiration = token.ExpiresAt
// 			first = false
// 		}
// 	}

// 	// If we found an expiration time, update the credentials
// 	if !first {
// 		creds.ExpiresAt = earliestExpiration
// 	}
// }
