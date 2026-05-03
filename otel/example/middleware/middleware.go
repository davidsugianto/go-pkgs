package middleware

import (
	"net/http"

	"github.com/davidsugianto/go-pkgs/logger"
	"github.com/davidsugianto/go-pkgs/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware creates HTTP middleware for tracing requests.
func TracingMiddleware(provider *otel.Provider) func(http.Handler) http.Handler {
	tracer := provider.Tracer("http-server")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.user_agent", r.UserAgent()),
				),
			)
			defer span.End()

			// Inject logger with trace context
			log := logger.WithContext(ctx)
			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("Incoming request")

			// Call next handler with context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
