package config

import (
	"fmt"
	"log"
	"os"

	"gitconnect-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

// ConnectDatabase initializes and connects to the database
func ConnectDatabase() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found, using system environment variables.")
	}

	// Get database environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Check if any variable is empty
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		return fmt.Errorf("❌ Missing required database environment variables")
	}

	// Database connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	// Connect to the database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("❌ Failed to connect to database: %w", err)
	}

	// 🔹 Auto-migrate tables (Ensure all necessary tables exist)
	if err := database.AutoMigrate(&models.User{}, &models.Profile{}, &models.Post{}, &models.Comment{}); err != nil {
		return fmt.Errorf("❌ Migration failed: %w", err)
	}

	DB = database
	log.Println("✅ Database connected and migrated successfully")
	return nil
}

