package models

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model struct {
	ID         string `gorm:"primary_key" json:"id"`
	CreatedAt  int    `json:"created_at"`
	ModifiedAt int    `json:"modified_at"`
	Status     bool   `json:"status"`
}

// Setup initializes the database instance
func Setup() {
	var err error
	dsn := "host=localhost user=things password=123456 dbname=things port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}

func MockSetup(conn *sql.DB) {
	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
}
