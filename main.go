package main

import (
	"bufio"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type LogView struct {
	sync.RWMutex
	view *tview.TextView
}

func NewLogView() *LogView {
	return &LogView{
		view: tview.NewTextView(),
	}
}

func main() {
	app := tview.NewApplication()
	listView := tview.NewList().ShowSecondaryText(false)
	// skip the column name
	for i, podName := range getPodNames() {
		if i > 0 {
			listView.AddItem(podName, "", rune('a'+i-1), nil)
		}
	}
	listView.SetBorder(true).SetTitle("Pod")
	// log := NewLogView()
	logView := tview.NewTextView()
	logView.SetBorder(true).SetTitle("Log")
	logText, _ := listView.GetItemText(listView.GetCurrentItem())
	writePodLogs(logView, logText)
	// listView.SetChangedFunc()
	flex := tview.NewFlex().
		AddItem(listView, 0, 1, false).
		AddItem(logView, 0, 3, false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		if event.Key() == tcell.KeyUp {
			newItem := listView.GetCurrentItem() - 1
			if newItem == -1 {
				newItem = listView.GetItemCount() - 1
			}
			podName, _ := listView.GetItemText(newItem)
			listView.SetCurrentItem(newItem)
			go writePodLogs(logView, podName)
		}
		if event.Key() == tcell.KeyDown {
			newItem := listView.GetCurrentItem() + 1
			if newItem == listView.GetItemCount() {
				newItem = 0
			}
			podName, _ := listView.GetItemText(newItem)
			listView.SetCurrentItem(newItem)
			go writePodLogs(logView, podName)
		}
		return event
	})
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func getPodNames() []string {
	cmd := "kubectl get pods |awk '{print $1}'"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run('%s') fail wtih %s\n", cmd, err)
	}
	pods, err := StringToLines(string(out))
	if err != nil {
		log.Fatalf("getPodNames() fail: %s\n", err)
	}
	return pods
}

func writePodLogs(target *tview.TextView, podName string) {
	_, _, _, height := target.GetRect()
	cmd := "kubectl logs " + podName + " --tail " + strconv.Itoa(height)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run('%s') fail wtih %s\n", cmd, err)
	}

	target.SetText(string(out))
}

func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}
