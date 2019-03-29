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
	handler func(string, io.Writer, func())
}

func NewLogView() *LogView {
	return &LogView{
		Content: tview.NewTextView(),
	}
}

func (logView *LogView) SetHandler(handler func(string, io.Writer, func())) {
	log.Debug("SetHandler called")
	logView.handler = handler
}

func (logView *LogView) Refresh(callback func()) {
	log.Debug("Refresh called")
	logView.Content.Clear()
	defer callback()
	go logView.handler(logView.PodName, logView.Content, callback)
}
