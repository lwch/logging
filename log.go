package logging

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/lwch/runtime"
)

func init() {
	log.SetOutput(os.Stdout)
}

// SetLevel set log level
func SetLevel(level Level) {
	DefaultLogger.setLevel(level)
}

// Debug debug log
func Debug(fmt string, a ...interface{}) {
	DefaultLogger.Debug(fmt, a...)
}

// Info info log
func Info(fmt string, a ...interface{}) {
	DefaultLogger.Info(fmt, a...)
}

// Error error log
func Error(fmt string, a ...interface{}) {
	DefaultLogger.Error(fmt, a...)
}

// Warning warning log
func Warning(fmt string, a ...interface{}) {
	DefaultLogger.Warning(fmt, a...)
}

// Printf print log
func Printf(fmt string, a ...interface{}) {
	DefaultLogger.Printf(fmt, a...)
}

// Flush flush log
func Flush() {
	DefaultLogger.flush()
}

type Logger struct {
	logger
	lastCheck time.Time
}

func (l *Logger) rateLimit() bool {
	if time.Since(l.lastCheck).Seconds() <= 1 {
		if rand.Intn(100) > 0 {
			return true
		}
	}
	return false
}

func (l *Logger) resetLastCheck() {
	l.lastCheck = time.Now()
}

// Debug debug log
func (l *Logger) Debug(fmt string, a ...interface{}) {
	defer l.resetLastCheck()
	if !l.rateLimit() {
		l.logger.rotate()
	}
	if l.currentLevel() >= LevelDebug {
		l.logger.printf("[DEBUG]"+fmt, a...)
	}
}

// Info info log
func (l *Logger) Info(fmt string, a ...interface{}) {
	defer l.resetLastCheck()
	if !l.rateLimit() {
		l.logger.rotate()
	}
	if l.currentLevel() >= LevelInfo {
		l.logger.printf("[INFO]"+fmt, a...)
	}
}

// Error error log
func (l *Logger) Error(fmt string, a ...interface{}) {
	defer l.resetLastCheck()
	if !l.rateLimit() {
		l.logger.rotate()
	}
	if l.currentLevel() >= LevelError {
		trace := strings.Join(runtime.Trace("  + "), separator)
		l.logger.printf("[ERROR]"+fmt+separator+trace, a...)
	}
}

// Warning warning log
func (l *Logger) Warning(fmt string, a ...interface{}) {
	defer l.resetLastCheck()
	if !l.rateLimit() {
		l.logger.rotate()
	}
	if l.currentLevel() >= LevelWarn {
		l.logger.printf("[WARN]"+fmt, a...)
	}
}

// Printf print log
func (l *Logger) Printf(fmt string, a ...interface{}) {
	defer l.resetLastCheck()
	if !l.rateLimit() {
		l.logger.rotate()
	}
	l.logger.printf(fmt, a...)
}

// Write write log
func (l *Logger) Write(data []byte) (int, error) {
	defer l.resetLastCheck()
	if !l.rateLimit() {
		l.logger.rotate()
	}
	str := string(data)
	str = strings.TrimSuffix(str, "\n")
	str = strings.TrimSuffix(str, "\r")
	l.logger.write(str)
	return len(data), nil
}

// Flush flush log
func (l *Logger) Flush() {
	l.logger.flush()
}
