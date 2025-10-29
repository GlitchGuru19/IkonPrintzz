package middleware

import (
	"context"
	"fileprintapp/internal/usecase"
	"net/http"
	"strings"
)

type contextKey string

const UsernameKey contextKey = "username"

// AuthMiddleware verifies JWT tokens for protected routes
type AuthMiddleware struct {
	authService *usecase.AuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *usecase.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate checks if the request has a valid JWT token
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		username, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add username to context
		ctx := context.WithValue(r.Context(), UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
