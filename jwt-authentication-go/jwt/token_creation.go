package jwt

import (
	//"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	//"net/http"
	//"strconv"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var secretKey []byte

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("env loaded...")
	secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func CreateToken(username string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "todo-app",                       // Issuer
		"aud": GetRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	//log.Fatal("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func GetRole(username string) string {
	if username == "senior" {
		return "senior"
	}
	return "employee"
}

// Function to verify JWT tokens
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}
