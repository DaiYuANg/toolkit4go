package main

import (
	"context"
	"fmt"

	"github.com/DaiYuANg/arcgo/authx"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	primaryProvider := authx.NewInMemoryIdentityProvider()
	secondaryProvider, err := authx.NewFuncIdentityProvider(func(ctx context.Context, principal string) (authx.UserDetails, error) {
		_ = ctx
		if principal != "alice" {
			return authx.UserDetails{}, authx.ErrUnauthorized
		}
		return authx.UserDetails{
			ID:           "u-1",
			Principal:    "alice",
			PasswordHash: string(hashedPassword),
			Name:         "Alice",
		}, nil
	})
	if err != nil {
		panic(err)
	}

	policySource := authx.NewMemoryPolicySource(authx.MemoryPolicySourceConfig{
		Name: "quickstart-policy",
		InitialPermissions: []authx.PermissionRule{
			authx.AllowPermission("u-1", "order:1001", "read"),
			authx.AllowPermission("role:admin", "order:1001", "write"),
		},
		InitialRoleBindings: []authx.RoleBinding{
			authx.NewRoleBinding("u-1", "role:admin"),
		},
	})

	manager, err := authx.NewManager(
		authx.WithProviders(primaryProvider, secondaryProvider),
		authx.WithSources(policySource),
	)
	if err != nil {
		panic(err)
	}

	version, err := manager.LoadPolicies(context.Background())
	if err != nil {
		panic(err)
	}

	ctx, authentication, err := manager.AuthenticatePassword(context.Background(), "alice", "secret")
	if err != nil {
		panic(err)
	}

	allowed, err := manager.Can(ctx, "write", "order:1001")
	if err != nil {
		panic(err)
	}

	fmt.Printf("authenticated=%v user=%s policyVersion=%d allowed=%v\n",
		authentication.IsAuthenticated(),
		authentication.Identity().ID(),
		version,
		allowed,
	)
}
