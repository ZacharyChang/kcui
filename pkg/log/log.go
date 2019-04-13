package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	DefaultCallerDepth = 2
)

var (
	currentLevel = InfoLevel
	logger       = NewLogger("root").AddHandler(NewFileHandler("test.log"))
)

// SetLogLevel set current log level, only the higher level log will be written to the file
func SetLogLevel(level logLevel) {
	currentLevel = level
}

func currentTime() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z07:00")
}

func log(level logLevel, v ...interface{}) {
	if currentLevel.level >= level.level {
		// runtime info
		file := "??"
		line := 0
		if _, _file, _line, ok := runtime.Caller(DefaultCallerDepth); ok {
			index := strings.LastIndex(_file, "/") + 1
			if index >= 0 {
				file, line = _file[index:], _line
			} else {
				file, line = _file, _line
			}
		}
		logger.Record(currentTime(), level.prefix, fmt.Sprintf("%s:%d\t| ", file, line), fmt.Sprint(v...), "\n")
	}
}

func logf(level logLevel, format string, v ...interface{}) {
	if currentLevel.level >= level.level {
		// runtime info
		file := "??"
		line := 0
		if _, _file, _line, ok := runtime.Caller(2); ok {
			index := strings.LastIndex(_file, "/") + 1
			if index >= 0 {
				file, line = _file[index:], _line
			} else {
				file, line = _file, _line
			}
		}
		logger.Record(currentTime(), level.prefix, fmt.Sprintf("%s:%d\t| ", file, line), fmt.Sprintf(format, v...), "\n")
	}
}

// Fatal print logs like log.Print, but has a prefix with "[FATAL]"
func Fatal(v ...interface{}) {
	log(FatalLevel, v...)
	os.Exit(1)
}

// Fatalf print logs like log.Printf, but has a prefix with "[FATAL]"
func Fatalf(format string, v ...interface{}) {
	logf(FatalLevel, format, v...)
	os.Exit(1)
}

// Error print logs like log.Print, but has a prefix with "[ERROR]"
func Error(v ...interface{}) {
	log(ErrorLevel, v...)
}

// Errorf print logs like log.Printf, but has a prefix with "[ERROR]"
func Errorf(format string, v ...interface{}) {
	logf(ErrorLevel, format, v...)
}

// Warn print logs like log.Print, but has a prefix with "[WARN]"
func Warn(v ...interface{}) {
	log(WarnLevel, v...)
}

// Warnf print logs like log.Printf, but has a prefix with "[WARN]"
func Warnf(foramt string, v ...interface{}) {
	logf(WarnLevel, foramt, v...)
}

// Info print logs like log.Print, but has a prefix with "[INFO]"
func Info(v ...interface{}) {
	log(InfoLevel, v...)
}

// Infof print logs like log.Printf, but has a prefix with "[INFO]"
func Infof(format string, v ...interface{}) {
	logf(InfoLevel, format, v...)
}

// Debug print logs like log.Print, but has a prefix with "[DEBUG]"
func Debug(v ...interface{}) {
	log(DebugLevel, v...)
}

// Debugf print logs like log.Printf, but has a prefix with "[DEBUG]"
func Debugf(format string, v ...interface{}) {
	logf(DebugLevel, format, v...)
}
