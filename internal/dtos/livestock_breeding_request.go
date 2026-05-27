package dtos

import "time"

type CreateLivestockBreedingRequest struct {
	MaleLivestockID     uint64 `json:"male_livestock_id" validate:"required"`
	FemaleLivestockID   uint64 `json:"female_livestock_id" validate:"required"`
	BreedingDate        string `json:"breeding_date" validate:"required"`
	BreedingType        string `json:"breeding_type" validate:"required"` // AI or Natural
	PregnancyCheckDate  string `json:"pregnancy_check_date"`
	PregnancyStatus     string `json:"pregnancy_status"`
	ExpectedCalvingDate string `json:"expected_calving_date"`
	ActualCalvingDate   string `json:"actual_calving_date"`
	Remarks             string `json:"remarks"`
}

type UpdateLivestockBreedingRequest struct {
	BreedingDate        string `json:"breeding_date"`
	BreedingType        string `json:"breeding_type"`
	MaleLivestockID     uint64 `json:"male_livestock_id" validate:"required"`
	FemaleLivestockID   uint64 `json:"female_livestock_id" validate:"required"`
	PregnancyCheckDate  string `json:"pregnancy_check_date"`
	PregnancyStatus     string `json:"pregnancy_status"`
	ExpectedCalvingDate string `json:"expected_calving_date"`
	ActualCalvingDate   string `json:"actual_calving_date"`
	Remarks             string `json:"remarks"`
}

type LivestockBreedingResponse struct {
	ID                  uint64     `json:"id"`
	LivestockID         uint64     `json:"livestock_id"`
	LivestockTagNo      string     `json:"livestock_tag_no"`
	BreedingDate        time.Time  `json:"breeding_date"`
	BreedingType        string     `json:"breeding_type"`
	SireID              *uint64    `json:"sire_id"`
	TechnicianName      string     `json:"technician_name"`
	PregnancyCheckDate  *time.Time `json:"pregnancy_check_date"`
	PregnancyStatus     string     `json:"pregnancy_status"`
	ExpectedCalvingDate *time.Time `json:"expected_calving_date"`
	ActualCalvingDate   *time.Time `json:"actual_calving_date"`
	Remarks             string     `json:"remarks"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
