package test

import (
	k8sclient "backend/pkg/client/kubernetes"
	"backend/pkg/kubernetes"
	"flag"
	"k8s.io/apimachinery/pkg/util/wait"
	"testing"
)

func TestNewKubernetesGetEvent(t *testing.T) {
	flag.Parse()
	client, err := k8sclient.KubeClient(kubeconfig)
	if err != nil {
		t.Fatalf("Create Client Error: %v", err)
	}
	kes := kubernetes.NewKubernetesEventSource(client)
	kes.Run(wait.NeverStop)
	_, eventError := kes.GetEvent()
	if eventError != nil {
		t.Fatalf("Failed to Event:%s", eventError)
	}
}
