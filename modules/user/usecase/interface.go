package usecase

import "talkspace-api/modules/user/entity"

type UserCommandUsecaseInterface interface {
	RegisterUser(user entity.User) (entity.User, error)
	LoginUser(email, password string) (entity.User, string, error)
	UpdateUserByID(id string, user entity.User) (entity.User, error)
	UpdateUserPassword(id string, password entity.User) (entity.User, error)
	NewUserPassword(email string, password entity.User) (entity.User, error)
	SendUserOTP(email string) (entity.User, error)
	VerifyUserOTP(email, otp string) (string, error)
}

type UserQueryUsecaseInterface interface {
	GetUserByID(id string) (entity.User, error)
}
