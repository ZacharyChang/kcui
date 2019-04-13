package view

import (
	"io"
	"sync"

	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/rivo/tview"
)

type LogView struct {
	sync.Mutex
	Content *tview.TextView
	PodName string
	handler func(string, io.Writer, func(), <-chan struct{})
	stopCh  chan struct{}
}

func NewLogView() *LogView {
	return &LogView{
		Content: tview.NewTextView().SetDynamicColors(true),
		stopCh:  make(chan struct{}),
	}
}

func (logView *LogView) SetHandler(handler func(string, io.Writer, func(), <-chan struct{})) {
	log.Debug("SetHandler called")
	logView.handler = handler
}

func (logView *LogView) Stop() {
	log.Debug("Stop called")
	logView.Lock()
	log.Debug(logView.stopCh)
	close(logView.stopCh)
	logView.Unlock()
}

func (logView *LogView) Refresh(callback func()) {
	log.Debug("Refresh called")

	logView.Content.Clear()

	go func() {
		logView.Stop()
		logView.stopCh = make(chan struct{})
		logView.handler(logView.PodName, tview.ANSIWriter(logView.Content), callback, logView.stopCh)
	}()
}
