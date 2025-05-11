// Package lease provides interfaces and implementations for acquiring and renewing authentication tokens
package lease

// CredentialCategory the Kind of credential retrieved
type CredentialCategory string

// CredentialMethod represents how the credential was retrieved
type CredentialMethod string

const (
	DeviceCode         CredentialCategory = "devicecode"
	ClientSecret       CredentialCategory = "clientsecret"
	InteractiveBrowser CredentialCategory = "interactivebrowser"

	LocalFile  CredentialMethod = "local"
	Config     CredentialMethod = "config"
	VaultLease CredentialMethod = "vault"
)
