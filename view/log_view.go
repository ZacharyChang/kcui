package view

import (
	"io"
	"sync"

	"github.com/ZacharyChang/kcui/log"
	"github.com/rivo/tview"
)

type LogView struct {
	sync.Mutex
	Content    *tview.TextView
	PodName    string
	handler    func(string, io.Writer, func()) io.ReadCloser
	closerList []io.ReadCloser
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
	logView.Lock()
	log.Debug("Refresh called")
	for _, closer := range logView.closerList {
		if closer != nil {
			err := closer.Close()
			if err != nil {
				log.Errorf("Fail to close last reader: %s", err.Error())
			}
			log.Debug("Last reader closed")
		}
	}

	logView.Content.Clear()
	go func() {
		logView.closerList = append(logView.closerList, logView.handler(logView.PodName, logView.Content, callback))
	}()
	logView.Unlock()
}
