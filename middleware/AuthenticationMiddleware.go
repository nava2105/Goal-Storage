package middleware

import (
	"context"
	"net/http"
)

// Middleware to extract the authorization header and attach it to the context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		// Check if it starts with "Bearer "
		if len(token) < 7 || token[:7] != "Bearer " {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Authorization header is valid, add it to the context
		ctx := context.WithValue(r.Context(), "Authorization", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
