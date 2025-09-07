package models

import (
	"database/sql"
	"fmt"
	"log"
	"nexus/pkg/setting"
	"strings"

	"github.com/google/uuid"
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
func Setup(config *setting.Database) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host, config.User, config.Password, config.Name, config.Port)
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

func GetDB() *gorm.DB {
	return db
}

func GenerateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GenerateDeviceID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GenerateSecretKey() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
