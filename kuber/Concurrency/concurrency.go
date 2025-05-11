package Concurrency

import (
	"sync"
)

type Pod struct {
	Name     string
	Ready    string
	Status   string
	Restarts int64
	Age      string
}

func ConcurrentFunctions(action string) {

	var wg sync.WaitGroup

	pods := make(chan []Pod)

	if action == "launch-console" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			launchConsole()
		}()

	} else if action == "get-pods" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			podsList := getAllPods()
			pods <- podsList
			close(pods)
		}()
	} else if action == "get-namespaces" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getAllNamespaces()
		}()
	} else if action == "get-deploy" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getAllDeploy()
		}()
	}
	wg.Wait()

}
