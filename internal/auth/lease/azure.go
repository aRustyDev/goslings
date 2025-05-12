package lease

//go:generate mockgen . AzureFactory
// -source=azure.go -destination=azure_mocks_test.go -package=lease

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	private "github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	managed "github.com/AzureAD/microsoft-authentication-library-for-go/apps/managedidentity"
	public "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/arustydev/goslings/internal/auth/shared"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Where is the app running from, client side, private server, or azure managed service?
type AppPosture string

const (
	Public  AppPosture = "public"
	Private AppPosture = "confidential"
	Managed AppPosture = "managed"

	Azure LeaseProvider = "azure"
	M365  LeaseProvider = "m365"
	D4iot LeaseProvider = "d4iot"

	DeviceCode         AcquisitionMethod = "devicecode"         // acquires a security token from the authority, by acquiring a device code and using that to acquire the token. Users need to create an AcquireTokenDeviceCodeParameters instance and pass it in.
	ClientSecret       AcquisitionMethod = "clientsecret"       // acquires a security token from the authority, using an authorization code. The specified redirect URI must be the same URI that was used when the authorization code was requested
	InteractiveBrowser AcquisitionMethod = "interactivebrowser" // acquires a security token from the authority using the default web browser to select the account
	UserPass           AcquisitionMethod = "userpass"           // acquires a security token from the authority, via Username/Password Authentication.
	Silent             AcquisitionMethod = "silent"             // acquires a token from either the cache or using a refresh token
	Credential         AcquisitionMethod = "credential"         //
	// invokes a callback to get assertions authenticating the application
)

type AzureCredentialFactory interface {
	DefaultCredentialFactory
	getCloudURL(params *shared.AuthParams) string
}

type AzureAuthFactory interface {
	DefaultAuthFactory
}

type AzureOptions struct {
	Posture  AppPosture
	Method   AcquisitionMethod
	ClientId string
}

// RealAzureAuthFactory implements CredentialFactoryfor Azure authentication
type RealAzureAuthFactory struct {
	DefaultAuthFactory
	CloudURL string
	Options  *AzureOptions
	Params   *shared.AuthParams
}

func (f *RealAzureAuthFactory) GetToken(
	ctx context.Context,
	opts policy.TokenRequestOptions,
) error {
	// Set the cloud URL based on parameters
	f.CloudURL = f.getCloudURL(f.Params)

	switch f.Options.Posture {
	case Public:
		client, err := public.New(f.Options.ClientId, public.WithHTTPClient(f.Client))
		if err != nil {
			log.Fatalf("%+v", err)
		}
		switch f.Options.Method {
		case DeviceCode:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenByDeviceCode(ctx, scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		case ClientSecret:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenByAuthCode(ctx, "code", "redirectURI", scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		case UserPass:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenByUsernamePassword(ctx, scopes, "username", "password")
			if err != nil {
				log.Fatalf("%+v", err)
			}
		case Silent:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenSilent(ctx, scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		default: // InteractiveBrowser:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenInteractive(ctx, scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}
	case Private:
		var cred private.Credential
		var err error
		switch f.Options.Method {
		case "assertion":
			// https://pkg.go.dev/github.com/AzureAD/microsoft-authentication-library-for-go@v1.4.2/apps/confidential#NewCredFromAssertionCallback
			// cred = private.NewCredFromAssertionCallback("secret") // takes a callback func
			log.Fatal("NOT IMPLEMENTED YET")
		case "token-provider":
			// https://pkg.go.dev/github.com/AzureAD/microsoft-authentication-library-for-go@v1.4.2/apps/confidential#NewCredFromTokenProvider
			// creates a Credential from a function that provides access tokens
			// cred = private.NewCredFromTokenProvider("secret") // takes a callback func
			log.Fatal("NOT IMPLEMENTED YET")
		case "certificate":
			b, err := os.ReadFile(viper.GetString("GOSLING_AZURE_CONFIDENTIAL_CERT_PATH"))
			if err != nil {
				log.Fatalf("%+v", err)
			}

			// This extracts our public certificates and private key from the PEM file. If it is
			// encrypted, the second argument must be password to decode.
			certs, priv, err := private.CertFromPEM(b, viper.GetString("GOSLING_AZURE_CONFIDENTIAL_CERT_PASS"))
			if err != nil {
				log.Fatalf("%+v", err)
			}
			cred, err = private.NewCredFromCert(certs, priv)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		default:
			cred, err = private.NewCredFromSecret("secret")
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}
		client, err := private.New("authority", f.Options.ClientId, cred, private.WithHTTPClient(f.Client))
		if err != nil {
			log.Fatalf("%+v", err)
		}
		switch f.Options.Method {
		case ClientSecret:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenByAuthCode(ctx, "code", "redirectURI", scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		case Credential:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenByCredential(ctx, scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		case UserPass:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenByUsernamePassword(ctx, scopes, viper.GetString("GOSLING_AZURE_CONFIDENTIAL_USER"), viper.GetString("GOSLING_AZURE_CONFIDENTIAL_PASS"))
			if err != nil {
				log.Fatalf("%+v", err)
			}
		case Silent:
			scopes := []string{"https://graph.microsoft.com/.default"}
			_, err = client.AcquireTokenSilent(ctx, scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		default: // InteractiveBrowser:
			scopes := []string{"https://graph.microsoft.com/.default"}
			// TODO: verify "userAssertion" values
			// https: //pkg.go.dev/github.com/AzureAD/microsoft-authentication-library-for-go@v1.4.2/apps/confidential#Client.AcquireTokenOnBehalfOf
			_, err = client.AcquireTokenOnBehalfOf(ctx, "userAssertion", scopes)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}
	default: // Managed
		client, err := managed.New(managed.SystemAssigned(), managed.WithHTTPClient(f.Client))
		if err != nil {
			log.Fatalf("%+v", err)
		}
		_, err = client.AcquireToken(ctx, "resource")
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
	f.Token = &azcore.AccessToken{}
	return nil
}

// getCloudURL returns the appropriate cloud URL based on params
func (f *RealAzureAuthFactory) getCloudURL(params *shared.AuthParams) string {
	if params.UsGovernment {
		return "https://login.microsoftonline.us"
	}
	return "https://login.microsoftonline.com"
}

// RealAzureAuthFactory implements CredentialFactoryfor Azure authentication
type RealAzureCredentialFactory struct {
	DefaultCredentialFactory
}

func (f *RealAzureCredentialFactory) GetCredential(
	ctx context.Context,
	AcquisitionMethod AcquisitionMethod,
	options *CredentialOptions,
) error {
	switch AcquisitionMethod {
	case DeviceCode:
		azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
			TenantID:   options.TenantID,
			ClientID:   options.ClientID,
			UserPrompt: options.UserPrompt,
		})
	case ClientSecret:
		azidentity.NewClientSecretCredential(
			options.TenantID,
			options.ClientID,
			options.ClientSecret,
			&azidentity.ClientSecretCredentialOptions{})
	case InteractiveBrowser:
		azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
			TenantID: options.TenantID,
			ClientID: options.ClientID,
		})
	}
	return nil
}
