package entity

import "talkspace-api/modules/user/model"

func UserEntityToUserModel(userEntity User) model.User {
	var gender *string
	if userEntity.Gender != "" {
		genderValue := userEntity.Gender
		gender = &genderValue
	}

	var bloodType *string
	if userEntity.BloodType != "" {
		bloodTypeValue := userEntity.BloodType
		bloodType = &bloodTypeValue
	}

	userModel := model.User{
		Fullname:       userEntity.Fullname,
		Email:          userEntity.Email,
		Password:       userEntity.Password,
		ProfilePicture: userEntity.ProfilePicture,
		Birthdate:      userEntity.Birthdate,
		Gender:         gender,
		BloodType:      bloodType,
		Height:         userEntity.Height,
		Weight:         userEntity.Weight,
		Role:           userEntity.Role,
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
	var gender string
	if userModel.Gender != nil {
		gender = *userModel.Gender
	}

	var bloodType string
	if userModel.BloodType != nil {
		bloodType = *userModel.BloodType
	}
	userEntity := User{
		ID:             userModel.ID,
		Fullname:       userModel.Fullname,
		Email:          userModel.Email,
		Password:       userModel.Password,
		ProfilePicture: userModel.ProfilePicture,
		Birthdate:      userModel.Birthdate,
		Gender:         gender,
		BloodType:      bloodType,
		Height:         userModel.Height,
		Weight:         userModel.Weight,
		Role:           userModel.Role,
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
