package router

import (
	"talkspace-api/middlewares"
	"talkspace-api/modules/user/handler"
	"talkspace-api/modules/user/repository"
	"talkspace-api/modules/user/usecase"

	"github.com/redis/go-redis/v9"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	userQueryRepository := repository.NewUserQueryRepository(db, rdb)
	userCommandRepository := repository.NewUserCommandRepository(db, rdb)

	userQueryUsecase := usecase.NewUserQueryUsecase(userCommandRepository, userQueryRepository)
	userCommandUsecase := usecase.NewUserCommandUsecase(userCommandRepository, userQueryRepository)

	userHandler := handler.NewUserHandler(userCommandUsecase, userQueryUsecase)

	account := e.Group("/account")
	account.POST("/register", userHandler.RegisterUser)
	account.POST("/login", userHandler.LoginUser)

	password := e.Group("/password")
	password.POST("/forgot-password", userHandler.ForgotUserPassword)
	password.POST("/verify-otp", userHandler.VerifyUserOTP)
	password.PATCH("/new-password", userHandler.NewUserPassword, middlewares.JWTMiddleware(true))
	password.PATCH("/change-password", userHandler.UpdateUserPassword, middlewares.JWTMiddleware(true))

	profile := e.Group("/profile", middlewares.JWTMiddleware(false))
	profile.GET("/:user_id", userHandler.GetUserByID)
	profile.PUT("/:user_id", userHandler.UpdateUserProfile)
}
