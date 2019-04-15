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
	handler func(string, io.Writer, <-chan struct{})
	stopCh  chan struct{}
}

func NewLogView() *LogView {
	return &LogView{
		Content: tview.NewTextView().SetDynamicColors(true),
		stopCh:  make(chan struct{}),
	}
}

func (logView *LogView) SetHandler(handler func(string, io.Writer, <-chan struct{})) {
	log.Debug("SetHandler called")
	logView.handler = handler
}

func (logView *LogView) Start() {
	logView.Content.Clear()
	logView.stopCh = make(chan struct{})
	go logView.handler(logView.PodName, tview.ANSIWriter(logView.Content), logView.stopCh)
}

func (logView *LogView) Stop() {
	log.Debug("Stop called")
	log.Debug(logView.stopCh)
	close(logView.stopCh)
}

func (logView *LogView) Refresh() {
	logView.Stop()
	logView.Start()
}
