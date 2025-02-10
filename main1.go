package main

import (
	"context"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 指定具体的 kubeconfig 文件路径
	kubeconfig := "/data/dev_tools/K8sPod/kubeconfig.yaml" // 这里替换成你的 kubeconfig 文件路径 

	// 加载指定的 kubeconfig 文件
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to build kubeconfig: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	for {
		// 获取所有命名空间中的 Pod 列表
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Failed to list pods: %v", err)
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// 示例：查找特定命名空间中的 Pod 是否存在
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "nginx-pod", metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Pod example-xxxxx not found or error occurred: %v\n", err)
		} else {
			fmt.Printf("Found example-xxxxx pod in default namespace\n")
		}

		time.Sleep(10 * time.Second)
	}
}

