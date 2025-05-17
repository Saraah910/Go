package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/k8s-client/getters"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube/config")
	fmt.Println(kubeconfig)

	//Fetches current context
	currentContext, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	client := kubernetes.NewForConfigOrDie(currentContext)
	pods := getters.GetPods("default", client)
	for _, pod := range pods {
		fmt.Println(pod)
	}
	newPodCreateTime, err := getters.CreatePod("default", client)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(newPodCreateTime)
	pods = getters.GetPods("default", client)
	for _, pod := range pods {
		fmt.Println(pod)
	}

}
