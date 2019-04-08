package k8s

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if client == nil {
		t.Fatalf("Error in createing kubernetes client using NewClient()")
	}
}

func TestSetNamespace(t *testing.T) {
	client := NewClient(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	client.SetNamespace("default")
	if client.GetNamespace() != "default" {
		t.Fatalf("Error in setting namespace for kubernetes client")
	}
}
