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
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DB")

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

	// // Drop the foreign key constraint
	// db.Exec("ALTER TABLE post_tags DROP FOREIGN KEY post_tags_ibfk_2;")

	// // Modify the columns to CHAR(36) for UUID compatibility
	// db.Exec("ALTER TABLE post_tags MODIFY COLUMN tag_id CHAR(36) NOT NULL;")
	// db.Exec("ALTER TABLE tags MODIFY COLUMN id CHAR(36) NOT NULL;")

	// // Re-add the foreign key constraint
	// db.Exec("ALTER TABLE post_tags ADD CONSTRAINT post_tags_ibfk_2 FOREIGN KEY (tag_id) REFERENCES tags(id);")

	// Check for success
	fmt.Println("Database tables created or migrated successfully.")
}
