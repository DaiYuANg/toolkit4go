package authx

import "context"

// AuthenticationRuntime defines authentication entrypoints.
type AuthenticationRuntime interface {
	Authenticate(ctx context.Context, credential Credential) (context.Context, Authentication, error)
	AuthenticatePassword(ctx context.Context, username, password string) (context.Context, Authentication, error)
}

// AuthorizationRuntime defines authorization entrypoints.
type AuthorizationRuntime interface {
	Can(ctx context.Context, action, resource string) (bool, error)
}

// PolicyRuntime defines policy loading and replacement operations.
type PolicyRuntime interface {
	LoadPolicies(ctx context.Context) (int64, error)
	LoadPoliciesFrom(ctx context.Context, source PolicySource) (int64, error)
	ReplacePolicies(ctx context.Context, snapshot PolicySnapshot) (int64, error)
	PolicyVersion() int64
}

// ProviderRegistry defines identity provider chain management operations.
type ProviderRegistry interface {
	SetIdentityProviders(providers ...IdentityProvider) error
	AddIdentityProvider(provider IdentityProvider) error
}

// PolicySourceRegistry defines policy source chain management operations.
type PolicySourceRegistry interface {
	SetPolicySources(sources ...PolicySource) error
	AddPolicySource(source PolicySource) error
}

// Manager is the public authx runtime contract.
// It intentionally hides concrete implementation details.
type Manager interface {
	AuthenticationRuntime
	AuthorizationRuntime
	PolicyRuntime
	ProviderRegistry
	PolicySourceRegistry
}
