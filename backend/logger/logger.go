package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	logFile     *os.File
	mu          sync.Mutex
)

// Init initializes the file logger
func Init(logPath string) error {
	mu.Lock()
	defer mu.Unlock()

	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create loggers
	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Close closes the log file
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
	// Also print to stdout
	log.Printf("INFO: "+format, v...)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
	// Also print to stderr
	log.Printf("ERROR: "+format, v...)
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	}
}
