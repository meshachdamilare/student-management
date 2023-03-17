package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/student-management/config"
	"github.com/meshachdamilare/student-management/model"
	"net/http"
)

var DB = config.GetDB()

func CreateStudent(c *gin.Context) {
	var studentPayload model.StudentRequest

	if err := c.ShouldBindJSON(&studentPayload); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": student})
}
func GetAllStudents(c *gin.Context) {
	var students []model.Student
	if err := DB.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}
	c.JSON(http.StatusOK, students)
}

func GetStudent(c *gin.Context) {
	var student model.Student
	studentId := c.Param("studentId")
	if err := DB.First(&student, "student_id=?", studentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, student)
}

func UpdateStudent(c *gin.Context) {
	var studentPayload model.StudentRequest
	if err := DB.First(&studentPayload, c.Param("studentId")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	if err := c.ShouldBindJSON(&studentPayload); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}
	c.JSON(http.StatusOK, student)
}

func DeleteStudent(c *gin.Context) {
	studentId := c.Param("studentId")
	// Retrieve the student from the database
	var student model.Student
	result := DB.Preload("Courses").First(&student, "student_id = ?", studentId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student"})
		return
	}
	// Delete the student's courses from the database
	for _, course := range student.Courses {
		result := DB.Model(&course).Association("Students").Delete(&student)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
			return
		}
	}
	// Delete the student from the database
	result = DB.Delete(&student, "student_id = ?", studentId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

func UpdateStudentCourses(c *gin.Context) {
	var student model.Student

	var updateCourseRequest model.UpdateCourseRequest
	studentId := c.Param("studentId")

	if err := c.BindJSON(&updateCourseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := DB.Where("student_id = ?", studentId).First(&student).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var courses []model.Course
	if err := DB.Where("course_id IN ?", updateCourseRequest.CourseIDs).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	student.Courses = courses
	if err := DB.Updates(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func ListCoursesTakenByStudent(c *gin.Context) {
	var student model.Student
	studentId := c.Param("studentId")
	if err := DB.Preload("Courses").Find(&student, "student_id", studentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, student.Courses)
}
