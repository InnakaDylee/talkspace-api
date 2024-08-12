package databases

import (
	"log"

	"gorm.io/gorm"

	um "talkspace-api/modules/user/model"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		&um.User{},
	)

	migrator := db.Migrator()
	tables := []string{"users"}
	for _, table := range tables {
		if !migrator.HasTable(table) {
			log.Fatalf("table %s was not successfully created", table)
		}
	}
	log.Println("all tables were successfully migrated")
}
