package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// UserIDContextKey is the key used to store the user ID in request context.
const UserIDContextKey contextKey = "userID"

// AuthMiddleware intercepts and validates JWT tokens.
type AuthMiddleware struct {
	jwtSecret []byte
}

// NewAuthMiddleware creates a new AuthMiddleware instance.
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: []byte(jwtSecret)}
}

// Handler wraps an http.Handler with JWT authorization checks.
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format. Must be 'Bearer <token>'")
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired access token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Retrieve user_id from token claims (parsed as float64 in JWT)
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "User ID not found in token")
			return
		}

		userID := int64(userIDFloat)

		// Inject User ID into context
		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext gets the user ID injected by the middleware.
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int64)
	return userID, ok
}
