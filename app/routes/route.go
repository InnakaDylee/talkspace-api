package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	ar "talkspace-api/modules/admin/router"
	dr "talkspace-api/modules/doctor/router"
	ur "talkspace-api/modules/user/router"
	tr "talkspace-api/modules/talkbot/router"
	cs "talkspace-api/modules/consultation/router"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB, rdb *redis.Client) {

	user := e.Group("/users")
	admin := e.Group("/admins")
	doctor := e.Group("/doctors")
	talkbot := e.Group("/talkbots")
	consultation := e.Group("/consultations")



	ur.UserRoutes(user, db, rdb)
	dr.DoctorRoutes(doctor, db, rdb)
	ar.AdminRoutes(admin, db, rdb)
	tr.TalkbotRoutes(talkbot, db)
	cs.ConsultationRoutes(consultation, db)


}
