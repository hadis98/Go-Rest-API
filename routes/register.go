package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadis98/rest-api/models"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id"})
		return
	}

	//* we have to be sure the eventId exists
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event"})
		return
	}

	//* add new registration to the registrationTable
	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not register user for event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "register created"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	fmt.Println("[cancelRegistration]: ", userId, "   ", eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	var event models.Event
	event.ID = eventId
	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not cancel registration"})
		return
	}

	fmt.Println("****************************************")
	context.JSON(http.StatusCreated, gin.H{"message": "registration was cancelled!"})

}
