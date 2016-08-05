package slog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	FATAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

var levelName = map[string]LogLevel{
	"FATAL": FATAL,
	"ERROR": ERROR,
	"WARN":  WARN,
	"INFO":  INFO,
	"DEBUG": DEBUG,
}

var (
	maxLogLevel = DEBUG
	logflag     = log.LstdFlags | log.Lshortfile
	cuted       = false
	dateFormat  = "2006-01-02"
	logStd      = newWrite(os.Stdout, "")
)

// log日志
type Logger struct {
	lock    sync.Mutex
	oldDate string
	logfile string
	fd      *os.File
	logger  *log.Logger
}

// 设置Access对象
func New(file, prefix string) *Logger {
	Writer := createLogger(file)
	return newWrite(Writer, prefix)
}

func newWrite(Writer io.Writer, prefix string) *Logger {
	l := &Logger{
		logfile: "",
		fd:      nil,
		oldDate: time.Now().Format(dateFormat),
	}
	l.logger = log.New(Writer, prefix, logflag)
	return l
}

func createLogger(accessFile string) *os.File {
	requestWriter, err := os.OpenFile(accessFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		Fatal(accessFile, err)
	}
	return requestWriter
}

// SetLogLevel sets MaxLogLevel based on the provided string
func SetLogLevel(level string) {
	level = strings.ToUpper(level)
	lev, ok := levelName[level]
	if !ok {
		log.Fatalf("Unknown log level requested: %v", level)
	}
	maxLogLevel = lev
}

func GetLogLevel() LogLevel {
	return maxLogLevel
}

func SetCut(cut bool) {
	cuted = cut
}

func GetStdLog() *Logger {
	return logStd
}

func (l *Logger) Fatal(args ...interface{}) {
	if FATAL > maxLogLevel {
		return
	}
	l.output("FATAL", fmt.Sprint(args...))
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	if FATAL > maxLogLevel {
		return
	}
	l.output("FATAL", fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Error logs a message to the 'standard' Logger (always)
func (l *Logger) Error(args ...interface{}) {
	if ERROR > maxLogLevel {
		return
	}
	l.output("ERROR", fmt.Sprint(args...))
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	if ERROR > maxLogLevel {
		return
	}
	l.output("ERROR", fmt.Sprintf(msg, args...))
}

// Warn logs a message to the 'standard' Logger if MaxLogLevel is >= WARN
func (l *Logger) Warn(args ...interface{}) {
	if WARN > maxLogLevel {
		return
	}
	l.output("WARN", fmt.Sprint(args...))
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	if WARN > maxLogLevel {
		return
	}
	l.output("WARN", fmt.Sprintf(msg, args...))
}

// Info logs a message to the 'standard' Logger if MaxLogLevel is >= INFO
func (l *Logger) Info(args ...interface{}) {
	if INFO > maxLogLevel {
		return
	}
	l.output("INFO", fmt.Sprint(args...))
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	if INFO > maxLogLevel {
		return
	}
	l.output("INFO", fmt.Sprintf(msg, args...))
}

// Trace logs a message to the 'standard' Logger if MaxLogLevel is >= DEBUG
func (l *Logger) Debug(args ...interface{}) {
	if DEBUG > maxLogLevel {
		return
	}
	l.output("DEBUG", fmt.Sprint(args...))
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	if DEBUG > maxLogLevel {
		return
	}
	l.output("DEBUG", fmt.Sprintf(msg, args...))
}

func (l *Logger) output(mode, msg string) {
	if cuted && l.fd != nil {
		nowDate := time.Now().Format(dateFormat)
		if l.oldDate != nowDate {
			l.lock.Lock()
			defer l.lock.Unlock()
			oldDate := l.oldDate
			l.oldDate = nowDate
			l.fd.Close()
			err := os.Rename(l.logfile, l.logfile+oldDate)
			if err != nil {
				Error(err)
			}
			requestWriter := createLogger(l.logfile)
			l.fd = requestWriter
			l.logger = log.New(requestWriter, "", log.LstdFlags)
		}
	}
	l.logger.Output(3, "["+mode+"] "+msg)
}

func Fatal(args ...interface{}) {
	if FATAL > maxLogLevel {
		return
	}
	logStd.output("FATAL", fmt.Sprint(args...))
	os.Exit(1)
}

func Fatalf(msg string, args ...interface{}) {
	if FATAL > maxLogLevel {
		return
	}
	logStd.output("FATAL", fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Error logs a message to the 'standard' Logger (always)
func Error(args ...interface{}) {
	if ERROR > maxLogLevel {
		return
	}
	logStd.output("ERROR", fmt.Sprint(args...))
}

func Errorf(msg string, args ...interface{}) {
	if ERROR > maxLogLevel {
		return
	}
	logStd.output("ERROR", fmt.Sprintf(msg, args...))
}

// Warn logs a message to the 'standard' Logger if MaxLogLevel is >= WARN
func Warn(args ...interface{}) {
	if WARN > maxLogLevel {
		return
	}
	logStd.output("WARN", fmt.Sprint(args...))
}

func Warnf(msg string, args ...interface{}) {
	if WARN > maxLogLevel {
		return
	}
	logStd.output("WARN", fmt.Sprintf(msg, args...))
}

// Info logs a message to the 'standard' Logger if MaxLogLevel is >= INFO
func Info(args ...interface{}) {
	if INFO > maxLogLevel {
		return
	}
	logStd.output("INFO", fmt.Sprint(args...))
}

func Infof(msg string, args ...interface{}) {
	if INFO > maxLogLevel {
		return
	}
	logStd.output("INFO", fmt.Sprintf(msg, args...))
}

// Trace logs a message to the 'standard' Logger if MaxLogLevel is >= DEBUG
func Debug(args ...interface{}) {
	if DEBUG > maxLogLevel {
		return
	}
	logStd.output("DEBUG", fmt.Sprint(args...))
}

func Debugf(msg string, args ...interface{}) {
	if DEBUG > maxLogLevel {
		return
	}
	logStd.output("DEBUG", fmt.Sprintf(msg, args...))
}
