package logging

import "log"

type logger interface {
	rotate()
	printf(string, ...interface{})
	write(string)
	flush()
}

var DefaultLogger Logger = Logger{dummyLogger{}}

type dummyLogger struct{}

func (l dummyLogger) rotate() {}
func (l dummyLogger) printf(fmt string, a ...interface{}) {
	log.Printf(fmt, a...)
}
func (l dummyLogger) write(str string) {
	log.Print(str)
}
func (l dummyLogger) flush() {}
