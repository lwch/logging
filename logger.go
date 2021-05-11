package logging

import "log"

type logger interface {
	rotate()
	write(string, ...interface{})
	flush()
}

var currentLogger Logger = Logger{dummyLogger{}}

type dummyLogger struct{}

func (l dummyLogger) rotate() {}
func (l dummyLogger) write(fmt string, a ...interface{}) {
	log.Printf(fmt, a...)
}
func (l dummyLogger) flush() {}
