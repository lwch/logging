package logging

import (
	"log"
	"time"
)

type logger interface {
	setLevel(Level)
	currentLevel() Level
	rotate()
	printf(string, ...interface{})
	write(string)
	flush()
}

var DefaultLogger = Logger{
	logger:    &dummyLogger{level: LevelInfo},
	lastCheck: time.Now(),
}

type dummyLogger struct {
	level Level
}

func (l *dummyLogger) setLevel(level Level) {
	l.level = level
}
func (l dummyLogger) currentLevel() Level {
	return l.level
}
func (l dummyLogger) rotate() {}
func (l dummyLogger) printf(fmt string, a ...interface{}) {
	log.Printf(fmt, a...)
}
func (l dummyLogger) write(str string) {
	log.Print(str)
}
func (l dummyLogger) flush() {}
