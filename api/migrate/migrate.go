package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"yeha-api/models" // Adjust the import path to your project
)

func main() {
	// Set up the connection to your test database
	user := os.Getenv("YEHA_MYSQL_USER")
	password := os.Getenv("YEHA_MYSQL_PWD")
	host := os.Getenv("YEHA_MYSQL_HOST")
	dbName := os.Getenv("YEHA_MYSQL_DB")

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)

	fmt.Println("Using database:", dbName) // Add this line to log the database name

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Ensure that the tables are created
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Report{},
		&models.Like{},
		&models.Content{},
		&models.Follower{},
		&models.Tag{},
		&models.PostTag{},
	) // Add any other models here as needed

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Check for success
	fmt.Println("Database tables created or migrated successfully.")
}
