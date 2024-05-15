package config

import (
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
	TursoUrl      string
	TursoToken    string
	TursoUrlToken string
	SMTP          SMTP
	Port          string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	env := os.Getenv("PROD")

	if env == "prod" {
		return Config{
			TursoUrl:      getEnv("TURSO_DATABASE_URL", "http://localhost:8080"),
			TursoToken:    getEnv("TURSO_AUTH_TOKEN", "mytoken"),
			TursoUrlToken: getEnv("TURSO_DATABASE_URL_TOKEN", "http://localhost:8080"),
			Port:          getEnv("PORT", "8080"),
			SMTP: SMTP{
				Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
				Port:     getEnv("SMTP_PORT", "587"),
				Author:   getEnv("SMTP_AUTHOR", "my_user@gmail.com"),
				Password: getEnv("SMTP_PASSWORD", "mypassword"),
			},
		}
	}

	return Config{
		TursoUrl:   getEnv("TURSO_DATABASE_URL", "http://localhost:8080"),
		TursoToken: getEnv("TURSO_AUTH_TOKEN", "mytoken"),
		Port:       getEnv("PORT", "8080"),
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
