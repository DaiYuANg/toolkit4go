package authxhttp

import (
	"errors"
	"net/http"
	"strings"

	"github.com/DaiYuANg/arcgo/authx"
)

const (
	// DefaultAPIKeyHeader is the default header key for API key extraction.
	DefaultAPIKeyHeader = "X-API-Key"
)

var (
	// ErrCredentialNotFound indicates no supported credential exists in request.
	ErrCredentialNotFound = errors.New("authxhttp: credential not found")
	// ErrForbidden indicates authorization denied in middleware layer.
	ErrForbidden = errors.New("authxhttp: forbidden")
)

// CredentialExtractor extracts auth credential from request.
type CredentialExtractor func(r *http.Request) (authx.Credential, error)

// ErrorHandler handles unauthorized/forbidden response mapping.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type config struct {
	optional       bool
	apiKeyHeader   string
	extractor      CredentialExtractor
	onUnauthorized ErrorHandler
	onForbidden    ErrorHandler
}

// Option configures middleware behavior.
type Option func(cfg *config)

// WithOptional marks authentication middleware as optional.
func WithOptional() Option {
	return func(cfg *config) {
		cfg.optional = true
	}
}

// WithAPIKeyHeader overrides default API key header key.
func WithAPIKeyHeader(header string) Option {
	return func(cfg *config) {
		trimmed := strings.TrimSpace(header)
		if trimmed != "" {
			cfg.apiKeyHeader = trimmed
		}
	}
}

// WithCredentialExtractor sets custom credential extractor.
func WithCredentialExtractor(extractor CredentialExtractor) Option {
	return func(cfg *config) {
		if extractor != nil {
			cfg.extractor = extractor
		}
	}
}

// WithUnauthorizedHandler sets custom 401 handler.
func WithUnauthorizedHandler(handler ErrorHandler) Option {
	return func(cfg *config) {
		if handler != nil {
			cfg.onUnauthorized = handler
		}
	}
}

// WithForbiddenHandler sets custom 403 handler.
func WithForbiddenHandler(handler ErrorHandler) Option {
	return func(cfg *config) {
		if handler != nil {
			cfg.onForbidden = handler
		}
	}
}

// Authenticate performs credential extraction and authx authentication.
func Authenticate(manager authx.Manager, opts ...Option) func(http.Handler) http.Handler {
	cfg := defaultConfig(opts...)

	return func(next http.Handler) http.Handler {
		if next == nil {
			next = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r == nil || manager == nil {
				cfg.onUnauthorized(w, r, authx.ErrInvalidAuthenticator)
				return
			}

			credential, err := cfg.extractor(r)
			if err != nil {
				if errors.Is(err, ErrCredentialNotFound) && cfg.optional {
					next.ServeHTTP(w, r)
					return
				}
				cfg.onUnauthorized(w, r, err)
				return
			}

			ctx, _, err := manager.Authenticate(r.Context(), credential)
			if err != nil {
				cfg.onUnauthorized(w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Require checks authorization decision by action/resource.
func Require(manager authx.Manager, action, resource string, opts ...Option) func(http.Handler) http.Handler {
	cfg := defaultConfig(opts...)

	normalizedAction := strings.TrimSpace(action)
	normalizedResource := strings.TrimSpace(resource)

	return func(next http.Handler) http.Handler {
		if next == nil {
			next = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r == nil || manager == nil {
				cfg.onUnauthorized(w, r, authx.ErrInvalidAuthorizer)
				return
			}

			allowed, err := manager.Can(r.Context(), normalizedAction, normalizedResource)
			if err != nil {
				if authx.IsForbidden(err) {
					cfg.onForbidden(w, r, err)
					return
				}
				cfg.onUnauthorized(w, r, err)
				return
			}
			if !allowed {
				cfg.onForbidden(w, r, ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// DefaultCredentialExtractor extracts credential in this order:
// 1) HTTP Basic auth
// 2) API key header
// 3) Authorization: Bearer <token> (mapped to APIKeyCredential)
func DefaultCredentialExtractor(apiKeyHeader string) CredentialExtractor {
	headerKey := strings.TrimSpace(apiKeyHeader)
	if headerKey == "" {
		headerKey = DefaultAPIKeyHeader
	}

	return func(r *http.Request) (authx.Credential, error) {
		if r == nil {
			return nil, ErrCredentialNotFound
		}

		username, password, ok := r.BasicAuth()
		if ok {
			username = strings.TrimSpace(username)
			password = strings.TrimSpace(password)
			if username != "" && password != "" {
				return authx.PasswordCredential{
					Username: username,
					Password: password,
				}, nil
			}
		}

		if apiKey := strings.TrimSpace(r.Header.Get(headerKey)); apiKey != "" {
			return authx.APIKeyCredential{Key: apiKey}, nil
		}

		bearer := extractBearerToken(r.Header.Get("Authorization"))
		if bearer != "" {
			return authx.APIKeyCredential{Key: bearer}, nil
		}

		return nil, ErrCredentialNotFound
	}
}

func defaultConfig(opts ...Option) config {
	cfg := config{
		apiKeyHeader:   DefaultAPIKeyHeader,
		onUnauthorized: defaultUnauthorized,
		onForbidden:    defaultForbidden,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if cfg.extractor == nil {
		cfg.extractor = DefaultCredentialExtractor(cfg.apiKeyHeader)
	}

	return cfg
}

func defaultUnauthorized(w http.ResponseWriter, _ *http.Request, _ error) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func defaultForbidden(w http.ResponseWriter, _ *http.Request, _ error) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func extractBearerToken(authorization string) string {
	parts := strings.Fields(strings.TrimSpace(authorization))
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
