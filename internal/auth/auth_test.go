package auth

// import (
// 	"context"
// 	"crypto/rand"
// 	"goslings/internal/auth/shared"
// 	"testing"
//  "github.com/stretchr/testify/assert"
// )

// // Mock implementation for testing
// type MockAzureService struct {
// 	// Mock data or behavior
// 	MockCreateResource func(ctx context.Context, resourceName string) (string, error)
// }

// func (m *MockAzureService) CreateResource(ctx context.Context, resourceName string) (string, error) {
// 	return m.MockCreateResource(ctx, resourceName)
// }

// func TestNewAuthManagerFileStore(t *testing.T) {
// 	ctx := t.Context()
// 	encryptionKey := make([]byte, 32)
// 	if _, err := rand.Read(encryptionKey); err != nil {
// 		t.Errorf("Failed to generate encryption key: %v", err)
// 	}

// 	want := map[string]any{
// 		"authp": &shared.AuthParams{
// 			Username:            "john",
// 			Password:            "john",
// 			TenantID:            "john",
// 			ClientID:            "john",
// 			ClientSecret:        "john",
// 			SubscriptionID:      "john",
// 			UsGovernment:        false,
// 			ExoUSGovernment:     false,
// 			M365Enabled:         false,
// 			MessageTraceEnabled: false,
// 		},
// 		"m365r": &shared.M365Resources{
// 			ExchangeCookies:  map[string]string{},
// 			ValidationKey:    "string",
// 			AdditionalTokens: map[string]string{},
// 		},
// 	}

// 	authmgr, err := NewAuthManager(Options{
// 		StoreType:     shared.FileStore,
// 		StorePath:     "./test/.credentials",
// 		EncryptionKey: encryptionKey,
// 	})

// 	authp := authmgr.GetAuthParams()

// 	if authp != want["authp"] {
// 		t.Errorf(`got "%+v", want "%+v"`, authp, want["authp"])
// 	}

// 	// TODO: test the returned struct
// 	m365r := authmgr.GetM365Resources()

// 	if m365r != want["m365r"] {
// 		t.Errorf(`got "%+v", want "%+v"`, m365r, want["m365r"])
// 	}

// 	if err = authmgr.Authenticate(ctx, authp); err != nil {
// 		t.Error(err)
// 	}

// 	// TODO: test the returned token
// 	if _, err = authmgr.GetToken("azure"); err != nil {
// 		t.Error(err)
// 	}

// 	if err = authmgr.saveToStore(ctx); err != nil {
// 		t.Error(err)
// 	}

// 	if err = authmgr.loadFromStore(ctx); err != nil {
// 		t.Error(err)
// 	}

// 	if err = authmgr.RenewTokens(ctx); err != nil {
// 		t.Error(err)
// 	}

// 	if err = authmgr.Clear(ctx); err != nil {
// 		t.Error(err)
// 	}

// }

// // func TestNewAuthManagerK8sStore(t *testing.T)   {}
// // func TestNewAuthManagerVaultStore(t *testing.T) {}
