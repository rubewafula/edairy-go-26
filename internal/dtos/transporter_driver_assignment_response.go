package dtos

import "time"

type TransporterDriverAssignmentResponse struct {
	ID                   uint64     `json:"id"`
	TransporterDriverID  uint64     `json:"transporter_driver_id"`
	DriverName           string     `json:"driver_name"`
	DriverNo             string     `json:"driver_no"`
	TransporterVehicleID uint64     `json:"transporter_vehicle_id"`
	VehicleRegNo         string     `json:"vehicle_reg_no"`
	AssignedFrom         time.Time  `json:"assigned_from"`
	AssignedTo           *time.Time `json:"assigned_to"`
	AssignmentType       string     `json:"assignment_type"`
	Active               bool       `json:"active"`
	Notes                string     `json:"notes"`
	CreatedAt            time.Time  `json:"created_at"`
}
