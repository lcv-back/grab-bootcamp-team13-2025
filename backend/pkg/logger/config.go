package logger

// DefaultConfig returns the default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		Pretty:     true,
		Outputs:    []string{"console"},
		Formatter:  &DefaultFormatter{TimeFormat: "2006-01-02T15:04:05Z07:00"},
	}
}

// ProductionConfig returns the production logger configuration
func ProductionConfig() Config {
	return Config{
		Level:      "info",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		Pretty:     false,
		Outputs:    []string{"file"},
		LogFile:    "/var/log/app.log",
		MaxSize:    100,    // 100MB
		MaxBackups: 3,      // Keep 3 backup files
		MaxAge:     30,     // Keep logs for 30 days
		Compress:   true,   // Compress rotated logs
		Formatter:  &JSONFormatter{TimeFormat: "2006-01-02T15:04:05Z07:00"},
	}
}

// DevelopmentConfig returns the development logger configuration
func DevelopmentConfig() Config {
	return Config{
		Level:      "debug",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		Pretty:     true,
		Outputs:    []string{"console", "file"},
		LogFile:    "logs/app.log",
		MaxSize:    10,     // 10MB
		MaxBackups: 3,      // Keep 3 backup files
		MaxAge:     7,      // Keep logs for 7 days
		Compress:   false,  // Don't compress in development
		Formatter:  &CustomFormatter{
			TimeFormat: "2006-01-02T15:04:05Z07:00",
			AppName:    "DevApp",
		},
	}
} 