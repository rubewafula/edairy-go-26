package dtos

import "time"

type ExchangeVisitAttendeeResponse struct {
	ID                   uint64    `json:"ID"`
	ExchangeVisitID      uint64    `json:"ExchangeVisitID"`
	Partner              string    `json:"Partner"`
	Attendee             string    `json:"Attendee"`
	AttendeeOrganization string    `json:"AttendeeOrganization"`
	AttendeeDesignation  string    `json:"AttendeeDesignation"`
	Attended             string    `json:"Attended"`
	Comments             string    `json:"Comments"`
	AttendanceEmployeeID uint64    `json:"AttendanceEmployeeID"`
	EmployeeFirstName    string    `json:"EmployeeFirstName"`
	EmployeeSurname      string    `json:"EmployeeSurname"`
	CreatedAt            time.Time `json:"CreatedAt"`
	UpdatedAt            time.Time `json:"UpdatedAt"`
}
