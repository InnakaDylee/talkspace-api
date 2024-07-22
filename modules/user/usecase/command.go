package usecase

import (
	"errors"
	"talkspace-api/middlewares"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/repository"
	"talkspace-api/utils/bcrypt"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/generator"
	"talkspace-api/utils/helper/email/mailer"
	"talkspace-api/utils/validator"
	"time"
)

type userCommandUsecase struct {
	userCommandRepository repository.UserCommandRepositoryInterface
	userQueryRepository   repository.UserQueryRepositoryInterface
}

func NewUserCommandUsecase(ucr repository.UserCommandRepositoryInterface, uqr repository.UserQueryRepositoryInterface) UserCommandUsecaseInterface {
	return &userCommandUsecase{
		userCommandRepository: ucr,
		userQueryRepository:   uqr,
	}
}
func (ucs *userCommandUsecase) RegisterUser(user entity.User) (entity.User, error) {

	errEmpty := validator.IsDataEmpty([]string{"fullname", "email", "password", "confirm_password"}, user.Fullname, user.Email, user.Password, user.ConfirmPassword)
	if errEmpty != nil {
		return entity.User{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(user.Email)
	if errEmailValid != nil {
		return entity.User{}, errEmailValid
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": user.Password})
	if errLength != nil {
		return entity.User{}, errLength
	}

	_, errGetEmail := ucs.userQueryRepository.GetUserByEmail(user.Email)
	if errGetEmail == nil {
		return entity.User{}, errors.New(constant.ERROR_EMAIL_EXIST)
	}

	if user.Password != user.ConfirmPassword {
		return entity.User{}, errors.New(constant.ERROR_PASSWORD_CONFIRM)
	}

	hashedPassword, errHash := bcrypt.HashPassword(user.Password)
	if errHash != nil {
		return entity.User{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}

	user.Password = hashedPassword

	token, errGenerateVerifyToken := generator.GenerateRandomBytes()
	if errGenerateVerifyToken != nil {
		return entity.User{}, errors.New(constant.ERROR_TOKEN_VERIFICATION)
	}
	user.VerifyAccount = token

	userEntity, errRegister := ucs.userCommandRepository.RegisterUser(user)
	if errRegister != nil {
		return entity.User{}, errRegister
	}

	mailer.SendEmailVerificationAccount(userEntity.Email, token)

	return userEntity, nil
}

func (ucs *userCommandUsecase) LoginUser(email, password string) (entity.User, string, error) {

	errEmpty := validator.IsDataEmpty([]string{"email", "password"}, email, password)
	if errEmpty != nil {
		return entity.User{}, "", errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.User{}, "", errEmailValid
	}

	userEntity, errGetEmail := ucs.userQueryRepository.GetUserByEmail(email)
	if errGetEmail != nil {
		return entity.User{}, "", errors.New(constant.ERROR_EMAIL_UNREGISTERED)
	}

	comparePassword := bcrypt.ComparePassword(userEntity.Password, password)
	if comparePassword != nil {
		return entity.User{}, "", errors.New(constant.ERROR_LOGIN)
	}

	token, errCreate := middlewares.CreateToken(userEntity.ID, userEntity.Role)
	if errCreate != nil {
		return entity.User{}, "", errors.New(constant.ERROR_TOKEN_GENERATE)
	}

	mailer.SendEmailNotificationLoginAccount(email)

	return userEntity, token, nil
}

func (ucs *userCommandUsecase) UpdateUserByID(id string, user entity.User) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New(constant.ERROR_ID_INVALID)
	}

	_, errGetID := ucs.userQueryRepository.GetUserByID(id)
	if errGetID != nil {
		return entity.User{}, errGetID
	}

	errEmailValid := validator.IsEmailValid(user.Email)
	if errEmailValid != nil {
		return entity.User{}, errEmailValid
	}

	errBirthdate := validator.IsDateValid(user.Birthdate)
	if errBirthdate != nil {
		return entity.User{}, errBirthdate
	}

	validGender := []interface{}{"male", "female"}
	errGender := validator.IsDataValid(user.Gender, validGender, true)
	if errGender != nil {
		return entity.User{}, errGender
	}

	validBloodType := []interface{}{"A", "B", "O", "AB"}
	errBloodType := validator.IsDataValid(user.BloodType, validBloodType, true)
	if errBloodType != nil {
		return entity.User{}, errBloodType
	}

	userEntity, errUpdate := ucs.userCommandRepository.UpdateUserByID(id, user)
	if errUpdate != nil {
		return entity.User{}, errUpdate
	}

	return userEntity, nil
}

func (ucs *userCommandUsecase) UpdateUserPassword(id string, password entity.User) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New(constant.ERROR_ID_INVALID)
	}

	result, errGetID := ucs.userQueryRepository.GetUserByID(id)
	if errGetID != nil {
		return entity.User{}, errGetID
	}

	errEmpty := validator.IsDataEmpty([]string{"password", "new_password", "confirm_password"}, password.Password, password.NewPassword, password.ConfirmPassword)
	if errEmpty != nil {
		return entity.User{}, errEmpty
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": password.NewPassword})
	if errLength != nil {
		return entity.User{}, errLength
	}

	comparePassword := bcrypt.ComparePassword(result.Password, password.Password)
	if comparePassword != nil {
		return entity.User{}, errors.New(constant.ERROR_OLDPASSWORD_INVALID)
	}

	if password.NewPassword != password.ConfirmPassword {
		return entity.User{}, errors.New(constant.ERROR_PASSWORD_CONFIRM)
	}

	HashPassword, errHash := bcrypt.HashPassword(password.NewPassword)
	if errHash != nil {
		return entity.User{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}
	password.Password = HashPassword

	userEntity, errUpdate := ucs.userCommandRepository.UpdateUserPassword(id, password)
	if errUpdate != nil {
		return entity.User{}, errUpdate
	}

	return userEntity, nil
}

func (ucs *userCommandUsecase) VerifyUser(token string) (bool, error) {
	if token == "" {
		return false, errors.New(constant.ERROR_TOKEN_INVALID)
	}

	user, errGetVerifyToken := ucs.userQueryRepository.GetUserByVerificationToken(token)
	if errGetVerifyToken != nil {
		return false, errors.New(constant.ERROR_DATA_RETRIEVED)
	}

	if user.IsVerified {
		return true, nil
	}

	_, errUpdate := ucs.userCommandRepository.UpdateUserIsVerified(user.ID, true)
	if errUpdate != nil {
		return false, errors.New(constant.ERROR_ACCOUNT_VERIFICATION)
	}

	return false, nil
}

func (ucs *userCommandUsecase) UpdateUserIsVerified(id string, isVerified bool) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New(constant.ERROR_ID_INVALID)
	}

	return ucs.userCommandRepository.UpdateUserIsVerified(id, isVerified)
}

