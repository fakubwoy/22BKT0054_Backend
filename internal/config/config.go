package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort      string
	ServerBaseURL   string
	JWTSecret       string
	JWTExpiration   time.Duration
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	RedisDB         int
	S3Enabled       bool
	S3Bucket        string
	S3Region        string
	LocalStorageDir string
}

func LoadConfig() *Config {
	jwtExp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		log.Fatalf("Failed to parse JWT expiration: %v", err)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatalf("Failed to parse Redis DB: %v", err)
	}

	s3Enabled, err := strconv.ParseBool(getEnv("S3_ENABLED", "false"))
	if err != nil {
		log.Fatalf("Failed to parse S3 enabled flag: %v", err)
	}

	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		ServerBaseURL:   getEnv("SERVER_BASE_URL", "http://localhost:8080"),
		JWTSecret:       getEnv("JWT_SECRET", "very-secret-key"),
		JWTExpiration:   jwtExp,
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "fileshare_user"),
		DBPassword:      getEnv("DB_PASSWORD", "securepassword"),
		DBName:          getEnv("DB_NAME", "fileshare"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         redisDB,
		S3Enabled:       s3Enabled,
		S3Bucket:        getEnv("S3_BUCKET", ""),
		S3Region:        getEnv("S3_REGION", ""),
		LocalStorageDir: getEnv("LOCAL_STORAGE_DIR", "./uploads"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
