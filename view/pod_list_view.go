package view

import (
	"github.com/ZacharyChang/kcui/k8s"
	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/ZacharyChang/kcui/pkg/option"
	"github.com/rivo/tview"
	"strings"
)

type PodListView struct {
	k8s.Client
	Content       *tview.List
	SelectedIndex int
	PodList       []string
}

func NewPodListView(opts *option.Options) *PodListView {
	view := PodListView{
		Content: tview.NewList(),
		Client:  *k8s.NewClient(opts),
	}
	view.Content.ShowSecondaryText(false)
	view.Content.SetBorder(true).SetTitle(" Pod ")
	return &view
}

func (view *PodListView) SetClient(client k8s.Client) *PodListView {
	view.Client = client
	return view
}

func (view *PodListView) SetPodList(list []string) *PodListView {
	log.Debugf("Get pods: %s", strings.Join(list, ","))
	view.PodList = list
	view.SelectedIndex = 0
	return view
}

func (view *PodListView) Refresh() {
	view.Content.Clear()
	view.SetPodList(view.ListPods())
	for i, podName := range view.PodList {
		view.Content.AddItem(podName, "", rune('A'+i), nil)
	}
}
