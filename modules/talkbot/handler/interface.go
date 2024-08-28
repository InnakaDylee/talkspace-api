package handler

import "github.com/labstack/echo/v4"

type UserHandlerInterface interface {
	// Command
	CreateTalkBot(c echo.Context) error
}
