package logging

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lwch/runtime"
)

type DateRotateConfig struct {
	Dir         string
	Name        string
	Rotate      int
	WriteStdout bool
	WriteFile   bool
}

type RotateDateLogger struct {
	sync.Mutex
	date string
	cfg  DateRotateConfig

	// runtime
	f *os.File
	w *writer
}

func NewRotateDateLogger(cfg DateRotateConfig) Logger {
	var ws []io.Writer
	if cfg.WriteStdout {
		ws = append(ws, os.Stdout)
	}
	var f *os.File
	if cfg.WriteFile {
		os.MkdirAll(cfg.Dir, 0755)
		var err error
		f, err = os.OpenFile(path.Join(cfg.Dir, cfg.Name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		runtime.Assert(err)
		ws = append(ws, f)
	}
	if len(ws) == 0 {
		panic(errors.New("no output"))
	}
	var w io.Writer
	if len(ws) == 1 {
		w = ws[0]
	} else {
		w = io.MultiWriter(ws[0], ws[1])
	}
	return Logger{
		logger: &RotateDateLogger{
			date: time.Now().Format("20060102"),
			cfg:  cfg,
			f:    f,
			w:    newWriter(w),
		},
		lastCheck: time.Now(),
	}
}

// SetDateRotate set log rotate by date
func SetDateRotate(cfg DateRotateConfig) {
	DefaultLogger = NewRotateDateLogger(cfg)
}

func (l *RotateDateLogger) rotate() {
	now := time.Now().Format("20060102")
	if l.date == now {
		return
	}
	if !l.cfg.WriteFile {
		return
	}
	l.Lock()
	defer l.Unlock()
	files, _ := filepath.Glob(path.Join(l.cfg.Dir, l.cfg.Name+"_*.log"))
	for _, file := range files {
		date := strings.TrimPrefix(path.Base(file), l.cfg.Name+"_")
		date = strings.TrimSuffix(date, ".log")
		t, _ := time.Parse("20060102", date)
		if time.Since(t).Hours() > float64(24*l.cfg.Rotate) {
			os.Remove(file)
		}
	}
	os.Rename(path.Join(l.cfg.Dir, l.cfg.Name+".log"),
		path.Join(l.cfg.Dir, l.cfg.Name+"_"+l.date+".log"))
	l.f.Close()
	l.f, _ = os.OpenFile(path.Join(l.cfg.Dir, l.cfg.Name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	var w io.Writer
	if l.cfg.WriteStdout {
		w = io.MultiWriter(os.Stdout, l.f)
	} else {
		w = l.f
	}
	l.w = newWriter(w)
	l.date = now
}

func (l *RotateDateLogger) printf(fmt string, a ...interface{}) {
	l.w.Printf(fmt, a...)
}

func (l *RotateDateLogger) write(str string) {
	l.w.Write(str)
}

func (l *RotateDateLogger) flush() {
	f := l.f
	if f != nil {
		f.Sync()
	}
}
