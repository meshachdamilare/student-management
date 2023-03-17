package model

import (
	"github.com/meshachdamilare/student-management/config"
)

func init() {
	config.Connection()
	db := config.GetDB()
	db.AutoMigrate(&Student{}, &Course{})
}

type Student struct {
	StudentID   uint     `gorm:"primaryKey"`
	FirstName   string   `json:"first_name" gorm:"not null"`
	LastName    string   `json:"last_name" gorm:"not null"`
	Email       string   `json:"email" gorm:"not null"`
	PhoneNumber string   `json:"phone_number" gorm:"not null""`
	Courses     []Course `gorm:"many2many:student_courses"`
}

type Course struct {
	CourseID    uint      `gorm:"primaryKey"`
	CourseName  string    `json:"course_name" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Instructor  string    `json:"instructor" gorm:"not null"`
	Students    []Student `gorm:"many2many:student_courses"`
}

type StudentRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateCourseRequest struct {
	CourseIDs []uint `json:"courseIDs"`
}

type CourseRequest struct {
	CourseName  string `json:"course_name"`
	Description string `json:"description"`
	Instructor  string `json:"instructor"`
}
