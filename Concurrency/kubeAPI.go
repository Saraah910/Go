package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func mainFunc() {
// 	var wg sync.WaitGroup
// 	podCount := make(chan int)

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		launchConsole()
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		pods := getAllPods("default")
// 		podCount <- pods
// 		close(podCount)
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		count := <-podCount
// 		var iPods sync.WaitGroup
// 		for i := count - 5; i <= count; i++ {
// 			iPods.Add(1)
// 			go func(podID int) {
// 				defer iPods.Done()
// 				restartPod(podID)
// 			}(i)
// 		}
// 		iPods.Wait()
// 	}()

// 	wg.Wait()
// }

// // func launchConsole() {
// // 	fmt.Println("Console launching...")
// // 	time.Sleep(time.Millisecond * 500)
// // 	fmt.Println("Console launched.")
// // }

// func getAllPods(namespace string) int {
// 	time.Sleep(time.Second * 2)
// 	fmt.Printf("Pod count in %v : %v\n", namespace, 10)
// 	return 10
// }

// func restartPod(id int) {
// 	time.Sleep(time.Millisecond * 500)
// 	fmt.Printf("Started pod with id: %v\n", id)
// }
