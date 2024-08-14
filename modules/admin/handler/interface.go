package handler

import "github.com/labstack/echo/v4"

type AdminHandlerInterface interface {
	// Query
	GetAdminByID(c echo.Context) error

	// Command
	RegisterAdmin(c echo.Context) error
	LoginAdmin(c echo.Context) error
	UpdateAdminPassword(c echo.Context) error
	ForgotAdminPassword(c echo.Context) error
	NewAdminPassword(c echo.Context) error
	VerifyAdminOTP(c echo.Context) error
}
