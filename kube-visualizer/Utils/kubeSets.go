package Utils

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetClientSet(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}
	return clientset, nil
}

func GetMetricsClient(kubeconfig string) *metrics.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}
	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Metrics clientset: %v", err)
	}
	return metricsClient
}

func GetDynamicClients(kubeconfig string) (*dynamic.DynamicClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	config.QPS = 50 // default is 5
	config.Burst = 100

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}

func GetDiscoveryAPIs(kubeconfig string, dynamicClient *dynamic.DynamicClient) ([]schema.GroupVersionResource, []schema.GroupVersionResource, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}
	config.QPS = 50 // default is 5
	config.Burst = 100
	discoveryClient := kubernetes.NewForConfigOrDie(config).Discovery()
	apiGroups, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		log.Println("Warning: some resources cannot be accessed:", err)
		return nil, nil, errors.New("some resources cannot be accissible.")
	}

	var watchableGVRs []schema.GroupVersionResource
	var clusterScoped []schema.GroupVersionResource
	for _, apiGroup := range apiGroups {
		group, version := parseGroupVersion(apiGroup.GroupVersion)
		for _, r := range apiGroup.APIResources {
			if !supportsWatch(r.Verbs) {
				continue
			}
			gvr := schema.GroupVersionResource{
				Group:    group,
				Version:  version,
				Resource: r.Name,
			}

			if r.Namespaced {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				_, err := dynamicClient.Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{Limit: 1})
				if err == nil {
					watchableGVRs = append(watchableGVRs, gvr)
				}
			} else {
				clusterScoped = append(clusterScoped, gvr)
			}

			// else skip silently; resource is not watchable now
		}
	}

	if len(watchableGVRs) == 0 && len(clusterScoped) == 0 {
		return nil, nil, errors.New("no resources found")
	}
	return watchableGVRs, clusterScoped, nil
}

func parseGroupVersion(groupVersion string) (string, string) {
	parts := strings.Split(groupVersion, "/")
	if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[0], parts[1]
}

func supportsWatch(verbs []string) bool {
	for _, v := range verbs {
		if v == "watch" {
			return true
		}
	}
	return false
}
