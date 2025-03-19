package main

import (
	"fmt"
	"log"

	"gitconnect-backend/config"
	"gitconnect-backend/models"
)

func main() {
	// Connect to the database
	err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Run migrations
	err = config.DB.AutoMigrate(&models.User{}, &models.Profile{}) // Add all models here
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	fmt.Println("✅ Database migration completed successfully!")
}

