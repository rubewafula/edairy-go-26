package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type ProductGradeService struct{}

func NewProductGradeService() *ProductGradeService {
	return &ProductGradeService{}
}

func (s *ProductGradeService) CreateProductGrade(req dtos.CreateProductGradeRequest) (*models.ProductGrade, error) {
	grade := &models.ProductGrade{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(grade).Error; err != nil {
		return nil, err
	}
	return grade, nil
}

func (s *ProductGradeService) GetProductGrades(page, limit int) ([]dtos.ProductGradeResponse, int64, error) {
	var grades []dtos.ProductGradeResponse
	var total int64
	db.DB.Model(&models.ProductGrade{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.ProductGrade{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&grades).Error
	return grades, total, err
}

func (s *ProductGradeService) GetProductGrade(id string) (*models.ProductGrade, error) {
	var grade models.ProductGrade
	if err := db.DB.First(&grade, id).Error; err != nil {
		return nil, err
	}
	return &grade, nil
}

func (s *ProductGradeService) UpdateProductGrade(id string, req dtos.UpdateProductGradeRequest) error {
	var grade models.ProductGrade
	if err := db.DB.First(&grade, id).Error; err != nil {
		return err
	}

	grade.Name = req.Name
	grade.Description = req.Description

	return db.DB.Save(&grade).Error
}

func (s *ProductGradeService) DeleteProductGrade(id string) error {
	return db.DB.Delete(&models.ProductGrade{}, id).Error
}
