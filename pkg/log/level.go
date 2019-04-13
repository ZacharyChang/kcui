package log

type logLevel struct {
	level  int
	prefix string
}

var (
	FatalLevel = logLevel{1, " [FATAL] "}
	ErrorLevel = logLevel{2, " [ERROR] "}
	WarnLevel  = logLevel{3, " [WARN]  "}
	InfoLevel  = logLevel{4, " [INFO]  "}
	DebugLevel = logLevel{5, " [DEBUG] "}
)
