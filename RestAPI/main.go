package main

import (
	"example.com/APIs/DB"
	"example.com/APIs/Routes"
	"github.com/gin-gonic/gin"
)

func main() {
	DB.InitDB()
	server := gin.Default()
	Routes.RegisterRoutes(server)

	server.Run("localhost:3000")
}
