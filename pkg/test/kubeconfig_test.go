package test

import (
	k8sclient "backend/pkg/client/kubernetes"
	"flag"
	"testing"
)

var (
	kubeconfig string
)

const (
	config = ""
)

func init() {
	flag.StringVar(&kubeconfig, "kubernetes", config, "Path to a kubernetes. Only required if out-of-cluster.")

}

func TestKubeClient(t *testing.T) {
	flag.Parse()
	_, err := k8sclient.KubeClient(kubeconfig)
	if err != nil {
		t.Fatalf("Create Client Error: %v", err)
	}
}