func (ucs *userCommandUsecase) SendUserOTP(email string) (entity.User, error) {

	errEmpty := validator.IsDataEmpty([]string{"email"}, email)
	if errEmpty != nil {
		return entity.User{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.User{}, errEmailValid
	}

	code, errGenerate := generator.GenerateRandomCode()
	if errGenerate != nil {
		return entity.User{}, errors.New(constant.ERROR_OTP_GENERATE)
	}

	expired := time.Now().Add(5 * time.Minute).Unix()

	userEntity, errSend := ucs.userCommandRepository.SendUserOTP(email, code, expired)
	if errSend != nil {
		return entity.User{}, errSend
	}

	mailer.SendEmailOTP(email, code)
	return userEntity, nil
}

func (ucs *userCommandUsecase) VerifyUserOTP(email, otp string) (string, error) {
	errEmpty := validator.IsDataEmpty([]string{"email", "otp"}, email, otp)
	if errEmpty != nil {
		return "", errEmpty
	}

	userEntity, err := ucs.userCommandRepository.VerifyUserOTP(email, otp)
	if err != nil {
		return "", errors.New(constant.ERROR_EMAIL_OTP)
	}

	if userEntity.OTPExpiration <= time.Now().Unix() {
		return "", errors.New(constant.ERROR_OTP_EXPIRED)
	}

	if userEntity.OTP != otp {
		return "", errors.New(constant.ERROR_OTP_INVALID)
	}

	token, err := middlewares.CreateVerifyToken(email)
	if err != nil {
		return "", errors.New(constant.ERROR_TOKEN_GENERATE)
	}

	_, errReset := ucs.userCommandRepository.ResetUserOTP(otp)
	if errReset != nil {
		return "", errors.New(constant.ERROR_OTP_RESET)
	}

	return token, nil
}

func (ucs *userCommandUsecase) NewUserPassword(email string, password entity.User) (entity.User, error) {
	errEmpty := validator.IsDataEmpty([]string{"email", "password", "confirm_passsword"}, email, password.Password, password.ConfirmPassword)
	if errEmpty != nil {
		return entity.User{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.User{}, errEmailValid
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": password.Password})
	if errLength != nil {
		return entity.User{}, errLength
	}

	if password.Password != password.ConfirmPassword {
		return entity.User{}, errors.New(constant.ERROR_OLDPASSWORD_INVALID)
	}

	HashPassword, errHash := bcrypt.HashPassword(password.Password)
	if errHash != nil {
		return entity.User{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}
	password.Password = HashPassword

	userEntity, errNewPass := ucs.userCommandRepository.NewUserPassword(email, password)
	if errNewPass != nil {
		return entity.User{}, errNewPass
	}

	return userEntity, nil
}
