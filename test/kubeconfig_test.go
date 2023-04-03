package test

import (
	k8sclient "backend/config/kubeconfig"
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
	flag.StringVar(&kubeconfig, "kubeconfig", config, "Path to a kubeconfig. Only required if out-of-cluster.")

}

func TestKubeClient(t *testing.T) {
	flag.Parse()
	_, err := k8sclient.KubeClient(kubeconfig)
	if err != nil {
		t.Fatalf("Create Client Error: %v", err)
	}
}
