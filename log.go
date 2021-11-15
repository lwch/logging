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
	rand.Seed(time.Now().UnixNano())
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
}

// Debug debug log
func (l Logger) Debug(fmt string, a ...interface{}) {
	l.logger.rotate()
	if rand.Intn(1000) < 1 {
		l.logger.printf("[DEBUG]"+fmt, a...)
	}
}

// Info info log
func (l Logger) Info(fmt string, a ...interface{}) {
	l.logger.rotate()
	l.logger.printf("[INFO]"+fmt, a...)
}

// Error error log
func (l Logger) Error(fmt string, a ...interface{}) {
	l.logger.rotate()
	trace := strings.Join(runtime.Trace("  + "), "\n")
	l.logger.printf("[ERROR]"+fmt+"\n"+trace, a...)
}

// Printf print log
func (l Logger) Printf(fmt string, a ...interface{}) {
	l.logger.rotate()
	l.logger.printf(fmt, a...)
}

// Write write log
func (l Logger) Write(data []byte) (int, error) {
	l.logger.rotate()
	l.logger.write(strings.TrimSuffix(string(data), "\n"))
	return len(data), nil
}

// Flush flush log
func (l Logger) Flush() {
	l.logger.flush()
}
