package handler

import "github.com/labstack/echo/v4"

type DoctorHandlerInterface interface {
	// Query
	GetDoctorByID(c echo.Context) error
	GetAllDoctors(c echo.Context) error

	// Command
	LoginDoctor(c echo.Context) error
	UpdateDoctorProfile(c echo.Context) error
	UpdateDoctorStatus(c echo.Context) error
	UpdateDoctorPassword(c echo.Context) error
	ForgotDoctorPassword(c echo.Context) error
	NewDoctorPassword(c echo.Context) error
	VerifyDoctorOTP(c echo.Context) error
}
