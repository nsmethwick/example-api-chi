package http

import (
	"context"
	"net/http"
)

// RouteContextKey is a custom type for the key used in context.Context
type routeContextKey struct{}

// RoutePatternContextKey is the key to store a route pattern in context.Context
var RoutePatternContextKey = routeContextKey{}

// ContextWithRoutePattern sets the request route into the context.Context
func ContextWithRoutePattern(mux *http.ServeMux, r *http.Request) context.Context {
	_, route := mux.Handler(r)
	return context.WithValue(r.Context(), RoutePatternContextKey, route)
}

// RoutePatternFromContext returns the request route from the context.Context
func RoutePatternFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(RoutePatternContextKey)
	if v != nil {
		if s, ok := v.(string); ok {
			return s, ok
		}
	}
	return "", false
}
