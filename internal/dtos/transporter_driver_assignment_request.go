package dtos

type CreateTransporterDriverAssignmentRequest struct {
	TransporterDriverID  uint64 `json:"transporter_driver_id" validate:"required"`
	TransporterVehicleID uint64 `json:"transporter_vehicle_id" validate:"required"`
	AssignedFrom         string `json:"assigned_from" validate:"required"`
	AssignedTo           string `json:"assigned_to"`
	AssignmentType       string `json:"assignment_type" validate:"omitempty,oneof=PRIMARY TEMPORARY RELIEF EMERGENCY"`
	Active               bool   `json:"active"`
	Notes                string `json:"notes"`
}

type UpdateTransporterDriverAssignmentRequest struct {
	AssignedTo     string `json:"assigned_to"`
	AssignmentType string `json:"assignment_type" validate:"required,oneof=PRIMARY TEMPORARY RELIEF EMERGENCY"`
	Active         bool   `json:"active"`
	Notes          string `json:"notes"`
}
