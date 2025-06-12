package routes

import (
	"example.com/kuber/DR"
	"example.com/kuber/models"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ClusterDRRequest struct {
	SourceClusterID      int64  `json:"source_cluster_id" binding:"required"`
	DestinationClusterID int64  `json:"destination_cluster_id" binding:"required"`
	DRtype               string `json:"dr_type" binding:"required" validate:"oneof=active-passive active-active"`
}

var gvrs = []schema.GroupVersionResource{
	{Group: "apps", Version: "v1", Resource: "deployments"},
	{Group: "", Version: "v1", Resource: "services"},
	{Group: "", Version: "v1", Resource: "configmaps"},
	{Group: "", Version: "v1", Resource: "secrets"},
	{Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
	{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
	// Add custom CRDs here if needed
}

func clusterDR(context *gin.Context) {
	var ClusterDRReq ClusterDRRequest
	if err := context.ShouldBindJSON(&ClusterDRReq); err != nil {
		context.JSON(400, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}
	if ClusterDRReq.SourceClusterID == ClusterDRReq.DestinationClusterID {
		context.JSON(400, gin.H{
			"error": "Source and destination clusters cannot be the same",
		})
		return
	}
	if ClusterDRReq.DRtype != "active-passive" && ClusterDRReq.DRtype != "active-active" {
		context.JSON(400, gin.H{
			"error": "Invalid DR type. Must be either 'active-passive' or 'active-active'",
		})
		return
	}
	// if ClusterDRReq.DRtype == "active-passive" {
	// 	context.JSON(400, gin.H{
	// 		"error": "DR type is required",
	// 	})
	// 	return
	// }
	// Check if source and destination clusters exist
	sourceCluster, err := models.GetClusterByID(ClusterDRReq.SourceClusterID)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to get source cluster: " + err.Error(),
		})
		return
	}
	if sourceCluster == nil {
		context.JSON(404, gin.H{
			"error": "Source cluster not found",
		})
		return
	}
	destinationCluster, err := models.GetClusterByID(ClusterDRReq.DestinationClusterID)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to get destination cluster: " + err.Error(),
		})
		return
	}
	if destinationCluster == nil {
		context.JSON(404, gin.H{
			"error": "Destination cluster not found",
		})
		return
	}

	context.JSON(200, gin.H{
		"message": "Cluster disaster recovery initiated",
	})
	SourceKubeconfig, err := models.GetKubeconfigFilePathByID(ClusterDRReq.SourceClusterID)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to get source kubeconfig: " + err.Error(),
		})
		return
	}
	DestinationKubeconfig, err := models.GetKubeconfigFilePathByID(ClusterDRReq.DestinationClusterID)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to get destination kubeconfig: " + err.Error(),
		})
		return
	}
	SourceClient, _, sourceK8sClient, err := DR.GetDynamicClient(SourceKubeconfig)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to create source dynamic client: " + err.Error(),
		})
		return
	}
	if sourceK8sClient == nil {
		context.JSON(500, gin.H{
			"error": "Failed to create source k8s client: client is nil",
		})
		return
	}
	if SourceClient == nil {
		context.JSON(500, gin.H{
			"error": "Failed to create source dynamic client: client is nil",
		})
		return
	}
	DestinationClient, _, _, err := DR.GetDynamicClient(DestinationKubeconfig)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to create destination dynamic client: " + err.Error(),
		})
		return
	}

	if DestinationClient == nil {
		context.JSON(500, gin.H{
			"error": "Failed to create destination dynamic client: client is nil",
		})
		return
	}
	povisioner, err := models.GetProvisionerByID(ClusterDRReq.SourceClusterID)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to get provisioner: " + err.Error(),
		})
		return
	}
	if povisioner == "" {
		context.JSON(500, gin.H{
			"error": "Provisioner not found for the source cluster",
		})
		return
	}

	err = DR.PerformClusterDR(SourceClient, DestinationClient, povisioner, ClusterDRReq.DRtype, gvrs, sourceK8sClient)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to perform cluster disaster recovery: " + err.Error(),
		})
		return
	}
	context.JSON(200, gin.H{
		"message":                "Cluster disaster recovery completed successfully",
		"source_cluster_id":      ClusterDRReq.SourceClusterID,
		"destination_cluster_id": ClusterDRReq.DestinationClusterID,
		"dr_type":                ClusterDRReq.DRtype,
	})
}
