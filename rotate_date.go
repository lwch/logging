package logging

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lwch/runtime"
)

type RotateDateLogger struct {
	sync.Mutex
	dir        string
	name       string
	date       string
	rotateDays int
	stdout     bool

	// runtime
	f *os.File
	w *writer
}

func NewRotateDateLogger(dir, name string, rotate int, stdout bool) Logger {
	os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(path.Join(dir, name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	runtime.Assert(err)
	var w io.Writer
	if stdout {
		w = io.MultiWriter(os.Stdout, f)
	} else {
		w = f
	}
	return Logger{&RotateDateLogger{
		dir:        dir,
		name:       name,
		date:       time.Now().Format("20060102"),
		rotateDays: rotate,
		stdout:     stdout,
		f:          f,
		w:          newWriter(w),
	}}
}

// SetDateRotate set log rotate by date
func SetDateRotate(dir, name string, rotate int, stdout bool) {
	DefaultLogger = NewRotateDateLogger(dir, name, rotate, stdout)
}

func (l *RotateDateLogger) rotate() {
	now := time.Now().Format("20060102")
	if l.date == now {
		return
	}
	l.Lock()
	defer l.Unlock()
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
	var w io.Writer
	if l.stdout {
		w = io.MultiWriter(os.Stdout, l.f)
	} else {
		w = l.f
	}
	l.w = newWriter(w)
	l.date = now
}

func (l *RotateDateLogger) write(fmt string, a ...interface{}) {
	l.w.Printf(fmt, a...)
}

func (l *RotateDateLogger) flush() {
	l.f.Sync()
}
