package models

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func fetchClusterDetailsFromKubeconfig(kubeconfigPath string) (controlPlanes []interface{}, workerNodes []interface{}, storageContainers []interface{}) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create Kubernetes client: %v\n", err)
		return
	}

	ctx := context.Background()

	// Get nodes and classify
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to list nodes: %v\n", err)
		return
	}

	for _, node := range nodes.Items {
		nodeInfo := map[string]string{
			"name":   node.Name,
			"status": string(getNodeReadyStatus(node)),
		}

		if isControlPlaneNode(node) {
			controlPlanes = append(controlPlanes, nodeInfo)
		} else {
			workerNodes = append(workerNodes, nodeInfo)
		}
	}

	// Get PVCs
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to list PVCs: %v\n", err)
		return
	}

	for _, pvc := range pvcs.Items {
		pvcInfo := map[string]string{
			"name":   pvc.Name,
			"status": string(pvc.Status.Phase),
		}
		storageContainers = append(storageContainers, pvcInfo)
	}

	return
}

func isControlPlaneNode(node v1.Node) bool {
	for _, taint := range node.Spec.Taints {
		if taint.Key == "node-role.kubernetes.io/control-plane" || taint.Key == "node-role.kubernetes.io/master" {
			return true
		}
	}
	return false
}

func getNodeReadyStatus(node v1.Node) v1.ConditionStatus {
	for _, cond := range node.Status.Conditions {
		if cond.Type == v1.NodeReady {
			return cond.Status
		}
	}
	return v1.ConditionUnknown
}

func GetServices(clusterID int64) ([]map[string]interface{}, error) {

	kubeconfigFilePath, err := GetKubeconfigFilePathByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig file path: %w", err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	ctx := context.Background()
	services, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var serviceList []map[string]interface{}
	for _, service := range services.Items {
		serviceMap := map[string]interface{}{
			"name":      service.Name,
			"namespace": service.Namespace,
			"clusterIP": service.Spec.ClusterIP,
			"ports":     service.Spec.Ports,
			"selector":  service.Spec.Selector,
			"type":      service.Spec.Type,
			"createdAt": service.CreationTimestamp,
		}
		serviceList = append(serviceList, serviceMap)
	}

	return serviceList, nil
}
