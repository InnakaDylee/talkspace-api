package routes

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	ur "talkspace-api/modules/user/router"
	ar "talkspace-api/modules/admin/router"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB, es *elasticsearch.Client, rdb *redis.Client) {

	user := e.Group("/users")
	admin := e.Group("/admins")
	// doctor := e.Group("/doctors")
	// transaction := e.Group("/transactions")
	// consultation := e.Group("/consultations")
	// chatbot := e.Group("/talkbots")

	ur.UserRoutes(user, db, es, rdb)
	ar.AdminRoutes(admin, db, es, rdb)

	// DoctorRoutes(doctor, db, rdb)
	// TransactionRoutes(transaction, db, rdb)
	// ConsultationRoutes(consultation, db, rdb)
	// ChatbotRoutes(chatbot)

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
