// Package lease provides interfaces and implementations for acquiring and renewing authentication tokens
package lease

import (
	"errors"
)

var (
	ErrNotAuthenticated   = errors.New("not authenticated")
	ErrCredentialsExpired = errors.New("credentials expired")
)
