package authx

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/samber/lo"
	"github.com/samber/oops"
)

// ProviderManager routes authentication credential to provider by credential concrete type.
type ProviderManager struct {
	mu        sync.RWMutex
	providers map[reflect.Type]AuthenticationProvider
}

// NewProviderManager constructs a ProviderManager and registers providers.
func NewProviderManager(providers ...AuthenticationProvider) *ProviderManager {
	manager := &ProviderManager{providers: make(map[reflect.Type]AuthenticationProvider)}
	manager.Register(providers...)
	return manager
}

// Register adds providers keyed by their credential type.
func (manager *ProviderManager) Register(providers ...AuthenticationProvider) {
	if manager == nil || len(providers) == 0 {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	lo.ForEach(providers, func(provider AuthenticationProvider, _ int) {
		if provider == nil {
			return
		}
		credentialType := provider.CredentialType()
		if credentialType == nil {
			return
		}
		manager.providers[credentialType] = provider
	})
}

// Authenticate dispatches credential to the registered provider for its concrete type.
func (manager *ProviderManager) Authenticate(
	ctx context.Context,
	credential any,
) (AuthenticationResult, error) {
	if credential == nil {
		return AuthenticationResult{}, oops.In("authx").
			With("op", "authenticate", "stage", "validate_credential").
			Wrapf(ErrInvalidAuthenticationCredential, "validate authentication credential")
	}
	if manager == nil {
		return AuthenticationResult{}, oops.In("authx").
			With("op", "authenticate", "stage", "validate_manager").
			Wrapf(ErrAuthenticationManagerNotConfigured, "validate authentication manager")
	}

	credentialType := reflect.TypeOf(credential)
	manager.mu.RLock()
	provider, ok := manager.providers[credentialType]
	providerCount := len(manager.providers)
	manager.mu.RUnlock()
	if !ok {
		return AuthenticationResult{}, oops.In("authx").
			With(
				"op", "authenticate",
				"stage", "resolve_provider",
				"credential_type", credentialType,
				"provider_count", providerCount,
			).
			Wrapf(ErrAuthenticationProviderNotFound, "resolve authentication provider")
	}

	result, err := provider.AuthenticateAny(ctx, credential)
	if err != nil {
		return AuthenticationResult{}, oops.In("authx").
			With("op", "authenticate", "stage", "provider_authenticate", "credential_type", credentialType).
			Wrapf(errors.Join(ErrUnauthenticated, err), "authenticate credential")
	}
	return result, nil
}
