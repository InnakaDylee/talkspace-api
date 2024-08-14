package databases

import (
	"log"

	"gorm.io/gorm"

	um "talkspace-api/modules/user/model"
	dm "talkspace-api/modules/doctor/model"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		&um.User{},
		&dm.Doctor{},
	)

	migrator := db.Migrator()
	tables := []string{"users", "doctors", "admins", "consultations", "transactions", "chatbots"}
	for _, table := range tables {
		if !migrator.HasTable(table) {
			log.Fatalf("table %s was not successfully created", table)
		}
	}
	log.Println("all tables were successfully migrated")
}