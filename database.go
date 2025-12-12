package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection and runs migrations
func InitDatabase() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed")

	// Seed initial data
	if err := seedData(); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}
}

// runMigrations creates all necessary tables
func runMigrations() error {
	return DB.AutoMigrate(
		&User{},
		&Session{},
		&Product{},
		&CartItem{},
		&Order{},
		&OrderItem{},
	)
}

// seedData adds initial products to the database
func seedData() error {
	// Check if products already exist
	var count int64
	DB.Model(&Product{}).Count(&count)
	if count > 0 {
		log.Println("Products already seeded, skipping...")
		return nil
	}

	products := []Product{
		{
			ID:          "1",
			Name:        "Premium Headphones",
			Description: "High-quality wireless headphones with noise cancellation",
			Price:       29900,
			ImageURL:    "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400",
		},
		{
			ID:          "2",
			Name:        "Smart Watch",
			Description: "Fitness tracking smartwatch with heart rate monitor",
			Price:       19900,
			ImageURL:    "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400",
		},
		{
			ID:          "3",
			Name:        "Laptop Stand",
			Description: "Ergonomic aluminum laptop stand for better posture",
			Price:       4900,
			ImageURL:    "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=400",
		},
		{
			ID:          "4",
			Name:        "Mechanical Keyboard",
			Description: "RGB mechanical keyboard with Cherry MX switches",
			Price:       12900,
			ImageURL:    "https://images.unsplash.com/photo-1511467687858-23d96c32e4ae?w=400",
		},
	}

	result := DB.Create(&products)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Seeded %d products", len(products))
	return nil
}

// CleanupExpiredSessions removes expired JWT sessions from database
func CleanupExpiredSessions() error {
	return DB.Where("expires_at < ?", time.Now()).Delete(&Session{}).Error
}
