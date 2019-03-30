package k8s

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZacharyChang/kcui/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	namespace  string
	kubeclient *kubernetes.Clientset
}

type Handler interface {
	Handle() *io.Writer
}

func NewClient() *Client {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Client{
		namespace:  "default",
		kubeclient: client,
	}
}

func (c *Client) SetNamespace(ns string) *Client {
	c.namespace = ns
	return c
}

func (c *Client) PodLogHandler(podName string, w io.Writer, callback func()) {
	log.Debugf("client: %v", c)
	log.Debugf("namespace: %s", c.namespace)
	lines := int64(100)
	req := c.kubeclient.CoreV1().Pods(c.namespace).GetLogs(podName, &corev1.PodLogOptions{
		TailLines: &lines,
		Follow:    true,
	})
	defer callback()

	podLogs, err := req.Stream()
	if err != nil {
		log.Errorf("error: fail to open stream %s", err.Error())
		_, err = fmt.Fprintf(w, "error: fail to open stream %s\n", err.Error())
		return
	}
	log.Debug("begin read")

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
	}

	defer podLogs.Close()

	log.Infof("stream finished: %s", podName)
	return
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