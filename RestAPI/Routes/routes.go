package Routes

import (
	"example.com/APIs/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	server.GET("/events", getEvents)

	authenticationGroup := server.Group("/")
	authenticationGroup.Use(middlewares.Authentication)
	authenticationGroup.POST("/events", postEvents)
	authenticationGroup.PUT("/events/:id", updateEvent)
	authenticationGroup.DELETE("/events/:id", deleteEvent)
	authenticationGroup.POST("/events/:id/register", registerEvent)
	authenticationGroup.DELETE("/events/:id/register", registerCancel)

	server.GET("/events/:id", getSingleEvent)
	server.POST("/signup", signupUser)
	server.POST("/login", login)
}
