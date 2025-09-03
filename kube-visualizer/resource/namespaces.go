package resource

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NsList struct {
	NamespaceList []string `json:"namespaces"`
	Count         int      `json:"count"`
}

func GetNamespaces(clientSet *kubernetes.Clientset) (*NsList, error) {
	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching namespaces: %w", err)
	}

	var nsNames []string
	for _, ns := range namespaces.Items {
		nsNames = append(nsNames, ns.Name)
	}

	return &NsList{
		NamespaceList: nsNames,
		Count:         len(nsNames),
	}, nil
}
