package resource

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"example.com/kube-visualizer/Utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type NodesUsage struct {
	Nodes []NodeData `json:"nodes"`
}
type NodeData struct {
	Name       string  `json:"name"`
	State      string  `json:"state"`
	Role       string  `json:"role"`
	CPU        int64   `json:"cpu"`    // in millicores
	Memory     int64   `json:"memory"` // in MB
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float64 `json:"mem_percent"`
}

// func GetClientSet(kubeconfig string) *kubernetes.Clientset {
// 	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
// 	if err != nil {
// 		log.Fatalf("Error building kubeconfig: %v", err)
// 	}
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		log.Fatalf("Error creating Kubernetes clientset: %v", err)
// 	}
// 	return clientset
// }

// func getMetricsClient(kubeconfig string) *metrics.Clientset {
// 	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
// 	if err != nil {
// 		log.Fatalf("Error building kubeconfig: %v", err)
// 	}
// 	metricsClient, err := metrics.NewForConfig(config)
// 	if err != nil {
// 		log.Fatalf("Error creating Metrics clientset: %v", err)
// 	}
// 	return metricsClient
// }

const metricsServerURL = "https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml"

func isMetricsServerAvailable(clientset *kubernetes.Clientset) bool {
	_, err := clientset.RESTClient().Get().AbsPath("/apis/metrics.k8s.io/v1beta1/nodes").DoRaw(context.TODO())
	if err != nil {
		return false
	}
	return true
}

func deployMetricsServer(config *rest.Config, clientset *kubernetes.Clientset) error {
	fmt.Println("⚠️ metrics-server not found, deploying now...")

	resp, err := http.Get(metricsServerURL)
	if err != nil {
		return fmt.Errorf("failed to download metrics-server manifest: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read metrics-server manifest: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(clientset.Discovery()))

	rbacYAML := `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: metrics-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:metrics-server
rules:
- apiGroups: [""]
  resources: ["nodes/metrics", "pods", "nodes", "namespaces", "configmaps"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:metrics-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:metrics-server
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: metrics-server:system:auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
`
	extraDocs := strings.Split(rbacYAML, "---")
	for _, doc := range extraDocs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}
		obj := &unstructured.Unstructured{}
		_, gvk, err := decUnstructured.Decode([]byte(doc), nil, obj)
		if err != nil {
			continue
		}
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}
		var ri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			ri = dynClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
		} else {
			ri = dynClient.Resource(mapping.Resource)
		}
		_, err = ri.Create(context.TODO(), obj, metav1.CreateOptions{})
		if err != nil && !errors.IsAlreadyExists(err) {
			return fmt.Errorf("failed creating RBAC %s: %w", obj.GetKind(), err)
		}
	}

	// Split multi-doc YAML and apply each
	docs := strings.Split(string(body), "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if len(doc) == 0 {
			continue
		}

		obj := &unstructured.Unstructured{}
		_, gvk, err := decUnstructured.Decode([]byte(doc), nil, obj)
		if err != nil {
			continue // skip invalid doc
		}

		if obj.GetKind() == "Deployment" && obj.GetName() == "metrics-server" {
			// force correct ServiceAccount
			_ = unstructured.SetNestedField(obj.Object, "metrics-server",
				"spec", "template", "spec", "serviceAccountName")

			containers, found, _ := unstructured.NestedSlice(obj.Object,
				"spec", "template", "spec", "containers")
			if found && len(containers) > 0 {
				c := containers[0].(map[string]interface{})
				args, _, _ := unstructured.NestedStringSlice(c, "args")
				args = append(args,
					"--kubelet-insecure-tls",
					"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
				)
				_ = unstructured.SetNestedStringSlice(c, args, "args")
				containers[0] = c
				_ = unstructured.SetNestedSlice(obj.Object, containers,
					"spec", "template", "spec", "containers")
			}
		}

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		var ri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			ns := obj.GetNamespace()
			if ns == "" {
				ns = "kube-system"
				obj.SetNamespace(ns)
			}
			ri = dynClient.Resource(mapping.Resource).Namespace(ns)
		} else {
			ri = dynClient.Resource(mapping.Resource)
		}

		_, err = ri.Create(context.TODO(), obj, metav1.CreateOptions{})
		if err != nil && !errors.IsAlreadyExists(err) {
			return fmt.Errorf("failed creating %s %s: %w", obj.GetKind(), obj.GetName(), err)
		}
	}

	fmt.Println("🚀 metrics-server deployed with correct RBAC")
	return nil
}

func waitForMetricsServerPod(clientset *kubernetes.Clientset, retries int, delay time.Duration) error {
	for i := 0; i < retries; i++ {
		pods, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{
			LabelSelector: "k8s-app=metrics-server",
		})
		if err != nil {
			return fmt.Errorf("error listing metrics-server pods: %w", err)
		}

		for _, pod := range pods.Items {
			for _, c := range pod.Status.Conditions {
				if c.Type == "Ready" && c.Status == "True" {
					fmt.Println("✅ metrics-server pod is running")
					return nil
				}
			}
		}

		fmt.Printf("⏳ Waiting for metrics-server pod to be ready... (%d/%d)\n", i+1, retries)
		time.Sleep(delay)
	}
	return fmt.Errorf("metrics-server pod not ready after %d retries", retries)
}

func MetricsServer(clientset *kubernetes.Clientset, kubeconfig string) {
	if isMetricsServerAvailable(clientset) {
		fmt.Println("✅ metrics-server is already running")
		return
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}
	err = deployMetricsServer(config, clientset)
	if err != nil {
		log.Fatalf("Error deploying metrics-server: %v", err)
	}
	err = waitForMetricsServerPod(clientset, 10, 6*time.Second)
	if err != nil {
		log.Fatalf("Error waiting for metrics-server API: %v", err)
	}
	fmt.Println("✅ metrics-server is now available")
}
func GetNodeResourceUsage(kubeconfig string, clientSet *kubernetes.Clientset) NodesUsage {
	clientset, err := Utils.GetClientSet(kubeconfig)
	MetricsServer(clientSet, kubeconfig)
	metricsClient := Utils.GetMetricsClient(kubeconfig)

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing nodes: %v", err)
	}

	nodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting node metrics: %v", err)
	}

	metricsMap := make(map[string]v1beta1.NodeMetrics)
	for _, m := range nodeMetrics.Items {
		metricsMap[m.Name] = m
	}

	NodesUsage := NodesUsage{}
	for _, node := range nodes.Items {
		role := "worker"
		if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
			role = "control-plane"
		} else if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
			role = "control-plane"
		}

		m := metricsMap[node.Name]
		usage := m.Usage
		cpu := usage.Cpu().MilliValue() // in millicores
		memory := usage.Memory().Value() / (1024 * 1024)
		allocCPU := node.Status.Allocatable.Cpu().MilliValue()               // millicores
		allocMem := node.Status.Allocatable.Memory().Value() / (1024 * 1024) // MB
		cpuPercent := float64(cpu) / float64(allocCPU) * 100
		memPercent := float64(memory) / float64(allocMem) * 100
		state := string(node.Status.Conditions[len(node.Status.Conditions)-1].Type)
		// in MB
		NodeData := NodeData{
			Name:       node.Name,
			State:      state,
			CPU:        cpu,
			Memory:     memory,
			Role:       role,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
		}
		NodesUsage.Nodes = append(NodesUsage.Nodes, NodeData)
		fmt.Printf("Node: %s | Role: %s | CPU: %dm | Memory: %dMi\n",
			node.Name, role, cpu, memory)
	}
	return NodesUsage

}
