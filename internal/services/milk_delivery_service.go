package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MilkDeliveryService struct{}

func NewMilkDeliveryService() *MilkDeliveryService {
	return &MilkDeliveryService{}
}

func (s *MilkDeliveryService) CreateDelivery(req dtos.CreateMilkDeliveryRequest) (*models.MilkDelivery, error) {
	delivery := &models.MilkDelivery{
		DeliveryNoteNumber: req.DeliveryNoteNumber,
		CustomerID:         req.CustomerID,
		QuantityAccepted:   req.QuantityAccepted,
		Cooler:             req.Cooler,
		Invoiced:           req.Invoiced,
		TransactionDate:    utils.ParseDate(req.TransactionDate),
		Amount:             req.Amount,
		AmountPaid:         req.AmountPaid,
		RouteID:            req.RouteID,
		Confirmed:          req.Confirmed,
		Processed:          req.Processed,
		TransporterID:      req.TransporterID,
	}

	if err := db.DB.Create(delivery).Error; err != nil {
		return nil, err
	}
	return delivery, nil
}

func (s *MilkDeliveryService) GetDeliveries(page, limit int) ([]dtos.MilkDeliveryResponse, int64, error) {
	var results []dtos.MilkDeliveryResponse
	var total int64
	db.DB.Model(&models.MilkDelivery{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			md.id, md.delivery_note_number, md.customer_id, c.full_names AS customer_name,
			md.quantity_accepted, md.cooler, md.invoiced, md.transaction_date,
			md.amount, md.amount_paid, md.route_id, r.route_name,
			md.confirmed, md.processed, md.transporter_id, t.transporter_no AS transporter_name,
			md.created_at, md.updated_at
		FROM milk_deliveries md
		LEFT JOIN customers c ON md.customer_id = c.id
		LEFT JOIN routes r ON md.route_id = r.id
		LEFT JOIN transporters t ON md.transporter_id = t.id
		WHERE md.deleted_at IS NULL
		ORDER BY md.transaction_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *MilkDeliveryService) GetDelivery(id string) (*dtos.MilkDeliveryResponse, error) {
	var result dtos.MilkDeliveryResponse
	query := `
		SELECT 
			md.id, md.delivery_note_number, md.customer_id, c.full_names AS customer_name,
			md.quantity_accepted, md.cooler, md.invoiced, md.transaction_date,
			md.amount, md.amount_paid, md.route_id, r.route_name,
			md.confirmed, md.processed, md.transporter_id, t.transporter_no AS transporter_name,
			md.created_at, md.updated_at
		FROM milk_deliveries md
		LEFT JOIN customers c ON md.customer_id = c.id
		LEFT JOIN routes r ON md.route_id = r.id
		LEFT JOIN transporters t ON md.transporter_id = t.id
		WHERE md.id = ? AND md.deleted_at IS NULL
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

func (s *MilkDeliveryService) UpdateDelivery(id string, req dtos.UpdateMilkDeliveryRequest) error {
	var delivery models.MilkDelivery
	if err := db.DB.First(&delivery, id).Error; err != nil {
		return err
	}

	delivery.DeliveryNoteNumber = req.DeliveryNoteNumber
	delivery.CustomerID = req.CustomerID
	delivery.QuantityAccepted = req.QuantityAccepted
	delivery.Cooler = req.Cooler
	delivery.Invoiced = req.Invoiced
	delivery.TransactionDate = utils.ParseDate(req.TransactionDate)
	delivery.Amount = req.Amount
	delivery.AmountPaid = req.AmountPaid
	delivery.RouteID = req.RouteID
	delivery.Confirmed = req.Confirmed
	delivery.Processed = req.Processed
	delivery.TransporterID = req.TransporterID

	return db.DB.Save(&delivery).Error
}

func (s *MilkDeliveryService) DeleteDelivery(id string) error {
	return db.DB.Delete(&models.MilkDelivery{}, id).Error
}
