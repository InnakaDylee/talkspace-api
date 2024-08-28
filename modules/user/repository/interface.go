package repository

import (
	"mime/multipart"
	"talkspace-api/modules/user/entity"
)

type UserCommandRepositoryInterface interface {
	RegisterUser(user entity.User) (entity.User, error)
	LoginUser(email, password string) (entity.User, error)
	UpdateUserProfile(id string, user entity.User, image *multipart.FileHeader) (entity.User, error)
	UpdateUserPassword(id string, password entity.User) (entity.User, error)
	NewUserPassword(email string, password entity.User) (entity.User, error)
	SendUserOTP(email string, otp string, expired int64) (entity.User, error)
	VerifyUserOTP(email, otp string) (entity.User, error)
	ResetUserOTP(otp string) (entity.User, error)
	RequestPremium(user entity.User, request_premium string) (entity.User, error)
	UpdateUserPremiumExpired(id string, status string) (entity.User, error)
}

type UserQueryRepositoryInterface interface {
	GetUserByID(id string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetRequestPremiumUsers() ([]entity.User, error)
}
