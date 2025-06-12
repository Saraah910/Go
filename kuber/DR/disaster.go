package DR

import (
	"context"
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func GetDynamicClient(kubeconfigPath string) (dynamic.Interface, *restmapper.DeferredDiscoveryRESTMapper, *kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, nil, nil, err
	}
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	return client, mapper, k8sClient, nil
}

func PerformClusterDR(SourceClient, DestinationClient dynamic.Interface, provisioner string, DRType string, gvrs []schema.GroupVersionResource, sourceK8sClient *kubernetes.Clientset) error {
	if SourceClient == nil || DestinationClient == nil {
		return fmt.Errorf("source or destination client is nil")
	}
	if provisioner == "" {
		return fmt.Errorf("provisioner cannot be empty")
	}
	if DRType == "" {
		return fmt.Errorf("DRType cannot be empty")
	}
	if DRType == "active-passive" {
		return fmt.Errorf("DRType %s not supported yet", DRType)
	}
	watchNamespacesAndSyncResources(SourceClient, DestinationClient, sourceK8sClient)

	return nil
}

func watchAndSync(gvr schema.GroupVersionResource, sourceClient, destClient dynamic.Interface, namespace string) {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(sourceClient, 0, namespace, nil)
	informer := factory.ForResource(gvr).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			u.SetResourceVersion("")
			_, err := destClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), u, metav1.CreateOptions{})
			if err != nil {
				log.Printf("Add error: %v", err)
			} else {
				log.Printf("✅ Created %s in %s", u.GetName(), namespace)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			u := newObj.(*unstructured.Unstructured)
			cleanObject(u)

			_, err := destClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), u, metav1.UpdateOptions{})
			if err != nil {
				// fallback to delete and create
				_ = destClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), u.GetName(), metav1.DeleteOptions{})
				_, err = destClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), u, metav1.CreateOptions{})
				if err != nil {
					log.Printf("❌ Fallback create failed for %s in %s: %v", u.GetName(), namespace, err)
				} else {
					log.Printf("✅ Fallback recreated %s in %s", u.GetName(), namespace)
				}
			} else {
				log.Printf("🔄 Updated %s in %s", u.GetName(), namespace)
			}
		},

		DeleteFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			err := destClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), u.GetName(), metav1.DeleteOptions{})
			if err != nil {
				log.Printf("Delete error: %v", err)
			} else {
				log.Printf("🗑️ Deleted %s from %s", u.GetName(), namespace)
			}
		},
	})

	stop := make(chan struct{})
	go informer.Run(stop)
	if !cache.WaitForCacheSync(stop, informer.HasSynced) {
		log.Fatalf("Failed to sync cache for %v in namespace %s", gvr, namespace)
	}
}

func watchNamespacesAndSyncResources(sourceClient, destClient dynamic.Interface, sourceK8sClient *kubernetes.Clientset) {
	nsWatcher, err := sourceK8sClient.CoreV1().Namespaces().Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to watch namespaces: %v", err)
	}

	// Define resources to sync
	gvrs := []schema.GroupVersionResource{
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "configmaps"},
		{Group: "", Version: "v1", Resource: "secrets"},
		{Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
		{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
	}

	log.Println("📡 Watching for namespace changes...")
	go func() {
		for event := range nsWatcher.ResultChan() {
			ns, ok := event.Object.(*corev1.Namespace)
			if !ok {
				log.Printf("Failed to cast namespace object, got type: %T", event.Object)
				continue
			}
			namespace := ns.GetName()
			fmt.Printf("Namespace event: %s (%s)\n", namespace, event.Type)
			if event.Type == watch.Added || event.Type == watch.Modified {
				log.Printf("🔍 Namespace change detected: %s (%s)", namespace, event.Type)
				for _, gvr := range gvrs {
					go watchAndSync(gvr, sourceClient, destClient, namespace)
				}
			}
		}
	}()
}

func cleanObject(u *unstructured.Unstructured) {
	unstructured.RemoveNestedField(u.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(u.Object, "metadata", "uid")
	unstructured.RemoveNestedField(u.Object, "metadata", "creationTimestamp")
	unstructured.RemoveNestedField(u.Object, "metadata", "managedFields")
	unstructured.RemoveNestedField(u.Object, "metadata", "selfLink")
	unstructured.RemoveNestedField(u.Object, "metadata", "generation")
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
			d.Spec.Replicas = ptr.To(int32(1)) // Scale to 1 replica
			_, err := dst.AppsV1().Deployments(ns.Name).Update(context.TODO(), &d, metav1.UpdateOptions{})
			if err != nil {
				fmt.Printf("⚠️ Failover update failed for %s: %v\n", d.Name, err)
			}
		}
	}
	return nil
}
