package main

import (
	"flag"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"path/filepath"

	"github.com/ZacharyChang/kcui/k8s"
	"github.com/ZacharyChang/kcui/log"
	"github.com/ZacharyChang/kcui/view"
)

var (
	kubeconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	namespace  = flag.String("namespace", "default", "(optional) k8s namespace")
)

func main() {
	flag.Parse()

	client := k8s.NewClient().SetNamespace(*namespace)

	app := tview.NewApplication()
	podListView := view.NewPodListView()
	podListView.Content.ShowSecondaryText(false)
	podListView.Content.SetBorder(true).SetTitle(" Pod ")
	podListView.SetPodList(client.GetPodNames())

	logView := view.NewLogView()
	logView.Content.SetBorder(true).SetTitle(" Log ")
	podListView.Refresh(func() {
		app.Draw()
	})
	podName, _ := podListView.Content.GetItemText(podListView.Content.GetCurrentItem())

	logView.PodName = podName
	logView.SetHandler(client.PodLogHandler)
	logView.Refresh(func() {
		app.Draw()
	})
	podListView.Content.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		logView.PodName = mainText
		logView.Refresh(func() {
			app.Draw()
		})
	})
	flex := tview.NewFlex().
		AddItem(podListView.Content, 0, 1, true).
		AddItem(logView.Content, 0, 3, false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		return event
	})
	log.Info("application started...")
	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal("application failed to start")
		panic(err)
	}
}
