// Package auth provides authentication functionality for Microsoft services
// This is a Go rewrite of the shared.py module from the Untitled Goose Tool
// Package auth provides authentication functionality for Microsoft services
package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/arustydev/goslings/internal/auth/lease"
	"github.com/arustydev/goslings/internal/auth/shared"
	"github.com/arustydev/goslings/internal/auth/store"
	log "github.com/sirupsen/logrus"
)

// Common errors
var (
	ErrStoreNotInitialized = errors.New("credential store not initialized")
	ErrLeaseNotInitialized = errors.New("authentication lease not initialized")
)

// Errors from external packages
var (
	ErrNotAuthenticated   = lease.ErrNotAuthenticated
	ErrCredentialsExpired = lease.ErrCredentialsExpired
)

// Service defines supported service types
type Service string

const (
	// AzureService represents Azure services
	AzureService Service = "azure"

	// M365Service represents Microsoft 365 services
	M365Service Service = "m365"

	// GraphService represents Microsoft Graph API
	GraphService Service = "graph"
)

// AuthManager is the main entry point for authentication functionality
type AuthManager struct {
	// Store handles credential storage and retrieval
	Store store.Store

	// Leases handles different authentication services
	Leases map[Service]lease.Lease

	// mu protects concurrent access to credentials
	mu sync.RWMutex

	// Current authentication state
	currentAuthParams *shared.AuthParams
	currentCreds      *shared.Credentials
	m365Resources     *shared.M365Resources
}

