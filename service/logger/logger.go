package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *Logger

func init() {
	logger = NewLogger()
}

// Logger will represent the logrus logger
type Logger struct {
	*logrus.Logger
}

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// NewLogger will return an instantiated logrus logger
func NewLogger() *Logger {
	var baseLogger = logrus.New()

	var logger = &Logger{
		baseLogger,
	}

	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		PadLevelText:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	// Log Level
	logger.SetLevel(logrus.InfoLevel)
	if os.Getenv("DEBUG") == "true" {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}

// Declare variables to store log messages as new Events
var (
	missingEnvMessage = Event{1, "Missing env key: %s"}
)

// MissingEnv is a standard error message
func MissingEnv(envName string) {
	logger.Panicf(missingEnvMessage.message, envName)
}

// Debug Log
func Debug(args ...interface{}) {
	logger.Debugln(args...)
}

// Debugf Log
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info Log
func Info(args ...interface{}) {
	logger.Infoln(args...)
}

// Infof Log
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn Log
func Warn(args ...interface{}) {
	logger.Warnln(args...)
}

// Warnf Log
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Panic Log
func Panic(args ...interface{}) {
	logger.Panicln(args...)
}

// Panicf Log
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Error Log
func Error(args ...interface{}) {
	logger.Errorln(args...)
}

// Errorf Log
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal Log
func Fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

// Fatalf Log
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
