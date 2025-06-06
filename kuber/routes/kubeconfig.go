package routes

import (
	"net/http"
	"strconv"

	"example.com/kuber/models"
	"github.com/gin-gonic/gin"
)

func getServices(context *gin.Context) {
	clusterIDStr := context.Param("id")
	clusterID, err := strconv.ParseInt(clusterIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid cluster ID.", "Error": err.Error()})
		return
	}
	services, err := models.GetServices(clusterID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch services.", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully fetched services", "services": services})
}
