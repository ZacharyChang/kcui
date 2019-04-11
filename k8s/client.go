package k8s

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ZacharyChang/kcui/log"
	"github.com/ZacharyChang/kcui/pkg/option"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	namespace  string
	kubeclient *kubernetes.Clientset
}

type Handler interface {
	Handle() *io.Writer
}

func NewClient(opts *option.Options) *Client {
	var config *rest.Config
	var err error
	if opts.ConfigValue != "" {
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(opts.ConfigValue))
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", opts.ConfigFile)
	}
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Client{
		namespace:  opts.Namespace,
		kubeclient: client,
	}
}

func (c *Client) SetNamespace(ns string) *Client {
	c.namespace = ns
	return c
}

func (c *Client) GetNamespace() string {
	return c.namespace
}

func (c *Client) PodLogHandler(podName string, w io.Writer, callback func()) io.ReadCloser {
	log.Debugf("client: %v", c)
	log.Debugf("namespace: %s", c.namespace)
	log.Debugf("getting log from %s:%s", c.namespace, podName)

	lines := int64(50)
	req := c.kubeclient.CoreV1().Pods(c.namespace).GetLogs(podName, &corev1.PodLogOptions{
		TailLines: &lines,
		Follow:    true,
	})
	defer callback()

	podLogs, err := req.Stream()
	if err != nil {
		log.Errorf("error: fail to open stream %s", err.Error())
		_, err = fmt.Fprintf(w, "error: fail to open stream %s\n", err.Error())
		return podLogs
	}
	log.Debug("begin read")

	go func() {
		reader := bufio.NewReader(podLogs)
		for {
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
			callback()
			time.Sleep(500)
		}
	}()

	return podLogs
}

func (c *Client) GetPodNames() (names []string) {
	log.Debug("getPodNames() called")
	pods, err := c.kubeclient.CoreV1().Pods(c.namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range pods.Items {
		names = append(names, v.ObjectMeta.Name)
	}
	log.Debugf("got pods: [ %s ]", strings.Join(names, " "))
	return
}
