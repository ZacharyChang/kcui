package log

import (
	"fmt"
	golog "log"
	"os"
	"time"
)

type logLevel struct {
	level  int
	prefix string
}

var (
	FatalLevel = logLevel{1, " [FATAL] "}
	ErrorLevel = logLevel{2, " [ERROR] "}
	WarnLevel  = logLevel{3, " [WARN] "}
	InfoLevel  = logLevel{4, " [INFO] "}
	DebugLevel = logLevel{5, " [DEBUG] "}
)

var (
	currentLevel = InfoLevel
)

func init() {
	f, err := os.OpenFile("kcui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		golog.Fatalf("error opening file: %v", err)
	}
	//defer f.Close()

	golog.SetFlags(0)
	golog.SetOutput(f)
}

// SetLogLevel set current log level, only the higher level log will be written to the file
func SetLogLevel(level logLevel) {
	currentLevel = level
}

func currentTime() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z07:00")
}

func log(level logLevel, v ...interface{}) {
	if currentLevel.level >= level.level {
		golog.Print(currentTime(), level.prefix, fmt.Sprint(v...))
	}
}

func logf(level logLevel, format string, v ...interface{}) {
	if currentLevel.level >= level.level {
		golog.Print(currentTime(), level.prefix, fmt.Sprintf(format, v...))
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
