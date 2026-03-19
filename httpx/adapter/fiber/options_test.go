//go:build !no_fiber

package fiber

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DaiYuANg/arcgo/httpx/adapter"
	fiberframework "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_UsesProvidedApp(t *testing.T) {
	external := fiberframework.New()
	a := New(external)

	assert.Same(t, external, a.Router())
}

func TestNew_AppliesDocsPaths(t *testing.T) {
	a := New(nil, adapter.HumaOptions{
		DocsPath:    "/reference",
		OpenAPIPath: "/spec",
	})

	docsResp := mustTest(t, a, httptest.NewRequest(http.MethodGet, "/reference", nil))
	assert.Equal(t, http.StatusOK, docsResp.StatusCode)

	oldDocsResp := mustTest(t, a, httptest.NewRequest(http.MethodGet, "/docs", nil))
	assert.Equal(t, http.StatusNotFound, oldDocsResp.StatusCode)

	specResp := mustTest(t, a, httptest.NewRequest(http.MethodGet, "/spec.json", nil))
	assert.Equal(t, http.StatusOK, specResp.StatusCode)
}

func TestNew_DisablesDocsRoutes(t *testing.T) {
	a := New(nil, adapter.HumaOptions{DisableDocsRoutes: true})

	docsResp := mustTest(t, a, httptest.NewRequest(http.MethodGet, "/docs", nil))
	assert.Equal(t, http.StatusNotFound, docsResp.StatusCode)

	specResp := mustTest(t, a, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))
	assert.Equal(t, http.StatusNotFound, specResp.StatusCode)
}

func mustTest(t *testing.T, a *Adapter, req *http.Request) *http.Response {
	t.Helper()

	resp, err := a.Router().Test(req, -1)
	require.NoError(t, err)
	t.Cleanup(func() {
		if resp != nil && resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	})
	return resp
}
