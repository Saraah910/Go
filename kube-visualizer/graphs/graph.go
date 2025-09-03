package graphs

import (
	"context"
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
)

type Node struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Namespace string `json:"namespace"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Rate   string `json:"rate"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
	Lock  sync.RWMutex
}

var stopCh chan struct{} = make(chan struct{})

func nodeID(obj *unstructured.Unstructured) string {
	return fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), obj.GetKind(), obj.GetName())
}

func (graph *Graph) AddNode(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	node := Node{
		Id:        nodeID(u),
		Type:      u.GetKind(),
		Namespace: u.GetNamespace(),
	}

	graph.Lock.Lock()
	defer graph.Lock.Unlock()
	// Append if not exists
	for _, n := range graph.Nodes {
		if n.Id == node.Id && n.Namespace == node.Namespace && n.Type == node.Type {
			return
		}
	}
	graph.Nodes = append(graph.Nodes, node)
	// Link resources dynamically
	graph.linkResources(u)
}

func (graph *Graph) UpdateNode(obj interface{}) {
	graph.AddNode(obj)
}

func (graph *Graph) RemoveNode(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	id := nodeID(u)

	graph.Lock.Lock()
	defer graph.Lock.Unlock()
	// Remove node
	newNodes := []Node{}
	for _, n := range graph.Nodes {
		if fmt.Sprintf("%s/%s/%s", n.Namespace, n.Type, n.Id) != id {
			newNodes = append(newNodes, n)
		}
	}
	graph.Nodes = newNodes

	newLinks := []Link{}
	for _, l := range graph.Links {
		if l.Source != id && l.Target != id {
			newLinks = append(newLinks, l)
		}
	}
	graph.Links = newLinks
}

func (g *Graph) addLink(src, tgt, rate string) {
	for _, l := range g.Links {
		if l.Source == src && l.Target == tgt && l.Rate == rate {
			return
		}
	}
	g.Links = append(g.Links, Link{Source: src, Target: tgt, Rate: rate})
}

func (g *Graph) linkResources(obj *unstructured.Unstructured) {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	srcID := nodeID(obj)

	// 1️⃣ OwnerReferences
	for _, owner := range obj.GetOwnerReferences() {
		targetID := fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), owner.Kind, owner.Name)
		g.addLink(srcID, targetID, "ownedBy")
	}

	// 2️⃣ Spec references (volumes, envFrom)
	spec, found, _ := unstructured.NestedMap(obj.Object, "spec")

	if found {

		if volumes, ok, _ := unstructured.NestedSlice(spec, "volumes"); ok {
			for _, v := range volumes {
				volMap := v.(map[string]interface{})
				if pvc, exists := volMap["persistentVolumeClaim"]; exists {
					pvcName := pvc.(map[string]interface{})["claimName"].(string)
					targetID := fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), "PersistentVolumeClaim", pvcName)
					g.addLink(srcID, targetID, "mounts")
				}
				if cm, exists := volMap["configMap"]; exists {
					cmName := cm.(map[string]interface{})["name"].(string)
					targetID := fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), "ConfigMap", cmName)
					g.addLink(srcID, targetID, "mounts")
				}
				if secret, exists := volMap["secret"]; exists {
					secName := secret.(map[string]interface{})["secretName"].(string)
					targetID := fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), "Secret", secName)
					g.addLink(srcID, targetID, "mounts")
				}
			}
		}

		// EnvFrom
		if containers, ok, _ := unstructured.NestedSlice(spec, "containers"); ok {
			for _, c := range containers {
				cMap := c.(map[string]interface{})
				if envFrom, ok, _ := unstructured.NestedSlice(cMap, "envFrom"); ok {
					for _, e := range envFrom {
						envMap := e.(map[string]interface{})
						if cmRef, exists := envMap["configMapRef"]; exists {
							cmName := cmRef.(map[string]interface{})["name"].(string)
							targetID := fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), "ConfigMap", cmName)
							g.addLink(srcID, targetID, "envForm")
						}
						if secRef, exists := envMap["secretRef"]; exists {
							secName := secRef.(map[string]interface{})["name"].(string)
							targetID := fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), "Secret", secName)
							g.addLink(srcID, targetID, "envForm")
						}
					}
				}
			}
		}
	}
}

