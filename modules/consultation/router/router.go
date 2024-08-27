package router

import (
	"talkspace-api/middlewares"
	"talkspace-api/modules/consultation/handler"
	"talkspace-api/modules/consultation/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ConsultationRoutes(e *echo.Group, db *gorm.DB) {
	hub := usecase.NewHub()

	consultationWebsocket := handler.NewHandler(hub,db)

	go hub.Run()

	e.POST("/createRoom", consultationWebsocket.CreateRoom, middlewares.JWTMiddleware(false))
	e.GET("/joinRoom/:roomId/:token", consultationWebsocket.JoinRoom)
	e.GET("/getRooms", consultationWebsocket.GetRooms, middlewares.JWTMiddleware(false))
	e.GET("/getDoctors", consultationWebsocket.GetDoctors, middlewares.JWTMiddleware(false))
}
