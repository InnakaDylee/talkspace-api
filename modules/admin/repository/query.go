package repository

import (
	"context"
	"encoding/json"
	"errors"
	"talkspace-api/modules/admin/entity"
	"talkspace-api/modules/admin/model"
	"talkspace-api/utils/constant"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type adminQueryRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewAdminQueryRepository(db *gorm.DB, rdb *redis.Client) AdminQueryRepositoryInterface {
	return &adminQueryRepository{
		db:  db,
		rdb: rdb,
	}
}

func (aqr *adminQueryRepository) GetAdminByID(id string) (entity.Admin, error) {
	cacheKey := "admin:id:" + id
	cachedAdmin, err := aqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedAdmin != "" {
		var admin entity.Admin
		if err := json.Unmarshal([]byte(cachedAdmin), &admin); err != nil {
			return entity.Admin{}, err
		}
		return admin, nil
	}

	adminModel := model.Admin{}
	result := aqr.db.Where("id = ?", id).First(&adminModel)
	if result.Error != nil {
		return entity.Admin{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.Admin{}, errors.New(constant.ERROR_ID_NOTFOUND)
	}

	adminEntity := entity.AdminModelToAdminEntity(adminModel)

	adminData, err := json.Marshal(adminEntity)
	if err == nil {
		aqr.rdb.Set(context.Background(), cacheKey, string(adminData), 10*time.Minute)
	}

	return adminEntity, nil
}

func (aqr *adminQueryRepository) GetAdminByEmail(email string) (entity.Admin, error) {
	cacheKey := "admin:email:" + email
	cachedAdmin, err := aqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedAdmin != "" {
		var admin entity.Admin
		if err := json.Unmarshal([]byte(cachedAdmin), &admin); err != nil {
			return entity.Admin{}, err
		}
		return admin, nil
	}

	adminModel := model.Admin{}
	result := aqr.db.Where("email = ?", email).First(&adminModel)

	if result.RowsAffected == 0 {
		return entity.Admin{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
	}

	if result.Error != nil {
		return entity.Admin{}, result.Error
	}

	adminEntity := entity.AdminModelToAdminEntity(adminModel)

	adminData, err := json.Marshal(adminEntity)
	if err == nil {
		aqr.rdb.Set(context.Background(), cacheKey, string(adminData), 10*time.Minute)
	}

	return adminEntity, nil
}
