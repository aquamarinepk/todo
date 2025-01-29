package am

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	ErrorLevel
)

const (
	colorReset     = "\033[0m"
	colorRed       = "\033[31m"
	colorGreen     = "\033[32m"
	colorYellow    = "\033[33m"
	colorBlue      = "\033[34m"
	colorTimestamp  = "\033[90m"
)

const (
	debugPrefix = colorBlue + "[DBG] " + colorReset
	infoPrefix  = colorGreen + "[INF] " + colorReset
	errorPrefix = colorRed + "[ERR] " + colorReset
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
		debug:    newCustomLogger(os.Stdout, debugPrefix, level, DebugLevel),
		info:     newCustomLogger(os.Stdout, infoPrefix, level, InfoLevel),
		error:    newCustomLogger(os.Stderr, errorPrefix, level, ErrorLevel),
		logLevel: level,
	}
}

func newCustomLogger(out *os.File, prefix string, currentLevel, level LogLevel) *log.Logger {
	return log.New(out, "", 0)
}

func (l *BaseLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

func (l *BaseLogger) Debug(v ...any) {
	if l.logLevel <= DebugLevel {
		l.debug.Output(2, formatLogMessage(debugPrefix, v...))
	}
}

func (l *BaseLogger) Debugf(format string, a ...any) {
	if l.logLevel <= DebugLevel {
		message := fmt.Sprintf(format, a...)
		l.debug.Output(2, formatLogMessage(debugPrefix, message))
	}
}

func (l *BaseLogger) Info(v ...any) {
	if l.logLevel <= InfoLevel {
		l.info.Output(2, formatLogMessage(infoPrefix, v...))
	}
}

func (l *BaseLogger) Infof(format string, a ...any) {
	if l.logLevel <= InfoLevel {
		message := fmt.Sprintf(format, a...)
		l.info.Output(2, formatLogMessage(infoPrefix, message))
	}
}

func (l *BaseLogger) Error(v ...any) {
	if l.logLevel <= ErrorLevel {
		message := fmt.Sprint(v...)
		l.error.Output(2, formatLogMessage(errorPrefix, message))
	}
}

func (l *BaseLogger) Errorf(format string, a ...any) {
	if l.logLevel <= ErrorLevel {
		message := fmt.Sprintf(format, a...)
		l.error.Output(2, formatLogMessage(errorPrefix, message))
	}
}

func formatLogMessage(prefix string, v ...any) string {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	return fmt.Sprintf("%s%s%s %s%s", prefix, colorTimestamp, timestamp, colorReset, fmt.Sprint(v...))
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
	l.debug = log.New(debug, "", 0)
}

// SetInfoOutput set the internal log.
// Used for package testing.
func (l *BaseLogger) SetInfoOutput(info *bytes.Buffer) {
	l.info = log.New(info, "", 0)
}

// SetErrorOutput set the internal log.
// Used for package testing.
func (l *BaseLogger) SetErrorOutput(error *bytes.Buffer) {
	l.error = log.New(error, "", 0)
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[1] = unicode.ToUpper(runes[0])
	return string(runes)
}
