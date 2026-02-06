package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name          string `mapstructure:"APP_NAME"`
	Port          string `mapstructure:"APP_PORT"`
	Env           string `mapstructure:"APP_ENV"`
	AdminUser     string `mapstructure:"ADMIN_USER"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	Expiry        int
	RefreshExpiry int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	expiryStr := getEnv("JWT_EXPIRY", "24")
	expiry, _ := strconv.Atoi(expiryStr)

	refreshExpiryStr := getEnv("JWT_REFRESH_EXPIRY", "168")
	refreshExpiry, _ := strconv.Atoi(refreshExpiryStr)

	return &Config{
		App: AppConfig{
			Name:          getEnv("APP_NAME", "post-api"),
			Port:          getEnv("APP_PORT", "8080"),
			Env:           getEnv("APP_ENV", "dev"),
			AdminUser:     getEnv("ADMIN_USER", "admin"),
			AdminPassword: getEnv("ADMIN_PASSWORD", "secret"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "post_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "supersecretkey"),
			Expiry:        expiry,
			RefreshExpiry: refreshExpiry,
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
