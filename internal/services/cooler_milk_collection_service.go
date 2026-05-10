package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type CoolerMilkCollectionService struct{}

func NewCoolerMilkCollectionService() *CoolerMilkCollectionService {
	return &CoolerMilkCollectionService{}
}

func (s *CoolerMilkCollectionService) CreateCollection(req dtos.CreateCoolerMilkCollectionRequest) (*models.CoolerMilkCollection, error) {
	collection := &models.CoolerMilkCollection{
		TransactionDate:     utils.ParseDate(req.TransactionDate),
		Quantity:            req.Quantity,
		TransportVehicleID:  req.TransportVehicleID,
		MilkDeliveryShiftID: req.MilkDeliveryShiftID,
		Confirmed:           req.Confirmed,
		SiteID:              req.SiteID,
		TransporterID:       req.TransporterID,
		RouteID:             req.RouteID,
	}

	if err := db.DB.Create(collection).Error; err != nil {
		return nil, err
	}
	return collection, nil
}

func (s *CoolerMilkCollectionService) GetCollections(page, limit int) ([]dtos.CoolerMilkCollectionResponse, int64, error) {
	var results []dtos.CoolerMilkCollectionResponse
	var total int64
	db.DB.Model(&models.CoolerMilkCollection{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			cmc.id, cmc.transaction_date, cmc.quantity, 
			cmc.transport_vehicle_id, tv.registration_no AS vehicle_reg_no,
			cmc.milk_delivery_shift_id, mds.name AS milk_delivery_shift,
			cmc.confirmed, cmc.site_id, 
			cmc.transporter_id, t.transporter_no,
			cmc.route_id, r.route_name,
			cmc.created_at, cmc.updated_at
		FROM cooler_milk_collections cmc
		LEFT JOIN routes r ON cmc.route_id = r.id
		LEFT JOIN transporters t ON cmc.transporter_id = t.id
		LEFT JOIN transporter_vehicles tv ON cmc.transport_vehicle_id = tv.id
		LEFT JOIN milk_delivery_shifts mds ON cmc.milk_delivery_shift_id = mds.id
		WHERE cmc.deleted_at IS NULL
		ORDER BY cmc.transaction_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CoolerMilkCollectionService) GetCollection(id string) (*dtos.CoolerMilkCollectionResponse, error) {
	var result dtos.CoolerMilkCollectionResponse
	query := `
		SELECT 
			cmc.id, cmc.transaction_date, cmc.quantity, 
			cmc.transport_vehicle_id, tv.registration_no AS vehicle_reg_no,
			cmc.milk_delivery_shift_id, mds.name AS milk_delivery_shift,
			cmc.confirmed, cmc.site_id, 
			cmc.transporter_id, t.transporter_no,
			cmc.route_id, r.route_name,
			cmc.created_at, cmc.updated_at
		FROM cooler_milk_collections cmc
		LEFT JOIN routes r ON cmc.route_id = r.id
		LEFT JOIN transporters t ON cmc.transporter_id = t.id
		LEFT JOIN transporter_vehicles tv ON cmc.transport_vehicle_id = tv.id
		LEFT JOIN milk_delivery_shifts mds ON cmc.milk_delivery_shift_id = mds.id
		WHERE cmc.id = ? AND cmc.deleted_at IS NULL
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

func (s *CoolerMilkCollectionService) UpdateCollection(id string, req dtos.UpdateCoolerMilkCollectionRequest) error {
	var collection models.CoolerMilkCollection
	if err := db.DB.First(&collection, id).Error; err != nil {
		return err
	}

	collection.TransactionDate = utils.ParseDate(req.TransactionDate)
	collection.Quantity = req.Quantity
	collection.TransportVehicleID = req.TransportVehicleID
	collection.MilkDeliveryShiftID = req.MilkDeliveryShiftID
	collection.Confirmed = req.Confirmed
	collection.SiteID = req.SiteID
	collection.TransporterID = req.TransporterID
	collection.RouteID = req.RouteID

	return db.DB.Save(&collection).Error
}

func (s *CoolerMilkCollectionService) DeleteCollection(id string) error {
	return db.DB.Delete(&models.CoolerMilkCollection{}, id).Error
}
