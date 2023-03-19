package main

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/config"
	_ "github.com/meshachdamilare/student-management/docs"
	"github.com/meshachdamilare/student-management/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

// @title Student Management API
// @description API for managing students and courses.
// @version 1
// @host localhost:8001
// @BasePath /api
// @schemes http

func main() {
	config.Connection()
	server := gin.Default()

	router := server.Group("/api")
	router.GET("/check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})
	routes.StudentRoutes(router)
	routes.CourseRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(server.Run(":8080"))
}
