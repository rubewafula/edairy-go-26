package dtos

import (
	"time"

	"gorm.io/datatypes"
)

type ActivityLogResponse struct {
	ID          uint64         `json:"id"`
	LogName     *string        `json:"log_name"`
	Description string         `json:"description"`
	SubjectType *string        `json:"subject_type"`
	BatchUUID   *string        `json:"batch_uuid"`
	SubjectID   *uint64        `json:"subject_id"`
	CauserType  *string        `json:"causer_type"`
	CauserID    *uint64        `json:"causer_id"`
	CauserName  string         `json:"causer_name"`
	Properties  datatypes.JSON `json:"properties"`
	Event       *string        `json:"event"`
	CreatedAt   *time.Time     `json:"created_at"`
}
