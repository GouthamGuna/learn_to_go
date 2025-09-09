package api

import (
	"net/http"
)

type HeaderConfig struct {
	ServerName    string
	Author        string
	PoweredBy     string
}

// DefaultHeaders returns a default header configuration
func DefaultHeaders() *HeaderConfig {
	return &HeaderConfig{
		ServerName:    "FileSync Server",
		Author:        "Gowtham Sankar Gunasekaran",
		PoweredBy:     "Go Language",
	}
}

type HeaderMiddleware struct {
	config *HeaderConfig
}

func NewHeaderMiddleware(cfg *HeaderConfig) *HeaderMiddleware {
	return &HeaderMiddleware{config: cfg}
}

// SetHeaders sets custom headers for all responses
func (hm *HeaderMiddleware) SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Set custom headers
		w.Header().Set("X-Server", hm.config.ServerName)
		w.Header().Set("X-Author", hm.config.Author)
		w.Header().Set("X-Powered-By", hm.config.PoweredBy)
		
		// Remove server identification
		w.Header().Del("Server")
		
		next.ServeHTTP(w, r)
	})
}

// SetCORSHeaders for cross-origin requests
func (hm *HeaderMiddleware) SetCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}