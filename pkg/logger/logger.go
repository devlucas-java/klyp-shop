// Package logger provides a simple leveled logger with optional component context
// for easier debugging and log tracing.
//
// Example usage:
//
//	log := logger.NewLogger(logger.DEBUG)
//	log = log.WithComponent("AuthService")
//	log.Info("user logged in")
//
// Or with explicit component creation:
//
//	log := logger.NewLoggerWithComponent(logger.INFO, "OrderService")
//	log.Warn("inventory low")
package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

type Level int32

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	level     atomic.Int32
	mu        sync.Mutex
	component string

	fatal *log.Logger
	err   *log.Logger
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger
	trace *log.Logger
}

var instance = newLogger(DEBUG, "")

// Instance returns the global singleton logger.
func Instance() *Logger {
	return instance
}

// SetLogLevel sets the level for the global singleton logger.
func SetLogLevel(level Level) {
	instance.level.Store(int32(level))
}

// NewLogger creates a logger with the provided level and no component prefix.
func NewLogger(level Level) *Logger {
	return newLogger(level, "")
}

// NewLoggerWithComponent creates a logger that prefixes every line with the
// provided component name, improving traceability in logs.
func NewLoggerWithComponent(level Level, component string) *Logger {
	return newLogger(level, component)
}

// WithComponent returns a derived logger with the same level but a new component prefix.
func (l *Logger) WithComponent(component string) *Logger {
	return newLogger(Level(l.level.Load()), component)
}

func newLogger(level Level, component string) *Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile
	componentPrefix := ""
	if component != "" {
		componentPrefix = fmt.Sprintf("[%s] ", component)
	}

	l := &Logger{
		component: component,
		fatal:     log.New(os.Stderr, fmt.Sprintf("%s[FATAL] ", componentPrefix), flags),
		err:       log.New(os.Stderr, fmt.Sprintf("%s[ERROR] ", componentPrefix), flags),
		warn:      log.New(os.Stdout, fmt.Sprintf("%s[WARN]  ", componentPrefix), flags),
		info:      log.New(os.Stdout, fmt.Sprintf("%s[INFO]  ", componentPrefix), flags),
		debug:     log.New(os.Stdout, fmt.Sprintf("%s[DEBUG] ", componentPrefix), flags),
		trace:     log.New(os.Stdout, fmt.Sprintf("%s[TRACE] ", componentPrefix), flags),
	}

	l.level.Store(int32(level))

	return l
}

func (l *Logger) canLog(level Level) bool {
	return level >= Level(l.level.Load())
}

func (l *Logger) Trace(v ...any) {
	if l.canLog(TRACE) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.trace.Println(v...)
	}
}

func (l *Logger) Debug(v ...any) {
	if l.canLog(DEBUG) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.debug.Println(v...)
	}
}

func (l *Logger) Info(v ...any) {
	if l.canLog(INFO) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.info.Println(v...)
	}
}

func (l *Logger) Warn(v ...any) {
	if l.canLog(WARN) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.warn.Println(v...)
	}
}

func (l *Logger) Error(v ...any) {
	if l.canLog(ERROR) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.err.Println(v...)
	}
}

func (l *Logger) Fatal(v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fatal.Fatalln(v...)
}

func (l *Logger) Tracef(format string, v ...any) {
	if l.canLog(TRACE) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.trace.Printf(format, v...)
	}
}

func (l *Logger) Debugf(format string, v ...any) {
	if l.canLog(DEBUG) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.debug.Printf(format, v...)
	}
}

func (l *Logger) Infof(format string, v ...any) {
	if l.canLog(INFO) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.info.Printf(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...any) {
	if l.canLog(WARN) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.warn.Printf(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	if l.canLog(ERROR) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.err.Printf(format, v...)
	}
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fatal.Fatalf(format, v...)
}
