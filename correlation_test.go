package correlation_id_traefik_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	correlation "github.com/saman-jafari/correlation-id-traefik"
)

func TestGeneratesCorrelationID(t *testing.T) {
	req := serve(t, correlation.CreateConfig(), nil)

	if got := req.Header.Get("X-Correlation-Id"); got == "" {
		t.Error("expected a generated correlation id on the default header, got empty")
	}
}

func TestUsesConfiguredHeaderName(t *testing.T) {
	cfg := correlation.CreateConfig()
	cfg.HeaderName = "correlation-id"

	req := serve(t, cfg, nil)

	if got := req.Header.Get("correlation-id"); got == "" {
		t.Error("expected a generated correlation id on the configured header, got empty")
	}
}

func TestPreservesIncomingCorrelationID(t *testing.T) {
	const existing = "already-set-123"

	req := serve(t, correlation.CreateConfig(), map[string]string{"X-Correlation-Id": existing})

	if got := req.Header.Get("X-Correlation-Id"); got != existing {
		t.Errorf("expected incoming id %q to be preserved, got %q", existing, got)
	}
}

// serve runs the middleware against a request carrying the given headers and
// returns that request after the handler mutated it.
func serve(t *testing.T, cfg *correlation.Config, headers map[string]string) *http.Request {
	t.Helper()

	ctx := context.Background()
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})

	handler, err := correlation.New(ctx, next, cfg, "plugin")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	handler.ServeHTTP(httptest.NewRecorder(), req)

	return req
}
