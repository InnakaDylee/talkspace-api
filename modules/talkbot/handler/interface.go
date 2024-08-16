package handler

import "github.com/labstack/echo/v4"

type UserHandlerInterface interface {
	// Query
	GetTalkBot(c echo.Context) error
	GetAllTalkBots(c echo.Context) error

	// Command
	CreateTalkBot(c echo.Context) error
}
