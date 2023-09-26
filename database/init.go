package database

import (
	"github.com/Quadra-hub/go-chatgpt/config"
	"github.com/Quadra-hub/go-chatgpt/gpt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Start() {
	urlConn := config.DBConnString()
	db, _ = gorm.Open(postgres.Open(urlConn), &gorm.Config{})
	db.AutoMigrate(&gpt.Embeddding{})
}

func GetDB() *gorm.DB {
	return db
}
