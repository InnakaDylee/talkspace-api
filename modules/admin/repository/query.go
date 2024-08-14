package repository

import (
	"context"
	"encoding/json"
	"errors"
	"talkspace-api/modules/admin/entity"
	"talkspace-api/modules/admin/model"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/validator"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type adminQueryRepository struct {
	db  *gorm.DB
	es  *elasticsearch.Client
	rdb *redis.Client
}

func NewAdminQueryRepository(db *gorm.DB, es *elasticsearch.Client, rdb *redis.Client) AdminQueryRepositoryInterface {
	return &adminQueryRepository{
		db:  db,
		es:  es,
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

    query := map[string]interface{}{
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                "id": id,
            },
        },
    }

    var admin entity.Admin
    res, err := aqr.es.Search(
        aqr.es.Search.WithContext(context.Background()),
        aqr.es.Search.WithIndex("admins"),
        aqr.es.Search.WithBody(validator.JSONReader(query)),
        aqr.es.Search.WithTrackTotalHits(true),
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

    var r map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
        return entity.Admin{}, err
    }

    if hits := r["hits"].(map[string]interface{})["hits"].([]interface{}); len(hits) > 0 {
        source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
        admin = entity.Admin{
            ID:              source["id"].(string),
            Fullname:        source["fullname"].(string),
            Email:           source["email"].(string),
            Password:        source["password"].(string),
            Role:            source["role"].(string),
            OTP:             source["otp"].(string),
            OTPExpiration:   int64(source["otpExpiration"].(float64)),
            CreatedAt:       time.Unix(int64(source["createdAt"].(float64)), 0),
            UpdatedAt:       time.Unix(int64(source["updatedAt"].(float64)), 0),
            DeletedAt:       validator.ConvertToTime(source["deletedAt"]),
        }

        adminData, err := json.Marshal(admin)
        if err == nil {
            aqr.rdb.Set(context.Background(), cacheKey, string(adminData), 10*time.Minute)
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

    query := map[string]interface{}{
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                "email": email,
            },
        },
    }

    var admin entity.Admin
    res, err := aqr.es.Search(
        aqr.es.Search.WithContext(context.Background()),
        aqr.es.Search.WithIndex("admins"),
        aqr.es.Search.WithBody(validator.JSONReader(query)),
        aqr.es.Search.WithTrackTotalHits(true),
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

    var r map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
        return entity.Admin{}, err
    }

    if hits := r["hits"].(map[string]interface{})["hits"].([]interface{}); len(hits) > 0 {
        source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
        admin = entity.Admin{
            ID:              source["id"].(string),
            Fullname:        source["fullname"].(string),
            Email:           source["email"].(string),
            Password:        source["password"].(string),
            Role:            source["role"].(string),
            OTP:             source["otp"].(string),
            OTPExpiration:   int64(source["otpExpiration"].(float64)),
            CreatedAt:       time.Unix(int64(source["createdAt"].(float64)), 0),
            UpdatedAt:       time.Unix(int64(source["updatedAt"].(float64)), 0),
            DeletedAt:       validator.ConvertToTime(source["deletedAt"]),
        }

        adminData, err := json.Marshal(admin)
        if err == nil {
            aqr.rdb.Set(context.Background(), cacheKey, string(adminData), 10*time.Minute)
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
