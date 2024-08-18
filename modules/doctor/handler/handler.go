package handler

import (
	"net/http"
	"strings"
	"talkspace-api/middlewares"
	"talkspace-api/modules/doctor/dto"
	"talkspace-api/modules/doctor/usecase"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/helper/cloud"
	"talkspace-api/utils/responses"

	"github.com/labstack/echo/v4"
)

type doctorHandler struct {
	doctorCommandUsecase usecase.DoctorCommandUsecaseInterface
	doctorQueryUsecase   usecase.DoctorQueryUsecaseInterface
}

func NewDoctorHandler(dcu usecase.DoctorCommandUsecaseInterface, dqu usecase.DoctorQueryUsecaseInterface) *doctorHandler {
	return &doctorHandler{
		doctorCommandUsecase: dcu,
		doctorQueryUsecase:   dqu,
	}
}

// Query
func (dh *doctorHandler) GetDoctorByID(c echo.Context) error {
	doctorIDParam := c.Param("doctor_id")
	if doctorIDParam == "" {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_ID_NOTFOUND))
	}

	tokenDoctorID, role, errExtract := middlewares.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	if role != constant.DOCTOR {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if doctorIDParam != tokenDoctorID {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	doctor, errGetID := dh.doctorQueryUsecase.GetDoctorByID(doctorIDParam)
	if errGetID != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errGetID.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorProfileResponse(doctor)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PROFILE_RETRIEVED, doctorResponse))
}

func (dh *doctorHandler) GetAllDoctors(c echo.Context) error {
	statusParam := c.QueryParam("status")         
	specializationParam := c.QueryParam("specialization") 

	var status *bool
	if statusParam != "" {
		b := statusParam == "true"
		status = &b
	}

	doctors, err := dh.doctorQueryUsecase.GetAllDoctors(status, specializationParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse(err.Error()))
	}

	doctorResponses := make([]dto.DoctorResponse, len(doctors))
	for i, doctor := range doctors {
		doctorResponses[i] = dto.DoctorEntityToDoctorResponse(doctor)
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_RETRIEVED, doctorResponses))
}

// Command
func (dh *doctorHandler) RegisterDoctor(c echo.Context) error {

	_, role, errExtractToken := middlewares.ExtractToken(c)
	if errExtractToken != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtractToken.Error()))
	}

	if role != constant.ADMIN {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	doctorRequest := dto.DoctorRegisterRequest{}

	if errBind := c.Bind(&doctorRequest); errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	image, errFile := c.FormFile("profile_picture")
	if errFile != nil && errFile != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_UPLOAD_IMAGE))
	}

	if image != nil {
		imageURL, errUpload := cloud.UploadImageToS3(image)
		if errUpload != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse(constant.ERROR_UPLOAD_IMAGE_S3))
		}
		doctorRequest.ProfilePicture = imageURL
	}
	
	doctorEntity := dto.DoctorRegisterRequestToDoctorEntity(doctorRequest)

	registeredDoctor, errRegister := dh.doctorCommandUsecase.RegisterDoctor(doctorEntity, image)
	if errRegister != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errRegister.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorRegisterResponse(registeredDoctor)

	return c.JSON(http.StatusCreated, responses.SuccessResponse(constant.SUCCESS_REGISTER, doctorResponse))
}

func (dh *doctorHandler) LoginDoctor(c echo.Context) error {
	doctorRequest := dto.DoctorLoginRequest{}

	errBind := c.Bind(&doctorRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	LoginDoctor, token, errLogin := dh.doctorCommandUsecase.LoginDoctor(doctorRequest.Email, doctorRequest.Password)
	if errLogin != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errLogin.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorLoginResponse(LoginDoctor, token)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_LOGIN, doctorResponse))
}

