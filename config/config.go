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
	// Load environment variables from .env file (for local development)
	_ = godotenv.Load()

	// Check for Railway's DATABASE_URL first.
	databaseURL := os.Getenv("DATABASE_URL")
	var dsn string

	if databaseURL != "" {
		dsn = databaseURL
		log.Println("🚀 Using DATABASE_URL from environment")
	} else {
		// Fallback to local database configuration
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")
		dbSSLMode := os.Getenv("DB_SSLMODE")
		if dbSSLMode == "" {
			dbSSLMode = "disable"
		}

		missingVars := []string{}
		if dbHost == "" { missingVars = append(missingVars, "DB_HOST") }
		if dbUser == "" { missingVars = append(missingVars, "DB_USER") }
		if dbPassword == "" { missingVars = append(missingVars, "DB_PASSWORD") }
		if dbName == "" { missingVars = append(missingVars, "DB_NAME") }
		if dbPort == "" { missingVars = append(missingVars, "DB_PORT") }
		if len(missingVars) > 0 {
			return fmt.Errorf("❌ Missing required database environment variables: %v", missingVars)
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)
		log.Println("🔧 Using local database configuration")
	}

	// Connect to the database with connection pooling
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to connect to database: %w", err)
	}

	// Auto-migrate tables (ensure necessary tables exist)
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

