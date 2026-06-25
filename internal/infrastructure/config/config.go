package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBSSLMode              string
	JWTSecret              string
	JWTAccessTokenExpiry   time.Duration
	JWTRefreshTokenExpiry  time.Duration
	ServerPort             string
	UploadBasePath         string
	UploadBooksPath        string
	UploadVideosPath       string
	UploadAudioPath        string
	UploadImagesPath       string
	UploadPhotoArchivePath string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}

	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m"))
	if err != nil {
		log.Fatalf("invalid JWT_ACCESS_TOKEN_EXPIRY: %v", err)
	}

	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		log.Fatalf("invalid JWT_REFRESH_TOKEN_EXPIRY: %v", err)
	}

	return &Config{
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBName:                 getEnv("DB_NAME", "hemra_siirow"),
		DBSSLMode:              getEnv("DB_SSLMODE", "disable"),
		JWTSecret:              getEnv("JWT_SECRET", "replace_with_secure_secret"),
		JWTAccessTokenExpiry:   accessExpiry,
		JWTRefreshTokenExpiry:  refreshExpiry,
		ServerPort:             getEnv("SERVER_PORT", "8080"),
		UploadBasePath:         getEnv("UPLOAD_BASE_PATH", "uploads"),
		UploadBooksPath:        getEnv("UPLOAD_BOOKS_PATH", "uploads/books"),
		UploadVideosPath:       getEnv("UPLOAD_VIDEOS_PATH", "uploads/videos"),
		UploadAudioPath:        getEnv("UPLOAD_AUDIO_PATH", "uploads/audio"),
		UploadImagesPath:       getEnv("UPLOAD_IMAGES_PATH", "uploads/images"),
		UploadPhotoArchivePath: getEnv("UPLOAD_PHOTOARCHIVE_PATH", "uploads/photoarchive"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
