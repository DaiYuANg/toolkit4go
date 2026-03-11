package authxhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DaiYuANg/arcgo/authx"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticateInjectsSecurityContext(t *testing.T) {
	manager := newTestManager(t)

	handler := Authenticate(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, err := authx.RequireIdentity(r.Context())
		if err != nil {
			t.Fatalf("require identity: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(identity.ID()))
	}))

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.SetBasicAuth("alice", "secret")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if rec.Body.String() != "u1" {
		t.Fatalf("unexpected identity id: %q", rec.Body.String())
	}
}

func TestAuthenticateMissingCredentialReturnsUnauthorized(t *testing.T) {
	manager := newTestManager(t)

	handler := Authenticate(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestAuthenticateOptionalPassesThrough(t *testing.T) {
	manager := newTestManager(t)

	handler := Authenticate(manager, WithOptional())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestRequireAllowsAuthorizedRequest(t *testing.T) {
	manager := newTestManager(t)

	protected := Authenticate(manager)(Require(manager, "read", "orders")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})))

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.SetBasicAuth("alice", "secret")
	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestRequireReturnsForbiddenWhenDenied(t *testing.T) {
	manager := newTestManager(t)

	protected := Authenticate(manager)(Require(manager, "write", "orders")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})))

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.SetBasicAuth("alice", "secret")
	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestRequireWithoutAuthenticationReturnsUnauthorized(t *testing.T) {
	manager := newTestManager(t)

	protected := Require(manager, "read", "orders")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func newTestManager(t *testing.T) authx.Manager {
	t.Helper()

	provider := authx.NewInMemoryIdentityProvider()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate hash: %v", err)
	}

	if err := provider.UpsertUser(authx.UserDetails{
		ID:           "u1",
		Principal:    "alice",
		PasswordHash: string(hashedPassword),
		Name:         "Alice",
	}); err != nil {
		t.Fatalf("upsert user: %v", err)
	}

	source := authx.NewMemoryPolicySource(authx.MemoryPolicySourceConfig{
		Name: "middleware-test-policy",
		InitialPermissions: []authx.PermissionRule{
			{
				Subject:  "u1",
				Resource: "orders",
				Action:   "read",
				Allowed:  true,
			},
		},
	})

	manager, err := authx.NewManager(
		authx.WithProvider(provider),
		authx.WithSource(source),
	)
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	if _, err := manager.LoadPolicies(context.Background()); err != nil {
		t.Fatalf("load policies: %v", err)
	}
	return manager
}
