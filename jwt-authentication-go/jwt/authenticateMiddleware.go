package jwt

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Function to verify JWT tokens
func AuthenticateMiddleware(c *gin.Context) {
	// Retrieve the token from the cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		log.Println("Token missing in cookie")
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	// Verify the token
	token, err := VerifyToken(tokenString)
	if err != nil {
		log.Printf("Token verification failed: %v\\n", err)
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	// Print information about the verified token
	log.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)

	// Continue with the next middleware or route handler
	c.Next()
}
