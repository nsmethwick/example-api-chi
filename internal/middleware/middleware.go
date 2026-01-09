package middleware

import (
	"net/http"

	h "github.com/fastly/api-boilerplate/internal/http"
)

// Decorator is a middleware function.
type Decorator func(http.Handler) http.Handler

// Pipeline is a type for wrapping an http.Handler and adding optional logging,
// metrics and tracing.
type Pipeline struct {
	middleware []Decorator
}

// New returns a middleware Pipeline.
func New(middleware ...Decorator) *Pipeline {
	p := &Pipeline{}
	p.middleware = append(p.middleware, middleware...)
	return p
}

// Decorate wraps next in the defined middleware and executes the pipeline.
func (p *Pipeline) Decorate(next http.Handler) http.Handler {
	for i := len(p.middleware) - 1; i >= 0; i-- {
		next = p.middleware[i](next)
	}
	return next
}

// AddRoutePattern adds the route pattern to the request context.
// NOTE: Other middleware, such as tracing and metrics, will lookup the route pattern.
func AddRoutePattern(mux *http.ServeMux) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := h.ContextWithRoutePattern(mux, r)
			next.ServeHTTP(w, r.WithContext(ctx)) //nolint:contextcheck
		})
	}
}
