package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type SMTP struct {
	Host     string
	Port     string
	Author   string
	Password string
}

type Config struct {
	Port     string
	TursoURl string
	SMTP     SMTP
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	tursoUrl := getEnv("TURSO_DATABASE_URL", "localhost:3306")
	tursoToken := getEnv("TURSO_AUTH_TOKEN", "mytoken")

	return Config{
		Port:     getEnv("PORT", "8080"),
		TursoURl: fmt.Sprintf("%s?authToken=%s", tursoUrl, tursoToken),
		SMTP: SMTP{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			Author:   getEnv("SMTP_AUTHOR", "myemail@gmail.com"),
			Password: getEnv("SMTP_PASSWORD", "mypassword"),
		},
	}

}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
