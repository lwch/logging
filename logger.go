package logging

import (
	"log"
	"time"
)

type logger interface {
	currentLevel() Level
	rotate()
	printf(string, ...interface{})
	write(string)
	flush()
}

var DefaultLogger Logger = Logger{
	logger:    dummyLogger{},
	lastCheck: time.Now(),
}

type dummyLogger struct{}

func (l dummyLogger) currentLevel() Level {
	return LevelInfo
}
func (l dummyLogger) rotate() {}
func (l dummyLogger) printf(fmt string, a ...interface{}) {
	log.Printf(fmt, a...)
}
func (l dummyLogger) write(str string) {
	log.Print(str)
}
func (l dummyLogger) flush() {}
