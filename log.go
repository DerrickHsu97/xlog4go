package dlog4go

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PUBLIC"}
	recordPool  *sync.Pool
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
	PUBLIC
)

const tunnel_size_default = 1024

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string {
	return fmt.Sprintf("[%s] [%s] [%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	writers   []Writer
	tunnel    chan *Record
	level     int
	c         chan bool
	layout    string
	levelFunc func(int) int
}

func NewLogger() *Logger {
	if logger_default != nil && !takeup {
		takeup = true
		return logger_default
	}

	l := new(Logger)
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.c = make(chan bool, 1)
	l.level = DEBUG
	l.layout = "2006-01-02T15:04:05.000+0800"

	go boostrapLogWriter(l)

	return l
}

func (l *Logger) SetLevelFunc(levelFunc func(int) int) {
	l.levelFunc = levelFunc
}

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

func (l *Logger) Public(fmt string, args ...interface{}) {
	l.DeliverRecordToWriter(PUBLIC, fmt, args...)
}

func (l *Logger) Trace(fmt string, args ...interface{}) {
	l.DeliverRecordToWriter(TRACE, fmt, args...)
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.DeliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.DeliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	l.DeliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	l.DeliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.DeliverRecordToWriter(FATAL, fmt, args...)
}

func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c

	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *Logger) getLevel(lvl int) bool {
	if l.levelFunc == nil {
		return lvl < l.level
	}
	level := l.levelFunc(lvl)
	if level <= PUBLIC && level >= TRACE && level > l.level {
		return lvl < level
	}
	return lvl < l.level
}

func (l *Logger) DeliverRecordToWriter(level int, format string, args ...interface{}) {
	var inf, code string

	if l.getLevel(level) {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	// source code, file and line num
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}

	r := recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = time.Now().Format(l.layout)
	r.level = level

	l.tunnel <- r
}

func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	var (
		r  *Record
		ok bool
	)

	if r, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}

	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}

			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}

			recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}

// default
var (
	logger_default *Logger
	takeup         = false
)

func SetLevel(lvl int) {
	logger_default.level = lvl
}

func SetLayout(layout string) {
	logger_default.layout = layout
}

func Public(fmt string, args ...interface{}) {
	logger_default.DeliverRecordToWriter(PUBLIC, fmt, args...)
}

func Trace(fmt string, args ...interface{}) {
	logger_default.DeliverRecordToWriter(TRACE, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	logger_default.DeliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	logger_default.DeliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	logger_default.DeliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	logger_default.DeliverRecordToWriter(ERROR, fmt, args...)
}

func SetLevelFunc(levelFunc func(int) int) {
	logger_default.SetLevelFunc(levelFunc)
}

func Fatal(fmt string, args ...interface{}) {
	logger_default.DeliverRecordToWriter(FATAL, fmt, args...)
}

func Register(w Writer) {
	logger_default.Register(w)
}

func Close() {
	logger_default.Close()
}

func GetLogger() *Logger {
	return logger_default
}

func init() {
	logger_default = NewLogger()
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
}
