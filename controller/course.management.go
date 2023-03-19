package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/model"
	"net/http"
)

// @Summary CreateCourse
// @Description Create a course with the given details
// @Accept  json
// @Produce  json
// @Param payload body model.CourseRequest true "Course details"
// @Success 200 {object} model.Course
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /courses [post]
func CreateCourse(c *gin.Context) {
	var coursePayload model.CourseRequest
	var errorResponse model.ErrorResponse
	if err := c.ShouldBindJSON(&coursePayload); err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
	}
	course := model.Course{
		CourseName:  coursePayload.CourseName,
		Description: coursePayload.Description,
		Instructor:  coursePayload.Instructor,
	}
	result := DB.Create(&course)
	if result.Error != nil {
		errorResponse.Message = result.Error.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": course})
}

// @Summary GetAllCourses
// @Description Get a list of all courses
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Course
// @Failure 500 {object} model.ErrorResponse
// @Router /courses [get]
func GetAllCourses(c *gin.Context) {
	var courses []model.Course
	var errorResponse model.ErrorResponse
	if err := DB.Find(&courses).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// @Summary GetCourse
// @Description Get the details of a course with the given ID
// @Produce  json
// @Param courseId path string true "Course ID"
// @Success 200 {object} model.Course
// @Failure 404 {object} model.ErrorResponse
// @Router /courses/{courseId} [get]
func GetCourse(c *gin.Context) {
	var course model.Course
	courseId := c.Param("courseId")
	if err := DB.First(&course, "course_id = ?", courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

// @Summary ListStudentTakenACourse
// @Description Get a list of all students taken a course
// @Produce  json
// @Param courseId path string true "Course ID"
// @Success 200 {array} model.Student
// @Failure 500 {object} model.ErrorResponse
// @Router /courses/{courseId}/students [get]
func ListStudentTakenACourse(c *gin.Context) {
	var course model.Course
	var errorResponse model.ErrorResponse
	id := c.Param("courseId")
	if err := DB.Preload("Students").First(&course, id).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course.Students)
}
