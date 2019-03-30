package view

import (
	"io"
	"sync"

	"github.com/ZacharyChang/kcui/log"
	"github.com/rivo/tview"
)

type LogView struct {
	sync.Mutex
	Content *tview.TextView
	PodName string
	handler func(string, io.Writer, func()) io.ReadCloser
	closer  io.ReadCloser
}

func NewLogView() *LogView {
	return &LogView{
		Content: tview.NewTextView(),
	}
}

func (logView *LogView) SetHandler(handler func(string, io.Writer, func()) io.ReadCloser) {
	log.Debug("SetHandler called")
	logView.handler = handler
}

func (logView *LogView) Refresh(callback func()) {
	log.Debug("Refresh called")
	if logView.closer != nil {
		err := logView.closer.Close()
		if err != nil {
			log.Errorf("Fail to close last reader: %s", err.Error())
		}
		log.Debug("Last reader closed")
	}
	logView.Content.Clear()
	go func() {
		logView.closer = logView.handler(logView.PodName, logView.Content, callback)
	}()
}
