package kubeconfig

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/url"
)

func getKubeConfig(path *url.URL) (*rest.Config, error) {
	if path == nil {
		return nil, fmt.Errorf("KubeConfig is NUll")
	}
	return nil, nil
}

func KubeClient(path *url.URL) (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(path)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	return client, err
}
