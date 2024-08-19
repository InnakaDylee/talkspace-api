package databases

import (
	"log"

	"gorm.io/gorm"

	um "talkspace-api/modules/user/model"
	// dm "talkspace-api/modules/doctor/model"
	// am "talkspace-api/modules/admin/model"
	// tm "talkspace-api/modules/transaction/model"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		&um.User{},
		// &dm.Doctor{},
		// &am.Admin{},
		// &tm.Transaction{},
	)

	migrator := db.Migrator()
	tables := []string{"users"} // "chatbots", "consultations"
	for _, table := range tables {
		if !migrator.HasTable(table) {
			log.Fatalf("table %s was not successfully created", table)
		}
	}
	log.Println("all tables were successfully migrated")
}