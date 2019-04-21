package k8s

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/ZacharyChang/kcui/pkg/log"
	"github.com/ZacharyChang/kcui/pkg/option"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	namespace  string
	kubeclient *kubernetes.Clientset
	Factory    informers.SharedInformerFactory
	podLister  listerv1.PodLister
}

func NewClient(opts *option.Options) *Client {
	config, err := clientcmd.BuildConfigFromFlags("", opts.Kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factoryOpts := informers.WithNamespace(opts.Namespace)
	f := informers.NewSharedInformerFactoryWithOptions(client, 0, factoryOpts)

	return &Client{
		namespace:  opts.Namespace,
		kubeclient: client,
		Factory:    f,
		podLister:  f.Core().V1().Pods().Lister(),
	}
}

func (c *Client) SetNamespace(ns string) *Client {
	c.namespace = ns
	return c
}

func (c *Client) GetNamespace() string {
	return c.namespace
}

func (c *Client) TailPodLog(podName string, w io.Writer, stopCh <-chan struct{}) {
	log.Debugf("client: %v", c)
	log.Debugf("namespace: %s", c.namespace)
	log.Debugf("getting log from %s:%s", c.namespace, podName)

	lines := int64(50)
	req := c.kubeclient.CoreV1().Pods(c.namespace).GetLogs(podName, &corev1.PodLogOptions{
		TailLines: &lines,
		Follow:    true,
	})

	podLogs, err := req.Stream()
	if err != nil {
		log.Errorf("error: fail to open stream %s", err.Error())
		_, err = fmt.Fprintf(w, "error: fail to open stream %s\n", err.Error())
		return
	}
	log.Debug("begin read")
	defer podLogs.Close()

	reader := bufio.NewReader(podLogs)
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("error: fail to read %s", err.Error())
			break
		}

		_, err = fmt.Fprint(w, line)
		if err != nil {
			log.Errorf("error: fail to output %s", err.Error())
			break
		}
		time.Sleep(500)
	}
	return
}
