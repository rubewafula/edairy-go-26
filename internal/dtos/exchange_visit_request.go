package dtos

type CreateExchangeVisitRequest struct {
	Partner    string `json:"exchange_visit_partner" validate:"required"`
	VisitDate  string `json:"exchange_visit_date" validate:"required"`
	Purpose    string `json:"purpose" validate:"required"`
	Venue      string `json:"venue" validate:"required"`
	EmployeeID uint64 `json:"exchange_visit_employee_id"`
	VisitNotes string `json:"visit_notes"`
}

type UpdateExchangeVisitRequest struct {
	Partner    string `json:"exchange_visit_partner" validate:"required"`
	VisitDate  string `json:"exchange_visit_date" validate:"required,datetime"`
	Purpose    string `json:"purpose" validate:"required"`
	Venue      string `json:"venue" validate:"required"`
	EmployeeID uint64 `json:"exchange_visit_employee_id"`
	VisitNotes string `json:"visit_notes"`
}
