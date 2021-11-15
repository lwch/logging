package logging

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

type writer struct {
	sync.Mutex
	w io.Writer
}

var separator = "\n"

func init() {
	if runtime.GOOS == "windows" {
		separator = "\r\n"
	}
}

func newWriter(w io.Writer) *writer {
	return &writer{w: w}
}

func (w *writer) Printf(format string, v ...interface{}) {
	w.log(fmt.Sprintf(format, v...))
}

func (w *writer) Write(str string) {
	w.log(str)
}

func (w *writer) log(s string) {
	s = time.Now().Format("2006/01/02 15:04:05 ") + s + separator
	w.Lock()
	defer w.Unlock()
	w.w.Write([]byte(s))
}
