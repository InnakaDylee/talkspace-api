package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"talkspace-api/modules/admin/entity"
	"talkspace-api/modules/admin/model"
	"talkspace-api/utils/bcrypt"
	"talkspace-api/utils/constant"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type adminCommandRepository struct {
	db  *gorm.DB
	es  *elasticsearch.Client
	rdb *redis.Client
}

func NewAdminCommandRepository(db *gorm.DB, es *elasticsearch.Client, rdb *redis.Client) AdminCommandRepositoryInterface {
	return &adminCommandRepository{
		db:  db,
		es:  es,
		rdb: rdb,
	}
}

func (acr *adminCommandRepository) RegisterAdmin(admin entity.Admin) (entity.Admin, error) {
	adminModel := entity.AdminEntityToAdminModel(admin)

	result := acr.db.Create(&adminModel)
	if result.Error != nil {
		return entity.Admin{}, result.Error
	}

	adminEntity := entity.AdminModelToAdminEntity(adminModel)
	data, err := json.Marshal(adminEntity)
	if err != nil {
		return entity.Admin{}, err
	}

	res, err := acr.es.Index(
		"admins",
		bytes.NewReader(data),
		acr.es.Index.WithContext(context.Background()),
		acr.es.Index.WithDocumentID(adminEntity.ID),
	)
	if err != nil {
		return entity.Admin{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return entity.Admin{}, err
		} else {
			return entity.Admin{}, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	cacheKey := "admin:" + adminEntity.ID
	err = acr.rdb.Set(context.Background(), cacheKey, data, 24*time.Hour).Err()
	if err != nil {
		return entity.Admin{}, err
	}

	return adminEntity, nil
}

func (acr *adminCommandRepository) LoginAdmin(email, password string) (entity.Admin, error) {
	cacheKey := "admin:" + email

	cachedData, err := acr.rdb.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		adminModel := model.Admin{}

		result := acr.db.Where("email = ?", email).First(&adminModel)
		if result.Error != nil {
			return entity.Admin{}, result.Error
		}

		if errComparePass := bcrypt.ComparePassword(adminModel.Password, password); errComparePass != nil {
			return entity.Admin{}, errors.New(constant.ERROR_PASSWORD_INVALID)
		}

		adminEntity := entity.AdminModelToAdminEntity(adminModel)

		data, _ := json.Marshal(adminEntity)
		err = acr.rdb.Set(context.Background(), cacheKey, data, 24*time.Hour).Err()
		if err != nil {
			return entity.Admin{}, err
		}

		return adminEntity, nil
	} else if err != nil {
		return entity.Admin{}, err
	}

	adminEntity := entity.Admin{}
	if err := json.Unmarshal([]byte(cachedData), &adminEntity); err != nil {
		return entity.Admin{}, err
	}

	return adminEntity, nil
}

func (acr *adminCommandRepository) SendAdminOTP(email string, otp string, expired int64) (entity.Admin, error) {
    adminModel := model.Admin{}

    result := acr.db.Where("email = ?", email).First(&adminModel)
    if result.Error != nil {
        if result.RowsAffected == 0 {
            return entity.Admin{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
        }
        return entity.Admin{}, result.Error
    }

    cacheKey := "otp:" + email
    err := acr.rdb.Set(context.Background(), cacheKey, otp, time.Duration(expired)*time.Second).Err()
    if err != nil {
        return entity.Admin{}, err
    }

    adminModel.OTP = otp
    adminModel.OTPExpiration = expired

    errUpdate := acr.db.Save(&adminModel).Error
    if errUpdate != nil {
        return entity.Admin{}, errUpdate
    }

    adminEntity := entity.AdminModelToAdminEntity(adminModel)

    return adminEntity, nil
}

func (acr *adminCommandRepository) VerifyAdminOTP(email, otp string) (entity.Admin, error) {
    cacheKey := "otp:" + email

    cachedOTP, err := acr.rdb.Get(context.Background(), cacheKey).Result()
    if err == redis.Nil || cachedOTP != otp {
        return entity.Admin{}, errors.New(constant.ERROR_EMAIL_OTP)
    } else if err != nil {
        return entity.Admin{}, err
    }

    adminModel := model.Admin{}
    result := acr.db.Where("otp = ? AND email = ?", otp, email).First(&adminModel)
    if result.Error != nil {
        return entity.Admin{}, result.Error
    }

    adminEntity := entity.AdminModelToAdminEntity(adminModel)

    acr.rdb.Del(context.Background(), cacheKey)

    return adminEntity, nil
}

func (acr *adminCommandRepository) ResetAdminOTP(otp string) (entity.Admin, error) {
	adminModel := model.Admin{}

	result := acr.db.Where("otp = ?", otp).First(&adminModel)
	if result.Error != nil {
		return entity.Admin{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.Admin{}, errors.New(constant.ERROR_OTP_NOTFOUND)
	}

	adminModel.OTP = ""
	adminModel.OTPExpiration = 0

	errUpdate := acr.db.Save(&adminModel).Error
	if errUpdate != nil {
		return entity.Admin{}, errUpdate
	}

	adminEntity := entity.AdminModelToAdminEntity(adminModel)

	return adminEntity, nil
}

func (acr *adminCommandRepository) UpdateAdminPassword(id string, password entity.Admin) (entity.Admin, error) {

	adminModel := entity.AdminEntityToAdminModel(password)

	result := acr.db.Where("id = ?", id).Updates(&adminModel)
	if result.Error != nil {
		return entity.Admin{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.Admin{}, errors.New(constant.ERROR_ID_NOTFOUND)
	}

	adminEntity := entity.AdminModelToAdminEntity(adminModel)

	return adminEntity, nil
}

func (acr *adminCommandRepository) NewAdminPassword(email string, password entity.Admin) (entity.Admin, error)  {
	adminModel := model.Admin{}

	result := acr.db.Where("email = ?", email).First(&adminModel)
	if result.Error != nil {
		return entity.Admin{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.Admin{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
	}

	errUpdate := acr.db.Model(&adminModel).Updates(entity.AdminEntityToAdminModel(password))
	if errUpdate != nil {
		return entity.Admin{}, errUpdate.Error
	}

	adminEntity := entity.AdminModelToAdminEntity(adminModel)
	data, err := json.Marshal(adminEntity)
	if err != nil {
		return entity.Admin{}, err
	}

	res, err := acr.es.Index(
		"admins",
		bytes.NewReader(data),
		acr.es.Index.WithContext(context.Background()),
		acr.es.Index.WithDocumentID(adminEntity.ID),
	)
	if err != nil {
		return entity.Admin{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return entity.Admin{}, err
		} else {
			return entity.Admin{}, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	return adminEntity, nil
}
