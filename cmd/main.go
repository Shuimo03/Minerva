package main

import (
	k8sclient "backend/config/kubeconfig"
	"backend/kubernetes"
	"backend/sink/logs"
	"flag"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
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
