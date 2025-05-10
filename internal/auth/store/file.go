// Package store provides interfaces and implementations for storing and retrieving credentials
package store

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/arustydev/goslings/internal/auth/shared"
	"golang.org/x/crypto/nacl/secretbox"
)

// Store is the interface that wraps the basic credential storage and retrieval methods
type Store interface {
	// StoreCredentials stores credentials in the backing store
	StoreCredentials(ctx context.Context, creds *shared.Credentials) error

	// LoadCredentials retrieves credentials from the backing store
	LoadCredentials(ctx context.Context) (*shared.Credentials, error)

	// StoreParams stores authentication parameters in the backing store
	StoreParams(ctx context.Context, params *shared.AuthParams) error

	// LoadParams retrieves authentication parameters from the backing store
	LoadParams(ctx context.Context) (*shared.AuthParams, error)

	// StoreM365Resources stores M365-specific resources in the backing store
	StoreM365Resources(ctx context.Context, resources *shared.M365Resources) error

	// LoadM365Resources retrieves M365-specific resources from the backing store
	LoadM365Resources(ctx context.Context) (*shared.M365Resources, error)

	// Clear removes all stored credentials and parameters
	Clear(ctx context.Context) error
}

// FileStore implements Store using encrypted local files
type FileStore struct {
	// BasePath is the directory where credential files will be stored
	BasePath string

	// EncryptionKey is the key used for encrypting sensitive data
	EncryptionKey [32]byte
}

// FileNames for different storage files
const (
	CredsFileName  = "credentials.enc"
	ParamsFileName = "params.enc"
	M365FileName   = "m365.enc"
)

// NewFileStore creates a new file-based credential store
func NewFileStore(basePath string, encryptionKey []byte) (*FileStore, error) {
	// Validate encryption key
	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key must be exactly 32 bytes")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0o700); err != nil {
		return nil, fmt.Errorf("failed to create directory for credentials: %w", err)
	}

	var key [32]byte
	copy(key[:], encryptionKey)

	return &FileStore{
		BasePath:      basePath,
		EncryptionKey: key,
	}, nil
}

// encrypt encrypts data using secretbox
func (fs *FileStore) encrypt(data []byte) ([]byte, error) {
	// Generate a random nonce
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	encrypted := secretbox.Seal(nonce[:], data, &nonce, &fs.EncryptionKey)
	return encrypted, nil
}

// decrypt decrypts data using secretbox
func (fs *FileStore) decrypt(data []byte) ([]byte, error) {
	if len(data) < 24 {
		return nil, errors.New("data too short to contain nonce")
	}

	// Extract the nonce
	var nonce [24]byte
	copy(nonce[:], data[:24])

	// Decrypt the data
	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, &fs.EncryptionKey)
	if !ok {
		return nil, errors.New("failed to decrypt data")
	}

	return decrypted, nil
}

// StoreCredentials implements Store.StoreCredentials for FileStore
func (fs *FileStore) StoreCredentials(ctx context.Context, creds *shared.Credentials) error {
	// Marshal credentials to JSON
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Encrypt the data
	encrypted, err := fs.encrypt(data)
	if err != nil {
		return err
	}

	// Write to file
	path := filepath.Join(fs.BasePath, CredsFileName)
	if err := os.WriteFile(path, encrypted, 0o600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// LoadCredentials implements Store.LoadCredentials for FileStore
func (fs *FileStore) LoadCredentials(ctx context.Context) (*shared.Credentials, error) {
	path := filepath.Join(fs.BasePath, CredsFileName)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("credentials file not found")
	}

	// Read file
	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Decrypt the data
	data, err := fs.decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	// Unmarshal credentials
	var creds shared.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &creds, nil
}

// StoreParams implements Store.StoreParams for FileStore
func (fs *FileStore) StoreParams(ctx context.Context, params *shared.AuthParams) error {
	// Marshal parameters to JSON
	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal parameters: %w", err)
	}

	// Encrypt the data
	encrypted, err := fs.encrypt(data)
	if err != nil {
		return err
	}

	// Write to file
	path := filepath.Join(fs.BasePath, ParamsFileName)
	if err := os.WriteFile(path, encrypted, 0o600); err != nil {
		return fmt.Errorf("failed to write parameters file: %w", err)
	}

	return nil
}

// LoadParams implements Store.LoadParams for FileStore
func (fs *FileStore) LoadParams(ctx context.Context) (*shared.AuthParams, error) {
	path := filepath.Join(fs.BasePath, ParamsFileName)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("parameters file not found")
	}

	// Read file
	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read parameters file: %w", err)
	}

	// Decrypt the data
	data, err := fs.decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	// Unmarshal parameters
	var params shared.AuthParams
	if err := json.Unmarshal(data, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}

	return &params, nil
}

// StoreM365Resources implements Store.StoreM365Resources for FileStore
func (fs *FileStore) StoreM365Resources(
	ctx context.Context,
	resources *shared.M365Resources,
) error {
	// Marshal M365 resources to JSON
	data, err := json.Marshal(resources)
	if err != nil {
		return fmt.Errorf("failed to marshal M365 resources: %w", err)
	}

	// Encrypt the data
	encrypted, err := fs.encrypt(data)
	if err != nil {
		return err
	}

	// Write to file
	path := filepath.Join(fs.BasePath, M365FileName)
	if err := os.WriteFile(path, encrypted, 0o600); err != nil {
		return fmt.Errorf("failed to write M365 resources file: %w", err)
	}

	return nil
}

// LoadM365Resources implements Store.LoadM365Resources for FileStore
func (fs *FileStore) LoadM365Resources(ctx context.Context) (*shared.M365Resources, error) {
	path := filepath.Join(fs.BasePath, M365FileName)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("M365 resources file not found")
	}

	// Read file
	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read M365 resources file: %w", err)
	}

	// Decrypt the data
	data, err := fs.decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	// Unmarshal M365 resources
	var resources shared.M365Resources
	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, fmt.Errorf("failed to unmarshal M365 resources: %w", err)
	}

	return &resources, nil
}

// Clear implements Store.Clear for FileStore
func (fs *FileStore) Clear(ctx context.Context) error {
	files := []string{
		filepath.Join(fs.BasePath, CredsFileName),
		filepath.Join(fs.BasePath, ParamsFileName),
		filepath.Join(fs.BasePath, M365FileName),
	}

	var firstErr error
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			if err := os.Remove(file); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

// Stub implementations for K8sStore and VaultStore would go here
// These would implement the Store interface with their respective backends
