package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/config"
	"github.com/meshachdamilare/student-management/model"
	"net/http"
)

var DB = config.GetDB()

// @Summary CreateStudent
// @Description Create a student with the given details
// @Accept  json
// @Produce  json
// @Param payload body model.StudentRequest true "Student details"
// @Success 200 {object} model.Student
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students [post]
func CreateStudent(c *gin.Context) {
	var studentPayload model.StudentRequest
	var errorResponse model.ErrorResponse

	if err := c.ShouldBindJSON(&studentPayload); err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
	}
	student := model.Student{
		FirstName:   studentPayload.FirstName,
		LastName:    studentPayload.LastName,
		Email:       studentPayload.Email,
		PhoneNumber: studentPayload.PhoneNumber,
	}
	result := DB.Create(&student)
	if result.Error != nil {
		errorResponse.Message = result.Error.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": student})
}

// @Summary GetAllStudents
// @Description Get a list of all students
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Student
// @Failure 500 {object} model.ErrorResponse
// @Router /students [get]
func GetAllStudents(c *gin.Context) {
	var students []model.Student
	var errorResponse model.ErrorResponse
	if err := DB.Find(&students).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}
	c.JSON(http.StatusOK, students)
}

// @Summary GetStudent
// @Description Get the details of a student with the given ID
// @Produce  json
// @Param studentId path string true "Student ID"
// @Success 200 {object} model.Student
// @Failure 404 {object} model.ErrorResponse
// @Router /students/{studentId} [get]
func GetStudent(c *gin.Context) {
	var student model.Student
	var errorResponse model.ErrorResponse
	studentId := c.Param("studentId")
	if err := DB.First(&student, "student_id=?", studentId).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, student)
}

// @Summary UpdateStudent
// @Description Update the details of a student with the given ID
// @Accept  json
// @Produce  json
// @Param studentId path string true "Student ID"
// @Param payload body model.StudentRequest true "Student details"
// @Success 200 {object} model.Student
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students/{studentId} [put]
func UpdateStudent(c *gin.Context) {
	var studentPayload model.StudentRequest
	var errorResponse model.ErrorResponse
	if err := DB.First(&studentPayload, c.Param("studentId")).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	if err := c.ShouldBindJSON(&studentPayload); err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student data"})
		return
	}

	student := model.Student{
		FirstName:   studentPayload.FirstName,
		LastName:    studentPayload.LastName,
		Email:       studentPayload.Email,
		PhoneNumber: studentPayload.PhoneNumber,
	}
	if err := DB.Save(&student).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}
	c.JSON(http.StatusOK, student)
}

// @Summary DeleteStudent
// @Description Delete a student with the given ID
// @Accept  json
// @Produce  json
// @Param studentId path string true "Student ID"
// @Success 200 {object} model.DeleteResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students/{studentId} [delete]
func DeleteStudent(c *gin.Context) {
	studentId := c.Param("studentId")
	var errorResponse model.ErrorResponse
	// Retrieve the student from the database
	var student model.Student

	result := DB.Preload("Courses").First(&student, "student_id = ?", studentId)
	if result.Error != nil {
		errorResponse.Message = result.Error.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student"})
		return
	}
	// Delete the student's courses from the database
	for _, course := range student.Courses {
		result := DB.Model(&course).Association("Students").Delete(&student)
		if result.Error != nil {
			errorResponse.Message = result.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
			return
		}
	}
	// Delete the student from the database
	result = DB.Delete(&student, "student_id = ?", studentId)
	if result.Error != nil {
		errorResponse.Message = result.Error.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete student"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

// @Summary UpdateStudentCourses
// @Description Update a student with the given ID courses
// @Accept  json
// @Produce  json
// @Param studentId path string true "Student ID"
// @Param payload body model.UpdateCourseRequest true "Courses IDs"
// @Success 200 {object} model.Student
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students/{studentId}/courses [put]
func UpdateStudentCourses(c *gin.Context) {
	var student model.Student
	var errorResponse model.ErrorResponse

	var updateCourseRequest model.UpdateCourseRequest
	studentId := c.Param("studentId")

	if err := c.BindJSON(&updateCourseRequest); err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := DB.Where("student_id = ?", studentId).First(&student).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var courses []model.Course
	if err := DB.Where("course_id IN ?", updateCourseRequest.CourseIDs).Find(&courses).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	student.Courses = courses
	if err := DB.Updates(&student).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// @Summary ListCoursesTakenByStudent
// @Description Get a list of all courses taken by student
// @Produce  json
// @Param studentId path string true "Student ID"
// @Success 200 {array} model.Course
// @Failure 500 {object} model.ErrorResponse
// @Router /students/{studentId}/courses [get]
func ListCoursesTakenByStudent(c *gin.Context) {
	var student model.Student
	var errorResponse model.ErrorResponse
	studentId := c.Param("studentId")
	if err := DB.Preload("Courses").Find(&student, "student_id", studentId).Error; err != nil {
		errorResponse.Message = err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, student.Courses)
}
