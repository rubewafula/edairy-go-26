package dtos

type CreateCustomerPayDateRangeRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	StartDate string `json:"start_date" validate:"required,datetime"`
	EndDate   string `json:"end_date" validate:"required,datetime"`
	PayMonth  string `json:"pay_month" validate:"required"`
	PayYear   string `json:"pay_year" validate:"required"`
}

type UpdateCustomerPayDateRangeRequest struct {
	Name      string `json:"name" validate:"required,max=255"`
	StartDate string `json:"start_date" validate:"required,datetime"`
	EndDate   string `json:"end_date" validate:"required,datetime"`
	PayMonth  string `json:"pay_month" validate:"required"`
	PayYear   string `json:"pay_year" validate:"required"`
}
