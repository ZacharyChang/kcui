package log

import (
	"io"
	"os"
	"time"
)

type Handler interface {
	Writer() io.Writer
	Close() error
	Flush() error
}

type FileHandler struct {
	filename string
	file     *os.File
}

func NewFileHandler(filename string) *FileHandler {
	return &FileHandler{
		filename: filename,
	}
}

func (handler *FileHandler) Writer() io.Writer {
	f, err := os.OpenFile(handler.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	handler.file = f
	return handler.file
}

func (handler *FileHandler) Close() error {
	return handler.file.Close()
}

func (handler *FileHandler) Flush() error {
	return handler.file.Sync()
}

type RotatingFileHandler struct {
	filename   string
	dateFormat string
	file       *os.File
}

func NewRotatingFileHandler(filename string) *RotatingFileHandler {
	return &RotatingFileHandler{
		filename:   filename,
		dateFormat: "2006-01-02-15",
	}
}

func (handler *RotatingFileHandler) SetDateFormat(f string) *RotatingFileHandler {
	handler.dateFormat = f
	return handler
}

func (handler *RotatingFileHandler) Writer() io.Writer {
	f, err := os.OpenFile(handler.filename+"."+time.Now().Format(handler.dateFormat), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	handler.file = f
	return handler.file
}

func (handler *RotatingFileHandler) Close() error {
	return handler.file.Close()
}

func (handler *RotatingFileHandler) Flush() error {
	return handler.file.Sync()
}

type StdoutHandler struct {
}

func NewStdoutHandler() *StdoutHandler {
	return &StdoutHandler{}
}

func (handler *StdoutHandler) Writer() io.Writer {
	return os.Stdout
}

func (handler *StdoutHandler) Close() error {
	return os.Stdout.Close()
}

func (handlere *StdoutHandler) Flush() error {
	return os.Stdout.Sync()
}
