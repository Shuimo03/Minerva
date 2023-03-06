package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	kubeConfig string
)

func init() {
	flag.StringVar(&kubeConfig, "kubeConfig", "", "Path to a kubeconfig")
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
	deploy, err := client.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})
	for _, v := range deploy.Items {
		fmt.Println("deployName:", v.Name)
	}
}
