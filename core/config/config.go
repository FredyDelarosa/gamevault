package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	// Cargar .env si existe (para desarrollo)
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "fredy"),
		DBPassword: getEnv("DB_PASSWORD", "125645"),
		DBName:     getEnv("DB_NAME", "gamevault"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
