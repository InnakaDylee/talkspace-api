package repository

import (
	"errors"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/model"
	"talkspace-api/utils/bcrypt"
	"talkspace-api/utils/constant"

	"gorm.io/gorm"
)

type userCommandRepository struct {
	db *gorm.DB
}

func NewUserCommandRepository(db *gorm.DB) UserCommandRepositoryInterface {
	return &userCommandRepository{
		db: db,
	}
}

func (ucr *userCommandRepository) RegisterUser(user entity.User) (entity.User, error) {
	userModel := entity.UserEntityToUserModel(user)

	result := ucr.db.Create(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) LoginUser(email, password string) (entity.User, error) {
	userModel := model.User{}

	result := ucr.db.Where("email = ?", email).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
	}

	if errComparePass := bcrypt.ComparePassword(userModel.Password, password); errComparePass != nil {
		return entity.User{}, errors.New(constant.ERROR_PASSWORD_INVALID)
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) UpdateUserByID(id string, user entity.User) (entity.User, error) {
	userModel := entity.UserEntityToUserModel(user)

	result := ucr.db.Where("id = ?", id).Updates(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_ID_NOTFOUND)
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) UpdateUserIsVerified(id string, isVerified bool) (entity.User, error) {
	userModel := model.User{}

	result := ucr.db.Where("id = ?", id).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_ID_NOTFOUND)
	}

	userModel.IsVerified = isVerified

	errSave := ucr.db.Save(&userModel)
	if errSave.Error != nil {
		return entity.User{}, errSave.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) SendUserOTP(email string, otp string, expired int64) (entity.User, error) {
	userModel := model.User{}

	result := ucr.db.Where("email = ?", email).First(&userModel)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return entity.User{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
		}
		return entity.User{}, result.Error
	}

	userModel.OTP = otp
	userModel.OTPExpiration = expired

	errUpdate := ucr.db.Save(&userModel).Error
	if errUpdate != nil {
		return entity.User{}, errUpdate
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) VerifyUserOTP(email, otp string) (entity.User, error){
	userModel := model.User{}

	result := ucr.db.Where("otp = ? AND email = ?", otp, email).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_EMAIL_OTP)
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) ResetUserOTP(otp string) (entity.User, error) {
	userModel := model.User{}

	result := ucr.db.Where("otp = ?", otp).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_OTP_NOTFOUND)
	}

	userModel.OTP = ""
	userModel.OTPExpiration = 0

	errUpdate := ucr.db.Save(&userModel).Error
	if errUpdate != nil {
		return entity.User{}, errUpdate
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) UpdateUserPassword(id string, password entity.User) (entity.User, error) {

	userModel := entity.UserEntityToUserModel(password)

	result := ucr.db.Where("id = ?", id).Updates(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_ID_NOTFOUND)
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (ucr *userCommandRepository) NewUserPassword(email string, password entity.User) (entity.User, error)  {
	userModel := model.User{}

	result := ucr.db.Where("email = ?", email).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
	}

	errUpdate := ucr.db.Model(&userModel).Updates(entity.UserEntityToUserModel(password))
	if errUpdate != nil {
		return entity.User{}, errUpdate.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}
