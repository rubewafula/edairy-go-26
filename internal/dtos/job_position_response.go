package dtos

import "time"

type JobPositionResponse struct {
	ID                uint64    `json:"id"`
	Code              string    `json:"code"`
	Name              string    `json:"name"`
	JobDescription    string    `json:"job_description"`
	DepartmentID      uint64    `json:"department_id"`
	DepartmentName    string    `json:"department_name"`
	GradeID           string    `json:"grade_id"` // Assuming grade_id is a string
	NoOfPosts         int       `json:"no_of_posts"`
	OccupiedPositions int       `json:"occupied_positions"`
	VaccantPositions  int       `json:"vaccant_positions"`
	CreatedBy         uint64    `json:"created_by"`
	UpdatedBy         uint64    `json:"updated_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
