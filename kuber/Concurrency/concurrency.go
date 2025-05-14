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

	wg.Add(1)
	go func() {
		defer wg.Done()
		launchConsole()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		podsList := getAllPods()
		pods <- podsList
		close(pods)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		getAllNamespaces()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		getAllDeploy()
	}()

	wg.Wait()

}
