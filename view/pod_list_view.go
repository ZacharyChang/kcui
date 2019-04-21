package view

import (
	"fmt"
	"sort"

	"github.com/ZacharyChang/kcui/k8s"
	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/ZacharyChang/kcui/pkg/option"
	"github.com/rivo/tview"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	informerv1 "k8s.io/client-go/informers/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type PodListView struct {
	k8s.Client
	Content       *tview.List
	SelectedIndex int
	PodMap        map[string]*v1.Pod
	factory       informers.SharedInformerFactory
	podLister     listerv1.PodLister
	podInformer   informerv1.PodInformer
}

func NewPodListView(opts *option.Options) *PodListView {
	client := k8s.NewClient(opts)
	view := PodListView{
		Content: tview.NewList(),
		Client:  *client,
		PodMap:  make(map[string]*v1.Pod, 0),
		factory: client.Factory,
	}
	view.Content.ShowSecondaryText(false)
	view.Content.SetBorder(true).SetTitle(" Pod ")
	view.podLister = view.factory.Core().V1().Pods().Lister()
	view.podInformer = view.factory.Core().V1().Pods()
	view.podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    view.addPod,
			UpdateFunc: view.updatePod,
			DeleteFunc: view.deletePod,
		},
	)
	return &view
}

func (view *PodListView) addPod(obj interface{}) {
	pod := obj.(*v1.Pod)
	view.PodMap[pod.Name] = pod
	view.Reload()
	log.Infof("Pod Created: %s", pod.Name)
}

func (view *PodListView) updatePod(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	view.PodMap[oldPod.Name] = newPod
	view.Reload()
	log.Infof("Pod Updated: %s -> %s", oldPod.Name, newPod.Name)
}

func (view *PodListView) deletePod(obj interface{}) {
	pod := obj.(*v1.Pod)
	delete(view.PodMap, pod.Name)
	view.Reload()
	log.Infof("Pod Deleted: %s", pod.Name)
}

func (view *PodListView) SetClient(client k8s.Client) *PodListView {
	view.Client = client
	return view
}

func (view *PodListView) Reload() {
	keys := make([]string, 0)
	for k := range view.PodMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	view.Content.Clear()
	for _, key := range keys {
		pod := view.PodMap[key]
		phase := pod.Status.Phase
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			phase = "Terminating"
		}
		log.Debug(view.PodMap[key].Status)
		view.Content.AddItem(fmt.Sprintf("%s:%s", phase, pod.Name), "", 'a', nil)
	}
}

func (view *PodListView) Run(stopCh chan struct{}) {
	view.factory.Start(stopCh)
}
