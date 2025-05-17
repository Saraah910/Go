package main

import (
	"fmt"
	"time"
)

type ClusterOps struct {
	ClusterName string
	Msg         string
}

func main() {
	clusters := []string{"thanosajksjhfk", "snowhitesjdfhskdfheig", "leet", "an", "saurabhsakshi"}
	done := make(chan ClusterOps, len(clusters))
	// done1 := make(chan string)

	for _, clusterName := range clusters {
		go launchConsole(clusterName, done)
	}
	launchedConsoles := 0
	for {
		select {
		case msg := <-done:
			fmt.Println(msg.Msg)
			go func(clusterName string) {
				op1 := make(chan ClusterOps)
				op2 := make(chan ClusterOps)
				getPods(msg.ClusterName, op1)
				getNamespaces(msg.ClusterName, op2)

				msg1 := <-op1
				msg2 := <-op2
				fmt.Println(msg1, msg2)
				// for i := 0; i < 2; i++ {
				// 	opsMsg := <-ops
				// 	fmt.Println(opsMsg.Msg)
				// }
			}(msg.ClusterName)

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

func launchConsole(clusterName string, done chan<- ClusterOps) {
	fmt.Printf("⚙️ Launching the console for cluster %v\n", clusterName)
	time.Sleep(time.Duration(1+len(clusterName)) * time.Second)
	data := ClusterOps{ClusterName: clusterName, Msg: fmt.Sprintf("✅ Launched console for cluster %v\n", clusterName)}
	done <- data

}

func getPods(clusterName string, done chan<- ClusterOps) {
	time.Sleep(time.Duration(1+len(clusterName)%2) * time.Millisecond)
	data := ClusterOps{Msg: fmt.Sprintf("The pod count is: %v for cluster: %v", len(clusterName), clusterName)}
	done <- data
}

func getNamespaces(clusterName string, done chan<- ClusterOps) {
	time.Sleep(time.Duration(1+len(clusterName)%5) * time.Millisecond)
	data := ClusterOps{Msg: fmt.Sprintf("The namespace count is: %v, for the cluster: %v", len(clusterName)*2, clusterName)}
	done <- data
}
