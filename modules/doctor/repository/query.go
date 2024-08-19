package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"talkspace-api/modules/doctor/entity"
	"talkspace-api/modules/doctor/model"
	"talkspace-api/utils/constant"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type doctorQueryRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewDoctorQueryRepository(db *gorm.DB, rdb *redis.Client) DoctorQueryRepositoryInterface {
	return &doctorQueryRepository{
		db:  db,
		rdb: rdb,
	}
}

func (dqr *doctorQueryRepository) GetDoctorByID(id string) (entity.Doctor, error) {
    cacheKey := "doctor:" + id
    cachedDoctor, err := dqr.rdb.Get(context.Background(), cacheKey).Result()
    if err == nil && cachedDoctor != "" {
        var doctor entity.Doctor
        if err := json.Unmarshal([]byte(cachedDoctor), &doctor); err != nil {
            return entity.Doctor{}, err
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

	var doctorModels []model.Doctor
	query := dqr.db.Model(&model.Doctor{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if specialization != "" {
		query = query.Where("specialization = ?", specialization)
	}

	result := query.Find(&doctorModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var doctors []entity.Doctor
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
