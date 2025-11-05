package logger

import (
	"log"
	"os"
)

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

// SimpleLogger is a simple implementation of Logger
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

// NewSimpleLogger creates a new simple logger
func NewSimpleLogger() Logger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	l.debugLogger.Printf(msg, args...)
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}
