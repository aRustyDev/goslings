package cmd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"goslings/internal/auth"
	"goslings/internal/auth/shared"
	"goslings/internal/conf"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// Common errors
var (
	ErrConfigNotSet = errors.New("A Config parameter was not set")
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "middle command for authenticating to the cloud",
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	ctx := cmd.Context()
	// 	value := ctx.Value("auth")
	// 	log.Infof("Auth(RunE): %s", value) // Output: value
	// 	return nil
	// },
	Run: func(cmd *cobra.Command, args []string) {
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

		// Authenticate
		log.Info("Starting authentication process")
		if err := authManager.Authenticate(cmd.Context(), conf.GetAuthConfig()); err != nil {
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
		if err := authManager.RenewTokens(cmd.Context()); err != nil {
			log.Fatalf("Failed to renew tokens: %v", err)
		}

		log.Info("Authentication example completed successfully")
	},
}
