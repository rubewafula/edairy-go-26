package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type OrganizationKybCommentService struct{}

func NewOrganizationKybCommentService() *OrganizationKybCommentService {
	return &OrganizationKybCommentService{}
}

func (s *OrganizationKybCommentService) CreateComment(req dtos.CreateOrganizationKybCommentRequest, userID uint64) (*models.OrganizationKybComment, error) {
	comment := &models.OrganizationKybComment{
		BaseModel: models.BaseModel{CreatedBy: userID},
		Issue:     req.Issue,
		Comment:   req.Comment,
		Iteration: req.Iteration,
	}
	if err := db.DB.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *OrganizationKybCommentService) GetComments(page, limit int) ([]dtos.OrganizationKybCommentResponse, int64, error) {
	var results []dtos.OrganizationKybCommentResponse
	var total int64
	db.DB.Model(&models.OrganizationKybComment{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.OrganizationKybComment{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *OrganizationKybCommentService) GetComment(id string) (*dtos.OrganizationKybCommentResponse, error) {
	var result dtos.OrganizationKybCommentResponse
	err := db.DB.Model(&models.OrganizationKybComment{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *OrganizationKybCommentService) UpdateComment(id string, req dtos.UpdateOrganizationKybCommentRequest, userID uint64) error {
	var comment models.OrganizationKybComment
	if err := db.DB.First(&comment, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"issue":      req.Issue,
		"comment":    req.Comment,
		"iteration":  req.Iteration,
		"updated_by": userID,
	}

	return db.DB.Model(&comment).Updates(updates).Error
}

func (s *OrganizationKybCommentService) DeleteComment(id string, userID uint64) error {
	var comment models.OrganizationKybComment
	if err := db.DB.First(&comment, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&comment).Update("updated_by", userID).Delete(&comment).Error
}

func (s *OrganizationKybCommentService) GetCommentsByIteration(iteration string) ([]dtos.OrganizationKybCommentResponse, error) {
	var results []dtos.OrganizationKybCommentResponse
	err := db.DB.Model(&models.OrganizationKybComment{}).Where("iteration = ?", iteration).Find(&results).Error
	return results, err
}
