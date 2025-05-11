package Routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/APIs/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch results."})
		return
	}
	context.JSON(http.StatusOK, events)
}

func postEvents(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Check parameters"})
		return
	}
	userID := context.GetInt64("userID")
	fmt.Printf("user id: %v", userID)
	event.UserID = userID

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not create results."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": event, "Message": "Successful"})

}

func getSingleEvent(context *gin.Context) {
	keyString := context.Param("id")
	eventID, err := strconv.ParseInt(keyString, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot parse the string to int"})
		return
	}

	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch event.", "id": eventID})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": event, "Message": "Successful"})
}

func updateEvent(context *gin.Context) {
	keyString := context.Param("id")
	eventID, err := strconv.ParseInt(keyString, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot parse the string to int", "Error": err})
		return
	}
	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not fetch event.", "id": eventID})
		return
	}
	userID := context.GetInt64("userID")
	if event.UserID != userID {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": "Not authorized to update the user.", "userID": userID})
		return
	}
	var updatedEvent models.Event
	Err := context.ShouldBindJSON(&updatedEvent)

	if Err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Check parameters"})
		return
	}

	updatedEvent.Id = eventID
	err = updatedEvent.UpdateEvent()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot update event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": updatedEvent, "Message": "Successful"})
}

func deleteEvent(context *gin.Context) {
	keyString := context.Param("id")
	eventID, err := strconv.ParseInt(keyString, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot parse the string to int", "Error": err})
		return
	}

	toDeleteEvent, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not fetch event.", "id": eventID})
		return
	}
	userID := context.GetInt64("userID")
	if toDeleteEvent.UserID != userID {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": "Not authorized to delete."})
		return
	}
	err = toDeleteEvent.DeleteEvent()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot delete event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"event": eventID, "Message": "Successful"})
}
