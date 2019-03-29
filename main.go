package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/ZacharyChang/kcui/k8s"
	"github.com/ZacharyChang/kcui/log"
	"github.com/ZacharyChang/kcui/view"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	namespace  = flag.String("namespace", "default", "(optional) k8s namespace")
	kubeclient *kubernetes.Clientset
)

func main() {
	flag.Parse()

	log.Debugf("flag namespace: %s", *namespace)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	kubeclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	app := tview.NewApplication()
	listView := tview.NewList().ShowSecondaryText(false)
	listView.SetBorder(true).SetTitle(" Pod ")
	logView := view.NewLogView()
	logView.Content.SetBorder(true).SetTitle(" Log ")
	for i, podName := range getPodNames() {
		listView.AddItem(podName, "", rune('A'+i), nil)
	}
	podName, _ := listView.GetItemText(listView.GetCurrentItem())

	logView.PodName = podName
	logView.SetHandler(k8s.NewClient().SetNamespace(*namespace).PodLogHandler)
	logView.Refresh(func() {
		app.Draw()
	})
	listView.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		logView.PodName = mainText
		logView.Refresh(func() {
			app.Draw()
		})
	})
	flex := tview.NewFlex().
		AddItem(listView, 0, 1, true).
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

func getPodNames() (names []string) {
	log.Debug("getPodNames() called")
	pods, err := kubeclient.CoreV1().Pods(*namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range pods.Items {
		names = append(names, v.ObjectMeta.Name)
	}
	log.Debugf("got pods: [ %s ]", strings.Join(names, " "))
	return
}
