package logger

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// ParseLevel converts a string to zerolog.Level
func ParseLevel(levelStr string) zerolog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// Logger interface defines methods
type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
}

// Struct to wrap zerolog
type zerologLogger struct {
	logger zerolog.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

type Config struct {
	AppName         string
	Environment     string
	Level           zerolog.Level
	FilePath        string
	LogstashAddress string
}

// Logstash writer
type LogstashWriter struct {
	conn net.Conn
}

func NewLogstashWriter(address string) (*LogstashWriter, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to logstash: %w", err)
	}
	return &LogstashWriter{conn: conn}, nil
}

func (w *LogstashWriter) Write(p []byte) (int, error) {
	return w.conn.Write(p)
}

func (w *LogstashWriter) Close() error {
	return w.conn.Close()
}

// NewLogger creates a logger with console, file, and logstash support
func NewLogger(config Config) (Logger, error) {
	var writers []io.Writer

	// Console writer
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	writers = append(writers, consoleWriter)

	// File writer
	if config.FilePath != "" {
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	// Logstash writer
	if config.LogstashAddress != "" {
		logstashWriter, err := NewLogstashWriter(config.LogstashAddress)
		if err != nil {
			return nil, err
		}
		writers = append(writers, logstashWriter)
	}

	// Create multi-writer (as io.Writer)
	multiWriter := zerolog.MultiLevelWriter(writers...)

	// Set global log level
	zerolog.SetGlobalLevel(config.Level)

	// Create logger
	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Str("app", config.AppName).
		Str("environment", config.Environment).
		Logger()

	return &zerologLogger{logger: logger}, nil
}

// Implement logging methods
func (l *zerologLogger) Info(msg string, fields ...Field) {
	event := l.logger.Info()
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}

func (l *zerologLogger) Error(msg string, fields ...Field) {
	event := l.logger.Error()
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}

func (l *zerologLogger) Debug(msg string, fields ...Field) {
	event := l.logger.Debug()
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}
