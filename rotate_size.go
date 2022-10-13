package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lwch/runtime"
)

type SizeRotateConfig struct {
	Dir         string
	Name        string
	Size        int64
	Rotate      int
	WriteStdout bool
	WriteFile   bool
}

type rotateSizeLogger struct {
	sync.Mutex
	cfg SizeRotateConfig

	// runtime
	f *os.File
	w *writer
}

func NewRotateSizeLogger(cfg SizeRotateConfig) Logger {
	var ws []io.Writer
	if cfg.WriteStdout {
		ws = append(ws, os.Stdout)
	}
	var f *os.File
	if cfg.WriteFile {
		os.MkdirAll(cfg.Dir, 0755)
		var err error
		f, err = os.OpenFile(filepath.Join(cfg.Dir, cfg.Name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
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
		logger: &rotateSizeLogger{
			cfg: cfg,
			f:   f,
			w:   newWriter(w),
		},
		lastCheck: time.Now(),
	}
}

// SetDateRotate set log rotate by date
func SetSizeRotate(cfg SizeRotateConfig) {
	DefaultLogger = NewRotateSizeLogger(cfg)
}

func (l *rotateSizeLogger) rotate() {
	if !l.cfg.WriteFile {
		return
	}
	l.Lock()
	defer l.Unlock()
	fi, err := l.f.Stat()
	if err != nil {
		return
	}
	if fi.Size() < l.cfg.Size {
		return
	}
	files, err := filepath.Glob(filepath.Join(l.cfg.Dir, l.cfg.Name+".log.*"))
	if err != nil {
		return
	}
	numbers := make([]int, 0, len(files))
	for _, file := range files {
		ver := strings.TrimPrefix(filepath.Base(file), l.cfg.Name+".log.")
		n, _ := strconv.ParseInt(ver, 10, 64)
		numbers = append(numbers, int(n))
	}
	sort.Ints(numbers)
	for i := 0; i < len(numbers)-l.cfg.Rotate+1; i++ {
		os.Remove(filepath.Join(l.cfg.Dir, fmt.Sprintf(l.cfg.Name+".log.%d", numbers[i])))
	}
	latest := 0
	if len(numbers) > 0 {
		latest = numbers[len(numbers)-1]
	}
	l.f.Close()
	os.Rename(filepath.Join(l.cfg.Dir, l.cfg.Name+".log"),
		filepath.Join(l.cfg.Dir, fmt.Sprintf(l.cfg.Name+".log.%d", latest+1)))
	l.f, _ = os.OpenFile(filepath.Join(l.cfg.Dir, l.cfg.Name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	var w io.Writer
	if l.cfg.WriteStdout {
		w = io.MultiWriter(os.Stdout, l.f)
	} else {
		w = l.f
	}
	l.w = newWriter(w)
}

func (l *rotateSizeLogger) printf(fmt string, a ...interface{}) {
	l.w.Printf(fmt, a...)
}

func (l *rotateSizeLogger) write(str string) {
	l.w.Write(str)
}

func (l *rotateSizeLogger) flush() {
	f := l.f
	if f != nil {
		f.Sync()
	}
}
