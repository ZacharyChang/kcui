package log

import (
	"fmt"
	"sync"
)

type Logger struct {
	sync.Mutex
	Name     string
	handlers []Handler
}

func NewLogger(name string) *Logger {
	return &Logger{
		Name: name,
	}
}

func (logger *Logger) AddHandler(handler Handler) *Logger {
	logger.handlers = append(logger.handlers, handler)
	return logger
}

func (logger *Logger) Record(record ...interface{}) error {
	logger.Lock()
	for _, handler := range logger.handlers {
		for _, v := range record {
			_, err := fmt.Fprint(handler.Writer(), v)
			if err != nil {
				return err
			}
		}
	}
	logger.Unlock()
	return nil
}

func (logger *Logger) Close() error {
	for _, v := range logger.handlers {
		err := v.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
