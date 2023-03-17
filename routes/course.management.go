package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/controller"
)

func CourseRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/courses")
	router.POST("/", controller.CreateCourse)
	router.GET("/", controller.GetAllCourses)
	router.GET("/:courseId", controller.GetCourse)

	// enrollment
	router.GET("/:courseId/students", controller.ListStudentTakenACourse)
}
