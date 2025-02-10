package main

import (
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // 导入v1包
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func main() {
	kubeconfig := "/data/dev_tools/K8sPod/kubeconfig.yaml" // 替换为您的 kubeconfig 路径

	// 加载配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载 kubeconfig 时出错: %v\n", err)
		os.Exit(1)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建 Kubernetes 客户端时出错: %v\n", err)
		os.Exit(1)
	}

	// Pod 和容器名称
	podName := "ubuntu-pod"
	containerName := "ubuntu"
	namespace := "default" // 如果有不同命名空间，请修改为相应的命名空间
	command := []string{"/bin/sh", "-c", "echo Hello from Pod!"}

	// 执行命令
	err = execCommand(clientset, config, namespace, podName, containerName, command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "执行命令时出错: %v\n", err)
		os.Exit(1)
	}
}

// execCommand 在指定的 Pod 中执行命令
func execCommand(clientset *kubernetes.Clientset, config *rest.Config, namespace, podName, containerName string, command []string) error {
	req := clientset.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec")

	req.VersionedParams(&metav1.PodExecOptions{
		Command:   command,
		Container: containerName,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		Tty:       false,
	}, metav1.ParameterCodec)

	// 创建执行器
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("创建执行器时出错: %w", err)
	}

	// 执行命令并流式传输输出
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		return fmt.Errorf("执行命令时出错: %w", err)
	}

	return nil
}
