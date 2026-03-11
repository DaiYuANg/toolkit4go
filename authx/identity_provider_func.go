package authx

import (
	"context"
	"fmt"
)

// IdentityProviderFunc adapts a function to IdentityProvider.
type IdentityProviderFunc func(ctx context.Context, principal string) (UserDetails, error)

// LoadByPrincipal implements IdentityProvider.
func (f IdentityProviderFunc) LoadByPrincipal(ctx context.Context, principal string) (UserDetails, error) {
	if f == nil {
		return UserDetails{}, fmt.Errorf("%w: identity provider function is nil", ErrInvalidAuthenticator)
	}
	return f(ctx, principal)
}

// NewFuncIdentityProvider creates an IdentityProvider from a function.
func NewFuncIdentityProvider(fn func(ctx context.Context, principal string) (UserDetails, error)) (IdentityProvider, error) {
	if fn == nil {
		return nil, fmt.Errorf("%w: identity provider function is nil", ErrInvalidAuthenticator)
	}
	return IdentityProviderFunc(fn), nil
}
