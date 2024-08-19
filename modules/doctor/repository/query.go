package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (dqr *doctorQueryRepository) GetAllDoctors(status *bool, specialization string, page, limit int) ([]entity.Doctor, int, error) {
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("doctors:all:status:%v:specialization:%s:page:%d:limit:%d", status, specialization, page, limit)
	cachedDoctors, err := dqr.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedDoctors != "" {
		var cacheResult struct {
			Doctors []entity.Doctor `json:"doctors"`
			Total   int             `json:"total"`
		}
		if err := json.Unmarshal([]byte(cachedDoctors), &cacheResult); err != nil {
			return nil, 0, err
		}
		return cacheResult.Doctors, cacheResult.Total, nil
	}

	var doctorModels []model.Doctor
	query := dqr.db.Model(&model.Doctor{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if specialization != "" {
		query = query.Where("specialization = ?", specialization)
	}

	var totalItems int64
	dqr.db.Model(&model.Doctor{}).Count(&totalItems)

	result := query.Offset(offset).Limit(limit).Find(&doctorModels)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var doctors []entity.Doctor
	for _, doctorModel := range doctorModels {
		doctorEntity := entity.DoctorModelToDoctorEntity(doctorModel)
		doctors = append(doctors, doctorEntity)
	}

	doctorData, err := json.Marshal(struct {
		Doctors []entity.Doctor `json:"doctors"`
		Total   int             `json:"total"`
	}{
		Doctors: doctors,
		Total:   int(totalItems),
	})
	if err == nil {
		dqr.rdb.Set(context.Background(), cacheKey, string(doctorData), 10*time.Minute)
	}

	return doctors, int(totalItems), nil
}
