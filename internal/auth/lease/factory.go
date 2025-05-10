package lease

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// DefaultCredentialFactory implements CredentialFactory using real Azure SDK
type DefaultCredentialFactory struct {
	// options CredentialOptions
}

type DefaultTokenCredential struct{}

// // NewDeviceCodeCredential creates a real device code credential
// func (f *DefaultCredentialFactory) NewDeviceCodeCredential(options *azidentity.DeviceCodeCredentialOptions) (TokenCredential, error) {
// 	return azidentity.NewDeviceCodeCredential(options)
// }

// // NewClientSecretCredential creates a real client secret credential
// func (f *DefaultCredentialFactory) NewClientSecretCredential(tenantID, clientID, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenCredential, error) {
// 	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
// }

// // NewInteractiveBrowserCredential creates a real interactive browser credential
// func (f *DefaultCredentialFactory) NewInteractiveBrowserCredential(options *azidentity.InteractiveBrowserCredentialOptions) (TokenCredential, error) {
// 	return azidentity.NewInteractiveBrowserCredential(options)
// }

func (f *DefaultCredentialFactory) GetCredential(ctx context.Context, CredentialCategory CredentialCategory, options *CredentialOptions) (TokenCredential, error) {
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
	return &DefaultTokenCredential{}, nil
}

func (f *DefaultTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{}, nil
}
