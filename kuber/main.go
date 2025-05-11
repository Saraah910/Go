package main

import (
	"example.com/kuber/db"
	"example.com/kuber/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()
	routes.Routes(server)
	server.Run("localhost:8000")
}
