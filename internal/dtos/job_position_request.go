package dtos

type CreateJobPositionRequest struct {
	Code           string `json:"code" validate:"required,max=255"`
	Name           string `json:"name" validate:"required,max=255"`
	JobDescription string `json:"job_description" validate:"max=255"`
	DepartmentID   uint64 `json:"department_id" validate:"required"`
	GradeID        string `json:"grade_id" validate:"max=255"`
	NoOfPosts      int    `json:"no_of_posts" validate:"min=0"`
}

type UpdateJobPositionRequest struct {
	Code              string `json:"code" validate:"required,max=255"`
	Name              string `json:"name" validate:"required,max=255"`
	JobDescription    string `json:"job_description" validate:"max=255"`
	DepartmentID      uint64 `json:"department_id" validate:"required"`
	GradeID           string `json:"grade_id" validate:"max=255"`
	NoOfPosts         int    `json:"no_of_posts" validate:"min=0"`
	OccupiedPositions int    `json:"occupied_positions" validate:"min=0"`
}
