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

// ConnectDatabase initializes and connects to the database.
func ConnectDatabase() error {
	// Load .env file for local development (Railway will inject env vars on deployment)
	_ = godotenv.Load()

	// Use Railway's DATABASE_URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("❌ DATABASE_URL is not set")
	}
	
	log.Println("🚀 Using DATABASE_URL from environment")
	
	dsn := databaseURL
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to connect to database: %w", err)
	}
	
	if err := database.AutoMigrate(&models.User{}, &models.Profile{}, &models.Post{}, &models.Comment{}); err != nil {
		return fmt.Errorf("❌ Migration failed: %w", err)
	}
	
	DB = database
	log.Println("✅ Database connected and migrated successfully")
	return nil
}

// CloseDatabase gracefully closes the DB connection.
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("⚠️ Warning: Unable to retrieve DB instance for closing.")
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("❌ Error closing the database:", err)
	} else {
		log.Println("✅ Database connection closed successfully")
	}
}

