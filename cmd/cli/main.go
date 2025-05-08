package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"goslings/internal/app/cli/cmd"

	"time"

	"goslings/internal/auth"
	"goslings/internal/auth/shared"

	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.Goodbye("from the cli!")
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Setup logging
	log.Info("Starting authentication example")

	// Generate a random encryption key (in a real app, you would store this securely)
	encryptionKey := make([]byte, 32)
	if _, err := rand.Read(encryptionKey); err != nil {
		log.Fatalf("Failed to generate encryption key: %v", err)
	}

	// Initialize the auth manager
	authManager, err := auth.NewAuthManager(auth.Options{
		StoreType:     shared.FileStore,
		StorePath:     "./.credentials",
		EncryptionKey: encryptionKey,
		Debug:         true,
	})
	if err != nil {
		log.Fatalf("Failed to create auth manager: %v", err)
	}

	// Get authentication parameters from config
	// In a real application, this would use your viper-based config package
	configParams := getConfigFromViper()

	// Authenticate
	log.Info("Starting authentication process")
	if err := authManager.Authenticate(ctx, configParams); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Successfully authenticated
	log.Info("Authentication successful!")

	// Example of getting a token for a service
	token, err := authManager.GetToken(auth.GraphService)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	log.Infof("Got token for Graph API, expires at: %v", token.ExpiresAt)

	// Example of using the token in an API call
	// In a real application, you would use this token to make API calls
	fmt.Printf("Token: %s %s\n", token.Type, token.Value[:10]+"...")

	// Simulate some time passing
	log.Info("Simulating time passing...")
	time.Sleep(5 * time.Second)

	// Renew tokens if needed
	log.Info("Checking if tokens need renewal")
	if err := authManager.RenewTokens(ctx); err != nil {
		log.Fatalf("Failed to renew tokens: %v", err)
	}

	log.Info("Authentication example completed successfully")
}

// getConfigFromViper simulates getting config from a viper-based config package
func getConfigFromViper() *shared.AuthParams {
	// In a real application, this would use your viper-based config package
	// For this example, we'll just return hardcoded values
	return &shared.AuthParams{
		Username:            "user@example.com",
		Password:            "password",
		TenantID:            "your-tenant-id",
		ClientID:            "your-client-id",
		ClientSecret:        "your-client-secret",
		SubscriptionID:      "your-subscription-id",
		UsGovernment:        false,
		ExoUSGovernment:     false,
		M365Enabled:         true,
		MessageTraceEnabled: false,
	}
}

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	cmd.Goodbye(name)
	return message
}
