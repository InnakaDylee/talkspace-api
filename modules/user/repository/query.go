package repository

import (
	"errors"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/model"
	"talkspace-api/utils/constant"

	"gorm.io/gorm"
)

type userQueryRepository struct {
	db *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) UserQueryRepositoryInterface {
	return &userQueryRepository{
		db: db,
	}
}

func (uqr *userQueryRepository) GetUserByID(id string) (entity.User, error) {
	userModel := model.User{}

	result := uqr.db.Where("id = ?", id).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_ID_NOTFOUND)
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (uqr *userQueryRepository) GetUserByEmail(email string) (entity.User, error) {
	userModel := model.User{}

	result := uqr.db.Where("email = ?", email).First(&userModel)

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
	}

	if result.Error != nil {
		return entity.User{}, result.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}

func (uqr *userQueryRepository) GetUserByVerificationToken(token string) (entity.User, error) {
	userModel := model.User{}

	result := uqr.db.Where("verification_token = ?", token).First(&userModel)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.User{}, errors.New(constant.ERROR_TOKEN_NOTFOUND)
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	return userEntity, nil
}
