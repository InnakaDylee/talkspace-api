package databases

import (
	"log"

	"gorm.io/gorm"

	am "talkspace-api/modules/admin/model"
	cm "talkspace-api/modules/consultation/model"
	dm "talkspace-api/modules/doctor/model"
	tm "talkspace-api/modules/talkbot/model"
	um "talkspace-api/modules/user/model"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		&um.User{},
		&dm.Doctor{},
		&am.Admin{},
		&cm.Consultation{},
		&cm.Message{},
		&tm.Talkbot{},
	)

	migrator := db.Migrator()
	tables := []string{"users", "admins", "doctors", "consultations", "messages", "talkbots"}
	for _, table := range tables {
		if !migrator.HasTable(table) {
			log.Fatalf("table %s was not successfully created", table)
		}
	}
	log.Println("all tables were successfully migrated")
}
