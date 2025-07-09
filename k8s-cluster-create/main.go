package main

import (
	"fmt"
	"os"
	"os/exec"
)

func InstallPrerequisites() {
	installKind()
	InstallHelm()
	InstallKubectl()
}

func CreateKindCluster(clusterName string) {
	fmt.Println("Creating Kind cluster...")
	cmd := exec.Command("kind", "create", "cluster", "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to create Kind cluster: %v\n", err)
		return
	}
	fmt.Println("Kind cluster created successfully.")
}

func main() {
	InstallPrerequisites()
	fmt.Println("All prerequisites installed successfully.")

	var clusterName string
	fmt.Printf("Enter the name for controller cluster: ")
	fmt.Scan(&clusterName)
	CreateKindCluster(clusterName)
}
