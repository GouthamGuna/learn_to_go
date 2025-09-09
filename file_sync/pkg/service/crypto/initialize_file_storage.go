package pkg

import (
	"fmt"
	"os"
)

func GetEncryptionKey() []byte {
	// Get key from environment variable or use a default
	keyString := os.Getenv("ENCRYPTION_KEY")
	if keyString == "" {
		keyString = "my-secret-encryption-key-32-chars-long!" // Change this in production!
		fmt.Println("WARNING: Using default encryption key. Set ENCRYPTION_KEY environment variable for production.")
	}

	return GenerateKeyFromPassword(keyString)
}
