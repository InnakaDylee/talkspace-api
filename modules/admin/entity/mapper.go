package entity

import "talkspace-api/modules/admin/model"

func AdminEntityToAdminModel(adminEntity Admin) model.Admin {
	adminModel := model.Admin{
		ID:             adminEntity.ID,
		Fullname:       adminEntity.Fullname,
		Email:          adminEntity.Email,
		Password:       adminEntity.Password,
		Role:           adminEntity.Role,
		OTP:            adminEntity.OTP,
		OTPExpiration:  adminEntity.OTPExpiration,
		CreatedAt:      adminEntity.CreatedAt,
		UpdatedAt:      adminEntity.UpdatedAt,
		DeletedAt:      adminEntity.DeletedAt,
	}
	return adminModel
}

func ListAdminEntityToAdminModel(adminEntities []Admin) []model.Admin {
	listAdminModel := []model.Admin{}
	for _, admin := range adminEntities {
		adminModel := AdminEntityToAdminModel(admin)
		listAdminModel = append(listAdminModel, adminModel)
	}
	return listAdminModel
}

func AdminModelToAdminEntity(adminModel model.Admin) Admin {

	adminEntity := Admin{
		ID:             adminModel.ID,
		Fullname:       adminModel.Fullname,
		Email:          adminModel.Email,
		Password:       adminModel.Password,
		Role:           adminModel.Role,
		OTP:            adminModel.OTP,
		OTPExpiration:  adminModel.OTPExpiration,
		CreatedAt:      adminModel.CreatedAt,
		UpdatedAt:      adminModel.UpdatedAt,
		DeletedAt:      adminModel.DeletedAt,
	}
	return adminEntity
}

func ListAdminModelToAdminEntity(adminModels []model.Admin) []Admin {
	listAdminEntity := []Admin{}
	for _, admin := range adminModels {
		adminEntity := AdminModelToAdminEntity(admin)
		listAdminEntity = append(listAdminEntity, adminEntity)
	}
	return listAdminEntity
}
