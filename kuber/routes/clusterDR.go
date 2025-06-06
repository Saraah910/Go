package routes

import (
	"example.com/kuber/DR"
	"example.com/kuber/models"
	"github.com/gin-gonic/gin"
)

type ClusterDRRequest struct {
	SourceClusterID      int64  `json:"source_cluster_id" binding:"required"`
	DestinationClusterID int64  `json:"destination_cluster_id" binding:"required"`
	DRtype               string `json:"dr_type" binding:"required" validate:"oneof=active-passive active-active"`
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
	SourceClient, err := DR.GetKubeClient(SourceKubeconfig)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to create source kube client: " + err.Error(),
		})
		return
	}
	DestinationClient, err := DR.GetKubeClient(DestinationKubeconfig)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "Failed to create destination kube client: " + err.Error(),
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
	err = DR.PerformClusterDR(SourceClient, DestinationClient, povisioner, ClusterDRReq.DRtype)
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
