package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: validate token (JWT, session, etc.)
		// userID, err := ValidateToken(token)
		// if err != nil {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		// Store user ID in context for handlers
		ctx := context.WithValue(r.Context(), UserIDKey, "omotadeogunmode@gmail.com")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
