package middleware

import (
	"context"
	"net/http"
	"strings"

	"lab-4/backend/internal/auth"
	"lab-4/backend/internal/service"
)

type contextKey string

const userIDKey contextKey = "user_id"

func UserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func RequireRoles(users service.UsersService, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := UserID(r.Context())
			if !ok {
				writeUnauthorized(w)
				return
			}
			user, err := users.GetByID(r.Context(), id)
			if err != nil {
				writeUnauthorized(w)
				return
			}
			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":{"code":"forbidden","message":"Insufficient permissions"}}`))
		})
	}
}

// RequirePermission keeps authorization policy out of handlers and prevents
// client-supplied roles from being trusted.
func RequirePermission(users service.UsersService, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := UserID(r.Context())
			if !ok {
				writeUnauthorized(w)
				return
			}
			user, err := users.GetByID(r.Context(), id)
			if err != nil || !auth.HasPermission(user.Role, permission) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error":{"code":"forbidden","message":"Insufficient permissions"}}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Fields(r.Header.Get("Authorization"))
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeUnauthorized(w)
				return
			}
			userID, err := auth.ParseAccessToken(parts[1], secret)
			if err != nil {
				writeUnauthorized(w)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, userID)))
		})
	}
}
func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"Authentication required"}}`))
}
