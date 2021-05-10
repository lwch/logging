package logging

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lwch/runtime"
)

type rotateDateLogger struct {
	sync.Mutex
	dir        string
	name       string
	date       string
	rotateDays int

	// runtime
	f *os.File
	l *log.Logger
}

func newRotateDateLogger(dir, name string, rotate int) *rotateDateLogger {
	os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(path.Join(dir, name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	runtime.Assert(err)
	return &rotateDateLogger{
		dir:        dir,
		name:       name,
		date:       time.Now().Format("20060102"),
		rotateDays: rotate,
		f:          f,
		l:          log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags),
	}
}

// SetDateRotate set log rotate by date
func SetDateRotate(dir, name string, rotate int) {
	currentLogger = newRotateDateLogger(dir, name, rotate)
}

func (l *rotateDateLogger) rotate() {
	now := time.Now().Format("20060102")
	if l.date == now {
		return
	}
	files, _ := filepath.Glob(path.Join(l.dir, l.name+"_*.log"))
	for _, file := range files {
		date := strings.TrimPrefix(path.Base(file), l.name+"_")
		date = strings.TrimSuffix(date, ".log")
		t, _ := time.Parse("20060102", date)
		if time.Since(t).Hours() > float64(24*l.rotateDays) {
			os.Remove(file)
		}
	}
	os.Rename(path.Join(l.dir, l.name+".log"),
		path.Join(l.dir, l.name+"_"+l.date+".log"))
	l.f.Close()
	l.f, _ = os.OpenFile(path.Join(l.dir, l.name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	l.l = log.New(io.MultiWriter(os.Stdout, l.f), "", log.LstdFlags)
	l.date = now
}

func (l *rotateDateLogger) write(fmt string, a ...interface{}) {
	l.l.Printf(fmt, a...)
}

func (l *rotateDateLogger) flush() {
	l.f.Sync()
}
