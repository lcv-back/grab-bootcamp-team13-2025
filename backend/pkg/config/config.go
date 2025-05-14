package config

import (
	"os"
)

type Config struct {
	ResetPasswordURL string
}

func NewConfig() *Config {
	return &Config{
		ResetPasswordURL: getEnv("RESET_PASSWORD_URL", "http://isymptom.vercel.app/reset-password"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 