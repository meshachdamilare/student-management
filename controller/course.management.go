package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/model"
	"net/http"
)

func CreateCourse(c *gin.Context) {
	var coursePayload model.CourseRequest
	if err := c.ShouldBindJSON(&coursePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
	}
	course := model.Course{
		CourseName:  coursePayload.CourseName,
		Description: coursePayload.Description,
		Instructor:  coursePayload.Instructor,
	}
	result := DB.Create(&course)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": course})
}

func GetAllCourses(c *gin.Context) {
	var courses []model.Course
	if err := DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func GetCourse(c *gin.Context) {
	var course model.Course
	courseId := c.Param("courseId")
	if err := DB.First(&course, "course_id = ?", courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func ListStudentTakenACourse(c *gin.Context) {
	var course model.Course
	id := c.Param("courseId")
	if err := DB.Preload("Students").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course.Students)
}
