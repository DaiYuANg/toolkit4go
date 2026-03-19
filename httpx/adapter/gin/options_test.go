package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DaiYuANg/arcgo/httpx/adapter"
	ginframework "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNew_UsesProvidedEngine(t *testing.T) {
	ginframework.SetMode(ginframework.TestMode)
	engine := ginframework.New()
	a := New(engine)

	assert.Same(t, engine, a.Router())
}

func TestNew_AppliesDocsPaths(t *testing.T) {
	ginframework.SetMode(ginframework.TestMode)
	a := New(nil, adapter.HumaOptions{
		DocsPath:    "/reference",
		OpenAPIPath: "/spec",
	})

	docsReq := httptest.NewRequest(http.MethodGet, "/reference", nil)
	docsRec := httptest.NewRecorder()
	a.Router().ServeHTTP(docsRec, docsReq)
	assert.Equal(t, http.StatusOK, docsRec.Code)

	oldDocsReq := httptest.NewRequest(http.MethodGet, "/docs", nil)
	oldDocsRec := httptest.NewRecorder()
	a.Router().ServeHTTP(oldDocsRec, oldDocsReq)
	assert.Equal(t, http.StatusNotFound, oldDocsRec.Code)

	specReq := httptest.NewRequest(http.MethodGet, "/spec.json", nil)
	specRec := httptest.NewRecorder()
	a.Router().ServeHTTP(specRec, specReq)
	assert.Equal(t, http.StatusOK, specRec.Code)
}

func TestNew_DisablesDocsRoutes(t *testing.T) {
	ginframework.SetMode(ginframework.TestMode)
	a := New(nil, adapter.HumaOptions{DisableDocsRoutes: true})

	docsReq := httptest.NewRequest(http.MethodGet, "/docs", nil)
	docsRec := httptest.NewRecorder()
	a.Router().ServeHTTP(docsRec, docsReq)
	assert.Equal(t, http.StatusNotFound, docsRec.Code)

	specReq := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	specRec := httptest.NewRecorder()
	a.Router().ServeHTTP(specRec, specReq)
	assert.Equal(t, http.StatusNotFound, specRec.Code)
}
