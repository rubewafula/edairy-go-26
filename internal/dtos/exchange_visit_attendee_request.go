package dtos

type CreateExchangeVisitAttendeeRequest struct {
	ExchangeVisitID      uint64 `json:"exchange_visit_id" validate:"required"`
	Attendee             string `json:"attendee" validate:"required"`
	AttendeeOrganization string `json:"attendee_organization" validate:"required"`
	AttendeeDesignation  string `json:"attendee_designation" validate:"required"`
	Attended             string `json:"attended" validate:"omitempty"`
	Comments             string `json:"comments"`
	AttendanceEmployeeID uint64 `json:"attendance_employee_id"`
}

type UpdateExchangeVisitAttendeeRequest struct {
	ExchangeVisitID      uint64 `json:"exchange_visit_id" validate:"required"`
	Attendee             string `json:"attendee" validate:"required"`
	AttendeeOrganization string `json:"attendee_organization" validate:"required"`
	AttendeeDesignation  string `json:"attendee_designation" validate:"required"`
	Attended             string `json:"attended" validate:"required"`
	Comments             string `json:"comments"`
	AttendanceEmployeeID uint64 `json:"attendance_employee_id"`
}
