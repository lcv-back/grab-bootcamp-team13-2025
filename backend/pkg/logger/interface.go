package logger

import "github.com/rs/zerolog"

// Logger defines the core logging interface
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// Formatter defines the interface for custom log formatting
type Formatter interface {
	Format(level zerolog.Level, msg string, fields map[string]interface{}) string
}

// Writer defines the interface for log output writers
type Writer interface {
	Write(p []byte) (n int, err error)
}

// WriterFactory defines the interface for creating log writers
type WriterFactory interface {
	CreateWriter(config WriterConfig) (Writer, error)
}

// WriterConfig holds configuration for log writers
type WriterConfig struct {
	Type       string
	Path       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Pretty     bool
} 