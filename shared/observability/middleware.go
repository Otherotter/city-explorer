package observability

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestLogger is a Chi-compatible middleware that logs
// every HTTP request as a structured JSON line via slog.
//
// Replaces Chi's built-in middleware.Logger which uses
// the standard log package and emits plain text.
//
// Usage in main():
//
//	r.Use(observability.RequestLogger)
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"bytes", ww.BytesWritten(),
			"remote_addr", r.RemoteAddr,
		)
	})
}
