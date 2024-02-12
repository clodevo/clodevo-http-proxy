// utils/logging.go
package utils

import (
	"log"
)

type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelDebug
	LogLevelTrace
)

var logger *Logger // Singleton instance of Logger

func init() {
	// Initialize the logger with a default level. This can be overridden.
	logger = NewLogger(LogLevelInfo)
}

type Logger struct {
	level LogLevel
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

func (l *Logger) SetLogLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level >= LogLevelInfo {
		log.Printf("[INFO] "+msg, args...)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level >= LogLevelDebug {
		log.Printf("[DEBUG] "+msg, args...)
	}
}

func (l *Logger) Trace(msg string, args ...interface{}) {
	if l.level >= LogLevelTrace {
		log.Printf("[TRACE] "+msg, args...)
	}
}

// GetLogger returns the singleton instance of the logger.
func GetLogger() *Logger {
	return logger
}
