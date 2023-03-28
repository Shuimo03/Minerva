package main

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

func main() {
	// 创建 Kubernetes 配置和客户端
	config, err := clientcmd.BuildConfigFromFlags("", "kubeconfig")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 创建 informer 工厂
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	// 创建 informer
	eventInformer := informerFactory.Core().V1().Events().Informer()

	// 创建一个工作队列
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// 定义事件处理程序
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event, ok := obj.(*corev1.Event)
			if !ok {
				klog.Errorf("object is not an Event")
				return
			}
			klog.Infof("kubernetes '%s' added", event.Name)
			queue.Add(event)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			event, ok := newObj.(*corev1.Event)
			if !ok {
				klog.Errorf("object is not an Event")
				return
			}
			klog.Infof("kubernetes '%s' updated", event.Name)
			queue.Add(event)
		},
		DeleteFunc: func(obj interface{}) {
			event, ok := obj.(*corev1.Event)
			if !ok {
				klog.Errorf("object is not an Event")
				return
			}
			klog.Infof("kubernetes '%s' deleted", event.Name)
			queue.Add(event)
		},
	}

	// 注册事件处理程序
	eventInformer.AddEventHandler(eventHandler)

	// 启动 informer
	informerFactory.Start(wait.NeverStop)

	// 处理工作队列中的事件
	for {
		item, quit := queue.Get()
		if quit {
			break
		}
		event, ok := item.(*corev1.Event)
		if !ok {
			klog.Errorf("item is not an Event")
			queue.Forget(item)
			continue
		}
		fmt.Printf("Processing kubernetes '%s'\n", event.Name)
		queue.Forget(item)
	}
}
