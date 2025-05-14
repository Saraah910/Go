package routes

import (
	"net/http"
	"strconv"

	"example.com/kuber/Concurrency"
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

func postCluster(context *gin.Context) {
	var cluster models.Kube
	err := context.ShouldBindJSON(&cluster)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Check parameters", "Error": err.Error()})
		return
	}
	userID := context.GetInt64("userID")
	cluster.Status = "Completed"
	cluster.UserID = userID

	err = cluster.SaveCluster()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot create cluster", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully saved clusters", "clusters": cluster})
}

func getClusterByID(context *gin.Context) {
	keyString := context.Param("id")
	clusterID, err := strconv.ParseInt(keyString, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cluster not exists", "Error": err.Error()})
		return
	}
	cluster, err := models.GetClusterById(clusterID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch cluster.", "id": clusterID})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": cluster, "Message": "Successful"})

}

func updateCluster(context *gin.Context) {
	keyString := context.Param("id")
	clusterID, err := strconv.ParseInt(keyString, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cluster id not exists", "Error": err.Error()})
		return
	}
	cluster, err := models.GetClusterById(clusterID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cluster id not exists", "Error": err.Error()})
		return
	}
	userId := context.GetInt64("userID")

	if userId != cluster.UserID {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": "User not authorized"})
		return
	}
	var newCluster models.Kube
	err = context.ShouldBindJSON(&newCluster)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Check parameters", "Error": err.Error()})
		return
	}
	newCluster.ClusterID = clusterID
	err = newCluster.UpdateCluster(clusterID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot update cluster"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Cluster": newCluster, "Message": "Successful"})

}

func performAction(context *gin.Context) {
	keyString := context.Param("id")
	action := context.Param("action")
	clusterID, err := strconv.ParseInt(keyString, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cluster not exists", "Error": err.Error()})
		return
	}
	Concurrency.ConcurrentFunctions(action)
	context.JSON(http.StatusOK, gin.H{"Cluster ID": clusterID, "Message": "Successful", "Action": action})
}
