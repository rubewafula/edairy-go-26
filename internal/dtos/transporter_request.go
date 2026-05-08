package dtos

type CreateTransporterRequest struct {
	Name         string `json:"name" validate:"required,max=255"`
	Phone        string `json:"phone" validate:"required,max=15"`
	IDNumber     string `json:"id_number" validate:"required,max=25"`
	VehicleRegNo string `json:"vehicle_reg_no" validate:"required,max=50"`
	Status       string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateTransporterRequest struct {
	Name         string `json:"name" validate:"required,max=255"`
	Phone        string `json:"phone" validate:"required,max=15"`
	IDNumber     string `json:"id_number" validate:"required,max=25"`
	VehicleRegNo string `json:"vehicle_reg_no" validate:"required,max=50"`
	Status       string `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
}
