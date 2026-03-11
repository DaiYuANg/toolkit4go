package authx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DaiYuANg/arcgo/observabilityx"
)

// ManagerOption configures manager construction.
type ManagerOption func(cfg *managerConfig) error

type managerConfig struct {
	providers      []IdentityProvider
	sources        []PolicySource
	policyMerger   PolicyMerger
	casbinOptions  []CasbinAuthorizerOption
	logger         *slog.Logger
	observability  observabilityx.Observability
	eventPublisher *EventPublisher
}

// WithProvider appends one authentication provider to provider chain.
func WithProvider(provider IdentityProvider) ManagerOption {
	return func(cfg *managerConfig) error {
		if cfg == nil {
			return fmt.Errorf("%w: manager config is nil", ErrInvalidAuthenticator)
		}
		if provider == nil {
			return fmt.Errorf("%w: identity provider is nil", ErrInvalidAuthenticator)
		}
		cfg.providers = append(cfg.providers, provider)
		return nil
	}
}

// WithProviderFunc appends one function-based authentication provider.
func WithProviderFunc(fn func(ctx context.Context, principal string) (UserDetails, error)) ManagerOption {
	return func(cfg *managerConfig) error {
		provider, err := NewFuncIdentityProvider(fn)
		if err != nil {
			return err
		}
		return WithProvider(provider)(cfg)
	}
}

// WithProviders appends multiple authentication providers to provider chain.
func WithProviders(providers ...IdentityProvider) ManagerOption {
	return func(cfg *managerConfig) error {
		for _, provider := range providers {
			if err := WithProvider(provider)(cfg); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithSource appends one policy source.
func WithSource(source PolicySource) ManagerOption {
	return func(cfg *managerConfig) error {
		if cfg == nil {
			return fmt.Errorf("%w: manager config is nil", ErrInvalidPolicy)
		}
		if source == nil {
			return fmt.Errorf("%w: policy source is nil", ErrInvalidPolicy)
		}
		cfg.sources = append(cfg.sources, source)
		return nil
	}
}

// WithSources appends multiple policy sources.
func WithSources(sources ...PolicySource) ManagerOption {
	return func(cfg *managerConfig) error {
		for _, source := range sources {
			if err := WithSource(source)(cfg); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithPolicyMerger sets the policy merger used by PolicySourceChain.
func WithPolicyMerger(merger PolicyMerger) ManagerOption {
	return func(cfg *managerConfig) error {
		if cfg == nil {
			return fmt.Errorf("%w: manager config is nil", ErrInvalidPolicy)
		}
		if merger == nil {
			return fmt.Errorf("%w: policy merger is nil", ErrInvalidPolicy)
		}
		cfg.policyMerger = merger
		return nil
	}
}

// WithCasbinOptions sets additional options used to build casbin authorizer.
func WithCasbinOptions(opts ...CasbinAuthorizerOption) ManagerOption {
	return func(cfg *managerConfig) error {
		if cfg == nil {
			return fmt.Errorf("%w: manager config is nil", ErrInvalidAuthorizer)
		}
		for _, opt := range opts {
			if opt != nil {
				cfg.casbinOptions = append(cfg.casbinOptions, opt)
			}
		}
		return nil
	}
}

// WithLogger sets slog logger used by authx runtime nodes.
func WithLogger(logger *slog.Logger) ManagerOption {
	return func(cfg *managerConfig) error {
		if cfg == nil {
			return fmt.Errorf("%w: manager config is nil", ErrInvalidAuthenticator)
		}
		cfg.logger = logger
		return nil
	}
}

// WithObservability sets optional observability integration.
func WithObservability(obs observabilityx.Observability) ManagerOption {
	return func(cfg *managerConfig) error {
		if cfg == nil {
			return fmt.Errorf("%w: manager config is nil", ErrInvalidAuthenticator)
		}
		cfg.observability = obs
		return nil
	}
}

// TypedProviderLoader loads typed domain payload by principal.
type TypedProviderLoader[T any] interface {
	LoadByPrincipal(ctx context.Context, principal string) (T, error)
}

// TypedProviderMapper maps typed domain payload into AuthX user details.
type TypedProviderMapper[T any] interface {
	MapToUserDetails(ctx context.Context, principal string, payload T) (UserDetails, error)
}

// MappedProvider is a generic provider contract for custom user model mapping.
type MappedProvider[T any] interface {
	TypedProviderLoader[T]
	TypedProviderMapper[T]
}

type mappedIdentityProvider[T any] struct {
	provider MappedProvider[T]
}

// NewMappedIdentityProvider creates a generic strong-typed provider adapter.
func NewMappedIdentityProvider[T any](provider MappedProvider[T]) (IdentityProvider, error) {
	if provider == nil {
		return nil, fmt.Errorf("%w: mapped provider is nil", ErrInvalidAuthenticator)
	}

	return &mappedIdentityProvider[T]{
		provider: provider,
	}, nil
}

// WithMappedProvider appends a generic strong-typed provider adapter.
func WithMappedProvider[T any](provider MappedProvider[T]) ManagerOption {
	return func(cfg *managerConfig) error {
		adaptedProvider, err := NewMappedIdentityProvider(provider)
		if err != nil {
			return err
		}
		return WithProvider(adaptedProvider)(cfg)
	}
}

func (p *mappedIdentityProvider[T]) LoadByPrincipal(ctx context.Context, principal string) (UserDetails, error) {
	if p == nil || p.provider == nil {
		return UserDetails{}, fmt.Errorf("%w: mapped provider is not configured", ErrInvalidAuthenticator)
	}

	payload, err := p.provider.LoadByPrincipal(ctx, principal)
	if err != nil {
		return UserDetails{}, err
	}

	details, err := p.provider.MapToUserDetails(ctx, principal, payload)
	if err != nil {
		return UserDetails{}, err
	}
	details.Payload = payload
	return details, nil
}

// WithEventPublisher sets a custom event publisher for authx events.
// If not provided, a default event publisher will be created.
func WithEventPublisher(publisher *EventPublisher) ManagerOption {
	return func(cfg *managerConfig) error {
		if cfg == nil {
			return fmt.Errorf("%w: manager config is nil", ErrInvalidAuthenticator)
		}
		cfg.eventPublisher = publisher
		return nil
	}
}
