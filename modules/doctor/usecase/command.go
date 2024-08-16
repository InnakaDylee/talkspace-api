package usecase

import (
	"errors"
	"mime/multipart"
	"talkspace-api/middlewares"
	"talkspace-api/modules/doctor/entity"
	"talkspace-api/modules/doctor/repository"
	"talkspace-api/utils/bcrypt"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/generator"
	"talkspace-api/utils/helper/email/mailer"
	"talkspace-api/utils/validator"
	"time"
)

type doctorCommandUsecase struct {
	doctorCommandRepository repository.DoctorCommandRepositoryInterface
	doctorQueryRepository   repository.DoctorQueryRepositoryInterface
}

func NewDoctorCommandUsecase(dcr repository.DoctorCommandRepositoryInterface, dqr repository.DoctorQueryRepositoryInterface) DoctorCommandUsecaseInterface {
	return &doctorCommandUsecase{
		doctorCommandRepository: dcr,
		doctorQueryRepository:   dqr,
	}
}

func (dcu *doctorCommandUsecase) RegisterDoctor(doctor entity.Doctor, image *multipart.FileHeader) (entity.Doctor, error) {

	errEmpty := validator.IsDataEmpty([]string{
		"fullname", "email", "password", "profile_picture",
		"gender", "specialization", "years_of_experience",
		"license_number", "alumnus", "about", "location"},
		doctor.Fullname, doctor.Email, doctor.Password,
		doctor.ProfilePicture, doctor.Gender, doctor.Specialization,
		doctor.YearsOfExperience, doctor.LicenseNumber, doctor.Alumnus,
		doctor.About, doctor.Location,
	)
	if errEmpty != nil {
		return entity.Doctor{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(doctor.Email)
	if errEmailValid != nil {
		return entity.Doctor{}, errEmailValid
	}

	_, errGetEmail := dcu.doctorQueryRepository.GetDoctorByEmail(doctor.Email)
	if errGetEmail == nil {
		return entity.Doctor{}, errors.New(constant.ERROR_EMAIL_EXIST)
	}

	if doctor.Password == "" {
		password, err := generator.GenerateRandomPassword(15) 
		if err != nil {
			return entity.Doctor{}, err
		}
		doctor.Password = password
	}

	hashedPassword, errHash := bcrypt.HashPassword(doctor.Password)
	if errHash != nil {
		return entity.Doctor{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}
	doctor.Password = hashedPassword

	doctorEntity, errRegister := dcu.doctorCommandRepository.RegisterDoctor(doctor, image)
	if errRegister != nil {
		return entity.Doctor{}, errRegister
	}

	mailer.SendEmailNotificationRegisterDoctor(
        doctorEntity.Fullname,
        doctorEntity.LicenseNumber,
        doctorEntity.Email,
        doctor.Password, 
    )

	return doctorEntity, nil
}

func (dcs *doctorCommandUsecase) LoginDoctor(email, password string) (entity.Doctor, string, error) {

	errEmpty := validator.IsDataEmpty([]string{"email", "password"}, email, password)
	if errEmpty != nil {
		return entity.Doctor{}, "", errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.Doctor{}, "", errEmailValid
	}

	doctorEntity, errGetEmail := dcs.doctorQueryRepository.GetDoctorByEmail(email)
	if errGetEmail != nil {
		return entity.Doctor{}, "", errors.New(constant.ERROR_EMAIL_UNREGISTERED)
	}

	comparePassword := bcrypt.ComparePassword(doctorEntity.Password, password)
	if comparePassword != nil {
		return entity.Doctor{}, "", errors.New(constant.ERROR_LOGIN)
	}

	token, errCreate := middlewares.GenerateToken(doctorEntity.ID, doctorEntity.Role)
	if errCreate != nil {
		return entity.Doctor{}, "", errors.New(constant.ERROR_TOKEN_GENERATE)
	}

	mailer.SendEmailNotificationLoginAccount(email)

	return doctorEntity, token, nil
}

func (dcs *doctorCommandUsecase) UpdateDoctorProfile(id string, doctor entity.Doctor, image *multipart.FileHeader) (entity.Doctor, error) {
	if id == "" {
		return entity.Doctor{}, errors.New(constant.ERROR_ID_INVALID)
	}

	_, errGetID := dcs.doctorQueryRepository.GetDoctorByID(id)
	if errGetID != nil {
		return entity.Doctor{}, errGetID
	}

	if doctor.Email != "" {
		errEmailValid := validator.IsEmailValid(doctor.Email)
		if errEmailValid != nil {
			return entity.Doctor{}, errEmailValid
		}
	}

	validGender := []interface{}{"male", "female"}
	errGender := validator.IsDataValid(doctor.Gender, validGender, true)
	if errGender != nil {
		return entity.Doctor{}, errGender
	}

	doctorEntity, errUpdate := dcs.doctorCommandRepository.UpdateDoctorProfile(id, doctor, image)
	if errUpdate != nil {
		return entity.Doctor{}, errUpdate
	}

	return doctorEntity, nil
}

func (dcs *doctorCommandUsecase) UpdateDoctorStatus(id string, status bool) (entity.Doctor, error) {
	if id == "" {
		return entity.Doctor{}, errors.New(constant.ERROR_ID_INVALID)
	}

	validStatuses := map[bool]bool{
		true:  true, // true represents "active"
		false: true, // false represents "inactive"
	}

	if _, isValidStatus := validStatuses[status]; !isValidStatus {
		return entity.Doctor{}, errors.New(constant.ERROR_STATUS_INVALID)
	}

	doctorEntity, err := dcs.doctorCommandRepository.UpdateDoctorStatus(id, status)
	if err != nil {
		return entity.Doctor{}, err
	}

	return doctorEntity, nil
}

func (dcs *doctorCommandUsecase) UpdateDoctorPassword(id string, password entity.Doctor) (entity.Doctor, error) {
	if id == "" {
		return entity.Doctor{}, errors.New(constant.ERROR_ID_INVALID)
	}

	result, errGetID := dcs.doctorQueryRepository.GetDoctorByID(id)
	if errGetID != nil {
		return entity.Doctor{}, errGetID
	}

	errEmpty := validator.IsDataEmpty([]string{"password", "new_password", "confirm_password"}, password.Password, password.NewPassword, password.ConfirmPassword)
	if errEmpty != nil {
		return entity.Doctor{}, errEmpty
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": password.NewPassword})
	if errLength != nil {
		return entity.Doctor{}, errLength
	}

	comparePassword := bcrypt.ComparePassword(result.Password, password.Password)
	if comparePassword != nil {
		return entity.Doctor{}, errors.New(constant.ERROR_OLDPASSWORD_INVALID)
	}

	if password.NewPassword != password.ConfirmPassword {
		return entity.Doctor{}, errors.New(constant.ERROR_PASSWORD_CONFIRM)
	}

	HashPassword, errHash := bcrypt.HashPassword(password.NewPassword)
	if errHash != nil {
		return entity.Doctor{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}
	password.Password = HashPassword

	doctorEntity, errUpdate := dcs.doctorCommandRepository.UpdateDoctorPassword(id, password)
	if errUpdate != nil {
		return entity.Doctor{}, errUpdate
	}

	return doctorEntity, nil
}

func (dcs *doctorCommandUsecase) SendDoctorOTP(email string) (entity.Doctor, error) {

	errEmpty := validator.IsDataEmpty([]string{"email"}, email)
	if errEmpty != nil {
		return entity.Doctor{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.Doctor{}, errEmailValid
	}

	code, errGenerate := generator.GenerateRandomCode()
	if errGenerate != nil {
		return entity.Doctor{}, errors.New(constant.ERROR_OTP_GENERATE)
	}

	expired := time.Now().Add(5 * time.Minute).Unix()

	doctorEntity, errSend := dcs.doctorCommandRepository.SendDoctorOTP(email, code, expired)
	if errSend != nil {
		return entity.Doctor{}, errSend
	}

	mailer.SendEmailOTP(email, code)
	return doctorEntity, nil
}

func (dcs *doctorCommandUsecase) VerifyDoctorOTP(email, otp string) (string, error) {
	errEmpty := validator.IsDataEmpty([]string{"email", "otp"}, email, otp)
	if errEmpty != nil {
		return "", errEmpty
	}

	doctorEntity, err := dcs.doctorCommandRepository.VerifyDoctorOTP(email, otp)
	if err != nil {
		return "", errors.New(constant.ERROR_EMAIL_OTP)
	}

	if doctorEntity.OTPExpiration <= time.Now().Unix() {
		return "", errors.New(constant.ERROR_OTP_EXPIRED)
	}

	if doctorEntity.OTP != otp {
		return "", errors.New(constant.ERROR_OTP_INVALID)
	}

	token, err := middlewares.GenerateVerifyToken(email)
	if err != nil {
		return "", errors.New(constant.ERROR_TOKEN_GENERATE)
	}

	_, errReset := dcs.doctorCommandRepository.ResetDoctorOTP(otp)
	if errReset != nil {
		return "", errors.New(constant.ERROR_OTP_RESET)
	}

	return token, nil
}

func (dcs *doctorCommandUsecase) NewDoctorPassword(email string, password entity.Doctor) (entity.Doctor, error) {
	errEmpty := validator.IsDataEmpty([]string{"email", "password", "confirm_passsword"}, email, password.Password, password.ConfirmPassword)
	if errEmpty != nil {
		return entity.Doctor{}, errEmpty
	}

	errEmailValid := validator.IsEmailValid(email)
	if errEmailValid != nil {
		return entity.Doctor{}, errEmailValid
	}

	errLength := validator.IsMinLengthValid(10, map[string]string{"password": password.Password})
	if errLength != nil {
		return entity.Doctor{}, errLength
	}

	if password.Password != password.ConfirmPassword {
		return entity.Doctor{}, errors.New(constant.ERROR_OLDPASSWORD_INVALID)
	}

	HashPassword, errHash := bcrypt.HashPassword(password.Password)
	if errHash != nil {
		return entity.Doctor{}, errors.New(constant.ERROR_PASSWORD_HASH)
	}
	password.Password = HashPassword

	doctorEntity, errNewPass := dcs.doctorCommandRepository.NewDoctorPassword(email, password)
	if errNewPass != nil {
		return entity.Doctor{}, errNewPass
	}

	return doctorEntity, nil
}
