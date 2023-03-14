package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/events/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"time"
)

var (
	kubeConfig string
)

func init() {
	flag.StringVar(&kubeConfig, "kubeConfig", "/root/.kube/config", "Path to a kubeconfig") //获取默认的config
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)

	if err != nil {
		klog.Fatal("Error building kubeconfig: ", err)
	}
	//获取客户端
	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		klog.Fatal("Error building kubernetes clientset: ", err)
	}
	eventInFormerFactory := informers.NewSharedInformerFactory(client, time.Minute)
	stopChan := make(chan struct{})
	defer close(stopChan)

	infromerEvent := eventInFormerFactory.Events().V1().Events().Informer()
	addChan := make(chan v1.Event)
	deleteChan := make(chan v1.Event)
	infromerEvent.AddEventHandlerWithResyncPeriod(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			unstructObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
			if err != nil {
				klog.Fatal("Error:", err)
			}
			event := &v1.Event{}
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj, event)
			if err != nil {
				klog.Fatal("Event Error:", err)
			}
			addChan <- *event
		},
		UpdateFunc: func(oldObj, newObj interface{}) {

		},
		DeleteFunc: func(obj interface{}) {

		},
	}, 0)

	go func() {
		for {
			select {
			case event := <-addChan:
				str, err := json.Marshal(&event)
				if err != nil {
					klog.Fatal("Json Error", err)
				}
				fmt.Println("Event:", string(str))
				break
			case <-deleteChan:
				break
			}
		}
	}()
	infromerEvent.Run(stopChan)
}
