package main

import (
	k8sclient "backend/pkg/client/kubernetes"
	"backend/pkg/kubernetes"
	"backend/pkg/sink/logs"
	"flag"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
)

var (
	kubeconfig string
)

const (
	config = "/Users/wucola/.kube/dev-rke.yaml"
)

func init() {
	flag.StringVar(&kubeconfig, "kubernetes", config, "Path to a kubernetes. e.g./.kube/config.")

}

func main() {
	flag.Parse()
	// 创建 Kubernetes 配置和客户端
	client, err := k8sclient.KubeClient(kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	es := kubernetes.NewKubernetesEventSource(client)
	es.Run(wait.NeverStop)
	log := logs.LogSinker{}
	for {
		event, err := es.GetEvent()
		if err != nil {
			klog.Error(err)
		}
		log.ExportEvent(event)
	}
}
