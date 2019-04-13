package view

import (
	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/rivo/tview"
	"strings"
)

type PodListView struct {
	Content       *tview.List
	SelectedIndex int
	PodList       []string
}

func NewPodListView() *PodListView {
	return &PodListView{
		Content: tview.NewList(),
	}
}

func (view *PodListView) SetPodList(list []string) *PodListView {
	log.Debugf("Get pods: %s", strings.Join(list, ","))
	view.PodList = list
	view.SelectedIndex = 0
	return view
}

func (view *PodListView) Refresh(callback func()) {
	view.Content.Clear()
	for i, podName := range view.PodList {
		view.Content.AddItem(podName, "", rune('A'+i), nil)
	}
	defer callback()
}
