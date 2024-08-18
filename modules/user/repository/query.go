package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		fmt.Println("Cache hit: ", cachedUser)
		var user entity.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			fmt.Println("Error unmarshalling cached user: ", err)
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

	fmt.Println("Query Elasticsearch: ", query)
	var user entity.User
	res, err := uqr.es.Search(
		uqr.es.Search.WithContext(context.Background()),
		uqr.es.Search.WithIndex("users"),
		uqr.es.Search.WithBody(validator.JSONReader(query)),
		uqr.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		fmt.Println("Elasticsearch query error: ", err)
		return entity.User{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
        // Decode the Elasticsearch error response into a map
        var e map[string]interface{}
        fmt.Println("Attempting to decode Elasticsearch error response...")
        if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
            fmt.Println("Error decoding Elasticsearch error response: ", err)
            return entity.User{}, err
        }
    
        // Check if the "error" key is present in the map and get the error reason
        fmt.Println("Checking for 'error' and 'reason' keys in response...")
        if errorDetail, ok := e["error"].(map[string]interface{}); ok {
            if reason, ok := errorDetail["reason"].(string); ok {
                fmt.Println("Elasticsearch error reason: ", reason)
                return entity.User{}, errors.New(reason)
            }
            fmt.Println("Error 'reason' key is not a string.")
        } else {
            fmt.Println("'error' key is not found or is not a map.")
        }
    
        // If the 'reason' key is not found, return a generic error message
        fmt.Println("No 'reason' key found in error response.")
        return entity.User{}, errors.New("unknown error from Elasticsearch")
    }
    

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println("Error decoding Elasticsearch response: ", err)
		return entity.User{}, err
	}

	fmt.Println("Elasticsearch response: ", r)
	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		fmt.Println("Unexpected hits format: ", r["hits"])
		return entity.User{}, errors.New("unexpected hits format")
	}

	if len(hits) == 0 {
		fmt.Println("No hits found, querying database")
		userModel := model.User{}
		result := uqr.db.Where("email = ?", email).First(&userModel)

		if result.RowsAffected == 0 {
			return entity.User{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
		}

		if result.Error != nil {
			fmt.Println("Database query error: ", result.Error)
			return entity.User{}, result.Error
		}

		userEntity := entity.UserModelToUserEntity(userModel)

		userData, err := json.Marshal(userEntity)
		if err == nil {
			uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute)
		}

		return userEntity, nil
	}

	fmt.Println("Hits data: ", hits)
	source, ok := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
	if !ok {
		fmt.Println("Unexpected source format: ", hits[0])
		return entity.User{}, errors.New("unexpected data format from Elasticsearch")
	}

	fmt.Println("Source data: ", source)
	user.ID = validator.GetStringFromMap(source, "id")
	user.Fullname = validator.GetStringFromMap(source, "fullname")
	user.Email = validator.GetStringFromMap(source, "email")
	user.Password = validator.GetStringFromMap(source, "password")
	user.NewPassword = validator.GetStringFromMap(source, "newPassword")
	user.ConfirmPassword = validator.GetStringFromMap(source, "confirmPassword")
	user.ProfilePicture = validator.GetStringFromMap(source, "profilePicture")
	user.Birthdate = validator.GetStringFromMap(source, "birthdate")
	user.Gender = validator.GetStringFromMap(source, "gender")
	user.BloodType = validator.GetStringFromMap(source, "bloodType")
	user.Height = validator.GetIntFromMap(source, "height")
	user.Weight = validator.GetIntFromMap(source, "weight")
	user.Role = validator.GetStringFromMap(source, "role")
	user.OTP = validator.GetStringFromMap(source, "otp")
	user.OTPExpiration = validator.GetInt64FromMap(source, "otpExpiration")
	user.CreatedAt = time.Unix(validator.GetInt64FromMap(source, "createdAt"), 0)
	user.UpdatedAt = time.Unix(validator.GetInt64FromMap(source, "updatedAt"), 0)
	user.DeletedAt = validator.ConvertToTime(source["deletedAt"])

	userData, err := json.Marshal(user)
	if err == nil {
		uqr.rdb.Set(context.Background(), cacheKey, string(userData), 10*time.Minute)
	}

	return user, nil
}
