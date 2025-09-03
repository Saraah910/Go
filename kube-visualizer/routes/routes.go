package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/", homepage)
	server.POST("/upload", uploadFile)
	server.POST("/health", nodeResourceUsage)
	server.GET("/data", getNamespaces)
	server.POST("/set-kubeconfig", setKubeConfig)
	server.GET("/get-graph-data/:ns", graphData)
}