func (dh *doctorHandler) UpdateDoctorProfile(c echo.Context) error {
	doctorIDParam := c.Param("doctor_id")
	if doctorIDParam == "" {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_ID_NOTFOUND))
	}

	tokenDoctorID, role, errExtract := middlewares.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	if role != constant.DOCTOR {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if doctorIDParam != tokenDoctorID {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	doctorRequest := dto.DoctorUpdateProfileRequest{}

	errBind := c.Bind(&doctorRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	image, errFile := c.FormFile("profile_picture")
	if errFile != nil && errFile != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_UPLOAD_IMAGE))
	}

	if image != nil {
		imageURL, errUpload := cloud.UploadImageToS3(image)
		if errUpload != nil {
			return c.JSON(http.StatusInternalServerError, responses.ErrorResponse(constant.ERROR_UPLOAD_IMAGE_S3))
		}
		doctorRequest.ProfilePicture = imageURL
	}

	doctorEntity := dto.DoctorUpdateProfileRequestToDoctorEntity(doctorRequest)

	doctor, errUpdate := dh.doctorCommandUsecase.UpdateDoctorProfile(doctorIDParam, doctorEntity, image)
	if errUpdate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errUpdate.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorUpdateProfileResponse(doctor)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PROFILE_UPDATED, doctorResponse))
}

func (dh *doctorHandler) UpdateDoctorStatus(c echo.Context) error {
	doctorIDParam := c.Param("doctor_id")
	if doctorIDParam == "" {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(constant.ERROR_ID_NOTFOUND))
	}

	tokenDoctorID, role, errExtract := middlewares.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	if role != constant.DOCTOR {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if doctorIDParam != tokenDoctorID {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	doctorRequest := dto.DoctorUpdateStatusRequest{}
	if err := c.Bind(&doctorRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(err.Error()))
	}

	doctorEntity := dto.DoctorUpdateStatusRequestToDoctorEntity(doctorRequest)

	updatedDoctor, errUpdate := dh.doctorCommandUsecase.UpdateDoctorStatus(doctorIDParam, doctorEntity.Status)
	if errUpdate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errUpdate.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorUpdateStatusResponse(updatedDoctor)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_STATUS_UPDATED, doctorResponse))
}

func (dh *doctorHandler) UpdateDoctorPassword(c echo.Context) error {
	doctorRequest := dto.DoctorUpdatePasswordRequest{}

	errBind := c.Bind(&doctorRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	doctorID, role, errExtractToken := middlewares.ExtractToken(c)

	if role != constant.DOCTOR {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if errExtractToken != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtractToken.Error()))
	}

	doctorEntity := dto.DoctorUpdatePasswordRequestToDoctorEntity(doctorRequest)

	password, errUpdate := dh.doctorCommandUsecase.UpdateDoctorPassword(doctorID, doctorEntity)
	if errUpdate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errUpdate.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorResponse(password)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PASSWORD_UPDATED, doctorResponse))
}

func (dh *doctorHandler) ForgotDoctorPassword(c echo.Context) error {
	doctorRequest := dto.DoctorSendOTPRequest{}

	errBind := c.Bind(&doctorRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	doctorEntity := dto.DoctorSendOTPRequestToDoctorEntity(doctorRequest)

	otp, errSendOTP := dh.doctorCommandUsecase.SendDoctorOTP(doctorEntity.Email)
	if errSendOTP != nil {
		if strings.Contains(errSendOTP.Error(), constant.ERROR_EMAIL_NOTFOUND) {
			return c.JSON(http.StatusNotFound, responses.ErrorResponse(errSendOTP.Error()))
		}
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errSendOTP.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorResponse(otp)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_OTP_SENT, doctorResponse))
}

func (dh *doctorHandler) VerifyDoctorOTP(c echo.Context) error {
	doctorRequest := dto.DoctorVerifyOTPRequest{}

	errBind := c.Bind(&doctorRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	doctorEntity := dto.DoctorVerifyOTPRequestToDoctorEntity(doctorRequest)

	token, errVerify := dh.doctorCommandUsecase.VerifyDoctorOTP(doctorEntity.Email, doctorEntity.OTP)
	if errVerify != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse(constant.ERROR_OTP_VERIFY+errVerify.Error()))
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_OTP_VERIFIED, token))
}

func (dh *doctorHandler) NewDoctorPassword(c echo.Context) error {
	doctorRequest := dto.DoctorNewPasswordRequest{}

	errBind := c.Bind(&doctorRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	email, errExtract := middlewares.ExtractVerifyToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtract.Error()))
	}

	doctorEntity := dto.DoctorNewPasswordRequestToDoctorEntity(doctorRequest)

	password, errCreate := dh.doctorCommandUsecase.NewDoctorPassword(email, doctorEntity)
	if errCreate != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errCreate.Error()))
	}

	doctorResponse := dto.DoctorEntityToDoctorResponse(password)

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_PASSWORD_UPDATED, doctorResponse))
}
