package getters

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreatePod(namespace string, client *kubernetes.Clientset) (time.Time, error) {
	podDefination := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "api-pod",
			Namespace:    "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "nginx-container-api",
					Image:           "nginx:latest",
					ImagePullPolicy: v1.PullIfNotPresent,
				},
			},
		},
	}
	pod, err := client.CoreV1().Pods(namespace).Create(context.Background(), podDefination, metav1.CreateOptions{})
	return pod.CreationTimestamp.Time, err

}
