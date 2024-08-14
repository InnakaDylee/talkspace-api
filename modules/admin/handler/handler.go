package handler

import (
	"net/http"
	"strings"
	"talkspace-api/middlewares"
	"talkspace-api/modules/admin/dto"
	"talkspace-api/modules/admin/usecase"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/responses"

	"github.com/labstack/echo/v4"
)

type adminHandler struct {
	adminCommandUsecase usecase.AdminCommandUsecaseInterface
	adminQueryUsecase   usecase.AdminQueryUsecaseInterface
}

func NewAdminHandler(acu usecase.AdminCommandUsecaseInterface, aqu usecase.AdminQueryUsecaseInterface) *adminHandler {
	return &adminHandler{
		adminCommandUsecase: acu,
		adminQueryUsecase:   aqu,
	}
}

// Query
func (ah *adminHandler) GetAdminByID(c echo.Context) error {
	adminIDParam := c.Param("admin_id")
	if adminIDParam == "" {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_ID_NOTFOUND))
	}

	tokenAdminID, role, errExtract := middlewares.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	if role != constant.ADMIN {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if adminIDParam != tokenAdminID {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	admin, errGetID := ah.adminQueryUsecase.GetAdminByID(adminIDParam)
	if errGetID != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errGetID.Error()))
	}

	adminResponse := dto.AdminEntityToAdminResponse(admin)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PROFILE_RETRIEVED, adminResponse))
}

// Command
func (ah *adminHandler) RegisterAdmin(c echo.Context) error {
	adminRequest := dto.AdminRegisterRequest{}

	errBind := c.Bind(&adminRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	adminEntity := dto.AdminRegisterRequestToAdminEntity(adminRequest)

	registeredAdmin, errRegister := ah.adminCommandUsecase.RegisterAdmin(adminEntity)
	if errRegister != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errRegister.Error()))
	}

	adminResponse := dto.AdminEntityToAdminRegisterResponse(registeredAdmin)

	return c.JSON(http.StatusCreated, responses.SuccessResponse(constant.SUCCESS_REGISTER, adminResponse))
}

func (ah *adminHandler) LoginAdmin(c echo.Context) error {
	adminRequest := dto.AdminLoginRequest{}

	errBind := c.Bind(&adminRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	loggedInAdmin, token, errLogin := ah.adminCommandUsecase.LoginAdmin(adminRequest.Email, adminRequest.Password)
	if errLogin != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errLogin.Error()))
	}

	adminResponse := dto.AdminEntityToAdminLoginResponse(loggedInAdmin, token)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_LOGIN, adminResponse))
}

func (ah *adminHandler) UpdateAdminPassword(c echo.Context) error {
	adminRequest := dto.AdminUpdatePasswordRequest{}

	errBind := c.Bind(&adminRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	adminID, role, errExtractToken := middlewares.ExtractToken(c)

	if role != constant.ADMIN {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if errExtractToken != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtractToken.Error()))
	}

	adminEntity := dto.AdminUpdatePasswordRequestToAdminEntity(adminRequest)

	password, errUpdate := ah.adminCommandUsecase.UpdateAdminPassword(adminID, adminEntity)
	if errUpdate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errUpdate.Error()))
	}

	adminResponse := dto.AdminEntityToAdminResponse(password)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PASSWORD_UPDATED, adminResponse))
}

func (ah *adminHandler) ForgotAdminPassword(c echo.Context) error {
	adminRequest := dto.AdminSendOTPRequest{}

	errBind := c.Bind(&adminRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	adminEntity := dto.AdminSendOTPRequestToAdminEntity(adminRequest)

	otp, errSendOTP := ah.adminCommandUsecase.SendAdminOTP(adminEntity.Email)
	if errSendOTP != nil {
		if strings.Contains(errSendOTP.Error(), constant.ERROR_EMAIL_NOTFOUND) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse(errSendOTP.Error()))
		}
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errSendOTP.Error()))
	}

	adminResponse := dto.AdminEntityToAdminResponse(otp)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_OTP_SENT, adminResponse))
}

func (ah *adminHandler) VerifyAdminOTP(c echo.Context) error {
	adminRequest := dto.AdminVerifyOTPRequest{}

	errBind := c.Bind(&adminRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	adminEntity := dto.AdminVerifyOTPRequestToAdminEntity(adminRequest)

	token, errVerify := ah.adminCommandUsecase.VerifyAdminOTP(adminEntity.Email, adminEntity.OTP)
	if errVerify != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse(constant.ERROR_OTP_VERIFY+errVerify.Error()))
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_OTP_VERIFIED, token))
}

func (ah *adminHandler) NewAdminPassword(c echo.Context) error {
	adminRequest := dto.AdminNewPasswordRequest{}

	errBind := c.Bind(&adminRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	email, errExtract := middlewares.ExtractVerifyToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	adminEntity := dto.AdminNewPasswordRequestToAdminEntity(adminRequest)

	password, errCreate := ah.adminCommandUsecase.NewAdminPassword(email, adminEntity)
	if errCreate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errCreate.Error()))
	}

	adminResponse := dto.AdminEntityToAdminResponse(password)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PASSWORD_UPDATED, adminResponse))
}
