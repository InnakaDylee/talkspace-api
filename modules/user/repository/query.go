package repository

import (
	"context"
	"encoding/json"
	"errors"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userQueryRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserQueryRepository(db *gorm.DB, rdb *redis.Client) UserQueryRepositoryInterface {
	return &userQueryRepository{
		db:  db,
		rdb: rdb,
	}
}
func (uqr *userQueryRepository) GetUserByID(id string) (entity.User, error) {
	cacheKey := "user:id:" + id

	cachedUser, err := uqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedUser != "" {
		var user entity.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			return entity.User{}, err
		}
		return user, nil
	}

	var userModel model.User
	result := uqr.db.Where("id = ?", id).First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.User{}, errors.New("user not found")
		}
		return entity.User{}, result.Error
	}

	userEntity := entity.UserModelToUserEntity(userModel)

	userData, err := json.Marshal(userEntity)
	if err != nil {
		return entity.User{}, err
	}

	err = uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute).Err()
	if err != nil {
		return entity.User{}, err
	}

	return userEntity, nil
}

func (uqr *userQueryRepository) GetUserByEmail(email string) (entity.User, error) {
	cacheKey := "user:email:" + email

	cachedUser, err := uqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedUser != "" {
		var user entity.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			return entity.User{}, err
		}
		return user, nil
	}

	var user entity.User
	result := uqr.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.User{}, errors.New("user not found")
		}
		return entity.User{}, result.Error
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return entity.User{}, err
	}

	err = uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute).Err()
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (uqr *userQueryRepository) GetRequestPremiumUsers() ([]entity.User, error) {
	var users []entity.User
	result := uqr.db.Where("request_premium != ?","").Find(&users)
	if result.Error != nil {
		return []entity.User{}, result.Error
	}

	return users, nil
}
