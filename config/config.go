package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	PostgresDSN string
	MongoURI   string
	MongoDB    string
	JWTSecret  string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // baca .env

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		PostgresDSN: getEnv("POSTGRES_DSN", ""),
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:    getEnv("MONGO_DB", "prestasi_db"),
		JWTSecret:  getEnv("JWT_SECRET", "changeme"),
	}

	if cfg.PostgresDSN == "" {
		log.Println("[WARN] POSTGRES_DSN is empty, Postgres may not connect")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
