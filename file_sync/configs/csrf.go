package configs

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	logger "github.com/GouthamGuna/file_sync/pkg/logger"
	response "github.com/GouthamGuna/file_sync/pkg/response"
)

var (
	// tokenLength is the length of the CSRF token in bytes
	tokenLength = 32

	// csrfTokens stores valid tokens with their expiration times
	csrfTokens     = make(map[string]time.Time)
	csrfTokensLock sync.RWMutex

	// tokenExpiration is the duration after which a token expires
	tokenExpiration = 24 * time.Hour
)

// GenerateCSRFToken creates a new CSRF token
func generateCSRFToken() string {
	// Create a random token
	bytes := make([]byte, tokenLength)
	_, err := rand.Read(bytes)
	if err != nil {
		logger.Error("Failed to generate CSRF token:", err)
		return ""
	}

	token := base64.StdEncoding.EncodeToString(bytes)

	// Store the token with expiration time
	csrfTokensLock.Lock()
	csrfTokens[token] = time.Now().Add(tokenExpiration)
	csrfTokensLock.Unlock()

	return token
}

// ValidateCSRFToken checks if a token is valid
func validateCSRFToken(token string) bool {
	if token == "" {
		return false
	}

	csrfTokensLock.RLock()
	expiration, exists := csrfTokens[token]
	csrfTokensLock.RUnlock()

	// Token exists and has not expired
	return exists && time.Now().Before(expiration)
}

// CleanupExpiredTokens removes expired tokens from the map
func cleanupExpiredTokens() {
	csrfTokensLock.Lock()
	defer csrfTokensLock.Unlock()

	now := time.Now()
	for token, expiration := range csrfTokens {
		if now.After(expiration) {
			delete(csrfTokens, token)
		}
	}
}

// CSRF checks for a valid CSRF token in requests that modify data
func CSRF(next http.Handler) http.Handler {
	logger.Info("CSRF middleware initialized")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip token validation for GET, HEAD, OPTIONS requests
		// as these should be safe operations that don't modify data
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// For all other methods (POST, PUT, DELETE, etc.), validate the token
		token := r.Header.Get("X-CSRF-Token")
		if !validateCSRFToken(token) {
			response.WriteError(w, http.StatusForbidden, "Invalid or missing CSRF token", "Forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CSRFTokenHandler provides a token for client-side use
func CSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("CSRF token requested from " + r.RemoteAddr + " for " + r.UserAgent())
	// Clean up expired tokens periodically
	go cleanupExpiredTokens()

	// Generate a new token
	token := generateCSRFToken()

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"csrf_token": "` + token + `"}`))
}
