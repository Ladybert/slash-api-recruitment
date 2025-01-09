package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Address string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	config := &Config{}

	config.Server.Address = os.Getenv("SERVER_ADDRESS")
	if config.Server.Address == "" {
		config.Server.Address = ":8080" // Default fallback
	}

	config.DB.Host = os.Getenv("DB_HOST")
	config.DB.Port = os.Getenv("DB_PORT")
	config.DB.User = os.Getenv("DB_USER")
	config.DB.Password = os.Getenv("DB_PASSWORD")
	config.DB.Name = os.Getenv("DB_NAME")

	return config
}
