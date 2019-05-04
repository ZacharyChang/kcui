package view

import (
	"sync"

	"github.com/ZacharyChang/kcui/k8s"
	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/ZacharyChang/kcui/pkg/option"

	"github.com/rivo/tview"
)

type LogView struct {
	sync.Mutex
	k8s.Client
	Content *tview.TextView
	PodName string
	stopCh  chan struct{}
}

func NewLogView(opts *option.Options) *LogView {
	view := LogView{
		Client:  *k8s.NewClient(opts),
		Content: tview.NewTextView().SetDynamicColors(true),
		stopCh:  make(chan struct{}),
	}
	view.Content.SetBorder(true).SetTitle(" Log ")
	return &view
}

func (logView *LogView) SetClient(client k8s.Client) *LogView {
	log.Debug("SetHandler called")
	logView.Client = client
	return logView
}

func (logView *LogView) Start() {
	logView.Content.Clear()
	logView.stopCh = make(chan struct{})
	go logView.Client.TailPodLog(logView.PodName, tview.ANSIWriter(logView.Content), logView.stopCh)
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
