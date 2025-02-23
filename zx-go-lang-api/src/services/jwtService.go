package services

import (
	"fmt"
	"log"
	"time"
	"github.com/GouthamGuna/zx/src/utils"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte(utils.GetJWTEncryptionKey())

func CreateToken(username string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "zx-app",                         // Issuer
		"aud": GetRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	log.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func GetRole(username string) string {
	if username == "admin@zx" {
		return "super admin"
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
