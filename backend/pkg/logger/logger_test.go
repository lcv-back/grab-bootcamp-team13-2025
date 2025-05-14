package logger

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create logs directory if not exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}

	// Test development config
	devLogger, err := NewLogger(DevelopmentConfig())
	if err != nil {
		t.Fatalf("Failed to create development logger: %v", err)
	}

	// Test logging at different levels
	devLogger.Debug("This is a debug message", "test", "debug")
	devLogger.Info("This is an info message", "test", "info")
	devLogger.Warn("This is a warning message", "test", "warn")
	devLogger.Error("This is an error message", "test", "error")

	// Test production config
	prodLogger, err := NewLogger(ProductionConfig())
	if err != nil {
		t.Fatalf("Failed to create production logger: %v", err)
	}

	// Test logging at different levels
	prodLogger.Debug("This is a debug message", "test", "debug")
	prodLogger.Info("This is an info message", "test", "info")
	prodLogger.Warn("This is a warning message", "test", "warn")
	prodLogger.Error("This is an error message", "test", "error")
}

func TestLoggerWithCustomConfig(t *testing.T) {
	// Create custom config
	config := Config{
		Level:      "debug",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		Pretty:     true,
		Outputs:    []string{"console", "file"},
		LogFile:    "logs/test.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
		Formatter: &CustomFormatter{
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			AppName:    "TestApp",
		},
	}

	// Create logger with custom config
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger with custom config: %v", err)
	}

	// Test logging
	logger.Info("This is a test message", "test", "custom")
}

func TestLoggerWithInvalidConfig(t *testing.T) {
	// Test with invalid config
	invalidConfig := Config{
		Level:   "invalid_level",
		Outputs: []string{},
	}

	// Should return error
	_, err := NewLogger(invalidConfig)
	if err == nil {
		t.Error("Expected error with invalid config, got nil")
	}
}

func TestLoggerWithFileOutput(t *testing.T) {
	// Create test log file
	testLogFile := "logs/test_file.log"
	
	// Clean up after test
	defer os.Remove(testLogFile)

	// Create config with file output
	config := Config{
		Level:      "info",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		Pretty:     false,
		Outputs:    []string{"file"},
		LogFile:    testLogFile,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}

	// Create logger
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Log some messages
	logger.Info("Test file output", "test", "file")
	logger.Error("Test error in file", "test", "file_error")

	// Check if file exists
	if _, err := os.Stat(testLogFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", testLogFile)
	}
} 