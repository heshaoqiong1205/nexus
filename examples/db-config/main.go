package main

import (
	"fmt"
	"log"

	"nexus/pkg/setting"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Testing Database Connection...")
	fmt.Println("=============================")

	// Load configuration
	setting.Setup()

	// Create database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Name,
		setting.DatabaseSetting.Port,
	)

	fmt.Printf("Connecting to: %s@%s:%d/%s\n",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Port,
		setting.DatabaseSetting.Name)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Database connection successful!")

	// Test a simple query
	var result struct {
		Version string
		Now     string
	}

	err = db.Raw("SELECT version() as version, now()::text as now").Scan(&result).Error
	if err != nil {
		log.Fatalf("Failed to execute test query: %v", err)
	}

	fmt.Printf("✓ Database version: %s\n", result.Version[:50]+"...")
	fmt.Printf("✓ Current time: %s\n", result.Now)

	// Check if tables exist
	var tableCount int64
	err = db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name LIKE 'things_%'").Scan(&tableCount).Error
	if err != nil {
		log.Fatalf("Failed to count tables: %v", err)
	}

	fmt.Printf("✓ Found %d tables with 'things_' prefix\n", tableCount)

	// List all things_ tables
	var tables []string
	err = db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name LIKE 'things_%' ORDER BY table_name").Scan(&tables).Error
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	if len(tables) > 0 {
		fmt.Println("\nTables found:")
		for _, table := range tables {
			fmt.Printf("  - %s\n", table)
		}
	}

	fmt.Println("\n✓ Database test completed successfully!")
}
