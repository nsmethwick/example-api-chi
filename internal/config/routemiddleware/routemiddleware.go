package routemiddlewware

// ==================
// Custom Middleware Example
// ==================

// func CustomerMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		slog.Info("==== middleware here ====")

// 		// proceed with request
// 		next.ServeHTTP(w, r)
// 	})
// }
