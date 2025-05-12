// Package lease provides interfaces and implementations for acquiring and renewing authentication tokens
package lease

// AcquisitionMethod the Kind of credential retrieved
type AcquisitionMethod string

// CredentialMethod represents how the credential was retrieved
type CredentialMethod string

const (
	LocalFile CredentialMethod = "local" // Load Creds from local
	EnvVars   CredentialMethod = "env"   // Load Creds from envvars
	VaultRead CredentialMethod = "vault" // Load Creds from vault
)
