package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hadis98/rest-api/db"
	"github.com/hadis98/rest-api/routes"
)

func main() {

	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8081") //localhost:8081
}
