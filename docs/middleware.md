# Middleware

The router we are using in this project includes support for custom middleware for HTTP Requests/Responses.  Middleware can be injected for entire route groups or for specific routes.

## Middleware Code Example

```go
func CustomerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("==== middleware here ====")

		// proceed with request
		next.ServeHTTP(w, r)
	})
}
```