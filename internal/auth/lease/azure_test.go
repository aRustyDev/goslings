package lease

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"goslings/internal/auth"
// 	"goslings/internal/auth/lease"
// )

// // TestAzureLeaseRenew
// // TestAzureLeaseAcquire
// // TestAzureLeaseIsExpired
// // TestAzureLease

// func TestAzureLeaseRenew(t *testing.T) {

// 	type testCase struct {
// 		name               string
// 		creds              *auth.Credentials
// 		params             *auth.AuthParams
// 		setupMock          func(*MockCredentialFactory)
// 		expectedTokenValue string
// 		wantErr            bool
// 	}

// 	// Test cases
// 	testCases := []testCase{}

// 	for _, tc := range testCases {
// 	}
// }

// func TestAzureLeaseAcquire(t *testing.T) {

// 	type testCase struct {
// 		name               string
// 		creds              *auth.Credentials
// 		params             *auth.AuthParams
// 		setupMock          func(*MockCredentialFactory)
// 		expectedTokenValue string
// 		wantErr            bool
// 	}

// 	// Test cases
// 	testCases := []testCase{}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Create a context with timeout
// 			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()

// 			// Create a mock credential factory
// 			mockFactory := NewMockCredentialFactory()

// 			// Setup the mock
// 			tc.setupMock(mockFactory)

// 			// Create the lease with the mock factory
// 			azureLease := lease.NewAzureLease(logger, true).WithCredentialFactory(mockFactory)

// 			// Add the Renew method implementation to the AzureLease

// 			// Call Renew
// 			renewedCreds, err := azureLease.Renew(ctx, tc.creds, tc.params)

// 			// Check error
// 			if (err != nil) != tc.wantErr {
// 				t.Errorf("Renew() error = %v, wantErr %v", err, tc.wantErr)
// 				return
// 			}

// 			if tc.wantErr {
// 				// If we expect an error, stop here
// 				return
// 			}

// 			// Check credentials
// 			if renewedCreds == nil {
// 				t.Errorf("Renew() returned nil credentials")
// 				return
// 			}

// 			// Check token value
// 			if token, ok := renewedCreds.Tokens["graph"]; !ok {
// 				t.Error("Renew() missing graph token")
// 			} else if token.Value != tc.expectedTokenValue {
// 				t.Errorf("Renew() token value = %v, want %v", token.Value, tc.expectedTokenValue)
// 			}

// 			// Check expiration
// 			if renewedCreds.ExpiresAt.Before(time.Now()) {
// 				t.Error("Renew() returned credentials that are already expired")
// 			}

// 			// Check last refreshed
// 			if time.Since(renewedCreds.LastRefreshed) > time.Minute {
// 				t.Error("Renew() returned credentials with old LastRefreshed time")
// 			}
// 		})
// 	}
// }

// func TestAzureLeaseIsExpired(t *testing.T) {

// 	// Create the lease
// 	azureLease := lease.NewAzureLease(true)

// 	type testCase struct {
// 		name        string
// 		creds       *auth.Credentials
// 		gracePeriod time.Duration
// 		want        bool
// 	}

// 	// Test cases
// 	testCases := []testCase{}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Create a context with timeout
// 			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()

// 			// Create a mock credential factory
// 			mockFactory := NewMockCredentialFactory()

// 			// Setup the mock
// 			tc.setupMock(mockFactory)

// 			// Create the lease with the mock factory
// 			azureLease := lease.NewAzureLease(logger, true).WithCredentialFactory(mockFactory)

// 			// Call Acquire
// 			creds, err := azureLease.Acquire(ctx, tc.params)

// 			// Check error
// 			if (err != nil) != tc.wantErr {
// 				t.Errorf("Acquire() error = %v, wantErr %v", err, tc.wantErr)
// 				return
// 			}

// 			if tc.wantErr {
// 				// If we expect an error, stop here
// 				return
// 			}

// 			// Check credentials
// 			if creds == nil {
// 				t.Errorf("Acquire() returned nil credentials")
// 				return
// 			}

// 			// Check auth type
// 			if creds.AuthType != tc.expectedAuthType {
// 				t.Errorf("Acquire() auth type = %v, want %v", creds.AuthType, tc.expectedAuthType)
// 			}

// 			// Check tokens
// 			for _, tokenName := range tc.expectedTokens {
// 				if _, ok := creds.Tokens[tokenName]; !ok {
// 					t.Errorf("Acquire() missing token %s", tokenName)
// 				}
// 			}

// 			// Check expiration
// 			if creds.ExpiresAt.IsZero() {
// 				t.Error("Acquire() returned credentials with zero expiration time")
// 			}

// 			// Check that the correct mock was called based on auth type
// 			switch creds.AuthType {
// 			case auth.DeviceCodeAuth:
// 				if mockFactory.DeviceCodeCallCount != 1 {
// 					t.Errorf("Expected device code to be called once, got %d", mockFactory.DeviceCodeCallCount)
// 				}
// 			case auth.ClientCredentialsAuth:
// 				if mockFactory.ClientSecretCallCount != 1 {
// 					t.Errorf("Expected client secret to be called once, got %d", mockFactory.ClientSecretCallCount)
// 				}
// 			case auth.InteractiveAuth:
// 				if mockFactory.InteractiveBrowserCallCount != 1 {
// 					t.Errorf("Expected interactive browser to be called once, got %d", mockFactory.InteractiveBrowserCallCount)
// 				}
// 			}
// 		})
// 	}
// }
