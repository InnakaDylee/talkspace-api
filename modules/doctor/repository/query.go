package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"talkspace-api/modules/doctor/entity"
	"talkspace-api/modules/doctor/model"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/validator"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type doctorQueryRepository struct {
	db  *gorm.DB
	es  *elasticsearch.Client
	rdb *redis.Client
}

func NewDoctorQueryRepository(db *gorm.DB, es *elasticsearch.Client, rdb *redis.Client) DoctorQueryRepositoryInterface {
	return &doctorQueryRepository{
		db:  db,
		es:  es,
		rdb: rdb,
	}
}

func (dqr *doctorQueryRepository) GetDoctorByID(id string) (entity.Doctor, error) {
	cacheKey := "doctor:id:" + id
	cachedDoctor, err := dqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedDoctor != "" {
		var doctor entity.Doctor
		if err := json.Unmarshal([]byte(cachedDoctor), &doctor); err != nil {
			return entity.Doctor{}, err
		}
		return doctor, nil
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": id,
			},
		},
	}

	var doctor entity.Doctor
	res, err := dqr.es.Search(
		dqr.es.Search.WithContext(context.Background()),
		dqr.es.Search.WithIndex("doctors"),
		dqr.es.Search.WithBody(validator.JSONReader(query)),
		dqr.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return entity.Doctor{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return entity.Doctor{}, err
		} else {
			return entity.Doctor{}, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return entity.Doctor{}, err
	}

	if hits := r["hits"].(map[string]interface{})["hits"].([]interface{}); len(hits) > 0 {
		source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
		doctor = entity.Doctor{
			ID:                source["id"].(string),
			Fullname:          source["fullname"].(string),
			Email:             source["email"].(string),
			Password:          source["password"].(string),
			ProfilePicture:    source["profilePicture"].(string),
			Status:            source["status"].(bool),
			Gender:            source["gender"].(string),
			Specialization:    source["specialization"].(string),
			YearsOfExperience: source["yearsOfExperience"].(string),
			Price:             source["price"].(float64),
			LicenseNumber:     source["licenseNumber"].(string),
			Alumnus:           source["alumnus"].(string),
			About:             source["about"].(string),
			Location:          source["location"].(string),
			Role:              source["role"].(string),
			OTP:               source["otp"].(string),
			OTPExpiration:     int64(source["otpExpiration"].(float64)),
			CreatedAt:         time.Unix(int64(source["createdAt"].(float64)), 0),
			UpdatedAt:         time.Unix(int64(source["updatedAt"].(float64)), 0),
			DeletedAt:         validator.ConvertToTime(source["deletedAt"]),
		}

		doctorData, err := json.Marshal(doctor)
		if err == nil {
			dqr.rdb.Set(context.Background(), cacheKey, string(doctorData), 10*time.Minute)
		}

		return doctor, nil
	}

	doctorModel := model.Doctor{}
	result := dqr.db.Where("id = ?", id).First(&doctorModel)
	if result.Error != nil {
		return entity.Doctor{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entity.Doctor{}, errors.New(constant.ERROR_ID_NOTFOUND)
	}

	doctorEntity := entity.DoctorModelToDoctorEntity(doctorModel)

	doctorData, err := json.Marshal(doctorEntity)
	if err == nil {
		dqr.rdb.Set(context.Background(), cacheKey, string(doctorData), 10*time.Minute)
	}

	return doctorEntity, nil
}

func (dqr *doctorQueryRepository) GetDoctorByEmail(email string) (entity.Doctor, error) {
	cacheKey := "doctor:email:" + email
	cachedDoctor, err := dqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedDoctor != "" {
		var doctor entity.Doctor
		if err := json.Unmarshal([]byte(cachedDoctor), &doctor); err != nil {
			return entity.Doctor{}, err
		}
		return doctor, nil
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"email": email,
			},
		},
	}

	var doctor entity.Doctor
	res, err := dqr.es.Search(
		dqr.es.Search.WithContext(context.Background()),
		dqr.es.Search.WithIndex("doctors"),
		dqr.es.Search.WithBody(validator.JSONReader(query)),
		dqr.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return entity.Doctor{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return entity.Doctor{}, err
		} else {
			return entity.Doctor{}, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return entity.Doctor{}, err
	}

	if hits := r["hits"].(map[string]interface{})["hits"].([]interface{}); len(hits) > 0 {
		source := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
		doctor = entity.Doctor{
			ID:                source["id"].(string),
			Fullname:          source["fullname"].(string),
			Email:             source["email"].(string),
			Password:          source["password"].(string),
			ProfilePicture:    source["profilePicture"].(string),
			Status:            source["status"].(bool),
			Gender:            source["gender"].(string),
			Specialization:    source["specialization"].(string),
			YearsOfExperience: source["yearsOfExperience"].(string),
			Price:             source["price"].(float64),
			LicenseNumber:     source["licenseNumber"].(string),
			Alumnus:           source["alumnus"].(string),
			About:             source["about"].(string),
			Location:          source["location"].(string),
			Role:              source["role"].(string),
			OTP:               source["otp"].(string),
			OTPExpiration:     int64(source["otpExpiration"].(float64)),
			CreatedAt:         time.Unix(int64(source["createdAt"].(float64)), 0),
			UpdatedAt:         time.Unix(int64(source["updatedAt"].(float64)), 0),
			DeletedAt:         validator.ConvertToTime(source["deletedAt"]),
		}

		doctorData, err := json.Marshal(doctor)
		if err == nil {
			dqr.rdb.Set(context.Background(), cacheKey, string(doctorData), 10*time.Minute)
		}

		return doctor, nil
	}

	doctorModel := model.Doctor{}
	result := dqr.db.Where("email = ?", email).First(&doctorModel)

	if result.RowsAffected == 0 {
		return entity.Doctor{}, errors.New(constant.ERROR_EMAIL_NOTFOUND)
	}

	if result.Error != nil {
		return entity.Doctor{}, result.Error
	}

	doctorEntity := entity.DoctorModelToDoctorEntity(doctorModel)

	doctorData, err := json.Marshal(doctorEntity)
	if err == nil {
		dqr.rdb.Set(context.Background(), cacheKey, string(doctorData), 10*time.Minute)
	}

	return doctorEntity, nil
}

