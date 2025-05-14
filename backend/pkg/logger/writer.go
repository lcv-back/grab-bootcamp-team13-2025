package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ConsoleWriter implements Writer interface for console output
type ConsoleWriter struct {
	pretty     bool
	timeFormat string
}

func NewConsoleWriter(pretty bool, timeFormat string) *ConsoleWriter {
	return &ConsoleWriter{
		pretty:     pretty,
		timeFormat: timeFormat,
	}
}

func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	if w.pretty {
		writer := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: w.timeFormat,
		}
		return writer.Write(p)
	}
	return os.Stdout.Write(p)
}

// FileWriter implements Writer interface for file output
type FileWriter struct {
	writer *lumberjack.Logger
}

func NewFileWriter(config WriterConfig) (*FileWriter, error) {
	if config.Path == "" {
		return nil, fmt.Errorf("file path is required for file writer")
	}
	return &FileWriter{
		writer: &lumberjack.Logger{
			Filename:   config.Path,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		},
	}, nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// MultiWriter implements Writer interface for multiple outputs
type MultiWriter struct {
	writers []Writer
}

func NewMultiWriter(writers ...Writer) *MultiWriter {
	return &MultiWriter{
		writers: writers,
	}
}

func (w *MultiWriter) Write(p []byte) (n int, err error) {
	for _, writer := range w.writers {
		if _, err := writer.Write(p); err != nil {
			return 0, fmt.Errorf("failed to write to writer: %w", err)
		}
	}
	return len(p), nil
}

// WriterFactory implementation
type DefaultWriterFactory struct{}

func (f *DefaultWriterFactory) CreateWriter(config WriterConfig) (Writer, error) {
	switch config.Type {
	case "console":
		return NewConsoleWriter(config.Pretty, "2006-01-02T15:04:05Z07:00"), nil
	case "file":
		return NewFileWriter(config)
	default:
		return nil, fmt.Errorf("unsupported writer type: %s", config.Type)
	}
} 