// Options contains options for creating a new auth manager
type Options struct {
	// StoreType is the type of credential store to use
	StoreType shared.StoreType

	// StorePath is the path for file-based stores
	StorePath string

	// EncryptionKey is the key for encrypting sensitive data
	EncryptionKey []byte
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(opts Options) (*AuthManager, error) {
	auth := &AuthManager{
		Leases: make(map[Service]lease.Lease),
	}

	// Initialize the credential store
	var err error
	switch opts.StoreType {
	case shared.FileStore:
		auth.Store, err = store.NewFileStore(opts.StorePath, opts.EncryptionKey)
	case shared.K8sStore:
		// TODO: implement Kubernetes store logic
		// auth.Store, err = store.NewK8sStore(opts.StorePath, opts.EncryptionKey)
		return nil, errors.New("kubernetes store not implemented")
	case shared.VaultStore:
		// TODO: implement Vault store logic
		// auth.Store, err = store.NewVaultStore(opts.StorePath, opts.EncryptionKey)
		return nil, errors.New("vault store not implemented")
	default:
		return nil, fmt.Errorf("unsupported store type: %s", opts.StoreType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize credential store: %w", err)
	}

	// Initialize the leases
	auth.Leases[AzureService] = lease.NewAzureLease(lease.CredentialOptions{})
	// auth.Leases[M365Service] = lease.NewM365Lease()

	// Load credentials from store
	if err := auth.loadFromStore(context.Background()); err != nil {
		log.Debugf("Failed to load credentials from store: %v", err)
		// Continue without credentials, we'll get them later
	}

	return auth, nil
}

// loadFromStore loads authentication state from the store
func (a *AuthManager) loadFromStore(ctx context.Context) error {
	if a.Store == nil {
		return ErrStoreNotInitialized
	}

	// Load auth params
	params, err := a.Store.LoadParams(ctx)
	if err == nil {
		a.currentAuthParams = params
	}

	// Load credentials
	creds, err := a.Store.LoadCredentials(ctx)
	if err == nil {
		a.currentCreds = creds
	}

	// Load M365 resources
	resources, err := a.Store.LoadM365Resources(ctx)
	if err == nil {
		a.m365Resources = resources
	}

	return nil
}

// saveToStore saves authentication state to the store
func (a *AuthManager) saveToStore(ctx context.Context) error {
	if a.Store == nil {
		return ErrStoreNotInitialized
	}

	// Save auth params
	if a.currentAuthParams != nil {
		if err := a.Store.StoreParams(ctx, a.currentAuthParams); err != nil {
			return fmt.Errorf("failed to store auth params: %w", err)
		}
	}

	// Save credentials
	if a.currentCreds != nil {
		if err := a.Store.StoreCredentials(ctx, a.currentCreds); err != nil {
			return fmt.Errorf("failed to store credentials: %w", err)
		}
	}

	// Save M365 resources
	if a.m365Resources != nil {
		if err := a.Store.StoreM365Resources(ctx, a.m365Resources); err != nil {
			return fmt.Errorf("failed to store M365 resources: %w", err)
		}
	}

	return nil
}

// Authenticate performs authentication using the provided parameters
func (a *AuthManager) Authenticate(ctx context.Context, params *shared.AuthParams) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	log.Debug("Starting authentication process")

	// Update current auth params
	a.currentAuthParams = params

	// // Authenticate to Azure/Graph
	// azureLease, ok := a.Leases[AzureService]
	// if !ok {
	// 	return fmt.Errorf("azure lease not found")
	// }

	// azureCreds, err := azureLease.Acquire(ctx, params)
	// if err != nil {
	// 	return fmt.Errorf("failed to authenticate to Azure: %w", err)
	// }

	// // Authenticate to M365 if enabled
	// if params.M365Enabled {
	// 	// m365Lease, ok := a.Leases[M365Service]
	// 	// if !ok {
	// 	// 	return fmt.Errorf("M365 lease not found")
	// 	// }

	// 	// m365Creds, err := m365Lease.Acquire(ctx, params)
	// 	// if err != nil {
	// 	// 	log.Infof("Warning: Failed to authenticate to M365: %v", err)
	// 	// 	// Continue anyway, as this is not critical
	// 	// } else {
	// 	// 	// Merge tokens from M365 into Azure credentials
	// 	// 	for name, token := range m365Creds.Tokens {
	// 	// 		azureCreds.Tokens[name] = token
	// 	// 	}

	// 	// 	// Update expiration if M365 tokens expire earlier
	// 	// 	if m365Creds.ExpiresAt.Before(azureCreds.ExpiresAt) {
	// 	// 		azureCreds.ExpiresAt = m365Creds.ExpiresAt
	// 	// 	}
	// 	// }
	// }

	// // Update current credentials
	// a.currentCreds = azureCreds

	// Save to store
	if err := a.saveToStore(ctx); err != nil {
		log.Infof("Warning: Failed to save authentication state: %v", err)
		// Continue anyway, as we've authenticated successfully
	}

	log.Info("Authentication completed successfully")

	return nil
}

// GetToken gets a token for the specified service
// TODO: This probably shouldn't be public?
func (a *AuthManager) GetToken(service Service) (*shared.Token, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Check if we have credentials
	if a.currentCreds == nil || len(a.currentCreds.Tokens) == 0 {
		return nil, ErrNotAuthenticated
	}

	// Map service to token name
	var tokenName string
	switch service {
	case AzureService:
		tokenName = "azure"
	case M365Service:
		tokenName = "exchange"
	case GraphService:
		tokenName = "graph"
	default:
		return nil, fmt.Errorf("unsupported service: %s", service)
	}

	// Get the token
	token, ok := a.currentCreds.Tokens[tokenName]
	if !ok {
		// Try graph token as fallback for Azure
		if service == AzureService {
			token, ok = a.currentCreds.Tokens["graph"]
			if !ok {
				return nil, fmt.Errorf("token not found for service: %s", service)
			}
		} else {
			return nil, fmt.Errorf("token not found for service: %s", service)
		}
	}

	// Check if token is expired
	if time.Now().After(token.ExpiresAt) {
		return nil, ErrCredentialsExpired
	}

	return token, nil
}

// RenewTokens renews all tokens
func (a *AuthManager) RenewTokens(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if we have credentials to renew
	if a.currentCreds == nil || a.currentAuthParams == nil {
		return ErrNotAuthenticated
	}

	log.Debug("Starting token renewal process")

	// // Renew Azure/Graph tokens
	// azureLease, ok := a.Leases[AzureService]
	// if !ok {
	// 	return fmt.Errorf("azure lease not found")
	// }

	// // Check if tokens are expired or about to expire
	// if !azureLease.IsExpired(a.currentCreds, 5*time.Minute) {
	// 	log.Debug("Tokens are not expired, skipping renewal")
	// 	return nil
	// }

	// // Renew the tokens
	// azureCreds, err := azureLease.Renew(ctx, a.currentCreds, a.currentAuthParams)
	// if err != nil {
	// 	return fmt.Errorf("failed to renew Azure tokens: %w", err)
	// }

	// // Renew M365 tokens if needed
	// if a.currentAuthParams.M365Enabled {
	// 	// m365Lease, ok := a.Leases[M365Service]
	// 	// if !ok {
	// 	// 	return fmt.Errorf("M365 lease not found")
	// 	// }

	// 	// // Renew if M365 tokens are present
	// 	// m365TokenFound := false
	// 	// for name := range a.currentCreds.Tokens {
	// 	// 	if name == "exchange" || name == "msgtrace" {
	// 	// 		m365TokenFound = true
	// 	// 		break
	// 	// 	}
	// 	// }

	// 	// if m365TokenFound {
	// 	// 	m365Creds, err := m365Lease.Renew(ctx, a.currentCreds, a.currentAuthParams)
	// 	// 	if err != nil {
	// 	// 		log.Infof("Warning: Failed to renew M365 tokens: %v", err)
	// 	// 		// Continue anyway, as this is not critical
	// 	// 	} else {
	// 	// 		// Merge tokens from M365 into Azure credentials
	// 	// 		for name, token := range m365Creds.Tokens {
	// 	// 			azureCreds.Tokens[name] = token
	// 	// 		}

	// 	// 		// Update expiration if M365 tokens expire earlier
	// 	// 		if m365Creds.ExpiresAt.Before(azureCreds.ExpiresAt) {
	// 	// 			azureCreds.ExpiresAt = m365Creds.ExpiresAt
	// 	// 		}
	// 	// 	}
	// 	// }
	// }

	// // Update current credentials
	// a.currentCreds = azureCreds

	// Save to store
	if err := a.saveToStore(ctx); err != nil {
		log.Infof("Warning: Failed to save renewed tokens: %v", err)
		// Continue anyway, as we've renewed successfully
	}

	log.Info("Token renewal completed successfully")

	return nil
}

// Clear clears all stored credentials and parameters
func (a *AuthManager) Clear(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Store == nil {
		return ErrStoreNotInitialized
	}

	// Clear the store
	if err := a.Store.Clear(ctx); err != nil {
		return fmt.Errorf("failed to clear credential store: %w", err)
	}

	// Clear current state
	a.currentAuthParams = nil
	a.currentCreds = nil
	a.m365Resources = nil

	log.Info("Cleared all authentication state")

	return nil
}

// GetAuthParams returns the current authentication parameters
func (a *AuthManager) GetAuthParams() *shared.AuthParams {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.currentAuthParams
}

// GetM365Resources returns the current M365 resources
func (a *AuthManager) GetM365Resources() *shared.M365Resources {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.m365Resources
}
