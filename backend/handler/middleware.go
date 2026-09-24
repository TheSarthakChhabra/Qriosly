package handler

import (
	"context"
	"net/http"
	"quiz-backend/apperror"
	"quiz-backend/service"
	"strings"
)

type contextKey string

const (
	roleContextKey   contextKey = "role"
	userIDContextKey contextKey = "userID"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, apperror.Unauthorized("MISSING_AUTH_HEADER", "missing authorization header"))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, apperror.Unauthorized("INVALID_AUTH_HEADER", "invalid authorization header format"))
				return
			}

			userID, role, err := service.ParseToken(parts[1], jwtSecret)
			if err != nil {
				writeError(w, apperror.Unauthorized("INVALIS_TOKEN", "invalidor expired token"))
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			ctx = context.WithValue(ctx, roleContextKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleContextKey).(string)
	return role, ok
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := RoleFromContext(r.Context())
			if !ok {
				writeError(w, apperror.Unauthorized("UNAUTHENTICATED", "authentication required"))
				return
			}
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeError(w, apperror.Forbidden("INSUFFICIENT_PERMISSIONS", "you donot have permission to perform this action"))
		})

	}
}
