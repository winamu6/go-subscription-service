package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/winnamu6/go-subscription-service/internal/logger"
)

type Config struct {
	AppPort string
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
}

func Load() *Config {
	log := logger.Get()

	if err := godotenv.Load(".env"); err != nil {
		log.Warn(".env file not found, using default values")
	} else {
		log.Info(".env file loaded successfully")
	}

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DBHost:  getEnv("DB_HOST", "localhost"),
		DBPort:  getEnv("DB_PORT", "5432"),
		DBUser:  getEnv("DB_USER", "postgres"),
		DBPass:  getEnv("DB_PASS", ""),
		DBName:  getEnv("DB_NAME", "subscription"),
	}

	log.Infof("Config loaded: %+v", cfg)
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
