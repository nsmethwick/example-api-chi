package http

import (
	"fmt"
	"net/http"

	"github.com/erp-api/internal/config"
)

// WithHeaders adds common HTTP response headers if not already set.
func WithHeaders(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			if headers.Get("Surrogate-Control") == "" {
				headers.Set("Surrogate-Control", fmt.Sprintf("max-age=%d", cfg.CacheDefaultTTL))
			}
			if headers.Get("Cache-Control") == "" {
				headers.Set("Cache-Control", "no-store")
			}
			if headers.Get("Content-Type") == "" {
				headers.Set("Content-Type", "application/json")
			}
			next.ServeHTTP(w, r)
		})
	}
}
