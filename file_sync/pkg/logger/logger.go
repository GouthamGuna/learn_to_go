package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	debugLogger  *log.Logger
	infoLogger   *log.Logger
	warnLogger   *log.Logger
	errorLogger  *log.Logger
	fatalLogger  *log.Logger
	currentLevel LogLevel
	mu           sync.Mutex
)

func init() {
	// Default to info level
	currentLevel = LevelInfo
	initLoggers(os.Stdout, os.Stderr)
}

func initLoggers(infoHandle, errorHandle io.Writer) {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	debugLogger = log.New(infoHandle, "DEBUG: ", flags)
	infoLogger = log.New(infoHandle, "INFO: ", flags)
	warnLogger = log.New(infoHandle, "WARN: ", flags)
	errorLogger = log.New(errorHandle, "ERROR: ", flags)
	fatalLogger = log.New(errorHandle, "FATAL: ", flags)
}

// SetLevel sets the log level
func SetLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	currentLevel = level
}

// SetOutput sets the output destination for all loggers
func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	initLoggers(w, w)
}

// SetOutputs sets different outputs for info and error logs
func SetOutputs(infoWriter, errorWriter io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	initLoggers(infoWriter, errorWriter)
}

// Debug logs debug message
func Debug(v ...interface{}) {
	if currentLevel <= LevelDebug {
		debugLogger.Output(2, fmt.Sprint(v...))
	}
}

// Debugf logs formatted debug message
func Debugf(format string, v ...interface{}) {
	if currentLevel <= LevelDebug {
		debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info logs info message
func Info(v ...interface{}) {
	if currentLevel <= LevelInfo {
		infoLogger.Output(2, fmt.Sprint(v...))
	}
}

// Infof logs formatted info message
func Infof(format string, v ...interface{}) {
	if currentLevel <= LevelInfo {
		infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn logs warning message
func Warn(v ...interface{}) {
	if currentLevel <= LevelWarn {
		warnLogger.Output(2, fmt.Sprint(v...))
	}
}

// Warnf logs formatted warning message
func Warnf(format string, v ...interface{}) {
	if currentLevel <= LevelWarn {
		warnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error logs error message
func Error(v ...interface{}) {
	if currentLevel <= LevelError {
		errorLogger.Output(2, fmt.Sprint(v...))
	}
}

// Errorf logs formatted error message
func Errorf(format string, v ...interface{}) {
	if currentLevel <= LevelError {
		errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal logs fatal message and exits
func Fatal(v ...interface{}) {
	fatalLogger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf logs formatted fatal message and exits
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// GetLevel returns current log level
func GetLevel() LogLevel {
	mu.Lock()
	defer mu.Unlock()
	return currentLevel
}

// StringToLogLevel converts string to LogLevel
func StringToLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}
