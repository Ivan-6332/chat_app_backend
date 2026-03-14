package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth0Domain               string
	Auth0Audience             string
	Auth0ClientID             string
	Auth0ClientSecret         string
	Auth0ManagementAudience   string
	Auth0SyncAllowedClientIDs []string
	Port                      string
	GinMode                   string
	MongoDBURI                string
	MongoDBName               string
	Auth0SyncPageSize         int
	Auth0SyncEnabled          bool
	Auth0SyncIntervalMinutes  int
	Auth0SyncMaxPages         int
	Auth0SyncServiceTokenOnly bool
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Auth0Domain:               getEnv("AUTH0_DOMAIN", ""),
		Auth0Audience:             getEnv("AUTH0_AUDIENCE", ""),
		Auth0ClientID:             getEnv("AUTH0_CLIENT_ID", ""),
		Auth0ClientSecret:         getEnv("AUTH0_CLIENT_SECRET", ""),
		Auth0ManagementAudience:   getEnv("AUTH0_MANAGEMENT_AUDIENCE", "https://"+getEnv("AUTH0_DOMAIN", "")+"/api/v2/"),
		Auth0SyncAllowedClientIDs: getEnvList("AUTH0_SYNC_ALLOWED_CLIENT_IDS"),
		Port:                      getEnv("PORT", "8080"),
		GinMode:                   getEnv("GIN_MODE", "debug"),
		MongoDBURI:                getEnv("MONGODB_URI", ""),
		MongoDBName:               getEnv("MONGODB_DATABASE", "chatapp_db"),
		Auth0SyncPageSize:         getEnvInt("AUTH0_SYNC_PAGE_SIZE", 100),
		Auth0SyncEnabled:          getEnvBool("AUTH0_SYNC_ENABLED", false),
		Auth0SyncIntervalMinutes:  getEnvInt("AUTH0_SYNC_INTERVAL_MINUTES", 30),
		Auth0SyncMaxPages:         getEnvInt("AUTH0_SYNC_MAX_PAGES", 3),
		Auth0SyncServiceTokenOnly: getEnvBool("AUTH0_SYNC_SERVICE_TOKEN_ONLY", true),
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

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var parsed int
		if _, err := fmt.Sscanf(value, "%d", &parsed); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := strings.TrimSpace(strings.ToLower(os.Getenv(key))); value != "" {
		return value == "1" || value == "true" || value == "yes" || value == "on"
	}
	return defaultValue
}

func getEnvList(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}

	return items
}
