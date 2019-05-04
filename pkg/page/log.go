package page

import (
	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/ZacharyChang/kcui/pkg/option"
	"github.com/ZacharyChang/kcui/pkg/view"
	"github.com/rivo/tview"
	"strings"
)

// Log returns the page that display the pod list and logs
func Log(opts *option.Options, stopCh chan struct{}) (title string, content tview.Primitive) {
	if opts.Debug {
		log.SetLogLevel(log.DebugLevel)
	}

	podListView := view.NewPodListView(opts)

	logView := view.NewLogView(opts)
	podListView.Run(stopCh)

	podListView.Content.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		logView.PodName = strings.Split(mainText, ":")[1]
		logView.Refresh()
	})
	flex := tview.NewFlex().
		AddItem(podListView.Content, 0, 1, true).
		AddItem(logView.Content, 0, 3, false)

	return "Log", flex
}
