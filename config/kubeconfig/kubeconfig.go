package kubeconfig

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func getKubeConfig(path string) (*rest.Config, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("KubeConfig is NUll")
	}
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		klog.Fatal("Error building kubernetes Config: ", err)
	}
	return config, nil
}

func KubeClient(path string) (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(path)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal("Error building kubernetes clientset: ", err)
	}
	return client, err
}
