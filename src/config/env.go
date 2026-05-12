package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
	JWTSecret   string
	RedisHost   string
	RedisPort   string
	RedisPass   string
	OpenAIKey   string
	GoogleID    string
	GoogleSecret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPass:       getEnv("DB_PASS", ""),
		DBName:       getEnv("DB_NAME", "chat_db"),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		RedisHost:    getEnv("REDIS_HOST", "localhost"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		RedisPass:    getEnv("REDIS_PASS", ""),
		OpenAIKey:    getEnv("OPENAI_API_KEY", ""),
		GoogleID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
