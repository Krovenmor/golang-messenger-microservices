package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func HTTPMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http.server")
}
