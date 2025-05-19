package main

import (
	"net/http"

	"example.com/postgres/DB"
	"github.com/gin-gonic/gin"
)

func main() {
	DB.IntiDB()
	server := gin.Default()
	server.GET("/hello", welcome)
	server.Run(":8000")
}

func welcome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Hello, welcome user",
	})
}
