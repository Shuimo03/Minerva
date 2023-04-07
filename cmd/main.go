package main

import (
	"backend/pkg/client/kafka"
	k8sclient "backend/pkg/client/kubernetes"
	"backend/pkg/kubernetes"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	"net/url"
)

var (
	kubeconfig string
)

//TODO 当具有多个集群的时候,就需要开启Kafka消息队列做缓存

const (
	config   = ""
	kafkaUrl = ""
)

var (
	sinks = ""
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
	uri, err := url.Parse(kafkaUrl)
	if err != nil {
		panic(err)
	}
	kafkaCLient, clientError := kafka.NewKafkaClient(uri)
	if clientError != nil {
		panic(clientError)
	}

	es.Run(wait.NeverStop)
	//log := logs.LogSinker{}
	for {
		event, err := es.GetEvent()
		if err != nil {
			klog.Error(err)
		}
		//log.ExportEvent(event)
		fmt.Println("Kafka sinks")
		isSucess, err := kafkaCLient.ProducerMessage(event)
		if err != nil {
			panic(err)
		}
		fmt.Println(isSucess)
	}

}
