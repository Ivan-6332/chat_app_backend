package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth0Domain   string
	Auth0Audience string
	Port          string
	GinMode       string
	MongoDBURI    string
	MongoDBName   string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Auth0Domain:   getEnv("AUTH0_DOMAIN", ""),
		Auth0Audience: getEnv("AUTH0_AUDIENCE", ""),
		Port:          getEnv("PORT", "8080"),
		GinMode:       getEnv("GIN_MODE", "debug"),
		MongoDBURI:    getEnv("MONGODB_URI", ""),
		MongoDBName:   getEnv("MONGODB_DATABASE", "chatapp_db"),
	}

	// Validate required config
	if AppConfig.Auth0Domain == "" || AppConfig.Auth0Audience == "" {
		log.Fatal("AUTH0_DOMAIN and AUTH0_AUDIENCE must be set")
	}

	if AppConfig.MongoDBURI == "" {
		log.Fatal("MONGODB_URI must be set")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
