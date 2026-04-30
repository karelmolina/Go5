package database

import (
	"fmt"

	"github.com/karelmolina/go5/config"
	"github.com/karelmolina/go5/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	sslMode := config.Config("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Config("DB_HOST"),
		config.Config("DB_PORT"),
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
		sslMode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.AutoMigrate(&model.User{}, &model.Event{}, &model.EventResponse{})
	fmt.Println("Database Migrated")
}
