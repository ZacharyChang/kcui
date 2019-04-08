package view

import (
	"io"
	"sync"

	"github.com/ZacharyChang/kcui/log"
	"github.com/rivo/tview"
)

type LogView struct {
	sync.Mutex
	Content   *tview.TextView
	PodName   string
	handler   func(string, io.Writer, func()) io.ReadCloser
	closerMap map[string]io.ReadCloser
}

func NewLogView() *LogView {
	return &LogView{
		Content:   tview.NewTextView().SetDynamicColors(true),
		closerMap: make(map[string]io.ReadCloser, 0),
	}
}

func (logView *LogView) SetHandler(handler func(string, io.Writer, func()) io.ReadCloser) {
	log.Debug("SetHandler called")
	logView.handler = handler
}

func (logView *LogView) Stop() {
	logView.Lock()
	for k, closer := range logView.closerMap {
		if closer != nil {
			err := closer.Close()
			if err != nil {
				log.Errorf("Fail to close [%s] last reader: %s", k, err.Error())
				continue
			}
			log.Debugf("Reader [%s] closed", k)
		}
	}
	logView.Unlock()
}

func (logView *LogView) Refresh(callback func()) {
	log.Debug("Refresh called")

	logView.Content.Clear()
	defer logView.Stop()

	go func() {
		logView.closerMap[logView.PodName] = logView.handler(logView.PodName, tview.ANSIWriter(logView.Content), callback)
	}()
}
