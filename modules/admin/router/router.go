package router

import (
	"talkspace-api/middlewares"
	"talkspace-api/modules/admin/handler"
	"talkspace-api/modules/admin/repository"
	"talkspace-api/modules/admin/usecase"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AdminRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	adminQueryRepository := repository.NewAdminQueryRepository(db, rdb)
	adminCommandRepository := repository.NewAdminCommandRepository(db, rdb)

	adminQueryUsecase := usecase.NewAdminQueryUsecase(adminCommandRepository, adminQueryRepository)
	adminCommandUsecase := usecase.NewAdminCommandUsecase(adminCommandRepository, adminQueryRepository)

	adminHandler := handler.NewAdminHandler(adminCommandUsecase, adminQueryUsecase)

	account := e.Group("/account")
	account.POST("/register", adminHandler.RegisterAdmin)
	account.POST("/login", adminHandler.LoginAdmin)

	password := e.Group("/password")
	password.POST("/forgot-password", adminHandler.ForgotAdminPassword)
	password.POST("/verify-otp", adminHandler.VerifyAdminOTP)
	password.PATCH("/new-password", adminHandler.NewAdminPassword, middlewares.JWTMiddleware(true))
	password.PATCH("/change-password", adminHandler.UpdateAdminPassword, middlewares.JWTMiddleware(true))

	profile := e.Group("/profile", middlewares.JWTMiddleware(false))
	profile.GET("/:admin_id", adminHandler.GetAdminByID)
}
