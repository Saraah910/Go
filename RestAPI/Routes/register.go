package Routes

import (
	"net/http"
	"strconv"

	"example.com/APIs/models"
	"github.com/gin-gonic/gin"
)

func registerEvent(context *gin.Context) {
	userID := context.GetInt64("userID")
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 14)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event ID"})
	}

	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch event.", "id": eventID})
		return
	}
	err = event.RegisterEvent(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not register event.", "id": eventID})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"Message": "Registered event successfully", "Event": event})
}

func registerCancel(context *gin.Context) {
	userID := context.GetInt64("userID")
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 14)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event ID"})
	}

	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch event.", "id": eventID})
		return
	}
	err = event.CancelEvent(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not delete event.", "id": eventID})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"Message": "Deleted event successfully", "Event": event})
}
