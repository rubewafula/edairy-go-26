package services

import (
	"fmt"
	"math"

	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type AirtelDisbursementService struct {
	clients *DisbursementProviderClients
}

func NewAirtelDisbursementService() *AirtelDisbursementService {
	return &AirtelDisbursementService{clients: GetDisbursementProviderClients()}
}

type AirtelInitiateResult struct {
	PhoneNumber       string
	ClientReference   string
	TransactionID     string
	ProviderReference string
}

func (s *AirtelDisbursementService) InitiateDisbursement(channelID, memberID uint64, amount float64, phoneNumber, clientReference string) (*AirtelInitiateResult, error) {
	api, err := s.clients.Airtel(channelID)
	if err != nil {
		return nil, err
	}
	if !api.Configured() {
		return nil, fmt.Errorf("Airtel disbursement is not configured (set AIRTEL_CLIENT_ID, AIRTEL_CLIENT_SECRET, AIRTEL_DISBURSEMENT_PIN, AIRTEL_PUBLIC_KEY_PATH)")
	}
	if clientReference == "" {
		return nil, fmt.Errorf("disbursement reference is required")
	}

	var member models.Member
	if err := db.DB.First(&member, memberID).Error; err != nil {
		return nil, fmt.Errorf("member not found")
	}
	phone, err := utils.NormalizeKenyaPhone(phoneNumber, member.PrimaryPhone)
	if err != nil {
		return nil, err
	}

	amountInt := int(math.Round(amount))
	if amountInt < 10 {
		return nil, fmt.Errorf("Airtel minimum disbursement is KES 10")
	}

	resp, err := api.InitiateDisbursement(phone, clientReference, amountInt)
	if err != nil {
		return nil, err
	}

	txnID := resp.Data.Transaction.ID
	if txnID == "" {
		txnID = clientReference
	}

	return &AirtelInitiateResult{
		PhoneNumber:       phone,
		ClientReference:   clientReference,
		TransactionID:     txnID,
		ProviderReference: txnID,
	}, nil
}

func airtelPayloadPatch(result *AirtelInitiateResult) map[string]interface{} {
	return map[string]interface{}{
		"phone_number":            result.PhoneNumber,
		"client_reference":        result.ClientReference,
		"provider":                "AIRTEL",
		"provider_transaction_id": result.TransactionID,
		"submitted_at":            time.Now().Format(time.RFC3339),
	}
}
