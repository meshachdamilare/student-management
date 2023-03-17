package main

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/config"
	"github.com/meshachdamilare/student-management/routes"
	"log"
	"net/http"
)

func main() {
	config.Connection()
	server := gin.Default()

	router := server.Group("/api")
	router.GET("/check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})
	routes.StudentRoutes(router)
	routes.CourseRoutes(router)

	log.Fatal(server.Run(":8080"))
}
