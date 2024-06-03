package auth

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/MatthewAraujo/notify/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func GenerateJWT() (string, error) {
	godotenv.Load()

	app_id := os.Getenv("APP_ID")
	payload := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Second).Unix(),
		"iss": app_id,
	}

	// Read RSA private key from file
	privateKeyBytes, err := utils.ReadFile("key.pem")
	if err != nil {
		return "", err
	}

	// Decode PEM block
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing RSA private key")
	}

	// Parse RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Create JWT
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, payload).SignedString(privateKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
