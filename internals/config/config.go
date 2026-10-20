package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	Port        string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("/Users/manasagarwal/Documents/go/todo_with_auth/.env")
	if err != nil {
		log.Println("Failed to load the config", err)
		return nil, err
	}

	var config *Config = &Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	return config, nil
}
