package services

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type LivestockService struct {
	notificationService *UINotificationService
}

func NewLivestockService() *LivestockService {
	return &LivestockService{
		notificationService: NewUINotificationService(),
	}
}

// Livestock CRUD
func (s *LivestockService) CreateLivestock(req dtos.CreateLivestockRequest, userID uint64) (*models.Livestock, error) {

	livestock := &models.Livestock{
		BaseModel:           models.BaseModel{CreatedBy: userID},
		MemberID:            req.MemberID,
		LivestockCategoryID: req.LivestockCategoryID,
		LivestockName:       req.LivestockName,
		TagNo:               req.TagNo,
		LivestockBreedID:    req.LivestockBreedID,
		Gender:              req.Gender,
		Color:               req.Color,
		BirthDate:           utils.ParseDate(req.BirthDate),
		PurchaseDate:        utils.ParseDate(req.BirthDate),
		Notes:               req.Notes,
		Source:              *req.Source,
		Weight:              req.Weight,
		InsuranceNumber:     req.InsuranceNumber,
		Status:              "ACTIVE",
	}
	if err := db.DB.Create(livestock).Error; err != nil {
		return nil, err
	}
	return livestock, nil
}

func (s *LivestockService) GetLivestocks(page, limit int) ([]dtos.LivestockResponse, int64, error) {
	var results []dtos.LivestockResponse
	var total int64
	db.DB.Model(&models.Livestock{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT l.*, lb.breed_name, lc.category_name
		FROM livestocks l
		LEFT JOIN livestock_breeds lb ON l.livestock_breed_id = lb.id
		LEFT JOIN livestock_categories lc ON lb.livestock_category_id = lc.id
		WHERE l.deleted_at IS NULL
		ORDER BY l.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetLivestock(id string) (*dtos.LivestockResponse, error) {
	var result dtos.LivestockResponse
	query := `
		SELECT l.*, lb.breed_name, lc.category_name
		FROM livestocks l
		LEFT JOIN livestock_breeds lb ON l.livestock_breed_id = lb.id
		LEFT JOIN livestock_categories lc ON lb.livestock_category_id = lc.id
		WHERE l.id = ? AND l.deleted_at IS NULL LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil || result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *LivestockService) UpdateLivestock(id string, req dtos.UpdateLivestockRequest, userID uint64) error {
	return db.DB.Model(&models.Livestock{}).Where("id = ?", id).Updates(map[string]interface{}{
		"tag_no": req.TagNo, "breed_id": req.BreedID, "gender": req.Gender,
		"date_of_birth": utils.ParseDate(req.DateOfBirth), "status": req.Status,
		"description": req.Description, "updated_by": userID,
	}).Error
}

func (s *LivestockService) DeleteLivestock(id string, userID uint64) error {
	return db.DB.Model(&models.Livestock{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.Livestock{}).Error
}

// Category CRUD
func (s *LivestockService) CreateCategory(req dtos.CreateLivestockCategoryRequest) (*models.LivestockCategory, error) {
	category := &models.LivestockCategory{
		CategoryName: req.CategoryName,
		Description:  req.Description,
	}
	if err := db.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *LivestockService) GetCategories(page, limit int) ([]dtos.LivestockCategoryResponse, int64, error) {
	var results []dtos.LivestockCategoryResponse
	var total int64

	db.DB.Model(&models.LivestockCategory{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.LivestockCategory{}).
		Limit(limit).Offset(offset).Order("id DESC").
		Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetCategory(id string) (*dtos.LivestockCategoryResponse, error) {
	var result dtos.LivestockCategoryResponse
	err := db.DB.Model(&models.LivestockCategory{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *LivestockService) UpdateCategory(id string, req dtos.UpdateLivestockCategoryRequest, userID uint64) error {
	var category models.LivestockCategory
	if err := db.DB.First(&category, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"category_name": req.CategoryName,
		"description":   req.Description,
		"updated_by":    userID,
	}
	return db.DB.Model(&category).Updates(updates).Error
}

func (s *LivestockService) DeleteCategory(id string, userID uint64) error {
	var category models.LivestockCategory
	if err := db.DB.First(&category, id).Error; err != nil {
		return err
	}
	// Check for associated breeds before deleting category
	var count int64
	db.DB.Model(&models.LivestockBreed{}).Where("livestock_category_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete category with associated breeds")
	}
	return db.DB.Model(&category).Update("updated_by", userID).Delete(&category).Error
}

// Breed CRUD
func (s *LivestockService) CreateBreed(req dtos.CreateLivestockBreedRequest) (*models.LivestockBreed, error) {
	breed := &models.LivestockBreed{
		LivestockCategoryID: req.LivestockCategoryID,
		BreedName:           req.BreedName,
		Description:         req.Description,
	}
	if err := db.DB.Create(breed).Error; err != nil {
		return nil, err
	}
	return breed, nil
}

func (s *LivestockService) GetBreeds(page, limit int) ([]dtos.LivestockBreedResponse, int64, error) {
	var results []dtos.LivestockBreedResponse
	var total int64

	db.DB.Model(&models.LivestockBreed{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT b.id, b.breed_name, b.description, b.created_at, c.category_name
		FROM livestock_breeds b
		JOIN livestock_categories c ON b.livestock_category_id = c.id
		WHERE b.deleted_at IS NULL
		ORDER BY b.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetBreed(id string) (*dtos.LivestockBreedResponse, error) {
	var result dtos.LivestockBreedResponse
	query := `
		SELECT b.id, b.breed_name, b.description, b.created_at, c.category_name
		FROM livestock_breeds b
		JOIN livestock_categories c ON b.livestock_category_id = c.id
		WHERE b.id = ? AND b.deleted_at IS NULL
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

func (s *LivestockService) UpdateBreed(id string, req dtos.UpdateLivestockBreedRequest, userID uint64) error {
	var breed models.LivestockBreed
	if err := db.DB.First(&breed, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"livestock_category_id": req.LivestockCategoryID,
		"breed_name":            req.BreedName,
		"description":           req.Description,
		"updated_by":            userID,
	}
	return db.DB.Model(&breed).Updates(updates).Error
}

func (s *LivestockService) DeleteBreed(id string, userID uint64) error {
	var breed models.LivestockBreed
	if err := db.DB.First(&breed, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&breed).Update("updated_by", userID).Delete(&breed).Error
}

// Events (Death, Weight, Feeding, Health, etc)
func (s *LivestockService) CreateDeath(req dtos.CreateLivestockDeathRequest, userID uint64) (*models.LivestockDeath, error) {
	death := &models.LivestockDeath{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		LivestockID:    req.LivestockID,
		DeathDate:      utils.ParseDate(req.DeathDate),
		CauseOfDeath:   req.CauseOfDeath,
		DisposalMethod: req.DisposalMethod,
		Remarks:        req.Remarks,
	}
	if err := db.DB.Create(death).Error; err != nil {
		return nil, err
	}
	return death, nil
}

func (s *LivestockService) GetDeaths(livestockID string, page, limit int) ([]dtos.LivestockDeathResponse, int64, error) {
	var results []dtos.LivestockDeathResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockDeath{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT ld.id, ld.livestock_id, l.tag_no as livestock_tag_no, ld.death_date,
		       ld.cause_of_death, ld.disposal_method, ld.remarks, ld.created_at, ld.updated_at
		FROM livestock_deaths ld
		JOIN livestocks l ON ld.livestock_id = l.id
		WHERE (? = '' OR ld.livestock_id = ?) AND ld.deleted_at IS NULL
		ORDER BY ld.death_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetDeath(id string) (*dtos.LivestockDeathResponse, error) {
	var result dtos.LivestockDeathResponse
	query := `
		SELECT ld.id, ld.livestock_id, l.tag_no as livestock_tag_no, ld.death_date,
		       ld.cause_of_death, ld.disposal_method, ld.remarks, ld.created_at, ld.updated_at
		FROM livestock_deaths ld
		JOIN livestocks l ON ld.livestock_id = l.id
		WHERE ld.id = ? AND ld.deleted_at IS NULL
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

func (s *LivestockService) UpdateDeath(id string, req dtos.UpdateLivestockDeathRequest, userID uint64) error {
	var death models.LivestockDeath
	if err := db.DB.First(&death, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"death_date":      utils.ParseDate(req.DeathDate),
		"cause_of_death":  req.CauseOfDeath,
		"disposal_method": req.DisposalMethod,
		"remarks":         req.Remarks,
		"updated_by":      userID,
	}
	return db.DB.Model(&death).Updates(updates).Error
}

func (s *LivestockService) DeleteDeath(id string, userID uint64) error {
	var death models.LivestockDeath
	if err := db.DB.First(&death, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&death).Update("updated_by", userID).Delete(&death).Error
}

func (s *LivestockService) CreateFeeding(req dtos.CreateLivestockFeedingRequest, userID uint64) (*models.LivestockFeeding, error) {
	feeding := &models.LivestockFeeding{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		LivestockID: req.LivestockID,
		FeedName:    req.FeedName,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		FeedingDate: utils.ParseDate(req.FeedingDate),
		Cost:        req.Cost,
		Notes:       req.Notes,
	}
	if err := db.DB.Create(feeding).Error; err != nil {
		return nil, err
	}
	return feeding, nil
}

func (s *LivestockService) GetFeedings(livestockID string, page, limit int) ([]dtos.LivestockFeedingResponse, int64, error) {
	var results []dtos.LivestockFeedingResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockFeeding{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT lf.id, lf.livestock_id, l.tag_no as livestock_tag_no, lf.feed_name,
		       lf.quantity, lf.unit, lf.feeding_date, lf.cost, lf.notes,
		       lf.created_at, lf.updated_at
		FROM livestock_feedings lf
		JOIN livestocks l ON lf.livestock_id = l.id
		WHERE (? = '' OR lf.livestock_id = ?) AND lf.deleted_at IS NULL
		ORDER BY lf.feeding_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetFeeding(id string) (*dtos.LivestockFeedingResponse, error) {
	var result dtos.LivestockFeedingResponse
	query := `
		SELECT lf.id, lf.livestock_id, l.tag_no as livestock_tag_no, lf.feed_name,
		       lf.quantity, lf.unit, lf.feeding_date, lf.cost, lf.notes,
		       lf.created_at, lf.updated_at
		FROM livestock_feedings lf
		JOIN livestocks l ON lf.livestock_id = l.id
		WHERE lf.id = ? AND lf.deleted_at IS NULL
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

func (s *LivestockService) UpdateFeeding(id string, req dtos.UpdateLivestockFeedingRequest, userID uint64) error {
	var feeding models.LivestockFeeding
	if err := db.DB.First(&feeding, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"feed_name":    req.FeedName,
		"quantity":     req.Quantity,
		"unit":         req.Unit,
		"feeding_date": utils.ParseDate(req.FeedingDate),
		"cost":         req.Cost,
		"notes":        req.Notes,
		"updated_by":   userID,
	}
	return db.DB.Model(&feeding).Updates(updates).Error
}

func (s *LivestockService) DeleteFeeding(id string, userID uint64) error {
	var feeding models.LivestockFeeding
	if err := db.DB.First(&feeding, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&feeding).Update("updated_by", userID).Delete(&feeding).Error
}

func (s *LivestockService) CreateHealth(req dtos.CreateLivestockHealthRequest, userID uint64) (*models.LivestockHealthRecord, error) {
	record := &models.LivestockHealthRecord{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		LivestockID:   req.LivestockID,
		RecordType:    req.RecordType,
		Diagnosis:     req.Diagnosis,
		Medication:    req.Medication,
		Dosage:        req.Dosage,
		Veterinarian:  req.Veterinarian,
		TreatmentDate: utils.ParseDate(req.TreatmentDate),
		Notes:         req.Notes,
	}
	if req.NextVisitDate != "" {
		t := utils.ParseDate(req.NextVisitDate)
		record.NextVisitDate = &t
	}
	if err := db.DB.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s *LivestockService) GetHealthRecords(livestockID string, page, limit int) ([]dtos.LivestockHealthResponse, int64, error) {
	var results []dtos.LivestockHealthResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockHealthRecord{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT lhr.id, lhr.livestock_id, l.tag_no as livestock_tag_no, lhr.record_type,
		       lhr.diagnosis, lhr.medication, lhr.dosage, lhr.veterinarian,
		       lhr.treatment_date, lhr.next_visit_date, lhr.notes,
		       lhr.created_at, lhr.updated_at
		FROM livestock_health_records lhr
		JOIN livestocks l ON lhr.livestock_id = l.id
		WHERE (? = '' OR lhr.livestock_id = ?) AND lhr.deleted_at IS NULL
		ORDER BY lhr.treatment_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetHealthRecord(id string) (*dtos.LivestockHealthResponse, error) {
	var result dtos.LivestockHealthResponse
	query := `
		SELECT lhr.id, lhr.livestock_id, l.tag_no as livestock_tag_no, lhr.record_type,
		       lhr.diagnosis, lhr.medication, lhr.dosage, lhr.veterinarian,
		       lhr.treatment_date, lhr.next_visit_date, lhr.notes,
		       lhr.created_at, lhr.updated_at
		FROM livestock_health_records lhr
		JOIN livestocks l ON lhr.livestock_id = l.id
		WHERE lhr.id = ? AND lhr.deleted_at IS NULL
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

func (s *LivestockService) UpdateHealthRecord(id string, req dtos.UpdateLivestockHealthRequest, userID uint64) error {
	var record models.LivestockHealthRecord
	if err := db.DB.First(&record, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"record_type":    req.RecordType,
		"diagnosis":      req.Diagnosis,
		"medication":     req.Medication,
		"dosage":         req.Dosage,
		"veterinarian":   req.Veterinarian,
		"treatment_date": utils.ParseDate(req.TreatmentDate),
		"notes":          req.Notes,
		"updated_by":     userID,
	}
	if req.NextVisitDate != "" {
		t := utils.ParseDate(req.NextVisitDate)
		updates["next_visit_date"] = &t
	} else {
		updates["next_visit_date"] = nil
	}
	return db.DB.Model(&record).Updates(updates).Error
}

func (s *LivestockService) DeleteHealthRecord(id string, userID uint64) error {
	var record models.LivestockHealthRecord
	if err := db.DB.First(&record, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&record).Update("updated_by", userID).Delete(&record).Error
}

func (s *LivestockService) CreateMovement(req dtos.CreateLivestockMovementRequest, userID uint64) (*models.LivestockMovement, error) {
	movement := &models.LivestockMovement{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		LivestockID:  req.LivestockID,
		FromLocation: req.FromLocation,
		ToLocation:   req.ToLocation,
		MovementDate: utils.ParseDate(req.MovementDate),
		MovementType: req.MovementType,
		Transporter:  req.Transporter,
		Remarks:      req.Remarks,
	}
	if err := db.DB.Create(movement).Error; err != nil {
		return nil, err
	}
	return movement, nil
}

func (s *LivestockService) GetMovements(livestockID string, page, limit int) ([]dtos.LivestockMovementResponse, int64, error) {
	var results []dtos.LivestockMovementResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockMovement{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT lm.id, lm.livestock_id, l.tag_no as livestock_tag_no, lm.from_location,
		       lm.to_location, lm.movement_date, lm.movement_type, lm.transporter,
		       lm.remarks, lm.created_at, lm.updated_at
		FROM livestock_movements lm
		JOIN livestocks l ON lm.livestock_id = l.id
		WHERE (? = '' OR lm.livestock_id = ?) AND lm.deleted_at IS NULL
		ORDER BY lm.movement_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetMovement(id string) (*dtos.LivestockMovementResponse, error) {
	var result dtos.LivestockMovementResponse
	query := `
		SELECT lm.id, lm.livestock_id, l.tag_no as livestock_tag_no, lm.from_location,
		       lm.to_location, lm.movement_date, lm.movement_type, lm.transporter,
		       lm.remarks, lm.created_at, lm.updated_at
		FROM livestock_movements lm
		JOIN livestocks l ON lm.livestock_id = l.id
		WHERE lm.id = ? AND lm.deleted_at IS NULL
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

func (s *LivestockService) UpdateMovement(id string, req dtos.UpdateLivestockMovementRequest, userID uint64) error {
	var movement models.LivestockMovement
	if err := db.DB.First(&movement, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"from_location": req.FromLocation,
		"to_location":   req.ToLocation,
		"movement_date": utils.ParseDate(req.MovementDate),
		"movement_type": req.MovementType,
		"transporter":   req.Transporter,
		"remarks":       req.Remarks,
		"updated_by":    userID,
	}
	return db.DB.Model(&movement).Updates(updates).Error
}

func (s *LivestockService) DeleteMovement(id string, userID uint64) error {
	var movement models.LivestockMovement
	if err := db.DB.First(&movement, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&movement).Update("updated_by", userID).Delete(&movement).Error
}

func (s *LivestockService) CreatePhoto(req dtos.CreateLivestockPhotoRequest, userID uint64) (*models.LivestockPhoto, error) {
	photo := &models.LivestockPhoto{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		LivestockID: req.LivestockID,
		PhotoURL:    req.PhotoURL,
		Description: req.Description,
	}
	if err := db.DB.Create(photo).Error; err != nil {
		return nil, err
	}
	return photo, nil
}

func (s *LivestockService) GetPhotos(livestockID string, page, limit int) ([]dtos.LivestockPhotoResponse, int64, error) {
	var results []dtos.LivestockPhotoResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockPhoto{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT lp.id, lp.livestock_id, l.tag_no as livestock_tag_no, lp.photo_url,
		       lp.description, lp.created_at, lp.updated_at
		FROM livestock_photos lp
		JOIN livestocks l ON lp.livestock_id = l.id
		WHERE (? = '' OR lp.livestock_id = ?) AND lp.deleted_at IS NULL
		ORDER BY lp.created_at DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetPhoto(id string) (*dtos.LivestockPhotoResponse, error) {
	var result dtos.LivestockPhotoResponse
	query := `
		SELECT lp.id, lp.livestock_id, l.tag_no as livestock_tag_no, lp.photo_url,
		       lp.description, lp.created_at, lp.updated_at
		FROM livestock_photos lp
		JOIN livestocks l ON lp.livestock_id = l.id
		WHERE lp.id = ? AND lp.deleted_at IS NULL
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

func (s *LivestockService) UpdatePhoto(id string, req dtos.UpdateLivestockPhotoRequest, userID uint64) error {
	var photo models.LivestockPhoto
	if err := db.DB.First(&photo, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"photo_url":   req.PhotoURL,
		"description": req.Description,
		"updated_by":  userID,
	}
	return db.DB.Model(&photo).Updates(updates).Error
}

func (s *LivestockService) DeletePhoto(id string, userID uint64) error {
	var photo models.LivestockPhoto
	if err := db.DB.First(&photo, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&photo).Update("updated_by", userID).Delete(&photo).Error
}

func (s *LivestockService) CreateSale(req dtos.CreateLivestockSaleRequest, userID uint64) (*models.LivestockSale, error) {
	sale := &models.LivestockSale{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		LivestockID:   req.LivestockID,
		CustomerID:    req.CustomerID,
		SaleDate:      utils.ParseDate(req.SaleDate),
		Quantity:      req.Quantity,
		SalePrice:     req.SalePrice,
		PaymentStatus: req.PaymentStatus,
		Notes:         req.Notes,
	}
	if err := db.DB.Create(sale).Error; err != nil {
		return nil, err
	}
	return sale, nil
}

func (s *LivestockService) GetSales(livestockID string, page, limit int) ([]dtos.LivestockSaleResponse, int64, error) {
	var results []dtos.LivestockSaleResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockSale{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT ls.id, ls.livestock_id, l.tag_no as livestock_tag_no, ls.customer_id,
		       ls.sale_date, ls.quantity, ls.sale_price, ls.payment_status, ls.notes,
		       ls.created_at, ls.updated_at
		FROM livestock_sales ls
		JOIN livestocks l ON ls.livestock_id = l.id
		WHERE (? = '' OR ls.livestock_id = ?) AND ls.deleted_at IS NULL
		ORDER BY ls.sale_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetSale(id string) (*dtos.LivestockSaleResponse, error) {
	var result dtos.LivestockSaleResponse
	query := `
		SELECT ls.id, ls.livestock_id, l.tag_no as livestock_tag_no, ls.customer_id,
		       ls.sale_date, ls.quantity, ls.sale_price, ls.payment_status, ls.notes,
		       ls.created_at, ls.updated_at
		FROM livestock_sales ls
		JOIN livestocks l ON ls.livestock_id = l.id
		WHERE ls.id = ? AND ls.deleted_at IS NULL
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

func (s *LivestockService) UpdateSale(id string, req dtos.UpdateLivestockSaleRequest, userID uint64) error {
	var sale models.LivestockSale
	if err := db.DB.First(&sale, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"customer_id":    req.CustomerID,
		"sale_date":      utils.ParseDate(req.SaleDate),
		"quantity":       req.Quantity,
		"sale_price":     req.SalePrice,
		"payment_status": req.PaymentStatus,
		"notes":          req.Notes,
		"updated_by":     userID,
	}
	return db.DB.Model(&sale).Updates(updates).Error
}

func (s *LivestockService) DeleteSale(id string, userID uint64) error {
	var sale models.LivestockSale
	if err := db.DB.First(&sale, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&sale).Update("updated_by", userID).Delete(&sale).Error
}

func (s *LivestockService) CreateWeight(req dtos.CreateLivestockWeightRequest, userID uint64) (*models.LivestockWeightRecord, error) {
	weight := &models.LivestockWeightRecord{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		LivestockID: req.LivestockID,
		Weight:      req.Weight,
		RecordedAt:  utils.ParseDate(req.RecordedAt),
		Remarks:     req.Remarks,
	}
	if err := db.DB.Create(weight).Error; err != nil {
		return nil, err
	}
	return weight, nil
}

func (s *LivestockService) GetWeightRecords(livestockID string, page, limit int) ([]dtos.LivestockWeightResponse, int64, error) {
	var results []dtos.LivestockWeightResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockWeightRecord{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT lwr.id, lwr.livestock_id, l.tag_no as livestock_tag_no, lwr.weight,
		       lwr.recorded_at, lwr.remarks, lwr.created_at, lwr.updated_at
		FROM livestock_weight_records lwr
		JOIN livestocks l ON lwr.livestock_id = l.id
		WHERE (? = '' OR lwr.livestock_id = ?) AND lwr.deleted_at IS NULL
		ORDER BY lwr.recorded_at DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetWeightRecord(id string) (*dtos.LivestockWeightResponse, error) {
	var result dtos.LivestockWeightResponse
	query := `
		SELECT lwr.id, lwr.livestock_id, l.tag_no as livestock_tag_no, lwr.weight,
		       lwr.recorded_at, lwr.remarks, lwr.created_at, lwr.updated_at
		FROM livestock_weight_records lwr
		JOIN livestocks l ON lwr.livestock_id = l.id
		WHERE lwr.id = ? AND lwr.deleted_at IS NULL
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

func (s *LivestockService) UpdateWeightRecord(id string, req dtos.UpdateLivestockWeightRequest, userID uint64) error {
	var weight models.LivestockWeightRecord
	if err := db.DB.First(&weight, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"weight":      req.Weight,
		"recorded_at": utils.ParseDate(req.RecordedAt),
		"remarks":     req.Remarks,
		"updated_by":  userID,
	}
	return db.DB.Model(&weight).Updates(updates).Error
}

func (s *LivestockService) DeleteWeightRecord(id string, userID uint64) error {
	var weight models.LivestockWeightRecord
	if err := db.DB.First(&weight, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&weight).Update("updated_by", userID).Delete(&weight).Error
}

// Breeding CRUD
func (s *LivestockService) CreateBreeding(req dtos.CreateLivestockBreedingRequest, userID uint64) (*models.LivestockBreedingRecord, error) {
	record := &models.LivestockBreedingRecord{
		BaseModel:         models.BaseModel{CreatedBy: userID},
		MaleLivestockID:   req.MaleLivestockID,
		FemaleLivestockID: req.FemaleLivestockID,
		BreedingDate:      utils.ParseDate(req.BreedingDate),
		BreedingType:      req.BreedingType,
		PregnancyStatus:   req.PregnancyStatus,
		Remarks:           req.Remarks,
	}

	if req.MaleLivestockID != 0 {
		record.MaleLivestockID = req.MaleLivestockID
	}
	if req.PregnancyCheckDate != "" {
		t := utils.ParseDate(req.PregnancyCheckDate)
		record.PregnancyCheckDate = &t
	}
	if req.ExpectedCalvingDate != "" {
		t := utils.ParseDate(req.ExpectedCalvingDate)
		record.ExpectedCalvingDate = &t
	}
	if req.ActualCalvingDate != "" {
		t := utils.ParseDate(req.ActualCalvingDate)
		record.ActualCalvingDate = &t
	}

	if err := db.DB.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s *LivestockService) GetBreedingRecords(livestockID string, page, limit int) ([]dtos.LivestockBreedingResponse, int64, error) {
	var results []dtos.LivestockBreedingResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockBreedingRecord{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT lbr.*, l.tag_no as livestock_tag_no
		FROM livestock_breeding_records lbr
		JOIN livestocks l ON lbr.male_livestock_id = l.id
		WHERE (lbr.female_livestock_id = ? OR lbr.male_livestock_id = ?) AND lbr.deleted_at IS NULL
		ORDER BY lbr.breeding_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockService) GetBreedingRecord(id string) (*dtos.LivestockBreedingResponse, error) {
	var result dtos.LivestockBreedingResponse
	query := `
		SELECT lbr.*, l.tag_no as livestock_tag_no
		FROM livestock_breeding_records lbr
		JOIN livestocks l ON lbr.male_livestock_id = l.id
		WHERE lbr.id = ? AND lbr.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil || result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *LivestockService) UpdateBreedingRecord(id string, req dtos.UpdateLivestockBreedingRequest, userID uint64) error {
	var record models.LivestockBreedingRecord
	if err := db.DB.First(&record, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"breeding_date":    utils.ParseDate(req.BreedingDate),
		"breeding_type":    req.BreedingType,
		"pregnancy_status": req.PregnancyStatus,
		"remarks":          req.Remarks,
		"updated_by":       userID,
	}

	if req.MaleLivestockID != 0 {
		updates["sire_id"] = req.MaleLivestockID
	}
	if req.PregnancyCheckDate != "" {
		t := utils.ParseDate(req.PregnancyCheckDate)
		updates["pregnancy_check_date"] = &t
	}
	if req.ExpectedCalvingDate != "" {
		t := utils.ParseDate(req.ExpectedCalvingDate)
		updates["expected_calving_date"] = &t
	}
	if req.ActualCalvingDate != "" {
		t := utils.ParseDate(req.ActualCalvingDate)
		updates["actual_calving_date"] = &t
	}

	return db.DB.Model(&record).Updates(updates).Error
}

func (s *LivestockService) DeleteBreedingRecord(id string, userID uint64) error {
	return db.DB.Model(&models.LivestockBreedingRecord{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.LivestockBreedingRecord{}).Error
}

func (s *LivestockService) GetLivestockGenericResponse(livestockID string, page, limit int) ([]dtos.LivestockGenericResponse, int64, error) {
	var results []dtos.LivestockGenericResponse
	var total int64

	// This is a generic example, you'd need to build specific queries for each type of generic response
	// For now, let's just return an empty slice and 0 total
	return results, total, nil
}

func (s *LivestockService) RecordFeeding(req dtos.CreateLivestockFeedingRequest) (*models.LivestockFeeding, error) {
	feeding := &models.LivestockFeeding{
		LivestockID: req.LivestockID,
		FeedName:    req.FeedName,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		FeedingDate: utils.ParseDate(req.FeedingDate),
		Cost:        req.Cost,
		Notes:       req.Notes,
	}
	return feeding, db.DB.Create(feeding).Error
}

func (s *LivestockService) RecordHealth(req dtos.CreateLivestockHealthRequest) (*models.LivestockHealthRecord, error) {
	record := &models.LivestockHealthRecord{
		LivestockID:   req.LivestockID,
		RecordType:    req.RecordType,
		Diagnosis:     req.Diagnosis,
		Medication:    req.Medication,
		Dosage:        req.Dosage,
		Veterinarian:  req.Veterinarian,
		TreatmentDate: utils.ParseDate(req.TreatmentDate),
		Notes:         req.Notes,
	}
	if req.NextVisitDate != "" {
		t := utils.ParseDate(req.NextVisitDate)
		record.NextVisitDate = &t
	}
	return record, db.DB.Create(record).Error
}

func (s *LivestockService) RecordMovement(req dtos.CreateLivestockMovementRequest) (*models.LivestockMovement, error) {
	movement := &models.LivestockMovement{
		LivestockID:  req.LivestockID,
		FromLocation: req.FromLocation,
		ToLocation:   req.ToLocation,
		MovementDate: utils.ParseDate(req.MovementDate),
		MovementType: req.MovementType,
		Transporter:  req.Transporter,
		Remarks:      req.Remarks,
	}
	return movement, db.DB.Create(movement).Error
}

func (s *LivestockService) AddPhoto(req dtos.CreateLivestockPhotoRequest) (*models.LivestockPhoto, error) {
	photo := &models.LivestockPhoto{
		LivestockID: req.LivestockID,
		PhotoURL:    req.PhotoURL,
		Description: req.Description,
	}
	return photo, db.DB.Create(photo).Error
}

func (s *LivestockService) RecordBreeding(req dtos.CreateLivestockBreedingRequest) (*models.LivestockBreedingRecord, error) {
	record := &models.LivestockBreedingRecord{
		MaleLivestockID:   req.MaleLivestockID,
		FemaleLivestockID: req.FemaleLivestockID,
		BreedingDate:      utils.ParseDate(req.BreedingDate),
		BreedingType:      req.BreedingType,
		PregnancyStatus:   req.PregnancyStatus,
		Remarks:           req.Remarks,
	}
	if req.MaleLivestockID != 0 {
		record.MaleLivestockID = req.MaleLivestockID
	}
	if req.PregnancyCheckDate != "" {
		t := utils.ParseDate(req.PregnancyCheckDate)
		record.PregnancyCheckDate = &t
	}
	if req.ExpectedCalvingDate != "" {
		t := utils.ParseDate(req.ExpectedCalvingDate)
		record.ExpectedCalvingDate = &t
	}
	return record, db.DB.Create(record).Error
}

func (s *LivestockService) RecordSale(req dtos.CreateLivestockSaleRequest) (*models.LivestockSale, error) {
	sale := &models.LivestockSale{
		LivestockID:   req.LivestockID,
		CustomerID:    req.CustomerID,
		SaleDate:      utils.ParseDate(req.SaleDate),
		Quantity:      req.Quantity,
		SalePrice:     req.SalePrice,
		PaymentStatus: req.PaymentStatus,
		Notes:         req.Notes,
	}
	return sale, db.DB.Create(sale).Error
}

func (s *LivestockService) RecordWeight(req dtos.CreateLivestockWeightRequest) (*models.LivestockWeightRecord, error) {
	weight := &models.LivestockWeightRecord{
		LivestockID: req.LivestockID,
		Weight:      req.Weight,
		RecordedAt:  utils.ParseDate(req.RecordedAt),
		Remarks:     req.Remarks,
	}
	return weight, db.DB.Create(weight).Error
}

func (s *LivestockService) ImportLivestock(file *multipart.FileHeader, userID uint64) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	var data [][]string

	if ext == ".csv" {
		reader := csv.NewReader(src)
		data, err = reader.ReadAll()
	} else if ext == ".xlsx" || ext == ".xls" {
		f, err := excelize.OpenReader(src)
		if err != nil {
			return err
		}
		sheets := f.GetSheetList()
		data, err = f.GetRows(sheets[0])
	} else {
		return fmt.Errorf("unsupported format")
	}

	if err != nil {
		return err
	}

	go s.processLivestockRowsInBackground(data, userID)
	return nil
}

func (s *LivestockService) processLivestockRowsInBackground(data [][]string, userID uint64) {
	var wg sync.WaitGroup
	jobs := make(chan []string, len(data)-1)
	errorChan := make(chan error, len(data)-1)
	numWorkers := runtime.NumCPU() * 2

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							db.DB.Create(&models.LivestockImportError{
								BaseModel: models.BaseModel{CreatedBy: userID},
								RowData:   strings.Join(row, ","),
								Error:     fmt.Sprintf("Panic: %v", r),
							})
						}
					}()

					err := db.DB.Transaction(func(tx *gorm.DB) error {
						if len(row) < 5 {
							return fmt.Errorf("insufficient columns")
						}

						tagNo := strings.TrimSpace(row[1])
						var count int64
						tx.Model(&models.Livestock{}).Where("tag_no = ?", tagNo).Count(&count)
						if count > 0 {
							return fmt.Errorf("tag %s exists", tagNo)
						}

						// Resolve Category and Breed
						var cat models.LivestockCategory
						tx.Where("category_name = ?", row[2]).First(&cat)
						var breed models.LivestockBreed
						tx.Where("breed_name = ?", row[3]).First(&breed)

						weight, _ := utils.ParseFloat(row[7])
						livestock := models.Livestock{
							BaseModel:           models.BaseModel{CreatedBy: userID},
							TagNo:               &tagNo,
							LivestockName:       &row[4],
							Gender:              strings.ToLower(row[5]),
							BirthDate:           utils.ParseFlexibleDate(row[6]),
							LivestockCategoryID: cat.ID,
							LivestockBreedID:    &breed.ID,
							Weight:              &weight,
							Status:              "active",
						}
						return tx.Create(&livestock).Error
					})

					if err != nil {
						db.DB.Create(&models.LivestockImportError{
							BaseModel: models.BaseModel{CreatedBy: userID},
							RowData:   strings.Join(row, ","),
							Error:     err.Error(),
						})
						errorChan <- err
					}
				}()
			}
		}()
	}

	for i := 1; i < len(data); i++ {
		jobs <- data[i]
	}
	close(jobs)
	wg.Wait()
	close(errorChan)

	failedCount := 0
	for range errorChan {
		failedCount++
	}

	msg := "Livestock import completed."
	if failedCount > 0 {
		msg = fmt.Sprintf("Livestock import finished with %d errors.", failedCount)
	}
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title: "Livestock Import Status", Message: msg, NotificationType: "IMPORT_STATUS",
	})
}
