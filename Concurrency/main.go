package main

import (
	"fmt"
	"time"
)

func main() {
	clusters := []string{"thanosajksjhfk", "snowhitesjdfhskdfheig", "leet", "an", "saurabhsakshi"}
	done := make(chan string, len(clusters))

	for _, clusterName := range clusters {
		go launchConsole(clusterName, done)
	}
	launchedConsoles := 0
	for {
		select {
		case msg := <-done:
			fmt.Println(msg)
			launchedConsoles++
			if launchedConsoles == len(clusters) {
				fmt.Println("All consoles launched successfully.")
			}
		default:
			fmt.Println("🛠️ Main function is running..")
			time.Sleep(500 * time.Millisecond)
		}
	}

}

func launchConsole(clusterName string, done chan<- string) {
	fmt.Printf("⚙️ Launching the console for cluster %v\n", clusterName)
	time.Sleep(time.Duration(1+len(clusterName)%3) * time.Second)

	done <- fmt.Sprintf("✅ Launched console for cluster %v", clusterName)

}
