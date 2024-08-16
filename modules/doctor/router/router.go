package router

import (
	"talkspace-api/middlewares"
	"talkspace-api/modules/doctor/handler"
	"talkspace-api/modules/doctor/repository"
	"talkspace-api/modules/doctor/usecase"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DoctorRoutes(e *echo.Group, db *gorm.DB, es *elasticsearch.Client, rdb *redis.Client) {
	doctorQueryRepository := repository.NewDoctorQueryRepository(db, es, rdb)
	doctorCommandRepository := repository.NewDoctorCommandRepository(db, es, rdb)

	doctorQueryUsecase := usecase.NewDoctorQueryUsecase(doctorCommandRepository, doctorQueryRepository)
	doctorCommandUsecase := usecase.NewDoctorCommandUsecase(doctorCommandRepository, doctorQueryRepository)

	doctorHandler := handler.NewDoctorHandler(doctorCommandUsecase, doctorQueryUsecase)

	account := e.Group("/account")
	account.POST("/login", doctorHandler.LoginDoctor)

	password := e.Group("/password")
	password.POST("/forgot-password", doctorHandler.ForgotDoctorPassword)
	password.POST("/verify-otp", doctorHandler.VerifyDoctorOTP)
	password.PATCH("/new-password", doctorHandler.NewDoctorPassword, middlewares.JWTMiddleware(true))
	password.PATCH("/change-password", doctorHandler.UpdateDoctorPassword, middlewares.JWTMiddleware(true))

	profile := e.Group("/profile", middlewares.JWTMiddleware(false))
	profile.GET("/:doctor_id", doctorHandler.GetDoctorByID)
	profile.PUT("/:doctor_id", doctorHandler.UpdateDoctorProfile)

	status := e.Group("/status", middlewares.JWTMiddleware(false))
	status.PUT("/:doctor_id", doctorHandler.UpdateDoctorStatus)

	e.GET("", doctorHandler.GetAllDoctors)
}
