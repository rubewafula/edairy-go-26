package dtos

import "time"

type ExchangeVisitResponse struct {
	ID                uint64    `json:"id"`
	Partner           string    `json:"partner"`
	VisitDate         time.Time `json:"visit_date"`
	Purpose           string    `json:"purpose"`
	Venue             string    `json:"venue"`
	EmployeeID        uint64    `json:"employee_id"`
	EmployeeFirstName string    `json:"employee_first_name"`
	EmployeeSurname   string    `json:"employee_surname"`
	VisitNotes        string    `json:"visit_notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
