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
	currentLogger.Debug(fmt, a...)
}

// Info info log
func Info(fmt string, a ...interface{}) {
	currentLogger.Info(fmt, a...)
}

// Error error log
func Error(fmt string, a ...interface{}) {
	currentLogger.Error(fmt, a...)
}

// Write write log
func Write(fmt string, a ...interface{}) {
	currentLogger.Write(fmt, a...)
}

// Flush flush log
func Flush() {
	currentLogger.flush()
}

type Logger struct {
	logger
}

// Debug debug log
func (l Logger) Debug(fmt string, a ...interface{}) {
	l.logger.rotate()
	if rand.Intn(1000) < 1 {
		l.logger.write("[DEBUG]"+fmt, a...)
	}
}

// Info info log
func (l Logger) Info(fmt string, a ...interface{}) {
	l.logger.rotate()
	l.logger.write("[INFO]"+fmt, a...)
}

// Error error log
func (l Logger) Error(fmt string, a ...interface{}) {
	l.logger.rotate()
	trace := strings.Join(runtime.Trace("  + "), "\n")
	l.logger.write("[ERROR]"+fmt+"\n"+trace, a...)
}

// Write write log
func (l Logger) Write(fmt string, a ...interface{}) {
	l.logger.rotate()
	l.logger.write(fmt, a...)
}

// Flush flush log
func (l Logger) Flush() {
	l.logger.flush()
}
