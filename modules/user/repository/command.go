package repository

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/model"
	"talkspace-api/utils/bcrypt"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/helper/cloud"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userCommandRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserCommandRepository(db *gorm.DB, rdb *redis.Client) UserCommandRepositoryInterface {
	return &userCommandRepository{
		db:  db,
		rdb: rdb,
	}
}

func (ucr *userCommandRepository) RegisterUser(user entity.User) (entity.User, error) {
	userModel := entity.UserEntityToUserModel(user)

	result := ucr.db.Create(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	cacheKey := "user:" + userEntity.ID
	data, err := json.Marshal(userEntity)
	if err != nil {
		return entity.User{}, err
	}

	err = ucr.rdb.Set(context.Background(), cacheKey, data, 24*time.Hour).Err()
	if err != nil {
		return entity.User{}, err
	}

	return userEntity, nil
}

func (ucr *userCommandRepository) LoginUser(email, password string) (entity.User, error) {
	cacheKey := "user:" + email

	cachedData, err := ucr.rdb.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		userModel := model.User{}

		result := ucr.db.Where("email = ?", email).First(&userModel)
		if result.Error != nil {
			return entity.User{}, result.Error
		}

		if errComparePass := bcrypt.ComparePassword(userModel.Password, password); errComparePass != nil {
			return entity.User{}, errors.New(constant.ERROR_PASSWORD_INVALID)
		}

		userEntity := entity.UserModelToUserEntity(userModel)

		data, _ := json.Marshal(userEntity)
		err = ucr.rdb.Set(context.Background(), cacheKey, data, 24*time.Hour).Err()
		if err != nil {
			return entity.User{}, err
		}

		return userEntity, nil
	} else if err != nil {
		return entity.User{}, err
	}

	userEntity := entity.User{}
	if err := json.Unmarshal([]byte(cachedData), &userEntity); err != nil {
		return entity.User{}, err
	}

	return userEntity, nil
}

func (ucr *userCommandRepository) UpdateUserProfile(id string, user entity.User, image *multipart.FileHeader) (entity.User, error) {
	userModel := entity.UserEntityToUserModel(user)

	if image != nil {
        imageURL, errUpload := cloud.UploadImageToS3(image)
        if errUpload != nil {
            return entity.User{}, errUpload
        }
        userModel.ProfilePicture = imageURL
    }

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

	cacheKey := "user:" + id
	err = ucr.rdb.Set(context.Background(), cacheKey, data, 24*time.Hour).Err()
	if err != nil {
		return entity.User{}, err
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

	cacheKey := "otp:" + email
	err := ucr.rdb.Set(context.Background(), cacheKey, otp, time.Duration(expired)*time.Second).Err()
	if err != nil {
		return entity.User{}, err
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
	cacheKey := "otp:" + email

	cachedOTP, err := ucr.rdb.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil || cachedOTP != otp {
		return entity.User{}, errors.New(constant.ERROR_EMAIL_OTP)
	} else if err != nil {
		return entity.User{}, err
	}

	userModel := model.User{}
	result := ucr.db.Where("otp = ? AND email = ?", otp, email).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	ucr.rdb.Del(context.Background(), cacheKey)

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

func (ucr *userCommandRepository) NewUserPassword(email string, password entity.User) (entity.User, error) {
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

	cacheKey := "user:" + userEntity.ID
	err = ucr.rdb.Set(context.Background(), cacheKey, data, 24*time.Hour).Err()
	if err != nil {
		return entity.User{}, err
	}

	return userEntity, nil
}
