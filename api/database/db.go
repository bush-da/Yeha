package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes and connects to the database
func ConnectDatabase() (*gorm.DB, error) {
	// Fetch database credentials from environment variables
	user := os.Getenv("YEHA_MYSQL_USER")
	password := os.Getenv("YEHA_MYSQL_PWD")
	host := os.Getenv("YEHA_MYSQL_HOST")
	dbName := os.Getenv("YEHA_MYSQL_DB")

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)

	// Open the database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err // Return error if connection fails
	}

	log.Println("MySQL database connected successfully!")
	DB = db // Store DB connection globally

	return db, nil // Return the DB instance
}
