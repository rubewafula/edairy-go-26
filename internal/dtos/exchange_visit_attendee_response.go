package dtos

import "time"

type ExchangeVisitAttendeeResponse struct {
	ID                   uint64    `json:"id"`
	ExchangeVisitID      uint64    `json:"exchange_visit_id"`
	Partner              string    `json:"partner"`
	Attendee             string    `json:"attendee"`
	AttendeeOrganization string    `json:"attendee_organization"`
	AttendeeDesignation  string    `json:"attendee_designation"`
	Attended             string    `json:"attended"`
	Comments             string    `json:"comments"`
	AttendanceEmployeeID uint64    `json:"attendance_employee_id"`
	EmployeeFirstName    string    `json:"employee_first_name"`
	EmployeeSurname      string    `json:"employee_surname"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
