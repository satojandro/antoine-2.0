// Package utils provides core utilities for Antoine CLI
// This file implements a structured logging system with multiple output formats
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel represents different logging levels
type LogLevel string

const (
	LogLevelTrace LogLevel = "trace"
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
	LogLevelPanic LogLevel = "panic"
)

// LogFormat represents different log output formats
type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

// LogOutput represents different log output destinations
type LogOutput string

const (
	LogOutputStdout LogOutput = "stdout"
	LogOutputStderr LogOutput = "stderr"
	LogOutputFile   LogOutput = "file"
)

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Level       LogLevel  `yaml:"level"`
	Format      LogFormat `yaml:"format"`
	Output      LogOutput `yaml:"output"`
	Development bool      `yaml:"development"`
	Caller      bool      `yaml:"caller"`
	StackTrace  bool      `yaml:"stack_trace"`

	// File logging configuration
	File struct {
		Path       string `yaml:"path"`
		MaxSizeMB  int    `yaml:"max_size_mb"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAgeDays int    `yaml:"max_age_days"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"file"`
}

// Logger wraps logrus with Antoine-specific functionality
type Logger struct {
	*logrus.Logger
	config LoggerConfig
}

// NewLogger creates a new configured logger instance
func NewLogger(config LoggerConfig) (*Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := parseLogLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)

	// Configure output destination
	output, err := getLogOutput(config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure log output: %w", err)
	}
	logger.SetOutput(output)

	// Configure formatter
	formatter := getLogFormatter(config)
	logger.SetFormatter(formatter)

	// Configure caller reporting
	logger.SetReportCaller(config.Caller)

	antoineLogger := &Logger{
		Logger: logger,
		config: config,
	}

	return antoineLogger, nil
}

// parseLogLevel converts string log level to logrus level
func parseLogLevel(level LogLevel) (logrus.Level, error) {
	switch level {
	case LogLevelTrace:
		return logrus.TraceLevel, nil
	case LogLevelDebug:
		return logrus.DebugLevel, nil
	case LogLevelInfo:
		return logrus.InfoLevel, nil
	case LogLevelWarn:
		return logrus.WarnLevel, nil
	case LogLevelError:
		return logrus.ErrorLevel, nil
	case LogLevelFatal:
		return logrus.FatalLevel, nil
	case LogLevelPanic:
		return logrus.PanicLevel, nil
	default:
		return logrus.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// getLogOutput configures the output destination
func getLogOutput(config LoggerConfig) (io.Writer, error) {
	switch config.Output {
	case LogOutputStdout:
		return os.Stdout, nil
	case LogOutputStderr:
		return os.Stderr, nil
	case LogOutputFile:
		return getFileOutput(config), nil
	default:
		return os.Stderr, nil
	}
}

// getFileOutput creates a rotating file output
func getFileOutput(config LoggerConfig) io.Writer {
	// Expand home directory if needed
	logPath := config.File.Path
	if strings.HasPrefix(logPath, "~/") {
		home, _ := os.UserHomeDir()
		logPath = filepath.Join(home, logPath[2:])
	}

	// Ensure directory exists
	dir := filepath.Dir(logPath)
	os.MkdirAll(dir, 0755)

	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    config.File.MaxSizeMB,
		MaxBackups: config.File.MaxBackups,
		MaxAge:     config.File.MaxAgeDays,
		Compress:   config.File.Compress,
	}
}

// getLogFormatter creates the appropriate formatter
func getLogFormatter(config LoggerConfig) logrus.Formatter {
	switch config.Format {
	case LogFormatJSON:
		return &logrus.JSONFormatter{
			TimestampFormat:   time.RFC3339,
			DisableHTMLEscape: true,
			PrettyPrint:       config.Development,
		}
	case LogFormatText:
		fallthrough
	default:
		return &AntoineTextFormatter{
			Development: config.Development,
			ForceColors: true,
		}
	}
}

// AntoineTextFormatter is a custom text formatter for Antoine CLI
type AntoineTextFormatter struct {
	Development bool
	ForceColors bool
}

// Format implements the logrus Formatter interface
func (f *AntoineTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 31 // red
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // cyan
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())

	var caller string
	if entry.HasCaller() {
		funcName := filepath.Base(entry.Caller.Function)
		fileName := filepath.Base(entry.Caller.File)
		caller = fmt.Sprintf(" [%s:%d %s]", fileName, entry.Caller.Line, funcName)
	}

	var fieldsStr string
	if len(entry.Data) > 0 {
		var fields []string
		for k, v := range entry.Data {
			fields = append(fields, fmt.Sprintf("%s=%v", k, v))
		}
		fieldsStr = " " + strings.Join(fields, " ")
	}

	var formatted string
	if f.ForceColors && (f.Development || isTerminal()) {
		formatted = fmt.Sprintf("\x1b[%dm[%s]\x1b[0m [%s]%s %s%s\n",
			levelColor, level, timestamp, caller, entry.Message, fieldsStr)
	} else {
		formatted = fmt.Sprintf("[%s] [%s]%s %s%s\n",
			level, timestamp, caller, entry.Message, fieldsStr)
	}

	return []byte(formatted), nil
}

// isTerminal checks if output is a terminal using environment detection
func isTerminal() bool {
	// Simple and reliable approach: check environment variables

	// If NO_COLOR is set, assume no terminal colors
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// If TERM is set and not "dumb", assume terminal supports colors
	term := os.Getenv("TERM")
	if term != "" && term != "dumb" {
		return true
	}

	// Check for common CI environments (usually no color support)
	ciEnvs := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "CIRCLECI", "JENKINS"}
	for _, env := range ciEnvs {
		if os.Getenv(env) != "" {
			return false
		}
	}

	// Default to true for local development
	return true
}

// ContextLogger provides structured logging with context
type ContextLogger struct {
	logger *Logger
	fields logrus.Fields
}

// WithFields creates a new context logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *ContextLogger {
	logrusFields := make(logrus.Fields)
	for k, v := range fields {
		logrusFields[k] = v
	}

	return &ContextLogger{
		logger: l,
		fields: logrusFields,
	}
}

// WithField creates a new context logger with a single field
func (l *Logger) WithField(key string, value interface{}) *ContextLogger {
	return l.WithFields(map[string]interface{}{key: value})
}

// WithError creates a new context logger with an error field
func (l *Logger) WithError(err error) *ContextLogger {
	return l.WithField("error", err.Error())
}

// WithComponent creates a new context logger for a component
func (l *Logger) WithComponent(component string) *ContextLogger {
	return l.WithField("component", component)
}

// WithOperation creates a new context logger for an operation
func (l *Logger) WithOperation(operation string) *ContextLogger {
	return l.WithField("operation", operation)
}

// Context logger methods
func (cl *ContextLogger) WithField(key string, value interface{}) *ContextLogger {
	newFields := make(logrus.Fields)
	for k, v := range cl.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &ContextLogger{
		logger: cl.logger,
		fields: newFields,
	}
}

func (cl *ContextLogger) WithFields(fields map[string]interface{}) *ContextLogger {
	newFields := make(logrus.Fields)
	for k, v := range cl.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &ContextLogger{
		logger: cl.logger,
		fields: newFields,
	}
}

func (cl *ContextLogger) WithError(err error) *ContextLogger {
	return cl.WithField("error", err.Error())
}

// Logging methods for ContextLogger
func (cl *ContextLogger) Trace(args ...interface{}) {
	cl.logger.WithFields(cl.fields).Trace(args...)
}

func (cl *ContextLogger) Debug(args ...interface{}) {
	cl.logger.WithFields(cl.fields).Debug(args...)
}

func (cl *ContextLogger) Info(args ...interface{}) {
	cl.logger.WithFields(cl.fields).Info(args...)
}

func (cl *ContextLogger) Warn(args ...interface{}) {
	cl.logger.WithFields(cl.fields).Warn(args...)
}

func (cl *ContextLogger) Error(args ...interface{}) {
	cl.logger.WithFields(cl.fields).Error(args...)
}

func (cl *ContextLogger) Fatal(args ...interface{}) {
	cl.logger.WithFields(cl.fields).Fatal(args...)
}

func (cl *ContextLogger) Panic(args ...interface{}) {
	cl.logger.WithFields(cl.fields).Panic(args...)
}

// Formatted logging methods for ContextLogger
func (cl *ContextLogger) Tracef(format string, args ...interface{}) {
	cl.logger.WithFields(cl.fields).Tracef(format, args...)
}

func (cl *ContextLogger) Debugf(format string, args ...interface{}) {
	cl.logger.WithFields(cl.fields).Debugf(format, args...)
}

func (cl *ContextLogger) Infof(format string, args ...interface{}) {
	cl.logger.WithFields(cl.fields).Infof(format, args...)
}

func (cl *ContextLogger) Warnf(format string, args ...interface{}) {
	cl.logger.WithFields(cl.fields).Warnf(format, args...)
}

func (cl *ContextLogger) Errorf(format string, args ...interface{}) {
	cl.logger.WithFields(cl.fields).Errorf(format, args...)
}

func (cl *ContextLogger) Fatalf(format string, args ...interface{}) {
	cl.logger.WithFields(cl.fields).Fatalf(format, args...)
}

func (cl *ContextLogger) Panicf(format string, args ...interface{}) {
	cl.logger.WithFields(cl.fields).Panicf(format, args...)
}

// Performance logging helpers
func (l *Logger) LogDuration(operation string, start time.Time) {
	duration := time.Since(start)
	l.WithFields(map[string]interface{}{
		"operation":   operation,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
	}).Info("Operation completed")
}

func (cl *ContextLogger) LogDuration(operation string, start time.Time) {
	duration := time.Since(start)
	cl.WithFields(map[string]interface{}{
		"operation":   operation,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
	}).Info("Operation completed")
}

// HTTP request logging
func (l *Logger) LogHTTPRequest(method, url string, statusCode int, duration time.Duration) {
	entry := l.WithFields(map[string]interface{}{
		"method":      method,
		"url":         url,
		"status_code": statusCode,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
	})

	message := "HTTP request completed"
	if statusCode >= 500 {
		entry.Error(message)
	} else if statusCode >= 400 {
		entry.Warn(message)
	} else {
		entry.Info(message)
	}
}

// MCP operation logging
func (l *Logger) LogMCPOperation(server, method string, success bool, duration time.Duration) {
	entry := l.WithFields(map[string]interface{}{
		"mcp_server":  server,
		"method":      method,
		"success":     success,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
	})

	message := "MCP operation completed"
	if success {
		entry.Info(message)
	} else {
		entry.Error(message)
	}
}

// Error logging with stack trace
func (l *Logger) LogError(err error, message string, fields map[string]interface{}) {
	entry := l.WithError(err)
	if fields != nil {
		entry = entry.WithFields(fields)
	}

	if l.config.StackTrace {
		// Get stack trace
		stack := make([]byte, 1024*8)
		stack = stack[:runtime.Stack(stack, false)]
		entry = entry.WithField("stack_trace", string(stack))
	}

	entry.Error(message)
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(config LoggerConfig) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// Fallback to default configuration
		config := LoggerConfig{
			Level:  LogLevelInfo,
			Format: LogFormatText,
			Output: LogOutputStderr,
		}
		globalLogger, _ = NewLogger(config)
	}
	return globalLogger
}

// Convenience functions using global logger
func Trace(args ...interface{}) {
	GetGlobalLogger().Trace(args...)
}

func Debug(args ...interface{}) {
	GetGlobalLogger().Debug(args...)
}

func Info(args ...interface{}) {
	GetGlobalLogger().Info(args...)
}

func Warn(args ...interface{}) {
	GetGlobalLogger().Warn(args...)
}

func Error(args ...interface{}) {
	GetGlobalLogger().Error(args...)
}

func Fatal(args ...interface{}) {
	GetGlobalLogger().Fatal(args...)
}

func Panic(args ...interface{}) {
	GetGlobalLogger().Panic(args...)
}

// Formatted convenience functions
func Tracef(format string, args ...interface{}) {
	GetGlobalLogger().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	GetGlobalLogger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GetGlobalLogger().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GetGlobalLogger().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetGlobalLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetGlobalLogger().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	GetGlobalLogger().Panicf(format, args...)
}

// Convenience functions for structured logging
func WithFields(fields map[string]interface{}) *ContextLogger {
	return GetGlobalLogger().WithFields(fields)
}

func WithField(key string, value interface{}) *ContextLogger {
	return GetGlobalLogger().WithField(key, value)
}

func WithError(err error) *ContextLogger {
	return GetGlobalLogger().WithError(err)
}

func WithComponent(component string) *ContextLogger {
	return GetGlobalLogger().WithComponent(component)
}

func WithOperation(operation string) *ContextLogger {
	return GetGlobalLogger().WithOperation(operation)
}

// Performance logging convenience functions
func LogDuration(operation string, start time.Time) {
	GetGlobalLogger().LogDuration(operation, start)
}

func LogHTTPRequest(method, url string, statusCode int, duration time.Duration) {
	GetGlobalLogger().LogHTTPRequest(method, url, statusCode, duration)
}

func LogMCPOperation(server, method string, success bool, duration time.Duration) {
	GetGlobalLogger().LogMCPOperation(server, method, success, duration)
}

func LogError(err error, message string, fields map[string]interface{}) {
	GetGlobalLogger().LogError(err, message, fields)
}

// PrettyPrintJSON formats JSON for pretty printing (helper function)
func PrettyPrintJSON(data interface{}) (string, error) {
	// This is a placeholder - you might want to implement actual JSON formatting
	return fmt.Sprintf("%+v", data), nil
}
