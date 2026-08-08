package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type EquitelDisbursementService struct {
	clients *DisbursementProviderClients
}

func NewEquitelDisbursementService() *EquitelDisbursementService {
	return &EquitelDisbursementService{clients: GetDisbursementProviderClients()}
}

type EquitelInitiateResult struct {
	PhoneNumber       string
	ClientReference   string
	TransactionRef    string
	ProviderReference string
}

func (s *EquitelDisbursementService) InitiatePayout(channelID, memberID uint64, amount float64, phoneNumber, contractNo, clientReference string) (*EquitelInitiateResult, error) {
	jenga, err := s.clients.Jenga(channelID)
	if err != nil {
		return nil, err
	}
	if !jenga.Configured() {
		return nil, fmt.Errorf("Equitel/Jenga is not configured (set JENGA_API_KEY, JENGA_MERCHANT_CODE, JENGA_CONSUMER_SECRET, JENGA_PRIVATE_KEY_PATH, JENGA_SOURCE_ACCOUNT)")
	}
	if clientReference == "" {
		return nil, fmt.Errorf("disbursement reference is required")
	}
	if amount < 100 {
		return nil, fmt.Errorf("Equitel minimum disbursement is KES 100")
	}
	if amount > 140000 {
		return nil, fmt.Errorf("Equitel maximum disbursement is KES 140,000")
	}

	var member models.Member
	if err := db.DB.First(&member, memberID).Error; err != nil {
		return nil, fmt.Errorf("member not found")
	}
	phone, err := utils.NormalizeKenyaPhone(phoneNumber, member.PrimaryPhone)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(strings.TrimSpace(member.FirstName) + " " + strings.TrimSpace(member.LastName))
	if name == "" {
		name = "Member"
	}
	amountStr := fmt.Sprintf("%.2f", amount)
	desc := fmt.Sprintf("Loan disbursement %s", contractNo)

	resp, err := jenga.SendMobileEquitel(name, phone, amountStr, clientReference, desc)
	if err != nil {
		return nil, err
	}

	txnRef := resp.Data.TransactionReference
	if txnRef == "" {
		txnRef = clientReference
	}

	return &EquitelInitiateResult{
		PhoneNumber:       phone,
		ClientReference:   clientReference,
		TransactionRef:    txnRef,
		ProviderReference: txnRef,
	}, nil
}

func equitelPayloadPatch(result *EquitelInitiateResult) map[string]interface{} {
	return map[string]interface{}{
		"phone_number":            result.PhoneNumber,
		"client_reference":        result.ClientReference,
		"provider":                "JENGA_EQUITEL",
		"provider_transaction_id": result.TransactionRef,
		"submitted_at":            time.Now().Format(time.RFC3339),
	}
}
