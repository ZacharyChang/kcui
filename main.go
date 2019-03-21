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
	listView.SetBorder(true).SetTitle(" Pod ")
	logView := tview.NewTextView()
	logView.SetBorder(true).SetTitle(" Log ")
	for i, podName := range getPodNames() {
		// skip the column name
		if i > 0 {
			listView.AddItem(podName, "", rune('a'+i-1), func() {
				writePodLogs(logView, podName)
			})
		}
	}
	podName, _ := listView.GetItemText(listView.GetCurrentItem())
	go writePodLogs(logView, podName)
	listView.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		writePodLogs(logView, mainText)
	})
	flex := tview.NewFlex().
		AddItem(listView, 0, 1, true).
		AddItem(logView, 0, 3, false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
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
