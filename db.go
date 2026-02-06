package main

import (
	"log"
	
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Podcast struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string
	URL         string `gorm:"uniqueIndex"`
	LastUpdated string
	ETag        string
	Episodes    []Episode `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Episode struct {
	ID            uint   `gorm:"primaryKey"`
	PodcastID     uint
	GUID          string `gorm:"index"`
	Title         string
	PublishedDate string
	EnclosureURL  string
}

var DB *gorm.DB

func InitDB(dbPath string) {
	var err error
	if dbPath == "" {
		dbPath = "podcast2m3u.db"
	}
	
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schema
	err = DB.AutoMigrate(&Podcast{}, &Episode{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}
