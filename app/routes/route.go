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
	// transaction := e.Group("/transactions")
	// consultation := e.Group("/consultations")


	ur.UserRoutes(user, db, rdb)
	dr.DoctorRoutes(doctor, db, rdb)
	ar.AdminRoutes(admin, db, rdb)
	tr.TalkbotRoutes(talkbot)
	cs.ConsultationRoutes(consultation, db)

	// DoctorRoutes(doctor, db, rdb)
	// TransactionRoutes(transaction, db, rdb)
	// ConsultationRoutes(consultation, db, rdb)


}

/*
	== user ==
	 https://talkspace.api.id/users/account/register
	 https://talkspace.api.id/users/account/verify-account
	 https://talkspace.api.id/users/account/login

	 https://talkspace.api.id/users/password/forgot-password
	 https://talkspace.api.id/users/password/verify-otp
	 https://talkspace.api.id/users/password/new-password
	 https://talkspace.api.id/users/password/change-password

	 https://talkspace.api.id/users/profile


	== doctor ==
	 https://talkspace.api.id/doctors/account/register
	 https://talkspace.api.id/doctors/account/verify-account
	 https://talkspace.api.id/doctors/account/login

	 https://talkspace.api.id/doctors/password/forgot-password
	 https://talkspace.api.id/doctors/password/verify-otp
	 https://talkspace.api.id/doctors/password/new-password
	 https://talkspace.api.id/doctors/password/change-password

	 https://talkspace.api.id/doctors/profile


	== transaction ==
	 https://talkspace.api.id/transactions
	 https://talkspace.api.id/transactions/:transactions_id


	== consultation ==
	 https://talkspace.api.id/consultations/doctor
	 https://talkspace.api.id/consultations/doctor/:doctor_id

	 https://talkspace.api.id/consultations/roomchat
	 https://talkspace.api.id/consultations/roomchat/:transaction_id
	 https://talkspace.api.id/consultations/roomchat/:roomchat_id

	 https://talkspace.api.id/consultations/message/:roomchat_id


	== chatbot ==
	 https://talkspace.api.id/talkbots

*/
