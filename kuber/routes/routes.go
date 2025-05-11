package routes

import (
	"example.com/kuber/Middleware"
	"github.com/gin-gonic/gin"
)

func Routes(server *gin.Engine) {
	server.GET("/users/list", getAllUsers)
	server.POST("/users/signup", signUpUser)
	server.POST("/users/login", loginUser)

	authenticationGroups := server.Group("/")
	authenticationGroups.Use(Middleware.Authentication)
	authenticationGroups.POST("/kube/clusters/list")

}