func watchResource(dynamicClient dynamic.Interface, gvr schema.GroupVersionResource, namespace string, graph *Graph, stopCh <-chan struct{}) {
	informer := cache.NewSharedInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (runtime.Object, error) {
				return dynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return dynamicClient.Resource(gvr).Namespace(namespace).Watch(context.TODO(), lo)
			},
		},
		&unstructured.Unstructured{},
		0,
	)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { graph.AddNode(obj) },
		UpdateFunc: func(_, newObj interface{}) { graph.UpdateNode(newObj) },
		DeleteFunc: func(obj interface{}) { graph.RemoveNode(obj) },
	})

	go informer.Run(stopCh)
}

func GenerateGraph(dynamicClient dynamic.Interface, gvrList []schema.GroupVersionResource, namespace string) (*Graph, error) {
	graph := &Graph{
		Nodes: []Node{},
		Links: []Link{},
	}
	podsRes := dynamicClient.Resource(schemaGroupVersionResource("v1", "pods")).Namespace(namespace)
	pods, err := podsRes.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}
	for _, pod := range pods.Items {
		podName := pod.GetName()
		podNode := Node{Id: podName, Type: "Pod", Namespace: namespace}
		graph.Nodes = append(graph.Nodes, podNode)

		// Check PVCs attached to Pod
		volumes, found, _ := unstructured.NestedSlice(pod.Object, "spec", "volumes")
		if found {
			for _, v := range volumes {
				vol := v.(map[string]interface{})
				if pvc, ok := vol["persistentVolumeClaim"]; ok {
					claimName := pvc.(map[string]interface{})["claimName"].(string)
					graph.Nodes = append(graph.Nodes, Node{Id: claimName, Type: "PVC", Namespace: namespace})
					graph.Links = append(graph.Links, Link{Source: podName, Target: claimName})
				}
				if cm, ok := vol["configMap"]; ok {
					cmName := cm.(map[string]interface{})["name"].(string)
					graph.Nodes = append(graph.Nodes, Node{Id: cmName, Type: "ConfigMap", Namespace: namespace})
					graph.Links = append(graph.Links, Link{Source: podName, Target: cmName})
				}
				if secret, ok := vol["secret"]; ok {
					secretName := secret.(map[string]interface{})["secretName"].(string)
					graph.Nodes = append(graph.Nodes, Node{Id: secretName, Type: "Secret", Namespace: namespace})
					graph.Links = append(graph.Links, Link{Source: podName, Target: secretName})
				}

			}
		}

		// Check Services that select this Pod
		servicesRes := dynamicClient.Resource(schemaGroupVersionResource("v1", "services")).Namespace(namespace)
		services, err := servicesRes.List(context.Background(), metav1.ListOptions{})
		if err == nil {
			podLabels := pod.GetLabels()
			for _, svc := range services.Items {
				selector, _, _ := unstructured.NestedStringMap(svc.Object, "spec", "selector")
				if matchesSelector(podLabels, selector) {
					svcName := svc.GetName()
					graph.Nodes = append(graph.Nodes, Node{Id: svcName, Type: "Service", Namespace: namespace})
					graph.Links = append(graph.Links, Link{Source: svcName, Target: podName})
				}
			}
		}
	}

	return graph, nil
}

func GetGraph(dynamicClient dynamic.Interface, gvrList []schema.GroupVersionResource, namespace string) (*Graph, error) {
	graph, err := GenerateGraph(dynamicClient, gvrList, namespace)
	if err != nil {
		fmt.Printf("Cannot generate graph. %v", err.Error())
		return nil, err
	}
	for _, gvr := range gvrList {
		watchResource(dynamicClient, gvr, namespace, graph, stopCh)
	}

	return graph, nil
}

func schemaGroupVersionResource(version, resource string) schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "", Version: version, Resource: resource}
}

func matchesSelector(podLabels, selector map[string]string) bool {
	for k, v := range selector {
		if podLabels[k] != v {
			return false
		}
	}
	return true
}
