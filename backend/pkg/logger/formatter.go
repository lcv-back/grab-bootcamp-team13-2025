package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

// DefaultFormatter implements Formatter interface with basic formatting
type DefaultFormatter struct {
	TimeFormat string
}

func (f *DefaultFormatter) Format(level zerolog.Level, msg string, fields map[string]interface{}) string {
	return msg
}

// CustomFormatter implements Formatter interface with custom formatting
type CustomFormatter struct {
	TimeFormat string
	AppName    string
}

func (f *CustomFormatter) Format(level zerolog.Level, msg string, fields map[string]interface{}) string {
	formatted := fmt.Sprintf("[%s] %s: %s", f.AppName, level.String(), msg)
	if len(fields) > 0 {
		formatted += " | "
		for k, v := range fields {
			formatted += fmt.Sprintf("%s=%v ", k, v)
		}
	}
	return formatted
}

// JSONFormatter implements Formatter interface with JSON formatting
type JSONFormatter struct {
	TimeFormat string
}

func (f *JSONFormatter) Format(level zerolog.Level, msg string, fields map[string]interface{}) string {
	// In a real implementation, this would return JSON formatted string
	// For now, we'll just return the message as zerolog handles JSON formatting
	return msg
} 