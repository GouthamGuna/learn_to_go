package api

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"

	logger "github.com/GouthamGuna/file_sync/pkg/logger"
	response "github.com/GouthamGuna/file_sync/pkg/response"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const (
	userRoleKey contextKey = "userRole"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	logger.Info("Handling authentication request")
	// This middleware checks for a valid JWT token in the Authorization header
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteError(w, http.StatusUnauthorized, "Missing Authorization Header", "Unauthorized")
			return
		}

		// Validate the token
		userRole, err := validateJWTToken(authHeader)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, err.Error(), "Invalid token")
			return
		}

		// Set the user role in the context using our custom key type
		ctx := context.WithValue(r.Context(), userRoleKey, userRole)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateToken(username string) (string, error) {
	logger.Info("Initializing token creation")
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                             // Subject (user identifier)
		"iss": "File Sync",                          // Issuer
		"aud": getRole(username),                    // Audience (user role)
		"exp": time.Now().Add(8 * time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                    // Issued at
	})

	var jwt_secret = os.Getenv("JWT_SECRET")
	tokenString, err := claims.SignedString([]byte(jwt_secret))
	if err != nil {
		return "", err
	}

	logger.Infof("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func getRole(username string) string {
	logger.Info("Initializing role retrieval")
	if strings.Contains(strings.ToLower(username), "admin") {
		return "admin"
	}

	return "user"
}

func validateJWTToken(tokenString string) (string, error) {
	logger.Info("Validating JWT token")
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	var jwt_secret = os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	role, ok := claims["aud"].(string)
	if !ok {
		return "", err
	}

	return role, nil
}

// Helper function to retrieve user role from context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	logger.Info("Retrieving user role from context")
	if role, ok := ctx.Value(userRoleKey).(string); ok {
		return role, true
	}
	return "", false
}
