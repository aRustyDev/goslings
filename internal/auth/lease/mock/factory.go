package mock

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/arustydev/goslings/internal/auth/lease"
)

// MockCredentialFactory implements CredentialFactory using real Azure SDK
type MockCredentialFactory struct{}

type MockTokenCredential struct{}

// // NewDeviceCodeCredential creates a real device code credential
// func (f *MockCredentialFactory) NewDeviceCodeCredential(options *azidentity.DeviceCodeCredentialOptions) (TokenCredential, error) {
// 	return azidentity.NewDeviceCodeCredential(options)
// }

// // NewClientSecretCredential creates a real client secret credential
// func (f *MockCredentialFactory) NewClientSecretCredential(tenantID, clientID, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (TokenCredential, error) {
// 	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
// }

// // NewInteractiveBrowserCredential creates a real interactive browser credential
// func (f *MockCredentialFactory) NewInteractiveBrowserCredential(options *azidentity.InteractiveBrowserCredentialOptions) (TokenCredential, error) {
// 	return azidentity.NewInteractiveBrowserCredential(options)
// }

func (f *MockCredentialFactory) GetCredential(
	ctx context.Context,
	CredentialCategory lease.CredentialCategory,
	options *lease.CredentialOptions,
) (lease.TokenCredential, error) {
	return &MockTokenCredential{}, nil
}

func (f *MockTokenCredential) GetToken(
	ctx context.Context,
	opts policy.TokenRequestOptions,
) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{}, nil
}
