package log

import (
	"io"
	"os"
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
