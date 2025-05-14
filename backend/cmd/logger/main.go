package main

import (
	"time"

	"grab-bootcamp-be-team13-2025/pkg/logger"
)

func main() {
	// Create logger with development config
	log, err := logger.NewLogger(logger.DevelopmentConfig())
	if err != nil {
		panic(err)
	}

	// Test different log levels
	log.Debug("This is a debug message", "timestamp", time.Now().Unix())
	log.Info("This is an info message", "timestamp", time.Now().Unix())
	log.Warn("This is a warning message", "timestamp", time.Now().Unix())
	log.Error("This is an error message", "timestamp", time.Now().Unix())

	// Test with custom fields
	log.Info("User logged in",
		"user_id", "123",
		"ip", "192.168.1.1",
		"timestamp", time.Now().Unix(),
	)

	// Test with error
	log.Error("Failed to connect to database",
		"error", "connection timeout",
		"retry_count", 3,
		"timestamp", time.Now().Unix(),
	)

	// Test with custom formatter
	customConfig := logger.Config{
		Level:      "debug",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		Pretty:     true,
		Outputs:    []string{"console", "file"},
		LogFile:    "logs/custom.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
		Formatter: &logger.CustomFormatter{
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			AppName:    "CustomApp",
		},
	}

	customLog, err := logger.NewLogger(customConfig)
	if err != nil {
		panic(err)
	}

	customLog.Info("This is a custom formatted message",
		"custom_field", "custom_value",
		"timestamp", time.Now().Unix(),
	)
} 