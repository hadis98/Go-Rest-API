package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadis98/rest-api/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized!"})
		return
		// * with AbortWithStatusJSON function, no other request handlers will be executed after it
		// * to make sure upon recieving an error, we stop and no other code on the server runs
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized!"})
		return
	}

	context.Set("userId", userId)
	context.Next()
	//* ensures that the next request handler will be executed correctly

}
