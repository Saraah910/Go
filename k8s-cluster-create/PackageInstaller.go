package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func AutoDetectOS() (string, string) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	if osName == "linux" && arch == "amd64" {
		return "linux", "amd64"
	} else if osName == "linux" && arch == "arm64" {
		return "linux", "arm64"
	} else if osName == "darwin" && arch == "amd64" {
		return "darwin", "amd64"
	} else if osName == "darwin" && arch == "arm64" {
		return "darwin", "arm64"
	}
	return fmt.Sprintf("Unsupported OS: %s/%s", osName, arch), ""
}

func installKind() {
	_, err := exec.LookPath("kind")
	if err == nil {
		fmt.Println("✅ Kind is already installed.")
		return
	}
	fmt.Println("Auto-detecting OS...")

	var url string
	var osArchBinMap = map[string]string{
		"linux/amd64":  "https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-amd64",
		"linux/arm64":  "https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-arm64",
		"darwin/amd64": "https://kind.sigs.k8s.io/dl/v0.22.0/kind-darwin-amd64",
		"darwin/arm64": "https://kind.sigs.k8s.io/dl/v0.22.0/kind-darwin-arm64",
	}
	osName, arch := AutoDetectOS()
	url = osArchBinMap[osName+"/"+arch]

	if url == "" {
		fmt.Printf("Unsupported architecture: %s/%s\n", osName, arch)
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory:", err)
		return
	}
	installDir := filepath.Join(homeDir, "go", "bin")
	output := filepath.Join(installDir, "kind")

	if err := os.MkdirAll(installDir, 0755); err != nil {
		fmt.Println("Failed to create install directory:", err)
		return
	}

	fmt.Println("Downloading kind binary from:", url)
	cmd := exec.Command("curl", "-Lo", output, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Download failed:", err)
		return
	}

	cmd = exec.Command("chmod", "+x", output)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to make binary executable:", err)
		return
	}

	fmt.Printf("\n✅ Kind installed to: ")
	addPathToShellConfig(installDir)
}

func addPathToShellConfig(pathToAdd string) {
	home, _ := os.UserHomeDir()
	shell := os.Getenv("SHELL")
	var shellConfig string

	if strings.Contains(shell, "zsh") {
		shellConfig = filepath.Join(home, ".zshrc")
	} else {
		shellConfig = filepath.Join(home, ".bashrc")
	}

	fmt.Printf("🔧 Adding %s to PATH in %s\n", pathToAdd, shellConfig)

	exportLine := fmt.Sprintf("export PATH=\"$PATH:%s\"", pathToAdd)

	file, err := os.Open(shellConfig)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), pathToAdd) {
				fmt.Println("✅ PATH already configured in shell config.")
				file.Close()
				return
			}
		}
		file.Close()
	}

	f, err := os.OpenFile(shellConfig, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open shell config: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString("\n# Added by kind installer\n" + exportLine + "\n"); err != nil {
		fmt.Printf("❌ Failed to write to shell config: %v\n", err)
		return
	}

	fmt.Println("✅ PATH updated. Please run:")
	exec.Command("source", shellConfig).Run()

	fmt.Printf("\n    source %s\n\n", shellConfig)
	fmt.Println("or open a new terminal to start using `kind`.")
}

func InstallHelm() {
	_, err := exec.LookPath("helm")
	if err == nil {
		fmt.Println("✅ Helm is already installed.")
		return
	}
	fmt.Println("Auto-detecting OS for Helm installation...")
	var url string
	var osArchBinMap = map[string]string{
		"linux/amd64":  "https://get.helm.sh/helm-v3.9.0-linux-amd64.tar.gz",
		"linux/arm64":  "https://get.helm.sh/helm-v3.9.0-linux-arm64.tar.gz",
		"darwin/amd64": "https://get.helm.sh/helm-v3.9.0-darwin-amd64.tar.gz",
		"darwin/arm64": "https://get.helm.sh/helm-v3.9.0-darwin-arm64.tar.gz",
	}
	osName, arch := AutoDetectOS()
	url = osArchBinMap[osName+"/"+arch]
	if url == "" {
		fmt.Printf("Unsupported architecture: %s/%s\n", osName, arch)
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory:", err)
		return
	}
	installDir := filepath.Join(homeDir, "go", "bin")
	output := filepath.Join(installDir, "helm.tar.gz")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		fmt.Println("Failed to create install directory:", err)
		return
	}
	fmt.Println("Downloading Helm binary from:", url)
	cmd := exec.Command("curl", "-Lo", output, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Download failed:", err)
		return
	}
	cmd = exec.Command("tar", "-xzf", output, "-C", installDir)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to extract Helm binary:", err)
		return
	}
	platformFolder := fmt.Sprintf("%s-%s", osName, arch)
	helmSrc := filepath.Join(installDir, platformFolder, "helm")

	cmd = exec.Command("chmod", "+x", helmSrc)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to make Helm binary executable:", err)
		return
	}

	if err := os.Rename(helmSrc, filepath.Join(installDir, "helm")); err != nil {
		fmt.Println("Failed to rename Helm binary:", err)
		return
	}
	fmt.Printf("\n✅ Helm installed to: %s\n", filepath.Join(installDir, "helm"))
	addPathToShellConfig(installDir)
	fmt.Println("🔧 Adding Helm to PATH in shell config.")
	fmt.Printf("Please run:\n    source %s\n\n", filepath.Join(os.Getenv("HOME"), ".bashrc"))
	fmt.Println("or open a new terminal to start using `helm`.")
}

func InstallKubectl() {
	_, err := exec.LookPath("kubectl")
	if err == nil {
		fmt.Println("✅ kubectl is already installed.")
		return
	}
	fmt.Println("Auto-detecting OS for kubectl installation...")
	var url string
	var osArchBinMap = map[string]string{
		"linux/amd64":  "https://dl.k8s.io/release/v1.23.0/bin/linux/amd64/kubectl",
		"linux/arm64":  "https://dl.k8s.io/release/v1.23.0/bin/linux/arm64/kubectl",
		"darwin/amd64": "https://dl.k8s.io/release/v1.23.0/bin/darwin/amd64/kubectl",
		"darwin/arm64": "https://dl.k8s.io/release/v1.23.0/bin/darwin/arm64/kubectl",
	}
	osName, arch := AutoDetectOS()
	url = osArchBinMap[osName+"/"+arch]
	if url == "" {
		fmt.Printf("Unsupported architecture: %s/%s\n", osName, arch)
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory:", err)
		return
	}
	installDir := filepath.Join(homeDir, "go", "bin")
	output := filepath.Join(installDir, "kubectl")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		fmt.Println("Failed to create install directory:", err)
		return
	}
	fmt.Println("Downloading kubectl binary from:", url)
	cmd := exec.Command("curl", "-Lo", output, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Download failed:", err)
		return
	}
	cmd = exec.Command("chmod", "+x", output)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to make kubectl binary executable:", err)
		return
	}
	fmt.Printf("\n✅ kubectl installed to: %s\n", output)
	addPathToShellConfig(installDir)
	fmt.Println("🔧 Adding kubectl to PATH in shell config.")
	fmt.Printf("Please run:\n    source %s\n\n", filepath.Join(os.Getenv("HOME"), ".bashrc"))
	fmt.Println("or open a new terminal to start using `kubectl`.")
}
