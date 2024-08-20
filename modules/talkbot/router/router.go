package router

import (
	"talkspace-api/middlewares"
	"talkspace-api/modules/talkbot/handler"
	"talkspace-api/modules/talkbot/usecase"

	"github.com/labstack/echo/v4"
)

func TalkbotRoutes(e *echo.Group) {

	talkbotQueryUsecase := usecase.NewTalkbotQueryUsecase()
	talkbotHandler := handler.NewTalkbotHandler(talkbotQueryUsecase)

	e.POST("", talkbotHandler.CreateTalkBotMessage, middlewares.JWTMiddleware(false))

}
