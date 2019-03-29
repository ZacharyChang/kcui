package log

import (
	"fmt"
	golog "log"
	"os"
	"time"
)

const (
	fatal = " [FATAL] "
	error = " [ERROR] "
	warn  = " [WARN] "
	info  = " [INFO] "
	debug = " [DEBUG] "
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

func currentTime() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z07:00")
}

func log(level string, v ...interface{}) {
	golog.Print(currentTime(), level, fmt.Sprint(v...))
}

func logf(level string, format string, v ...interface{}) {
	golog.Print(currentTime(), level, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	log(fatal, v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	logf(fatal, format, v...)
	os.Exit(1)
}

func Error(v ...interface{}) {
	log(error, v...)
}

func Errorf(format string, v ...interface{}) {
	logf(error, format, v...)
}

func Warn(v ...interface{}) {
	log(warn, v...)
}

func Warnf(foramt string, v ...interface{}) {
	logf(warn, foramt, v...)
}

func Info(v ...interface{}) {
	log(info, v...)
}

func Infof(format string, v ...interface{}) {
	logf(info, format, v...)
}

func Debug(v ...interface{}) {
	log(debug, v...)
}

func Debugf(format string, v ...interface{}) {
	logf(debug, format, v...)
}
