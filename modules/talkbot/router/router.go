package router

import (
	"talkspace-api/middlewares"
	"talkspace-api/modules/talkbot/handler"
	"talkspace-api/modules/talkbot/repository"
	"talkspace-api/modules/talkbot/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func TalkbotRoutes(e *echo.Group, db *gorm.DB) {
	talkbotQueryRepository := repository.NewTalkbotQueryRepository(db)
	talkbotCommandRepository := repository.NewTalkbotCommandRepository(db)

	talkbotQueryUsecase := usecase.NewTalkbotQueryUsecase(talkbotCommandRepository, talkbotQueryRepository)
	
	talkbotHandler := handler.NewTalkbotHandler(talkbotQueryUsecase)

	e.POST("", talkbotHandler.CreateTalkBotMessage, middlewares.JWTMiddleware(false))

}
