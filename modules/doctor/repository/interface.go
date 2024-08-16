package repository

import "talkspace-api/modules/doctor/entity"

type DoctorCommandRepositoryInterface interface {
	LoginDoctor(email, password string) (entity.Doctor, error)
	UpdateDoctorProfile(id string, doctor entity.Doctor) (entity.Doctor, error)
	UpdateDoctorStatus(id string, status bool) (entity.Doctor, error)
	UpdateDoctorPassword(id string, password entity.Doctor) (entity.Doctor, error)
	NewDoctorPassword(email string, password entity.Doctor) (entity.Doctor, error)
	SendDoctorOTP(email string, otp string, expired int64) (entity.Doctor, error)
	VerifyDoctorOTP(email, otp string) (entity.Doctor, error)
	ResetDoctorOTP(otp string) (entity.Doctor, error)
}

type DoctorQueryRepositoryInterface interface {
	GetDoctorByID(id string) (entity.Doctor, error)
	GetDoctorByEmail(email string) (entity.Doctor, error)
	GetAllDoctors(status *bool, specialization string) ([]entity.Doctor, error)
}
