package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/model"
	"talkspace-api/utils/bcrypt"
	"talkspace-api/utils/constant"

	"github.com/elastic/go-elasticsearch/v8"

	"gorm.io/gorm"
)

type userCommandRepository struct {
	db *gorm.DB
	es *elasticsearch.Client
}

func NewUserCommandRepository(db *gorm.DB, es *elasticsearch.Client) UserCommandRepositoryInterface {
	return &userCommandRepository{
		db: db,
		es: es,
	}
}
func (ucr *userCommandRepository) RegisterUser(user entity.User) (entity.User, error) {
	userModel := entity.UserEntityToUserModel(user)

	result := ucr.db.Create(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)
	data, err := json.Marshal(userEntity)
	if err != nil {
		return entity.User{}, err
	}

	res, err := ucr.es.Index(
		"users",
		bytes.NewReader(data),
		ucr.es.Index.WithContext(context.Background()),
		ucr.es.Index.WithDocumentID(userEntity.ID),
	)
	if err != nil {
		return entity.User{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return entity.User{}, err
		} else {
			return entity.User{}, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

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
	data, err := json.Marshal(userEntity)
	if err != nil {
		return entity.User{}, err
	}

	res, err := ucr.es.Index(
		"users",
		bytes.NewReader(data),
		ucr.es.Index.WithContext(context.Background()),
		ucr.es.Index.WithDocumentID(userEntity.ID),
	)
	if err != nil {
		return entity.User{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return entity.User{}, err
		} else {
			return entity.User{}, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

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

func (ucr *userCommandRepository) VerifyUserOTP(email, otp string) (entity.User, error) {
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
	data, err := json.Marshal(userEntity)
	if err != nil {
		return entity.User{}, err
	}

	res, err := ucr.es.Index(
		"users",
		bytes.NewReader(data),
		ucr.es.Index.WithContext(context.Background()),
		ucr.es.Index.WithDocumentID(userEntity.ID),
	)
	if err != nil {
		return entity.User{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return entity.User{}, err
		} else {
			return entity.User{}, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	return userEntity, nil
}
