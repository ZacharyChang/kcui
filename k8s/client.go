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

var instance *Client

type Client struct {
	Namespace  string
	KubeClient *kubernetes.Clientset
	KubeConfig clientcmd.ClientConfig
	Factory    informers.SharedInformerFactory
	podLister  listerv1.PodLister
}

// NewClient returns a singleton instance of k8s client
func NewClient(opts *option.Options) *Client {
	if instance != nil {
		return instance
	}
	config, err := clientcmd.BuildConfigFromFlags("", opts.Kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// if namespace is not set from command, read from config file
	ns := opts.Namespace
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	if ns == "" {
		ns, _, _ = kubeconfig.Namespace()
	}

	factoryOpts := informers.WithNamespace(ns)
	f := informers.NewSharedInformerFactoryWithOptions(client, 0, factoryOpts)

	instance = &Client{
		Namespace:  ns,
		KubeClient: client,
		KubeConfig: kubeconfig,
		Factory:    f,
		podLister:  f.Core().V1().Pods().Lister(),
	}
	return instance
}

func (c *Client) TailPodLog(podName string, w io.Writer, stopCh <-chan struct{}) {
	log.Debugf("client: %v", c)
	log.Debugf("namespace: %s", c.Namespace)
	log.Debugf("getting log from %s:%s", c.Namespace, podName)

	lines := int64(50)
	req := c.KubeClient.CoreV1().Pods(c.Namespace).GetLogs(podName, &corev1.PodLogOptions{
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
