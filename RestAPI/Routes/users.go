package Routes

import (
	"net/http"

	"example.com/APIs/Utils"
	"example.com/APIs/models"
	"github.com/gin-gonic/gin"
)

func signupUser(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Check parameters"})
		return
	}
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not create user."})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"Message": "Successful"})

}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not parse user."})
		return
	}
	err = user.ValidateCreds()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": err})
		return
	}
	token, err := Utils.GenerateToken(user.Email, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not generate token", "error": err})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"Message": "Successful user login", "token": token})
}
