package dtos

import "time"

type TransporterDriverAssignmentResponse struct {
	ID                   uint64     `json:"ID"`
	TransporterDriverID  uint64     `json:"TransporterDriverID"`
	DriverName           string     `json:"DriverName"`
	DriverNo             string     `json:"DriverNo"`
	TransporterVehicleID uint64     `json:"TransporterVehicleID"`
	VehicleRegNo         string     `json:"VehicleRegNo"`
	AssignedFrom         time.Time  `json:"AssignedFrom"`
	AssignedTo           *time.Time `json:"AssignedTo"`
	AssignmentType       string     `json:"AssignmentType"`
	Active               bool       `json:"Active"`
	Notes                string     `json:"Notes"`
	CreatedAt            time.Time  `json:"CreatedAt"`
}
