# How to extend `Lease`

## New External Service

```go
// CredentialService strings are used to determine what service to use for gathering credentials
const ExampleService CredentialService = "myService"

// AuthService strings are used to determine what service to use for external authentication
const ExampleService AuthService = "myService"

type ExampleCredentialFactory interface {
	DefaultCredentialFactory 								// Allows using defaults if nothing needs to change
	getSomethingSpecial(params *shared.AuthParams) string 	// Extend your CredentialFactory if needed
}

type ExampleAuthFactory interface {
	DefaultAuthFactory // Allows using defaults if nothing needs to change
}

type ExampleOptions struct {
	Posture AppPosture
	Method  AcquisitionMethod
}

// RealExampleCredentialFactory implements ExampleCredentialFactory for example authentication
type RealExampleCredentialFactory struct {
	DefaultCredentialFactory  // Implements all Methods in DefaultCredentialFactory for RealExampleCredentialFactory
}

// RealExampleAuthFactory implements ExampleAuthFactory for example authentication
type RealExampleAuthFactory struct {
	DefaultAuthFactory // Implements all Methods in DefaultAuthFactory for RealExampleAuthFactory
	CloudURL string
}
```
