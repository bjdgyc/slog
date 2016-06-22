package slog

import (
	"fmt"
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
	logout      *log.Logger
	logAccess   *AccessLog
)

// request 按天分割日志
type AccessLog struct {
	lock    sync.Mutex
	oldDate string
	logfile string
	fd      *os.File
	access  *log.Logger
}

func init() {
	logout = log.New(os.Stdout, "", logflag)
}

func SetLogfile(outfile string) {
	fileWriter, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		Fatal(outfile, err)
	}
	logout = log.New(fileWriter, "", logflag)
}

// 设置request对象
func SetRequestfile(requestfile string) {
	dateFormat := "2006-01-02"
	requestWriter := createRequestlogger(requestfile)
	logAccess = &AccessLog{
		logfile: requestfile,
		fd:      requestWriter,
		oldDate: time.Now().Format(dateFormat),
	}
	logAccess.access = log.New(logAccess.fd, "", log.LstdFlags)

	//异步切割日志
	go func() {
		ticker := time.NewTicker(time.Minute * 3)
		for range ticker.C {
			nowDate := time.Now().Format(dateFormat)
			if logAccess.oldDate != nowDate {
				logAccess.lock.Lock()
				logAccess.oldDate = nowDate
				logAccess.fd.Close()
				err := os.Rename(requestfile, requestfile+nowDate)
				if err != nil {
					Error(err)
				}
				requestWriter := createRequestlogger(requestfile)
				logAccess.fd = requestWriter
				logAccess.access = log.New(logAccess.fd, "", log.LstdFlags)
				logAccess.lock.Unlock()
			}
		}
	}()

}

func createRequestlogger(requestfile string) *os.File {
	requestWriter, err := os.OpenFile(requestfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		Fatal(requestfile, err)
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

func Fatal(args ...interface{}) {
	if FATAL > maxLogLevel {
		return
	}
	output("FATAL", fmt.Sprint(args...))
	os.Exit(1)
}

func Fatalf(msg string, args ...interface{}) {
	if FATAL > maxLogLevel {
		return
	}
	output("FATAL", fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Error logs a message to the 'standard' Logger (always)
func Error(args ...interface{}) {
	if ERROR > maxLogLevel {
		return
	}
	output("ERROR", fmt.Sprint(args...))
}

func Errorf(msg string, args ...interface{}) {
	if ERROR > maxLogLevel {
		return
	}
	output("ERROR", fmt.Sprintf(msg, args...))
}

// Warn logs a message to the 'standard' Logger if MaxLogLevel is >= WARN
func Warn(args ...interface{}) {
	if WARN > maxLogLevel {
		return
	}
	output("WARN", fmt.Sprint(args...))
}

func Warnf(msg string, args ...interface{}) {
	if WARN > maxLogLevel {
		return
	}
	output("WARN", fmt.Sprintf(msg, args...))
}

// Info logs a message to the 'standard' Logger if MaxLogLevel is >= INFO
func Info(args ...interface{}) {
	if INFO > maxLogLevel {
		return
	}
	output("INFO", fmt.Sprint(args...))
}

func Infof(msg string, args ...interface{}) {
	if INFO > maxLogLevel {
		return
	}
	output("INFO", fmt.Sprintf(msg, args...))
}

// Trace logs a message to the 'standard' Logger if MaxLogLevel is >= DEBUG
func Debug(args ...interface{}) {
	if DEBUG > maxLogLevel {
		return
	}
	output("DEBUG", fmt.Sprint(args...))
}

func Debugf(msg string, args ...interface{}) {
	if DEBUG > maxLogLevel {
		return
	}
	output("DEBUG", fmt.Sprintf(msg, args...))
}

func Access(args ...interface{}) {
	if logAccess == nil {
		return
	}
	outputRequest(fmt.Sprint(args...))
}

func output(mode, msg string) {
	logout.Output(3, "["+mode+"] "+msg)
}

func outputRequest(msg string) {
	logAccess.access.Output(3, msg)
}
