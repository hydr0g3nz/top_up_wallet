package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents the application logger
type Logger struct {
	zap *zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger(isProduction bool) (*Logger, error) {
	var config zap.Config

	if isProduction {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create log directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("can't create log directory: %w", err)
	}

	// Set up file rotation
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("can't open log file: %w", err)
	}

	// Create core for writing to both file and stdout
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(file),
			config.Level,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(config.EncoderConfig),
			zapcore.AddSync(os.Stdout),
			config.Level,
		),
	)

	// Create logger
	zapLogger := zap.New(core)

	return &Logger{zapLogger}, nil
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.zap.Sugar().Infow(msg, fields...)

}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.zap.Sugar().Errorw(msg, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.zap.Sugar().Debugw(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.zap.Sugar().Warnw(msg, fields...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.zap.Sugar().Fatalw(msg, fields...)
}

// With creates a child logger with additional fields
func (l *Logger) With(fields ...interface{}) *Logger {
	f := make([]zapcore.Field, len(fields))
	for i, field := range fields {
		f[i] = zap.Any("field", field)
	}
	return &Logger{l.zap.With(f...)}
}

// Close flushes any buffered log entries
func (l *Logger) Close() error {
	return l.zap.Sync()
}
