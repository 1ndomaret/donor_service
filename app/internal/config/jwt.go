package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func JWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		env, err := godotenv.Read("secret.env")
		if err != nil {
			return nil, errors.New("JWT_SECRET is not configured")
		}

		secret = env["JWT_SECRET"]
	}

	return []byte(secret), nil
}
