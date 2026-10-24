package dtos

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CreateDefaultMilkRateRequest struct {
	Rate        float64  `json:"rate" validate:"required"`
	AdvanceRate *FlexFloat `json:"advance_rate"`
	RouteID     uint64     `json:"route_id"`
	MemberID    uint64   `json:"member_id"`
}

type UpdateDefaultMilkRateRequest struct {
	Rate        float64  `json:"rate" validate:"required"`
	AdvanceRate *FlexFloat `json:"advance_rate"`
	RouteID     uint64     `json:"route_id" `
	MemberID    uint64   `json:"member_id"`
}

type DefaultMilkRateResponse struct {
	ID         uint64    `json:"id"`
	Rate        float64   `json:"rate"`
	AdvanceRate *float64  `json:"advance_rate"`
	RouteID     uint64    `json:"route_id"`
	MemberID    uint64    `json:"member_id"`
	RouteName  string    `json:"route_name"`
	MemberName string    `json:"member_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FlexFloat accepts a JSON number, a numeric string ("40") or ""/null.
type FlexFloat float64

func (f *FlexFloat) UnmarshalJSON(b []byte) error {
	s := strings.Trim(strings.TrimSpace(string(b)), `"`)
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("advance_rate must be a number")
	}
	*f = FlexFloat(v)
	return nil
}

// Ptr returns the value as *float64 (nil when unset).
func (f *FlexFloat) Ptr() *float64 {
	if f == nil {
		return nil
	}
	v := float64(*f)
	return &v
}
