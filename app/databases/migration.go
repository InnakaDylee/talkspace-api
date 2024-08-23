package databases

import (
	"log"

	"gorm.io/gorm"

	am "talkspace-api/modules/admin/model"
	dm "talkspace-api/modules/doctor/model"
	um "talkspace-api/modules/user/model"
	cm "talkspace-api/modules/consultation/model"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		&um.User{},
		&dm.Doctor{},
		&am.Admin{},
		&cm.Consultation{},
		&cm.Message{},
	)

	migrator := db.Migrator()
	tables := []string{"users", "admins", "doctors"} // "chatbots", "consultations"
	for _, table := range tables {
		if !migrator.HasTable(table) {
			log.Fatalf("table %s was not successfully created", table)
		}
	}
	log.Println("all tables were successfully migrated")
}
