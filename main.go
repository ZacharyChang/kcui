package main

import (
	"bytes"
	"flag"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"io"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "github.com/ZacharyChang/kcui/log"
)

var (
	kubeconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	namespace  = flag.String("namespace", "default", "(optional) kubernetes namespace")
	kubeclient *kubernetes.Clientset
)

func main() {
	flag.Parse()

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
	logView := tview.NewTextView()
	logView.SetBorder(true).SetTitle(" Log ")
	for i, podName := range getPodNames() {
		listView.AddItem(podName, "", rune('a'+i-1), nil)
	}
	podName, _ := listView.GetItemText(listView.GetCurrentItem())
	go writePodLogs(logView, podName, func() {
		app.Draw()
	})
	listView.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		go writePodLogs(logView, mainText, func() {
			app.Draw()
		})
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
	log.Print("application started...")
	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal("application failed to start")
		panic(err)
	}
}

func getPodNames() (names []string) {
	pods, err := kubeclient.CoreV1().Pods(*namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range pods.Items {
		names = append(names, v.ObjectMeta.Name)
	}
	return
}

func writePodLogs(target *tview.TextView, podName string, callback func()) {
	_, _, _, height := target.GetRect()
	h := int64(height)
	req := kubeclient.CoreV1().Pods(*namespace).GetLogs(podName, &corev1.PodLogOptions{
		TailLines: &h,
		//Follow:    true,
	})

	if callback != nil {
		defer callback()
	}

	podLogs, err := req.Stream()
	if err != nil {
		target.SetText("error: fail to open stream\n")
		return
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		target.SetText("error: fail to copy logs\n")
		return
	}
	target.SetText(buf.String())
}
