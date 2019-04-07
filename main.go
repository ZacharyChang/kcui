package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ZacharyChang/kcui/k8s"
	"github.com/ZacharyChang/kcui/log"
	"github.com/ZacharyChang/kcui/view"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "kcui"
	app.Compiled = time.Now()
	app.Usage = "k8s log tail tool"

	opts := NewOptions()
	opts.AddFlags(app)

	app.Action = func(c *cli.Context) error {
		if err := startView(opts); err != nil {
			return cli.NewExitError(fmt.Sprintf("application failed to start: %v\n", err), 1)
		}
		log.Info("application started...")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Errorf("fail to run: %s", err.Error())
	}

}

func startView(opts *Options) error {
	client := k8s.NewClient(opts.Kubeconfig).SetNamespace(opts.Namespace)
	if opts.Debug {
		log.SetLogLevel(log.DebugLevel)
	}

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
	if err := app.SetRoot(flex, true).Run(); err != nil {
		return err
	}
	return nil
}
