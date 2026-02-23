package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	MongoDB   MongoDBConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type RateLimitConfig struct {
	Limit  int
	Window time.Duration
}

type CORSConfig struct {
	Origins []string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	rateLimitWindow, _ := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1m"))
	mongoTimeout, _ := time.ParseDuration(getEnv("MONGO_TIMEOUT", "10s"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Env:  getEnv("ENV", "development"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "vinodhini_api"),
			Timeout:  mongoTimeout,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
			Expiry: jwtExpiry,
		},
		RateLimit: RateLimitConfig{
			Limit:  100,
			Window: rateLimitWindow,
		},
		CORS: CORSConfig{
			Origins: []string{getEnv("CORS_ORIGINS", "http://localhost:3000")},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
