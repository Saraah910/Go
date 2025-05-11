package routes

import (
	"net/http"

	"example.com/kuber/Utils"
	"example.com/kuber/models"
	"github.com/gin-gonic/gin"
)

func getAllUsers(context *gin.Context) {
	users, err := models.FetchAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot fetch all users."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Succesfully fetched all users", "users": users})
}

func signUpUser(context *gin.Context) {
	var user models.Users
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot bind user parameters."})
		return
	}
	err = user.CreateUser()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot signup", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Succesfully created user", "userID": user.Id})
}

func loginUser(context *gin.Context) {
	var user models.Users
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot bind user parameters."})
		return
	}
	err = user.ValidateCreds()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot validate user credentials.", "Error": err.Error()})
		return
	}
	token, err := Utils.GenerateToken(user.Id, user.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot generate token.", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Succesfully logged in user", "token": token})

}
