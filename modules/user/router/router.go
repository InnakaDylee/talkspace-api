package router

import (
	"talkspace-api/middlewares"
	"talkspace-api/modules/user/handler"
	"talkspace-api/modules/user/repository"
	"talkspace-api/modules/user/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserRoutes(e *echo.Group, db *gorm.DB) {
	userQueryRepository := repository.NewUserQueryRepository(db)
	userCommandRepository := repository.NewUserCommandRepository(db)

	userQueryUsecase := usecase.NewUserQueryUsecase(userCommandRepository, userQueryRepository)
	userCommandUsecase := usecase.NewUserCommandUsecase(userCommandRepository, userQueryRepository)

	userHandler := handler.NewUserHandler(userCommandUsecase, userQueryUsecase)

	account := e.Group("/account")
	account.POST("register", userHandler.RegisterUser)
	account.GET("verify-account", userHandler.VerifyUser)
	account.POST("login", userHandler.LoginUser)

	password := e.Group("/password")
	password.POST("forgot-password", userHandler.ForgotUserPassword)
	password.POST("verify-otp", userHandler.VerifyUserOTP)
	password.PATCH("new-password", userHandler.NewUserPassword, middlewares.JWTMiddleware())
	password.PATCH("/change-password", userHandler.UpdateUserPassword, middlewares.JWTMiddleware())

	profile := e.Group("/profile", middlewares.JWTMiddleware())
	profile.GET("", userHandler.GetUserByID)
	profile.PUT("", userHandler.UpdateUserByID)
}
