package usecase

import "talkspace-api/modules/admin/entity"

type AdminCommandUsecaseInterface interface {
	RegisterAdmin(admin entity.Admin) (entity.Admin, error)
	LoginAdmin(email, password string) (entity.Admin, string, error)
	UpdateAdminPassword(id string, password entity.Admin) (entity.Admin, error)
	NewAdminPassword(email string, password entity.Admin) (entity.Admin, error)
	SendAdminOTP(email string) (entity.Admin, error)
	VerifyAdminOTP(email, otp string) (string, error)
}

type AdminQueryUsecaseInterface interface {
	GetAdminByID(id string) (entity.Admin, error)
}
