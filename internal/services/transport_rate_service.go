package services

import (
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type TransportRateService struct{}

func NewTransportRateService() *TransportRateService {
	return &TransportRateService{}
}

func (s *TransportRateService) CreateTransportRate(req dtos.CreateTransportRateRequest) (*models.TransportRate, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	rate := &models.TransportRate{
		RouteID:       req.RouteID,
		TransporterID: req.TransporterID,
		Rate:          req.TransportRate,
		MemberID:      req.MemberID,
		Status:        status,
	}

	if err := db.DB.Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *TransportRateService) GetTransportRates(page, limit int, transporterNo, routeID, memberNo string) ([]dtos.TransportRateResponse, int64, error) {
	var results []dtos.TransportRateResponse
	var total int64

	offset := (page - 1) * limit

	baseQuery := `
		SELECT 
			tr.id, m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			r.route_name,
			t.transporter_no,
			tr.transport_rate AS rate, tr.status, 
			tr.created_at, tr.updated_at
		FROM transport_rates tr
		LEFT JOIN member_registrations m ON tr.member_id = m.id
		LEFT JOIN transporters t ON tr.transporter_id = t.id
		LEFT JOIN routes r ON tr.route_id = r.id
	`

	baseCountQuery := `
		SELECT COUNT(*)
		FROM transport_rates tr
		LEFT JOIN member_registrations m ON tr.member_id = m.id
		LEFT JOIN transporters t ON tr.transporter_id = t.id
		LEFT JOIN routes r ON tr.route_id = r.id
	`

	var args []interface{}
	whereClauses := []string{"tr.deleted_at IS NULL"}

	if transporterNo != "" {
		whereClauses = append(whereClauses, "t.transporter_no LIKE ?")
		args = append(args, "%"+transporterNo+"%")
	}
	if routeID != "" {
		whereClauses = append(whereClauses, "tr.route_id = ?")
		args = append(args, routeID)
	}
	if memberNo != "" {
		whereClauses = append(whereClauses, "m.member_no LIKE ?")
		args = append(args, "%"+memberNo+"%")
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	if err := db.DB.Raw(baseCountQuery+whereSql, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	fullQuery := baseQuery + whereSql + " ORDER BY tr.id DESC LIMIT ? OFFSET ?"
	if err := db.DB.Raw(fullQuery, append(args, limit, offset)...).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *TransportRateService) GetTransportRate(id string) (*dtos.TransportRateResponse, error) {
	var result dtos.TransportRateResponse
	query := `
		SELECT 
			tr.id, m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			r.route_name,
			t.transporter_no,
			tr.transport_rate AS rate, tr.status, 
			tr.created_at, tr.updated_at
		FROM transport_rates tr
		LEFT JOIN member_registrations m ON tr.member_id = m.id
		LEFT JOIN transporters t ON tr.transporter_id = t.id
		LEFT JOIN routes r ON tr.route_id = r.id
		WHERE tr.id = ? AND tr.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *TransportRateService) UpdateTransportRate(id string, req dtos.UpdateTransportRateRequest) error {
	var rate models.TransportRate
	if err := db.DB.First(&rate, id).Error; err != nil {
		return err
	}

	rate.RouteID = req.RouteID
	rate.TransporterID = req.TransporterID
	rate.Rate = req.TransportRate
	rate.MemberID = req.MemberID
	rate.Status = req.Status

	return db.DB.Save(&rate).Error
}

func (s *TransportRateService) DeleteTransportRate(id string) error {
	return db.DB.Delete(&models.TransportRate{}, id).Error
}
