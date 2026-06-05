package dtos

type CreateShareTypeRequest struct {
	ShareCode         string  `json:"share_code"`
	ShareType         string  `json:"share_type"`
	Description       string  `json:"description"`
	Rate              float64 `json:"rate"`
	Mandatory         int     `json:"mandatory"`
	HasShareValue     int     `json:"has_share_value"`
	RepayMethod       string  `json:"repay_method"`
	CalculatingMethod string  `json:"calculating_method"`
	ShareValue        float64 `json:"share_value"`
	DeductionTypeID   uint64  `json:"deduction_type_id"`
	Priority          int     `json:"priority"`
}

type UpdateShareTypeRequest struct {
	ShareCode         string  `json:"share_code"`
	ShareType         string  `json:"share_type"`
	Description       string  `json:"description"`
	Rate              float64 `json:"rate"`
	Mandatory         int     `json:"mandatory"`
	HasShareValue     int     `json:"has_share_value"`
	RepayMethod       string  `json:"repay_method"`
	CalculatingMethod string  `json:"calculating_method"`
	ShareValue        float64 `json:"share_value"`
	DeductionTypeID   uint64  `json:"deduction_type_id"`
	Priority          int     `json:"priority"`
}
