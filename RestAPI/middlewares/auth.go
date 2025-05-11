package middlewares

import (
	"fmt"
	"net/http"

	"example.com/APIs/Utils"
	"github.com/gin-gonic/gin"
)

func Authentication(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Authentication error. Empty token"})
		return
	}
	userID, err := Utils.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": err})
		return
	}
	fmt.Printf("The user id is: %d", userID)
	context.Set("userID", userID)
	context.Next()
}
