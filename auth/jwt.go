package auth

import (
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/db"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func GenerateJWT() (string, error) {
	token, err := getJwt()
	if err != nil {
		if err.Error() == "token not found" {
			log.Printf("Token not found, generating new token")
		}
	} else {
		log.Printf("Token found in database")
		isExpired, err := IsTokenExpired(token)
		if err != nil {
			return "", err
		}

		if !isExpired {
			return token, nil
		}

		log.Printf("Token expired, generating new token")
	}

	godotenv.Load()

	app_id := os.Getenv("APP_ID")
	payload := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
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

	err = InsertJwtToken(tokenString)
	if err != nil {
		return "", err
	}

	log.Printf("Token generated and saved to database: %s", token)

	return tokenString, nil
}

// IsTokenExpired checks if the provided token is expired
func IsTokenExpired(tokenString string) (bool, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return true, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) > time.Now().Unix() {
				return false, nil // Token is not expired
			}
			return true, nil // Token is expired
		}
	}

	return true, fmt.Errorf("invalid token claims")
}

func InsertJwtToken(token string) error {
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO JwtToken (token) VALUES (?)", token)
	if err != nil {
		return err
	}

	return nil
}

func getJwt() (string, error) {
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		return "", err
	}

	var token string
	err = db.QueryRow("SELECT token FROM JwtToken").Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("token not found")
		}
		return "", err
	}

	if token == "" {
		log.Printf("Token not found")
		return "", fmt.Errorf("token not found")
	}

	log.Printf("Token found in database")

	return token, nil
}
