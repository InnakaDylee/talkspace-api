package entity

import "talkspace-api/modules/user/model"

func UserEntityToUserModel(userEntity User) model.User {
	userModel := model.User{
		Fullname:       userEntity.Fullname,
		Email:          userEntity.Email,
		Password:       userEntity.Password,
		ProfilePicture: userEntity.ProfilePicture,
		Birthdate:      userEntity.Birthdate,
		Gender:         userEntity.Gender,
		BloodType:      userEntity.BloodType,
		Height:         userEntity.Height,
		Weight:         userEntity.Weight,
		OTP:            userEntity.OTP,
		OTPExpiration:  userEntity.OTPExpiration,
	}
	return userModel
}

func ListUserEntityToUserModel(userEntity []User) []model.User {
	listUserModel := []model.User{}
	for _, user := range userEntity {
		userModel := UserEntityToUserModel(user)
		listUserModel = append(listUserModel, userModel)
	}
	return listUserModel
}

func UserModelToUserEntity(userModel model.User) User {
	userEntity := User{
		ID:             userModel.ID,
		Fullname:       userModel.Fullname,
		Email:          userModel.Email,
		Password:       userModel.Password,
		ProfilePicture: userModel.ProfilePicture,
		Birthdate:      userModel.Birthdate,
		Gender:         userModel.Gender,
		BloodType:      userModel.BloodType,
		Height:         userModel.Height,
		Weight:         userModel.Weight,
		OTP:            userModel.OTP,
		OTPExpiration:  userModel.OTPExpiration,
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
		DeletedAt:      userModel.DeletedAt,
	}
	return userEntity
}

func ListUserModelToUserEntity(userModel []model.User) []User {
	listUserEntity := []User{}
	for _, user := range userModel {
		userEntity := UserModelToUserEntity(user)
		listUserEntity = append(listUserEntity, userEntity)
	}
	return listUserEntity
}
