package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadis98/rest-api/models"
	"github.com/hadis98/rest-api/utils"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user) //we want to fill the user data with the data that is attached to the incoming request(json)
	// the incoming request should contain both email and password
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	// validate the sent credentilas that are attached to the sent request
	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not authenticate user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})

}

func getUsers(context *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch users. Try again later"})
		return
	}
	// context.JSON(http.StatusOK, gin.H{"message": "Hello"})
	context.JSON(http.StatusOK, users)
}
