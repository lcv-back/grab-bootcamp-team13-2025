package main

import (
	"grab-bootcamp-be-team13-2025/internal/app"
	applogger "grab-bootcamp-be-team13-2025/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Khởi tạo logger custom
	logConfig := applogger.Config{
		AppName:     "GrabBootcampApp",
		Environment: "development",
		Level:       applogger.ParseLevel("debug"),
		FilePath:    "logs/app.log",
	}
	customLogger, err := applogger.NewLogger(logConfig)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	customLogger.Info("Test log hoạt động", applogger.Field{Key: "test", Value: true})

	// initialize the application
	appInstance, err := app.NewApp("config/config.yaml", customLogger)
	if err != nil {
		customLogger.Error("Failed to initialize app", applogger.Field{Key: "error", Value: err})
		return
	}

	// start the server
	if err := appInstance.Run(); err != nil {
		customLogger.Error("App exited with error", applogger.Field{Key: "error", Value: err})
	}
}
