package services

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthenticateMiddleware(c *gin.Context) {
	// Retrieve the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return
	}

	// Extract the token from the Authorization header
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		log.Println("Invalid Authorization header format")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	tokenString := parts[1]

	// Verify the token
	token, err := VerifyToken(tokenString)
	if err != nil {
		log.Printf("Token verification failed: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed"})
		c.Abort()
		return
	}

	// Print information about the verified token
	log.Printf("Token verified successfully. Claims: %+v\n", token)

	// Continue with the next middleware or route handler
	c.Next()
}
