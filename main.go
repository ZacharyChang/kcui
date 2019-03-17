package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	list := tview.NewList().ShowSecondaryText(false).
		AddItem("List item 1", "", '1', nil).
		AddItem("List item 2", "", '2', nil).
		AddItem("List item 3", "", '3', nil).
		AddItem("List item 4", "", '4', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	list.SetBorder(true).SetTitle("Pod")
	flex := tview.NewFlex().
		AddItem(list, 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Log"), 0, 3, false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		if event.Key() == tcell.KeyUp {
			list.SetCurrentItem(list.GetCurrentItem() - 1)
		}
		if event.Key() == tcell.KeyDown {
			newItem := list.GetCurrentItem() + 1
			if newItem < list.GetItemCount() {
				list.SetCurrentItem(newItem)
			} else {
				list.SetCurrentItem(0)
			}
		}
		return event
	})
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
