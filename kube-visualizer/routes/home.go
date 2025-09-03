package routes

import (
	"fmt"

	"example.com/kube-visualizer/Utils"
	"example.com/kube-visualizer/graphs"
	"example.com/kube-visualizer/resource"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Kubesets struct {
	KubeconfigPath string
	ClientSet      *kubernetes.Clientset
	MetricsSet     metrics.Clientset
}

type UploadRequest struct {
	FileName  string `json:"file_name" binding:"required"`
	SavedPath string `json:"file_path"`
}

var latestFile UploadRequest
var Sets Kubesets

func homepage(c *gin.Context) {
	c.String(200, "Welcome to the Kube Visualizer!")
}

type KubeRequest struct {
	FileName string `json:"file_name"`
}

func setKubeConfig(c *gin.Context) {
	var req KubeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "file_name is required"})
		return
	}

	kubeconfig := "./uploads/" + req.FileName

	// Recreate clientset & metricsset
	newSets, err := CreateKubeSets(kubeconfig)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to initialize kubeconfig"})
		return
	}

	Sets = newSets // ✅ update global

	c.JSON(200, gin.H{"message": "kubeconfig reinitialized", "file": req.FileName})
}

func uploadFile(c *gin.Context) {
	var req UploadRequest
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File is required"})
		return
	}
	req.FileName = file.Filename
	req.SavedPath = "./uploads/" + file.Filename
	err = c.SaveUploadedFile(file, req.SavedPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}
	latestFile = req
	c.JSON(200, gin.H{"message": "File received", "file_name": req.FileName, "file_path": req.SavedPath})
	Sets, err = CreateKubeSets(req.SavedPath)
	if err != nil {
		c.JSON(200, gin.H{"message": "Cannot create clients"})
		return
	}

}

func CreateKubeSets(kubdeconfigPath string) (Kubesets, error) {
	clientSets, err := Utils.GetClientSet(kubdeconfigPath)
	if err != nil {
		fmt.Print("cannot create clientSet")
	}
	metricSets := Utils.GetMetricsClient(kubdeconfigPath)
	return Kubesets{
		KubeconfigPath: kubdeconfigPath,
		ClientSet:      clientSets,
		MetricsSet:     *metricSets,
	}, nil
}

func nodeResourceUsage(c *gin.Context) {
	var req UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "file_name is required"})
		return
	}
	kubeconfig := "./uploads/" + req.FileName
	clientSet, err := Utils.GetClientSet(kubeconfig)
	if err != nil {
		c.JSON(200, gin.H{"Message": "Cannot find file"})
	}
	usage := resource.GetNodeResourceUsage(kubeconfig, clientSet)
	c.JSON(200, gin.H{"node_resource_usage": usage})
}

func getNamespaces(c *gin.Context) {
	clientSet := Sets.ClientSet
	namespaces, err := resource.GetNamespaces(clientSet)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"Message": "Successfully fetched namespaces", "namespaces": namespaces})
}

func graphData(c *gin.Context) {
	ns := c.Param("ns")
	fmt.Print(ns)
	dynamicClienet, err := Utils.GetDynamicClients(Sets.KubeconfigPath)
	if err != nil {
		c.JSON(403, gin.H{"Message": "Cannot create the dynamic client for the provided kubeconfig", "eror": err.Error()})
		return
	}

	gvrList, _, err := Utils.GetDiscoveryAPIs(Sets.KubeconfigPath, dynamicClienet)
	if err != nil {
		c.JSON(403, gin.H{"Message": err.Error()})
		return
	}
	graph, err := graphs.GetGraph(dynamicClienet, gvrList, ns)
	if err != nil {
		c.JSON(403, gin.H{"Message": "Error generating graph", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Message": "Successfully fetched graph", "graph": graph})
}
