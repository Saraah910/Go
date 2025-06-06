package routes

import (
	"net/http"
	"strconv"
	"time"

	"example.com/kuber/models"
	"github.com/gin-gonic/gin"
)

func getAllClusters(context *gin.Context) {
	clusters, err := models.GetAllClusters()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch results.", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully fetched clusters", "clusters": clusters})
}

func getClusterByID(context *gin.Context) {
	id := context.Param("id")
	clusterID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid cluster ID format.", "Error": err.Error()})
		return
	}
	cluster, err := models.GetClusterByID(clusterID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch cluster by ID.", "Error": err.Error()})
		return
	}
	if cluster == nil {
		context.JSON(http.StatusNotFound, gin.H{"Message": "Cluster not found."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully fetched cluster", "cluster": cluster})
}

func addCluster(context *gin.Context) {
	var cluster models.Cluster
	err := context.ShouldBindJSON(&cluster)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Check parameters", "Error": err.Error()})
		return
	}
	userID := context.GetInt64("userID")
	role, permission, err := models.GetPermission(userID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot get roles and permissions"})
		return
	}
	if role != "admin" {
		if permission != "write" && permission != "full" {
			context.JSON(http.StatusBadRequest, gin.H{"Message": "No permission to add cluster."})
			return
		}
	}
	if cluster.Workspace == "" {
		cluster.Workspace = "default"
	}
	cluster.Status = "Completed"
	cluster.UserID = userID
	cluster.CreatedAt = time.Now().Format(time.RFC3339)
	cluster.UpdatedAt = time.Now().Format(time.RFC3339)

	err = cluster.SaveCluster()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot create cluster", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully saved clusters", "clusters": cluster})
}

func getClustersForUser(context *gin.Context) {
	userID := context.GetInt64("userID")
	clusters, err := models.GetClustersForUser(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch clusters for user.", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully fetched clusters for user", "clusters": clusters})
}

func getClustersByWorkspaceName(context *gin.Context) {
	workspace := context.Param("workspace")
	clusters, err := models.GetClustersByWorkspaceName(workspace)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch clusters by workspace name.", "Error": err.Error()})
		return
	}
	if len(clusters) == 0 {
		context.JSON(http.StatusNotFound, gin.H{"Message": "No clusters found for the specified workspace."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully fetched clusters by workspace name", "clusters": clusters})
}
