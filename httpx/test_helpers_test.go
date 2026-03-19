package httpx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	adapterecho "github.com/DaiYuANg/arcgo/httpx/adapter/echo"
	adapterfiber "github.com/DaiYuANg/arcgo/httpx/adapter/fiber"
	adaptergin "github.com/DaiYuANg/arcgo/httpx/adapter/gin"
	adapterstd "github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/stretchr/testify/require"
)

func serveRequest(t testing.TB, server ServerRuntime, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()

	impl := unwrapServer(server)
	require.NotNil(t, impl)
	require.NotNil(t, impl.adapter)

	switch host := impl.adapter.(type) {
	case *adapterstd.Adapter:
		rec := httptest.NewRecorder()
		host.Router().ServeHTTP(rec, req)
		return rec
	case *adaptergin.Adapter:
		rec := httptest.NewRecorder()
		host.Router().ServeHTTP(rec, req)
		return rec
	case *adapterecho.Adapter:
		rec := httptest.NewRecorder()
		host.Router().ServeHTTP(rec, req)
		return rec
	case *adapterfiber.Adapter:
		resp, err := host.Router().Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		rec := httptest.NewRecorder()
		for key, values := range resp.Header {
			for _, value := range values {
				rec.Header().Add(key, value)
			}
		}
		rec.WriteHeader(resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		_, _ = rec.Body.Write(body)
		return rec
	default:
		t.Fatalf("unsupported adapter type %T", host)
		return nil
	}
}
