package usecase

import (
	"mime/multipart"
	"talkspace-api/modules/doctor/entity"
)

type DoctorCommandUsecaseInterface interface {
	RegisterDoctor(doctor entity.Doctor, image *multipart.FileHeader) (entity.Doctor, error)
	LoginDoctor(email, password string) (entity.Doctor, string, error)
	UpdateDoctorProfile(id string, doctor entity.Doctor, image *multipart.FileHeader) (entity.Doctor, error)
	UpdateDoctorStatus(id string, status bool) (entity.Doctor, error)
	UpdateDoctorPassword(id string, password entity.Doctor) (entity.Doctor, error)
	NewDoctorPassword(email string, password entity.Doctor) (entity.Doctor, error)
	SendDoctorOTP(email string) (entity.Doctor, error)
	VerifyDoctorOTP(email, otp string) (string, error)
}

type DoctorQueryUsecaseInterface interface {
	GetDoctorByID(id string) (entity.Doctor, error)
	GetAllDoctors(status *bool, specialization string, page, limit int) ([]entity.Doctor, int, error)
}
