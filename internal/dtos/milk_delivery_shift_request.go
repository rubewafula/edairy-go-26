package dtos

type CreateMilkDeliveryShiftRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

type UpdateMilkDeliveryShiftRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}
