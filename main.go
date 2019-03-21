package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	kubeclient *kubernetes.Clientset
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

func getPodNames() (names []string) {
	pods, err := kubeclient.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range pods.Items {
		names = append(names, v.ObjectMeta.Name)
	}
	return
}

func writePodLogs(target *tview.TextView, podName string) {
	_, _, _, height := target.GetRect()
	cmd := "kubectl logs " + podName + " --all-containers --tail " + strconv.Itoa(height)
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
