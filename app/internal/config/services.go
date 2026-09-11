package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServicesConfig struct {
	BaseURL        string
	UserPort       string
	UserURL        string
	GeoapifyAPIKey string
	ServiceToken   string
}

func ServicesSecret() ServicesConfig {
	err := godotenv.Load("secret.env")
	if err != nil {
		log.Println("using system env")
	}

	return ServicesConfig{
		BaseURL:        os.Getenv("BASE_URL"),
		UserPort:       os.Getenv("USER_PORT"),
		UserURL:        os.Getenv("USER_URL"),
		GeoapifyAPIKey: os.Getenv("GEOAPIFY_API_KEY"),
		ServiceToken:   os.Getenv("SERVICE_TOKEN"),
	}
}
