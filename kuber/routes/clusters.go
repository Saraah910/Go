package routes

import (
	"fmt"
	"net/http"
	"time"

	"example.com/kuber/db"
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
	// clusterID, err := strconv.ParseInt(id, 10, 64)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid cluster ID format.", "Error": err.Error()})
	// 	return
	// }
	cluster, err := models.GetClusterByID(id)
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
	if cluster.WorkspaceName == "" {
		cluster.WorkspaceName = "default"
	}
	workspaceID, err := models.GetWorkspaceIDByName(cluster.WorkspaceName)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Workspace not found", "Error": err.Error()})
		return
	}
	uuid, err := db.GetUUID()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot generate UUID", "Error": err.Error()})
		return
	}
	cluster.ID = uuid
	cluster.WorkspaceID = workspaceID
	cluster.Status = "Completed"
	cluster.UserID = userID
	cluster.CreatedAt = time.Now()
	cluster.UpdatedAt = time.Now()

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

func getClustersByWorkspace(context *gin.Context) {
	workspaceID := context.Param("workspace_uuid")
	workspace, err := models.GetWorkspaceByUUID(workspaceID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch clusters by workspace name.", "Error": err.Error()})
		return
	}
	if workspace.ClusterCount == 0 {
		context.JSON(http.StatusNotFound, gin.H{"Message": "No clusters found for the specified workspace."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully fetched clusters by workspace UUID", "clusters": workspace})
}

func createCluster(context *gin.Context) {
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
			context.JSON(http.StatusBadRequest, gin.H{"Message": "No permission to create cluster."})
			return
		}
	}
	if cluster.WorkspaceName == "" {
		cluster.WorkspaceName = "default"
	}
	workspaceID, err := models.GetWorkspaceIDByName(cluster.WorkspaceName)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Workspace not found", "Error": err.Error()})
		return
	}
	infraList, err := models.GetInfraByName(cluster.Provisioner)
	if err != nil {
		return
	}
	fmt.Println(infraList)
	uuid, err := db.GetUUID()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot generate UUID", "Error": err.Error()})
		return
	}
	cluster.ID = uuid
	cluster.WorkspaceID = workspaceID
	cluster.Status = "Completed"
	cluster.UserID = userID
	cluster.CreatedAt = time.Now()
	cluster.UpdatedAt = time.Now()
	err = cluster.SaveCluster()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot create cluster", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully saved cluster", "cluster": cluster})
}
