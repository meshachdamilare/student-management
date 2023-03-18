package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func Connection() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DB_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	DB_USER := os.Getenv("MYSQL_USER")
	DB_NAME := os.Getenv("MYSQL_DATABASE")
	DB_PORT := os.Getenv("MYSQL_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(mysql:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_USER, DB_PASSWORD, DB_PORT, DB_NAME)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}
func GetDB() *gorm.DB {
	return db
}
