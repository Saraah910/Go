package DR

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
)

func GetKubeClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create Kubernetes client: %v\n", err)
		return nil, err
	}
	return clientset, nil
}

func PerformClusterDR(SourceClient, DestinationClient *kubernetes.Clientset, provisioner string, DRType string) error {
	if DRType != "active-passive" {
		return fmt.Errorf("DRType %s not supported yet", DRType)
	}
	if err := syncNamespaces(SourceClient, DestinationClient); err != nil {
		return fmt.Errorf("failed syncing namespaces: %v", err)
	}
	if err := syncWorkloads(SourceClient, DestinationClient); err != nil {
		return fmt.Errorf("failed syncing workloads: %v", err)
	}
	if err := syncPersistentVolumes(SourceClient, DestinationClient, provisioner); err != nil {
		return fmt.Errorf("failed syncing PVs: %v", err)
	}
	go monitorClusterHealth(SourceClient, func() {
		fmt.Println("❗ Source cluster unhealthy. Initiating failover...")
		if err := initiateFailover(DestinationClient); err != nil {
			fmt.Printf("⚠️ Failover error: %v\n", err)
		}
	})

	return nil
}

func syncNamespaces(src, dst *kubernetes.Clientset) error {
	namespaces, err := src.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ns := range namespaces.Items {
		if errors.IsNotFound(err) {
			ns.ResourceVersion = "" // Clear RV for creation
			if _, err := dst.CoreV1().Namespaces().Create(context.TODO(), &ns, metav1.CreateOptions{}); err != nil {
				return fmt.Errorf("failed to create ns %s: %v", ns.Name, err)
			}
			fmt.Printf("✅ Namespace %s synced\n", ns.Name)
		}
		if errors.IsAlreadyExists(err) {
			// If namespace already exists, we can skip or update if needed
			fmt.Printf("ℹ️ Namespace %s already exists, skipping\n", ns.Name)
		} else {
			fmt.Printf("⚠️ Error syncing namespace %s: %v\n", ns.Name, err)
		}
	}
	return nil
}
func syncWorkloads(src, dst *kubernetes.Clientset) error {
	nsList, err := src.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ns := range nsList.Items {
		deployments, err := src.AppsV1().Deployments(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			continue
		}
		for _, deploy := range deployments.Items {
			deploy.ResourceVersion = ""
			_, err := dst.AppsV1().Deployments(ns.Name).Create(context.TODO(), &deploy, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				fmt.Printf("⚠️ Deployment %s failed to sync: %v\n", deploy.Name, err)
			} else {
				fmt.Printf("✅ Deployment %s synced to %s\n", deploy.Name, ns.Name)
			}
		}
	}
	return nil
}

// Services → CoreV1().Services(ns)
// ConfigMaps → CoreV1().ConfigMaps(ns)
// Secrets → CoreV1().Secrets(ns)
func syncPersistentVolumes(src, dst *kubernetes.Clientset, provisioner string) error {
	pvs, err := src.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, pv := range pvs.Items {
		if strings.Contains(pv.Spec.StorageClassName, provisioner) {
			pv.ResourceVersion = ""
			_, err := dst.CoreV1().PersistentVolumes().Create(context.TODO(), &pv, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				fmt.Printf("⚠️ PV %s replication failed: %v\n", pv.Name, err)
			}
		}
	}
	return nil
}

func monitorClusterHealth(client *kubernetes.Clientset, onFailure func()) {
	for {
		_, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			onFailure()
			return
		}
		time.Sleep(10 * time.Second)
	}
}
func initiateFailover(dst *kubernetes.Clientset) error {
	// Sample: scale replicas of all deployments to 1
	nsList, err := dst.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ns := range nsList.Items {
		deployments, _ := dst.AppsV1().Deployments(ns.Name).List(context.TODO(), metav1.ListOptions{})
		for _, d := range deployments.Items {
			d.Spec.Replicas = pointer.Int32Ptr(1)
			_, err := dst.AppsV1().Deployments(ns.Name).Update(context.TODO(), &d, metav1.UpdateOptions{})
			if err != nil {
				fmt.Printf("⚠️ Failover update failed for %s: %v\n", d.Name, err)
			}
		}
	}
	return nil
}
