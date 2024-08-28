package handler

import "github.com/labstack/echo/v4"

type UserHandlerInterface interface {
	// Query
	GetUserByID(c echo.Context) error
	GetRequestPremiumUsers(c echo.Context) error

	// Command
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
	UpdateUserByID(c echo.Context) error
	UpdateUserPassword(c echo.Context) error
	ForgotUserPassword(c echo.Context) error
	NewUserPassword(c echo.Context) error
	VerifyUserOTP(c echo.Context) error
	RequestPremium(c echo.Context) error
	UpdateUserPremiumExpired(c echo.Context) error
}
