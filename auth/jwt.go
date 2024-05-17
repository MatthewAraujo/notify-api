package auth

import (
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/db"
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
		}
	} else {

		tokenIsExpired, err = isJwtExpired(token)
		if err != nil {
			return "", err
		}

		if !tokenIsExpired {
			log.Printf("Token is not expired")
			return token, nil
		}

		log.Printf("Token is expired, generating new token")
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

func isJwtExpired(tokenString string) (bool, error) {
	keyPath := "key.pem"
	// Carrega a chave pública do arquivo PEM
	keyData, err := utils.ReadFile(keyPath)
	if err != nil {
		return false, fmt.Errorf("erro ao ler o arquivo de chave: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return false, fmt.Errorf("erro ao analisar a chave pública: %v", err)
	}

	// Análise e verificação do token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifique se o método de assinatura é RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("erro ao analisar o token JWT: %v", err)
	}

	// Verifica se o token é válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verifica a expiração
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			return expirationTime.Before(time.Now()), nil
		} else {
			return false, errors.New("o campo 'exp' não está presente no token")
		}
	} else {
		return false, errors.New("token JWT inválido")
	}
}
