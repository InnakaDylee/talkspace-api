package repository

import (
	"context"
	"encoding/json"
	"errors"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/model"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/validator"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userQueryRepository struct {
	db  *gorm.DB
	es  *elasticsearch.Client
	rdb *redis.Client
}

func NewUserQueryRepository(db *gorm.DB, es *elasticsearch.Client, rdb *redis.Client) UserQueryRepositoryInterface {
	return &userQueryRepository{
		db:  db,
		es:  es,
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

    query := map[string]interface{}{
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                "id": id,
            },
        },
    }

    var user entity.User
    res, err := uqr.es.Search(
        uqr.es.Search.WithContext(context.Background()),
        uqr.es.Search.WithIndex("users"),
        uqr.es.Search.WithBody(validator.JSONReader(query)),
        uqr.es.Search.WithTrackTotalHits(true),
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

    var r map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
        return entity.User{}, err
    }

    if hits := r["hits"].(map[string]interface{})["hits"].([]interface{}); len(hits) > 0 {
        source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
        user = entity.User{
            ID:              source["id"].(string),
            Fullname:        source["fullname"].(string),
            Email:           source["email"].(string),
            Password:        source["password"].(string),
            NewPassword:     source["newPassword"].(string),
            ConfirmPassword: source["confirmPassword"].(string),
            ProfilePicture:  source["profilePicture"].(string),
            Birthdate:       source["birthdate"].(string),
            Gender:          source["gender"].(string),
            BloodType:       source["bloodType"].(string),
            Height:          int(source["height"].(float64)),
            Weight:          int(source["weight"].(float64)),
            Role:            source["role"].(string),
            OTP:             source["otp"].(string),
            OTPExpiration:   int64(source["otpExpiration"].(float64)),
            CreatedAt:       time.Unix(int64(source["createdAt"].(float64)), 0),
            UpdatedAt:       time.Unix(int64(source["updatedAt"].(float64)), 0),
            DeletedAt:       validator.ConvertToTime(source["deletedAt"]),
        }

        userData, err := json.Marshal(user)
        if err == nil {
            uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute)
        }

        return user, nil
    }

    userModel := model.User{}
    result := uqr.db.Where("id = ?", id).First(&userModel)
    if result.Error != nil {
        return entity.User{}, result.Error
    }

    if result.RowsAffected == 0 {
        return entity.User{}, errors.New(constant.ERROR_ID_NOTFOUND)
    }

    userEntity := entity.UserModelToUserEntity(userModel)

    userData, err := json.Marshal(userEntity)
    if err == nil {
        uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute)
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

    query := map[string]interface{}{
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                "email": email,
            },
        },
    }

    var user entity.User
    res, err := uqr.es.Search(
        uqr.es.Search.WithContext(context.Background()),
        uqr.es.Search.WithIndex("users"),
        uqr.es.Search.WithBody(validator.JSONReader(query)),
        uqr.es.Search.WithTrackTotalHits(true),
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

    var r map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
        return entity.User{}, err
    }

    if hits := r["hits"].(map[string]interface{})["hits"].([]interface{}); len(hits) > 0 {
        source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
        user = entity.User{
            ID:              source["id"].(string),
            Fullname:        source["fullname"].(string),
            Email:           source["email"].(string),
            Password:        source["password"].(string),
            NewPassword:     source["newPassword"].(string),
            ConfirmPassword: source["confirmPassword"].(string),
            ProfilePicture:  source["profilePicture"].(string),
            Birthdate:       source["birthdate"].(string),
            Gender:          source["gender"].(string),
            BloodType:       source["bloodType"].(string),
            Height:          int(source["height"].(float64)),
            Weight:          int(source["weight"].(float64)),
            Role:            source["role"].(string),
            OTP:             source["otp"].(string),
            OTPExpiration:   int64(source["otpExpiration"].(float64)),
            CreatedAt:       time.Unix(int64(source["createdAt"].(float64)), 0),
            UpdatedAt:       time.Unix(int64(source["updatedAt"].(float64)), 0),
            DeletedAt:       validator.ConvertToTime(source["deletedAt"]),
        }

        userData, err := json.Marshal(user)
        if err == nil {
            uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute)
        }

        return user, nil
    }

    userModel := model.User{}
    result := uqr.db.Where("email = ?", email).First(&userModel)

    if result.RowsAffected == 0 {
        return entity.User{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
    }

    if result.Error != nil {
        return entity.User{}, result.Error
    }

    userEntity := entity.UserModelToUserEntity(userModel)

    userData, err := json.Marshal(userEntity)
    if err == nil {
        uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute)
    }

    return userEntity, nil
}

