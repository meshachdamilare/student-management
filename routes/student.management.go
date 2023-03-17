package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/controller"
)

func StudentRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/students")
	router.POST("/", controller.CreateStudent)
	router.GET("/", controller.GetAllStudents)
	router.GET("/:studentId", controller.GetStudent)
	router.PUT("/:studentId", controller.UpdateStudent)
	router.DELETE("/:studentId", controller.DeleteStudent)

	// enrollments
	router.PUT("/:studentId/courses", controller.UpdateStudentCourses)
	router.GET("/:studentId/courses", controller.ListCoursesTakenByStudent)
}
