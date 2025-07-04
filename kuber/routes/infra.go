package routes

import (
	"net/http"

	"example.com/kuber/models"
	"github.com/gin-gonic/gin"
)

func AWSInfra(context *gin.Context) {
	var InfraInput models.AWSInfra
	err := context.ShouldBindJSON(&InfraInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Check parameters", "Error": err.Error()})
		return
	}
	userID := context.GetInt64("userID")
	IsAuthorized, err := isAuthorized(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if !IsAuthorized {
		context.JSON(http.StatusUnauthorized, gin.H{"Error": err})
		return
	}
	InfraInput.UserID = userID
	_, err = InfraInput.SaveAWSInfra()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot create Infrastructure", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully saved infrastructure", "Infrastructure": InfraInput})
}

func NutanixInfra(context *gin.Context) {
	var InfraInput models.NutanixInfra
	err := context.ShouldBindJSON(&InfraInput)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Check parameters", "Error": err.Error()})
		return
	}
	userID := context.GetInt64("userID")
	IsAuthorized, err := isAuthorized(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if !IsAuthorized {
		context.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		return
	}
	InfraInput.UserID = userID
	_, err = InfraInput.SaveNutanixInfra()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot create Infrastructure", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully saved infrastructure", "Infrastructure": InfraInput})

}

// func getInfraForUser(ctx *gin.Context) {
// 	KeyString := ctx.Param("id")
// 	userID, err := strconv.ParseInt(KeyString, 10, 64)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"Msg": err.Error()})
// 	}
// 	models.GetInfraByUserID(userID)
// }

func getAllInfrastructures(context *gin.Context) {
	infrastructures, err := models.GetInfrastructures()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch results.", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Successfully fetched infrastructures", "infrastructures": infrastructures})
}
