package page

import (
	"strings"

	"github.com/ZacharyChang/kcui/pkg/option"
	"github.com/ZacharyChang/kcui/pkg/view"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type LogPage struct {
	PodListView  *view.PodListView
	LogView      *view.LogView
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
}

func NewLogPage(opts *option.Options) *LogPage {
	return &LogPage{
		PodListView: view.NewPodListView(opts),
		LogView:     view.NewLogView(opts),
	}
}

func (page *LogPage) SetInputCapture(event func(event *tcell.EventKey) *tcell.EventKey) *LogPage {
	page.inputCapture = event
	return page
}

// Show returns the page that display the pod list and logs
func (page *LogPage) Show(app *tview.Application, stopCh chan struct{}) (title string, content tview.Primitive) {

	page.PodListView.Run(stopCh)

	page.PodListView.Content.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		page.LogView.PodName = strings.Split(mainText, ":")[1]
		page.LogView.Refresh()
	})
	flex := tview.NewFlex().
		AddItem(page.PodListView.Content, 0, 1, true).
		AddItem(page.LogView.Content, 0, 3, false)
	page.PodListView.Content.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		app.SetFocus(page.LogView.Content)
	})
	page.LogView.Content.SetDoneFunc(func(key tcell.Key) {
		app.SetFocus(page.PodListView.Content)
	})
	return "Log", flex
}
