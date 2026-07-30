package main

import (
	"fmt"
	"os"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

func main() {
	db.ConnectToDatabase()

	var dtoResults []dtos.DailyMilkVarianceResponse
	dtoErr := db.DB.Model(&models.DailyMilkVariance{}).Limit(2).Order("day DESC").Scan(&dtoResults).Error
	fmt.Printf("DTO Scan err: %v\n", dtoErr)

	var modelResults []models.DailyMilkVariance
	modelErr := db.DB.Model(&models.DailyMilkVariance{}).Limit(2).Order("day DESC").Find(&modelResults).Error
	fmt.Printf("Model Find err: %v\n", modelErr)
	if modelErr == nil && len(modelResults) > 0 {
		fmt.Printf("Model[0] day=%v transporter=%q\n", modelResults[0].Day, modelResults[0].Transporter)
	}
	os.Exit(0)
}
