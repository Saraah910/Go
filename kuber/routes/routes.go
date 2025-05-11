package routes

import (
	"example.com/kuber/Middleware"
	"github.com/gin-gonic/gin"
)

func Routes(server *gin.Engine) {
	server.GET("/users/list", getAllUsers)
	server.POST("/users/signup", signUpUser)
	server.POST("/users/login", loginUser)
	server.GET("/kube/clusters/list", getAllClusters)
	server.GET("/kube/cluster/:id/actions", getClusterByID)

	authenticationGroups := server.Group("/")
	authenticationGroups.Use(Middleware.Authentication)
	authenticationGroups.POST("/kube/clusters/list", postCluster)
	authenticationGroups.GET("/kube/cluster/:id/actions/:action", performAction)
}
