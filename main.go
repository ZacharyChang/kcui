package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

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
	logView := tview.NewTextView()
	logView.SetBorder(true).SetTitle("Log")
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
			writePodLogs(logView, podName)
		}
		if event.Key() == tcell.KeyDown {
			newItem := listView.GetCurrentItem() + 1
			if newItem == listView.GetItemCount() {
				newItem = 0
			}
			podName, _ := listView.GetItemText(newItem)
			listView.SetCurrentItem(newItem)
			writePodLogs(logView, podName)
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
		log.Fatalf("cmd.Run() fail wtih %s\n", err)
	}
	pods, err := StringToLines(string(out))
	if err != nil {
		log.Fatalf("getPodNames() fail: %s\n", err)
	}
	return pods
}

func writePodLogs(target *tview.TextView, podName string) {
	cmd := "kubectl logs " + podName + " --all-containers --tail 100"

	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() fail wtih %s\n", err)
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
