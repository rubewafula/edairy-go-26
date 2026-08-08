package services

import (
	"fmt"
	"math"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type MpesaB2cDisbursementService struct {
	clients *DisbursementProviderClients
}

func NewMpesaB2cDisbursementService() *MpesaB2cDisbursementService {
	return &MpesaB2cDisbursementService{clients: GetDisbursementProviderClients()}
}

type MpesaB2CInitiateResult struct {
	PhoneNumber              string
	ClientReference          string
	ConversationID           string
	OriginatorConversationID string
	ProviderReference        string
}

func (s *MpesaB2cDisbursementService) InitiateB2C(channelID, memberID uint64, amount float64, phoneNumber, contractNo, clientReference string) (*MpesaB2CInitiateResult, error) {
	daraja, err := s.clients.Mpesa(channelID)
	if err != nil {
		return nil, err
	}
	if !daraja.Configured() {
		return nil, fmt.Errorf("M-Pesa Daraja is not configured (set MPESA_CONSUMER_KEY, MPESA_CONSUMER_SECRET, MPESA_B2C_SHORTCODE, MPESA_INITIATOR_NAME, MPESA_INITIATOR_PASSWORD, MPESA_CERT_PATH)")
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
		return nil, fmt.Errorf("M-Pesa B2C minimum disbursement is KES 10")
	}

	remarks := fmt.Sprintf("Loan disbursement %s", contractNo)
	resp, err := daraja.InitiateB2C(clientReference, phone, amountInt, remarks)
	if err != nil {
		return nil, err
	}

	return &MpesaB2CInitiateResult{
		PhoneNumber:              phone,
		ClientReference:          clientReference,
		ConversationID:           resp.ConversationID,
		OriginatorConversationID: resp.OriginatorConversationID,
		ProviderReference:        resp.ConversationID,
	}, nil
}

func mpesaB2CPayloadPatch(result *MpesaB2CInitiateResult) map[string]interface{} {
	return map[string]interface{}{
		"phone_number":               result.PhoneNumber,
		"client_reference":           result.ClientReference,
		"conversation_id":            result.ConversationID,
		"originator_conversation_id": result.OriginatorConversationID,
		"provider":                   "MPESA_DARAJA",
		"provider_transaction_id":    result.ConversationID,
		"submitted_at":               time.Now().Format(time.RFC3339),
	}
}
