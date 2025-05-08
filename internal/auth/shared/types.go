// Package auth provides authentication functionality for Microsoft services
package shared

import (
	"time"
)

// AuthType represents the type of authentication being used
type AuthType string

const (
	// DeviceCodeAuth represents authentication via device code flow
	DeviceCodeAuth AuthType = "device_code"

	// ClientCredentialsAuth represents authentication via client credentials
	ClientCredentialsAuth AuthType = "client_credentials"

	// InteractiveAuth represents authentication via interactive browser
	InteractiveAuth AuthType = "interactive"

	// M365Auth represents authentication specific to M365 services
	M365Auth AuthType = "m365"
)

// AuthParams contains all the parameters needed for authentication
type AuthParams struct {
	// Username is the user's email/username
	Username string `mapstructure:"GOSLING_USER"`

	// Password is the user's password (should be handled securely)
	Password string `mapstructure:"GOSLING_PASS"`

	// TenantID is the Azure/Microsoft 365 tenant ID
	TenantID string `mapstructure:"GOSLING_TENANT"`

	// ClientID (also known as AppID) is the application's client ID
	ClientID string `mapstructure:"GOSLING_APP_ID"`

	// ClientSecret is the application's client secret
	ClientSecret string `mapstructure:"GOSLING_APP_SECRET"`

	// SubscriptionID is the Azure subscription ID
	SubscriptionID string `mapstructure:"GOSLING_SUBSCRIPTION"`

	// UsGovernment indicates whether to use US Government cloud endpoints
	UsGovernment bool `mapstructure:"GOSLING_USGOV_CLOUD"`

	// ExoUSGovernment indicates whether to use US Government Exchange Online endpoints
	ExoUSGovernment bool `mapstructure:"GOSLING_USGOV_EXO"`

	// M365Enabled indicates whether M365 authentication is enabled
	M365Enabled bool `mapstructure:"GOSLING_M365_AUTH"`

	// MessageTraceEnabled indicates whether Exchange message trace is enabled
	MessageTraceEnabled bool `mapstructure:"GOSLING_EXO_MSG_TRACE"`
}

// Token represents an authentication token
type Token struct {
	// Value is the actual token string
	Value string

	// Type is the token type (e.g., "Bearer")
	Type string

	// ExpiresAt is when the token expires
	ExpiresAt time.Time

	// RefreshToken is the token that can be used to get a new access token
	RefreshToken string

	// Scopes contains the granted scopes for this token
	Scopes []string

	// Resource is the resource this token grants access to
	Resource string
}

// Credentials holds a collection of tokens for different services
type Credentials struct {
	// Tokens is a map of service names to their respective tokens
	Tokens map[string]*Token

	// AuthType indicates how these credentials were obtained
	AuthType AuthType

	// LastRefreshed indicates when these credentials were last refreshed
	LastRefreshed time.Time

	// ExpiresAt is when these credentials expire (usually the earliest token expiration)
	ExpiresAt time.Time
}

// AuthConfig contains authentication configuration and credentials
type AuthConfig struct {
	// Params contains the authentication parameters
	Params *AuthParams

	// Credentials contains the authentication tokens
	Credentials *Credentials
}

// StoreType represents the type of credential store being used
type StoreType string

const (
	// FileStore represents a file-based credential store
	FileStore StoreType = "file"

	// K8sStore represents a Kubernetes Secret-based credential store
	K8sStore StoreType = "kubernetes"

	// VaultStore represents a HashiCorp Vault-based credential store
	VaultStore StoreType = "vault"
)

// M365Resources holds M365-specific authentication resources
type M365Resources struct {
	// Cookies for Exchange Online sessions
	ExchangeCookies map[string]string

	// ValidationKey for certain M365 operations
	ValidationKey string

	// Additional tokens or session identifiers specific to M365
	AdditionalTokens map[string]string
}
