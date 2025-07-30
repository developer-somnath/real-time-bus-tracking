package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	*log.Logger
	file *os.File
}

func Init(serviceName string) *Logger {
	logDir := "../../build/tmp/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	logFilePath := filepath.Join(logDir, serviceName+".log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	logger := log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)
	logger.Printf("[INFO] Logger initialized for %s", serviceName)
	return &Logger{Logger: logger, file: file}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *Logger) Error(msg string, err error, keysAndValues ...interface{}) {
	l.Printf("[ERROR] %s: %v %v", msg, err, keysAndValues)
}

func (l *Logger) Fatal(msg string, err error, keysAndValues ...interface{}) {
	l.Printf("[FATAL] %s: %v %v", msg, err, keysAndValues)
	l.file.Close()
	os.Exit(1)
}

func (l *Logger) Close() {
	l.file.Close()
}
