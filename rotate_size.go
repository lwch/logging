package logging

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lwch/runtime"
)

type rotateSizeLogger struct {
	sync.Mutex
	dir         string
	name        string
	rotateSize  int64
	rotateCount int

	// runtime
	currentSize int
	f           *os.File
	l           *log.Logger
	lastCheck   time.Time
}

func NewRotateSizeLogger(dir, name string, size, rotate int) Logger {
	os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(path.Join(dir, name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	runtime.Assert(err)
	fi, err := f.Stat()
	runtime.Assert(err)
	return Logger{&rotateSizeLogger{
		dir:         dir,
		name:        name,
		rotateSize:  int64(size),
		rotateCount: rotate,
		currentSize: int(fi.Size()),
		f:           f,
		l:           log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags),
		lastCheck:   time.Now(),
	}}
}

// SetDateRotate set log rotate by date
func SetSizeRotate(dir, name string, size, rotate int) {
	currentLogger = NewRotateSizeLogger(dir, name, size, rotate)
}

func (l *rotateSizeLogger) rotate() {
	defer func() {
		l.lastCheck = time.Now()
	}()
	// 1% probability to rotate in high rate
	if time.Since(l.lastCheck).Seconds() <= 1 {
		if rand.Intn(100) > 0 {
			return
		}
	}
	fi, err := l.f.Stat()
	if err != nil {
		return
	}
	if fi.Size() < l.rotateSize {
		return
	}
	files, err := filepath.Glob(path.Join(l.dir, l.name+".log.*"))
	if err != nil {
		return
	}
	numbers := make([]int, 0, len(files))
	for _, file := range files {
		ver := strings.TrimPrefix(path.Base(file), l.name+".log.")
		n, _ := strconv.ParseInt(ver, 10, 64)
		numbers = append(numbers, int(n))
	}
	sort.Ints(numbers)
	for i := 0; i < len(numbers)-l.rotateCount+1; i++ {
		os.Remove(path.Join(l.dir, fmt.Sprintf(l.name+".log.%d", numbers[i])))
	}
	latest := 0
	if len(numbers) > 0 {
		latest = numbers[len(numbers)-1]
	}
	os.Rename(path.Join(l.dir, l.name+".log"),
		path.Join(l.dir, fmt.Sprintf(l.name+".log.%d", latest+1)))
	l.f.Close()
	l.f, _ = os.OpenFile(path.Join(l.dir, l.name+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	l.l = log.New(io.MultiWriter(os.Stdout, l.f), "", log.LstdFlags)
}

func (l *rotateSizeLogger) write(fmt string, a ...interface{}) {
	l.l.Printf(fmt, a...)
}

func (l *rotateSizeLogger) flush() {
	l.f.Sync()
}
