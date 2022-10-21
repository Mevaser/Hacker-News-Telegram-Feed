package db

import (
	"hn_feed/config"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	log.Printf("[*] Initializing %s.", config.Config.SqliteDb)

	var err error
	DB, err = gorm.Open(sqlite.Open(config.Config.SqliteDb), &gorm.Config{})
	if err != nil {
		log.Fatal("[!] Failed to connect to database.")
	}
	DB.AutoMigrate(
		&Channel{},
		&Post{},
		&Keyword{},
	)
}
