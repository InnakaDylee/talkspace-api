package dto

type ConsultationRequest struct {
	DoctorID	  string `json:"doctor_id" validate:"required"`
	Status        bool   `json:"status" validate:"required"`
}