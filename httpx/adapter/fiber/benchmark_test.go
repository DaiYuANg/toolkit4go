package fiber

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

type benchmarkOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func benchmarkAdapterWithRoute(b *testing.B) *Adapter {
	b.Helper()

	a := New(nil)
	huma.Register(a.HumaAPI(), huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/ping",
	}, func(ctx context.Context, input *struct{}) (*benchmarkOutput, error) {
		out := &benchmarkOutput{}
		out.Body.Message = "pong"
		return out, nil
	})
	return a
}

func BenchmarkAdapterTestRequest(b *testing.B) {
	a := benchmarkAdapterWithRoute(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		resp, err := a.Router().Test(req, -1)
		if err != nil {
			b.Fatalf("fiber test request failed: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("unexpected status code: %d", resp.StatusCode)
		}
	}
}
