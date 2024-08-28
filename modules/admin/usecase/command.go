package usecase

import (
	"errors"
	"talkspace-api/middlewares"
	"talkspace-api/modules/admin/entity"
	"talkspace-api/modules/admin/repository"
	"talkspace-api/utils/bcrypt"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/generator"
	"talkspace-api/utils/helper/email/mailer"
	"talkspace-api/utils/validator"
	"time"
)

type adminCommandUsecase struct {
	adminCommandRepository repository.AdminCommandRepositoryInterface
	adminQueryRepository   repository.AdminQueryRepositoryInterface
}

func NewAdminCommandUsecase(acr repository.AdminCommandRepositoryInterface, aqr repository.AdminQueryRepositoryInterface) AdminCommandUsecaseInterface {
	return &adminCommandUsecase{
		adminCommandRepository: acr,
		adminQueryRepository:   aqr,
	}
}

func (acu *adminCommandUsecase) RegisterAdmin(admin entity.Admin) (entity.Admin, error) {

	errEmpty := validator.IsDataEmpty([]string{"fullname", "email", "password", "confirm_password"}, admin.Fullname, admin.Email, admin.Password, admin.ConfirmPassword)
	if errEmpty != nil {
		return entity.Admin{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(admin.Email)
	if errEmailValid != nil {
		return entity.Admin{}, errEmailValid
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": admin.Password})
	if errLength != nil {
		return entity.Admin{}, errLength
	}

	_, errGetEmail := acu.adminQueryRepository.GetAdminByEmail(admin.Email)
	if errGetEmail == nil {
		return entity.Admin{}, errors.New(constant.ERROR_EMAIL_EXIST)
	}

	if admin.Password != admin.ConfirmPassword {
		return entity.Admin{}, errors.New(constant.ERROR_PASSWORD_CONFIRM)
	}

	hashedPassword, errHash := bcrypt.HashPassword(admin.Password)
	if errHash != nil {
		return entity.Admin{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}

	admin.Password = hashedPassword

	adminEntity, errRegister := acu.adminCommandRepository.RegisterAdmin(admin)
	if errRegister != nil {
		return entity.Admin{}, errRegister
	}

	mailer.SendEmailNotificationRegisterAccount(adminEntity.Email)

	return adminEntity, nil
}

func (acu *adminCommandUsecase) LoginAdmin(email, password string) (entity.Admin, string, error) {

	errEmpty := validator.IsDataEmpty([]string{"email", "password"}, email, password)
	if errEmpty != nil {
		return entity.Admin{}, "", errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.Admin{}, "", errEmailValid
	}

	adminEntity, errGetEmail := acu.adminQueryRepository.GetAdminByEmail(email)
	if errGetEmail != nil {
		return entity.Admin{}, "", errors.New(constant.ERROR_EMAIL_UNREGISTERED)
	}

	comparePassword := bcrypt.ComparePassword(adminEntity.Password, password)
	if comparePassword != nil {
		return entity.Admin{}, "", errors.New(constant.ERROR_LOGIN)
	}

	token, errCreate := middlewares.GenerateToken(adminEntity.ID, adminEntity.Role)
	if errCreate != nil {
		return entity.Admin{}, "", errors.New(constant.ERROR_TOKEN_GENERATE)
	}

	mailer.SendEmailNotificationLoginAccount(email)

	return adminEntity, token, nil
}

func (acu *adminCommandUsecase) UpdateAdminPassword(id string, password entity.Admin) (entity.Admin, error) {
	if id == "" {
		return entity.Admin{}, errors.New(constant.ERROR_ID_INVALID)
	}

	result, errGetID := acu.adminQueryRepository.GetAdminByID(id)
	if errGetID != nil {
		return entity.Admin{}, errGetID
	}

	errEmpty := validator.IsDataEmpty([]string{"password", "new_password", "confirm_password"}, password.Password, password.NewPassword, password.ConfirmPassword)
	if errEmpty != nil {
		return entity.Admin{}, errEmpty
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": password.NewPassword})
	if errLength != nil {
		return entity.Admin{}, errLength
	}

	comparePassword := bcrypt.ComparePassword(result.Password, password.Password)
	if comparePassword != nil {
		return entity.Admin{}, errors.New(constant.ERROR_OLDPASSWORD_INVALID)
	}

	if password.NewPassword != password.ConfirmPassword {
		return entity.Admin{}, errors.New(constant.ERROR_PASSWORD_CONFIRM)
	}

	HashPassword, errHash := bcrypt.HashPassword(password.NewPassword)
	if errHash != nil {
		return entity.Admin{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}
	password.Password = HashPassword

	adminEntity, errUpdate := acu.adminCommandRepository.UpdateAdminPassword(id, password)
	if errUpdate != nil {
		return entity.Admin{}, errUpdate
	}

	return adminEntity, nil
}

func (acu *adminCommandUsecase) SendAdminOTP(email string) (entity.Admin, error) {

	errEmpty := validator.IsDataEmpty([]string{"email"}, email)
	if errEmpty != nil {
		return entity.Admin{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.Admin{}, errEmailValid
	}

	code, errGenerate := generator.GenerateRandomCode()
	if errGenerate != nil {
		return entity.Admin{}, errors.New(constant.ERROR_OTP_GENERATE)
	}

	expired := time.Now().Add(5 * time.Minute).Unix()

	adminEntity, errSend := acu.adminCommandRepository.SendAdminOTP(email, code, expired)
	if errSend != nil {
		return entity.Admin{}, errSend
	}

	mailer.SendEmailOTP(email, code)
	return adminEntity, nil
}

func (acu *adminCommandUsecase) VerifyAdminOTP(email, otp string) (string, error) {
	errEmpty := validator.IsDataEmpty([]string{"email", "otp"}, email, otp)
	if errEmpty != nil {
		return "", errEmpty
	}

	adminEntity, err := acu.adminCommandRepository.VerifyAdminOTP(email, otp)
	if err != nil {
		return "", errors.New(constant.ERROR_EMAIL_OTP)
	}

	if adminEntity.OTPExpiration <= time.Now().Unix() {
		return "", errors.New(constant.ERROR_OTP_EXPIRED)
	}

	if adminEntity.OTP != otp {
		return "", errors.New(constant.ERROR_OTP_INVALID)
	}

	token, err := middlewares.GenerateVerifyToken(email)
	if err != nil {
		return "", errors.New(constant.ERROR_TOKEN_GENERATE)
	}

	_, errReset := acu.adminCommandRepository.ResetAdminOTP(otp)
	if errReset != nil {
		return "", errors.New(constant.ERROR_OTP_RESET)
	}

	return token, nil
}

func (acu *adminCommandUsecase) NewAdminPassword(email string, password entity.Admin) (entity.Admin, error) {
	errEmpty := validator.IsDataEmpty([]string{"email", "password", "confirm_passsword"}, email, password.Password, password.ConfirmPassword)
	if errEmpty != nil {
		return entity.Admin{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.Admin{}, errEmailValid
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": password.Password})
	if errLength != nil {
		return entity.Admin{}, errLength
	}

	if password.Password != password.ConfirmPassword {
		return entity.Admin{}, errors.New(constant.ERROR_OLDPASSWORD_INVALID)
	}

	HashPassword, errHash := bcrypt.HashPassword(password.Password)
	if errHash != nil {
		return entity.Admin{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}
	password.Password = HashPassword

	adminEntity, errNewPass := acu.adminCommandRepository.NewAdminPassword(email, password)
	if errNewPass != nil {
		return entity.Admin{}, errNewPass
	}

	return adminEntity, nil
}
