package dtos

import "time"

type ExchangeVisitResponse struct {
	ID                uint64    `json:"ID"`
	Partner           string    `json:"Partner"`
	VisitDate         time.Time `json:"VisitDate"`
	Purpose           string    `json:"Purpose"`
	Venue             string    `json:"Venue"`
	EmployeeID        uint64    `json:"EmployeeID"`
	EmployeeFirstName string    `json:"EmployeeFirstName"`
	EmployeeSurname   string    `json:"EmployeeSurname"`
	VisitNotes        string    `json:"VisitNotes"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
}
