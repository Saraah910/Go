package routes

import (
	"fmt"
	"net/http"
	"strconv"

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
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot bind user parameters.", "Error": err.Error()})
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

func logoutUser(context *gin.Context) {
	// userID := context.GetInt64("userID")

}

func getPermissions(context *gin.Context) {
	userId := context.GetInt64("userID")

	role, permission, err := models.GetPermission(userId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "User not found.", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Succesfully fetched permissions nd roles", "Roles": role, "Permission": permission})
}

func deleteUser(context *gin.Context) {
	keyString := context.Param("id")
	tbdID, err := strconv.ParseInt(keyString, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "User not found.", "Error": err.Error()})
		return
	}
	userID := context.GetInt64("userID")
	fmt.Println(tbdID, userID)
	if tbdID != userID {
		role, permission, err := models.GetPermission(userID)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized to delete.", "Error": err.Error()})
			return
		}
		if role != "admin" && permission != "full" {
			context.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized", "Error": err})
			return
		}
	}

	err = models.DeleteUser(tbdID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot delete user.", "Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Message": "Succesfully deleted user"})
}
