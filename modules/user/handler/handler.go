package handler

import (
	"net/http"
	"strings"
	"talkspace-api/middlewares"
	"talkspace-api/modules/user/dto"
	"talkspace-api/modules/user/usecase"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/responses"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userCommandUsecase usecase.UserCommandUsecaseInterface
	userQueryUsecase   usecase.UserQueryUsecaseInterface
}

func NewUserHandler(ucu usecase.UserCommandUsecaseInterface, uqu usecase.UserQueryUsecaseInterface) *userHandler {
	return &userHandler{
		userCommandUsecase: ucu,
		userQueryUsecase:   uqu,
	}
}

// Query
func (uh *userHandler) GetUserByID(c echo.Context) error {
	userIDParam := c.Param("user_id")
	if userIDParam == "" {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_ID_NOTFOUND))
	}

	tokenUserID, role, errExtract := middlewares.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	if role != constant.USER {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if userIDParam != tokenUserID {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	user, errGetID := uh.userQueryUsecase.GetUserByID(userIDParam)
	if errGetID != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errGetID.Error()))
	}

	userResponse := dto.UserEntityToUserUpdateResponse(user)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PROFILE_RETRIEVED, userResponse))
}

// Command
func (uh *userHandler) RegisterUser(c echo.Context) error {
	userRequest := dto.UserRegisterRequest{}

	errBind := c.Bind(&userRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	userEntity := dto.UserRegisterRequestToUserEntity(userRequest)

	registeredUser, errRegister := uh.userCommandUsecase.RegisterUser(userEntity)
	if errRegister != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errRegister.Error()))
	}

	userResponse := dto.UserEntityToUserRegisterResponse(registeredUser)

	return c.JSON(http.StatusCreated, responses.SuccessResponse(constant.SUCCESS_REGISTER, userResponse))
}

func (uh *userHandler) LoginUser(c echo.Context) error {
	userRequest := dto.UserLoginRequest{}

	errBind := c.Bind(&userRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	LoginUser, token, errLogin := uh.userCommandUsecase.LoginUser(userRequest.Email, userRequest.Password)
	if errLogin != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errLogin.Error()))
	}

	userResponse := dto.UserEntityToUserLoginResponse(LoginUser, token)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_LOGIN, userResponse))
}

func (uh *userHandler) UpdateUserByID(c echo.Context) error {
	userIDParam := c.Param("user_id")
	if userIDParam == "" {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_ID_NOTFOUND))
	}

	tokenUserID, role, errExtract := middlewares.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	if role != constant.USER {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if userIDParam != tokenUserID {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	userRequest := dto.UserUpdateRequest{}

	errBind := c.Bind(&userRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	userEntity := dto.UserUpdateRequestToUserEntity(userRequest)

	user, errUpdate := uh.userCommandUsecase.UpdateUserByID(userIDParam, userEntity)
	if errUpdate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errUpdate.Error()))
	}

	userResponse := dto.UserEntityToUserUpdateResponse(user)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PROFILE_UPDATED, userResponse))
}

func (uh *userHandler) UpdateUserPassword(c echo.Context) error {
	userRequest := dto.UserUpdatePasswordRequest{}

	errBind := c.Bind(&userRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	userID, role, errExtractToken := middlewares.ExtractToken(c)

	if role != constant.USER {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if errExtractToken != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtractToken.Error()))
	}

	userEntity := dto.UserUpdatePasswordRequestToUserEntity(userRequest)

	password, errUpdate := uh.userCommandUsecase.UpdateUserPassword(userID, userEntity)
	if errUpdate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errUpdate.Error()))
	}

	userResponse := dto.UserEntityToUserResponse(password)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PASSWORD_UPDATED, userResponse))
}

func (uh *userHandler) ForgotUserPassword(c echo.Context) error {
	userRequest := dto.UserSendOTPRequest{}

	errBind := c.Bind(&userRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	userEntity := dto.UserSendOTPRequestToUserEntity(userRequest)

	otp, errSendOTP := uh.userCommandUsecase.SendUserOTP(userEntity.Email)
	if errSendOTP != nil {
		if strings.Contains(errSendOTP.Error(), constant.ERROR_EMAIL_NOTFOUND) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse(errSendOTP.Error()))
		}
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errSendOTP.Error()))
	}

	userResponse := dto.UserEntityToUserResponse(otp)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_OTP_SENT, userResponse))
}

func (uh *userHandler) VerifyUserOTP(c echo.Context) error {
	userRequest := dto.UserVerifyOTPRequest{}

	errBind := c.Bind(&userRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	userEntity := dto.UserVerifyOTPRequestToUserEntity(userRequest)

	token, errVerify := uh.userCommandUsecase.VerifyUserOTP(userEntity.Email, userEntity.OTP)
	if errVerify != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse(constant.ERROR_OTP_VERIFY+errVerify.Error()))
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_OTP_VERIFIED, token))
}

func (uh *userHandler) NewUserPassword(c echo.Context) error {
	userRequest := dto.UserNewPasswordRequest{}

	errBind := c.Bind(&userRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	email, errExtract := middlewares.ExtractVerifyToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	userEntity := dto.UserNewPasswordRequestToUserEntity(userRequest)

	password, errCreate := uh.userCommandUsecase.NewUserPassword(email, userEntity)
	if errCreate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errCreate.Error()))
	}

	userResponse := dto.UserEntityToUserResponse(password)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PASSWORD_UPDATED, userResponse))
}
