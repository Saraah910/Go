package Concurrency

import (
	"fmt"
	"time"
)

func launchConsole() {
	fmt.Println("Launching console")
	time.Sleep(time.Millisecond * 500)
	fmt.Print("Completed console launch")
}

func getAllPods() []Pod {
	fmt.Println("Getting pods in default namespace")
	time.Sleep(time.Millisecond * 500)
	var podsList []Pod
	pod := Pod{
		Name:     "new",
		Ready:    "1/1",
		Status:   "Running",
		Restarts: 0,
		Age:      "2h",
	}
	podsList = append(podsList, pod)
	return podsList
}

func getAllNamespaces() {
	fmt.Println("Getting all namespaces")
	time.Sleep(time.Millisecond * 500)
}

func getAllDeploy() {
	fmt.Println("Getting all deployments in default namespace.")
	time.Sleep(time.Millisecond * 500)

}
