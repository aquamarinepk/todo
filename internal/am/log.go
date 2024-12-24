package am

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	ErrorLevel
)

type Logger interface {
	SetLogLevel(level LogLevel)
	Debug(v ...any)
	Debugf(format string, a ...any)
	Info(v ...any)
	Infof(format string, a ...any)
	Error(v ...any)
	Errorf(format string, a ...any)
}

type BaseLogger struct {
	debug    *log.Logger
	info     *log.Logger
	error    *log.Logger
	logLevel LogLevel
}

func NewLogger(logLevel string) *BaseLogger {
	level := ToValidLevel(logLevel)
	return &BaseLogger{
		debug:    log.New(os.Stdout, "[DBG] ", log.LstdFlags),
		info:     log.New(os.Stdout, "[INF] ", log.LstdFlags),
		error:    log.New(os.Stderr, "[ERR] ", log.LstdFlags),
		logLevel: level,
	}
}

func (l *BaseLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

func (l *BaseLogger) Debug(v ...any) {
	if l.logLevel <= DebugLevel {
		l.debug.Println(v...)
	}
}

func (l *BaseLogger) Debugf(format string, a ...any) {
	if l.logLevel <= DebugLevel {
		message := fmt.Sprintf(format, a...)
		l.debug.Println(message)
	}
}

func (l *BaseLogger) Info(v ...any) {
	if l.logLevel <= InfoLevel {
		l.info.Println(v...)
	}
}

func (l *BaseLogger) Infof(format string, a ...any) {
	if l.logLevel <= InfoLevel {
		message := fmt.Sprintf(format, a...)
		l.info.Println(message)
	}
}

func (l *BaseLogger) Error(v ...interface{}) {
	if l.logLevel <= ErrorLevel {
		message := fmt.Sprint(v...)
		l.error.Println(message)
	}
}

func (l *BaseLogger) Errorf(format string, a ...interface{}) {
	if l.logLevel <= ErrorLevel {
		message := fmt.Sprintf(format, a...)
		l.error.Println(message)
	}
}

func ToValidLevel(level string) LogLevel {
	level = strings.ToLower(level)

	switch level {
	case "debug", "dbg":
		return DebugLevel
	case "info", "inf":
		return InfoLevel
	case "error", "err":
		return ErrorLevel
	default:
		return ErrorLevel
	}
}

// SetDebugOutput set the internal log.
// Used for package testing.
func (l *BaseLogger) SetDebugOutput(debug *bytes.Buffer) {
	l.debug = log.New(debug, "", 1)
}

// SetInfoOutput set the internal log.
// Used for package testing.
func (l *BaseLogger) SetInfoOutput(info *bytes.Buffer) {
	l.info = log.New(info, "", 1)
}

// SetErrorOutput set the internal log.
// Used for package testing.
func (l *BaseLogger) SetErrorOutput(error *bytes.Buffer) {
	l.error = log.New(error, "", 1)
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[1] = unicode.ToUpper(runes[0])
	return string(runes)
}
