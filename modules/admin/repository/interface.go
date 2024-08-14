package repository

import "talkspace-api/modules/admin/entity"

type AdminCommandRepositoryInterface interface {
	RegisterAdmin(admin entity.Admin) (entity.Admin, error)
	LoginAdmin(email, password string) (entity.Admin, error)
	UpdateAdminPassword(id string, password entity.Admin) (entity.Admin, error)
	NewAdminPassword(email string, password entity.Admin) (entity.Admin, error)
	SendAdminOTP(email string, otp string, expired int64) (entity.Admin, error)
	VerifyAdminOTP(email, otp string) (entity.Admin, error)
	ResetAdminOTP(otp string) (entity.Admin, error)
}

type AdminQueryRepositoryInterface interface {
	GetAdminByID(id string) (entity.Admin, error)
	GetAdminByEmail(email string) (entity.Admin, error)
}