func (dqr *doctorQueryRepository) GetAllDoctors(status *bool, specialization string) ([]entity.Doctor, error) {
	cacheKey := "doctors:all"
	if status != nil {
		cacheKey += ":status:" + strconv.FormatBool(*status)
	}
	if specialization != "" {
		cacheKey += ":specialization:" + specialization
	}
	cachedDoctors, err := dqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedDoctors != "" {
		var doctors []entity.Doctor
		if err := json.Unmarshal([]byte(cachedDoctors), &doctors); err != nil {
			return nil, err
		}
		return doctors, nil
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"match_all": map[string]interface{}{},
					},
				},
				"filter": []interface{}{
					func() interface{} {
						if status != nil {
							return map[string]interface{}{
								"term": map[string]interface{}{
									"status": *status,
								},
							}
						}
						return nil
					}(),
					func() interface{} {
						if specialization != "" {
							return map[string]interface{}{
								"match": map[string]interface{}{
									"specialization": specialization,
								},
							}
						}
						return nil
					}(),
				},
			},
		},
	}

	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = validator.RemoveNilValues(query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{}))

	var doctors []entity.Doctor
	res, err := dqr.es.Search(
		dqr.es.Search.WithContext(context.Background()),
		dqr.es.Search.WithIndex("doctors"),
		dqr.es.Search.WithBody(validator.JSONReader(query)),
		dqr.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			return nil, errors.New(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if hits := r["hits"].(map[string]interface{})["hits"].([]interface{}); len(hits) > 0 {
		for _, hit := range hits {
			source := hit.(map[string]interface{})["_source"].(map[string]interface{})
			doctor := entity.Doctor{
				ID:                source["id"].(string),
				Fullname:          source["fullname"].(string),
				Email:             source["email"].(string),
				Password:          source["password"].(string),
				ProfilePicture:    source["profilePicture"].(string),
				Status:            source["status"].(bool),
				Gender:            source["gender"].(string),
				Specialization:    source["specialization"].(string),
				YearsOfExperience: source["yearsOfExperience"].(string),
				Price:             source["price"].(float64),
				LicenseNumber:     source["licenseNumber"].(string),
				Alumnus:           source["alumnus"].(string),
				About:             source["about"].(string),
				Location:          source["location"].(string),
				Role:              source["role"].(string),
				OTP:               source["otp"].(string),
				OTPExpiration:     int64(source["otpExpiration"].(float64)),
				CreatedAt:         time.Unix(int64(source["createdAt"].(float64)), 0),
				UpdatedAt:         time.Unix(int64(source["updatedAt"].(float64)), 0),
				DeletedAt:         validator.ConvertToTime(source["deletedAt"]),
			}
			doctors = append(doctors, doctor)
		}

		doctorData, err := json.Marshal(doctors)
		if err == nil {
			dqr.rdb.Set(context.Background(), cacheKey, string(doctorData), 10*time.Minute)
		}

		return doctors, nil
	}

	var doctorModels []model.Doctor
	result := dqr.db.Find(&doctorModels)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, doctorModel := range doctorModels {
		doctorEntity := entity.DoctorModelToDoctorEntity(doctorModel)
		doctors = append(doctors, doctorEntity)
	}

	doctorData, err := json.Marshal(doctors)
	if err == nil {
		dqr.rdb.Set(context.Background(), cacheKey, string(doctorData), 10*time.Minute)
	}

	return doctors, nil
}
