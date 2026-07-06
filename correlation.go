// Package correlation_id_traefik is a Traefik middleware plugin that ensures
// every request carries a correlation ID header, generating a UUID v7 when one
// is absent.
package correlation_id_traefik

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const defaultHeaderName = "X-Correlation-Id"

// Config the plugin configuration.
type Config struct {
	HeaderName string `json:"header_name,omitempty" yaml:"headerName,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName: "",
	}
}

// Correlation is the middleware handler.
type Correlation struct {
	next       http.Handler
	name       string
	headerName string
}

// New creates a new Correlation middleware.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	headerName := config.HeaderName
	if headerName == "" {
		headerName = defaultHeaderName
	}

	return &Correlation{
		next:       next,
		name:       name,
		headerName: headerName,
	}, nil
}

// ServeHTTP preserves an incoming correlation ID, or generates a UUID v7
// (Unix-timestamp-ordered, see https://github.com/google/uuid/blob/master/version7.go),
// then forwards the request to the next handler.
func (c *Correlation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	correlationID := request.Header.Get(c.headerName)

	if correlationID == "" {
		if id, err := uuid.NewV7(); err == nil {
			correlationID = id.String()
		}
	}

	if correlationID != "" {
		request.Header.Set(c.headerName, correlationID)
	}

	c.next.ServeHTTP(writer, request)
}
