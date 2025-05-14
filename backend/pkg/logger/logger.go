package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds the logger configuration
type Config struct {
	Level      string
	TimeFormat string
	Pretty     bool
	// Log rotation config
	LogFile    string
	MaxSize    int    // megabytes
	MaxBackups int
	MaxAge     int    // days
	Compress   bool
	// Multiple outputs
	Outputs    []string // "console", "file", "json"
	// Custom formatter
	Formatter  Formatter
}

// NewLogger creates a new logger instance
func NewLogger(cfg Config) (Logger, error) {
	// Validate config
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure time format
	zerolog.TimeFieldFormat = cfg.TimeFormat

	// Create writers
	writerFactory := &DefaultWriterFactory{}
	var writers []Writer

	// Add console output if specified
	if contains(cfg.Outputs, "console") {
		consoleWriter, err := writerFactory.CreateWriter(WriterConfig{
			Type:   "console",
			Pretty: cfg.Pretty,
		})
		if err != nil {
			return nil, err
		}
		writers = append(writers, consoleWriter)
	}

	// Add file output if specified
	if contains(cfg.Outputs, "file") && cfg.LogFile != "" {
		fileWriter, err := writerFactory.CreateWriter(WriterConfig{
			Type:       "file",
			Path:       cfg.LogFile,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		if err != nil {
			return nil, err
		}
		writers = append(writers, fileWriter)
	}

	// Create multi-writer
	multiWriter := NewMultiWriter(writers...)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	return &logger{
		formatter: cfg.Formatter,
	}, nil
}

type logger struct {
	formatter Formatter
}

func (l *logger) formatMessage(level zerolog.Level, msg string, fields ...interface{}) string {
	if l.formatter != nil {
		fieldMap := make(map[string]interface{})
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				fieldMap[fields[i].(string)] = fields[i+1]
			}
		}
		return l.formatter.Format(level, msg, fieldMap)
	}
	return msg
}

func (l *logger) Debug(msg string, fields ...interface{}) {
	formattedMsg := l.formatMessage(zerolog.DebugLevel, msg, fields...)
	log.Debug().Fields(fields).Msg(formattedMsg)
}

func (l *logger) Info(msg string, fields ...interface{}) {
	formattedMsg := l.formatMessage(zerolog.InfoLevel, msg, fields...)
	log.Info().Fields(fields).Msg(formattedMsg)
}

func (l *logger) Warn(msg string, fields ...interface{}) {
	formattedMsg := l.formatMessage(zerolog.WarnLevel, msg, fields...)
	log.Warn().Fields(fields).Msg(formattedMsg)
}

func (l *logger) Error(msg string, fields ...interface{}) {
	formattedMsg := l.formatMessage(zerolog.ErrorLevel, msg, fields...)
	log.Error().Fields(fields).Msg(formattedMsg)
}

func (l *logger) Fatal(msg string, fields ...interface{}) {
	formattedMsg := l.formatMessage(zerolog.FatalLevel, msg, fields...)
	log.Fatal().Fields(fields).Msg(formattedMsg)
}

// Helper function to check if a string is in a slice
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// validateConfig validates the logger configuration
func validateConfig(cfg Config) error {
	if cfg.Level == "" {
		return fmt.Errorf("log level is required")
	}
	if len(cfg.Outputs) == 0 {
		return fmt.Errorf("at least one output is required")
	}
	if contains(cfg.Outputs, "file") && cfg.LogFile == "" {
		return fmt.Errorf("log file path is required when file output is enabled")
	}
	return nil
}