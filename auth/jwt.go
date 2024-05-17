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
	"github.com/MatthewAraujo/notify/types"
	"github.com/MatthewAraujo/notify/utils"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func GenerateJWT() (string, error) {
	token, err := getJwt()
	tokenIsExpired := false
	if err != nil {
		if err.Error() == "token not found" {
			log.Printf("Token not found, generating new token")
		} else {
			isExpired, err := IsTokenExpired(token.Token)

			if err != nil {
				return "", err
			}

			if !isExpired {
				fmt.Println("Token is not expired")
				return token.Token, nil
			}
			tokenIsExpired = true

			log.Printf("Token is expired, generating new token")
		}
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

	if tokenIsExpired {
		err = UpdateJwtToken(tokenString)
		if err != nil {
			return "", err
		}
	} else {
		err = InsertJwtToken(tokenString)
		if err != nil {
			return "", err
		}
	}

	return tokenString, nil
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

func UpdateJwtToken(token string) error {
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE JwtToken SET token = ?", token)
	if err != nil {
		return err
	}

	return nil
}

func getJwt() (*types.JwtToken, error) {
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * from JwtToken")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	token := new(types.JwtToken)

	for rows.Next() {
		token, err = scanRowIntoJwtToken(rows)
		if err != nil {
			return nil, err
		}
	}

	if token.Token == "" {
		return nil, fmt.Errorf("token not found")
	}

	return token, nil
}

func IsTokenExpired(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	return time.Now().After(expirationTime), nil
}

func scanRowIntoJwtToken(rows *sql.Rows) (*types.JwtToken, error) {
	var token types.JwtToken
	if err := rows.Scan(&token.Token); err != nil {
		return nil, err
	}
	return &token, nil
}
