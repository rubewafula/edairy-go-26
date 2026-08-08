package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type DtbDisbursementService struct {
	clients *DisbursementProviderClients
}

func NewDtbDisbursementService() *DtbDisbursementService {
	return &DtbDisbursementService{clients: GetDisbursementProviderClients()}
}

type DtbMpesaInitiateResult struct {
	SCAIntentID       string
	TransferData      map[string]interface{}
	ExternalUniqueID  string
	PhoneNumber       string
	ProviderReference string
}

func (s *DtbDisbursementService) InitiateMpesaTransfer(channelID, memberID uint64, amount float64, phoneNumber, contractNo string) (*DtbMpesaInitiateResult, error) {
	astra, err := s.clients.Astra(channelID)
	if err != nil {
		return nil, err
	}
	if !astra.Configured() {
		return nil, fmt.Errorf("DTB/Astra API is not configured (set ASTRA_API_BASE_URL, ASTRA_AUTHENTICATION_API_BASE_URL, AUTH_IDENTITY, ASTRA_TENANT_PASSWORD)")
	}

	var member models.Member
	if err := db.DB.First(&member, memberID).Error; err != nil {
		return nil, fmt.Errorf("member not found")
	}

	var memberWallet models.MemberWallet
	if err := db.DB.Where("member_id = ?", memberID).Order("id DESC").First(&memberWallet).Error; err != nil {
		return nil, fmt.Errorf("member Astra wallet not found — member must be onboarded on DTB/Astra first")
	}

	phone, err := utils.NormalizeKenyaPhone(phoneNumber, member.PrimaryPhone)
	if err != nil {
		return nil, err
	}

	cfg, err := s.clients.configSvc.BuildAstraConfig(channelID)
	if err != nil {
		return nil, err
	}

	externalUniqueID := uuid.NewString()
	desc := fmt.Sprintf("Loan disbursement %s to %s %s", contractNo, member.FirstName, member.LastName)
	payload := map[string]interface{}{
		"amount":           amount,
		"type":             "KE_DTB_MPESA",
		"description":      desc,
		"externalUniqueId": externalUniqueID,
		"deliverToPhone":   phone,
		"callbackUrl":      cfg.WithdrawalCallback,
		"reference":        fmt.Sprintf("loan-contract#%s", contractNo),
	}

	endpoint := fmt.Sprintf("wallets/%s/withdrawals", memberWallet.WalletID)
	resp, err := astra.InitiateCashTransfer(endpoint, payload)
	if err != nil {
		return nil, err
	}

	scaIntentID := parseSCAIntentID(resp.Headers)
	if scaIntentID == "" {
		msg := "DTB M-Pesa initiation did not return SCA header (OTP)"
		if resp.Body != nil {
			if details, ok := resp.Body["details"]; ok {
				msg = fmt.Sprintf("%s: %v", msg, details)
			} else if message, ok := resp.Body["message"].(string); ok && message != "" {
				msg = message
			}
		}
		return nil, fmt.Errorf(msg)
	}

	return &DtbMpesaInitiateResult{
		SCAIntentID:       scaIntentID,
		TransferData:      payload,
		ExternalUniqueID:  externalUniqueID,
		PhoneNumber:       phone,
		ProviderReference: scaIntentID,
	}, nil
}

type DtbMpesaCompleteResult struct {
	ExternalReference string
	Raw               map[string]interface{}
}

func (s *DtbDisbursementService) CompleteMpesaTransfer(channelID, memberID uint64, scaIntentID, otp string, transferData map[string]interface{}) (*DtbMpesaCompleteResult, error) {
	if scaIntentID == "" || otp == "" {
		return nil, fmt.Errorf("sca_intent_id and otp are required")
	}
	if len(transferData) == 0 {
		return nil, fmt.Errorf("transfer_data is missing")
	}

	astra, err := s.clients.Astra(channelID)
	if err != nil {
		return nil, err
	}

	var memberWallet models.MemberWallet
	if err := db.DB.Where("member_id = ?", memberID).Order("id DESC").First(&memberWallet).Error; err != nil {
		return nil, fmt.Errorf("member Astra wallet not found")
	}

	scaJWT, err := astra.ExchangeSCAForJWT(scaIntentID, otp)
	if err != nil {
		return nil, fmt.Errorf("OTP validation failed: %w", err)
	}

	endpoint := fmt.Sprintf("wallets/%s/withdrawals", memberWallet.WalletID)
	resp, err := astra.PostWithAuth(endpoint, transferData, scaJWT)
	if err != nil {
		return nil, fmt.Errorf("M-Pesa transfer submission failed: %w", err)
	}
	if resp == nil {
		return &DtbMpesaCompleteResult{}, nil
	}
	if errVal, ok := resp["error"]; ok && errVal != nil && errVal != false {
		return nil, fmt.Errorf("M-Pesa transfer rejected: %v", resp)
	}

	extRef := ""
	for _, key := range []string{"withdrawalId", "withdrawal_id", "transactionId", "transaction_id", "reference"} {
		if v, ok := resp[key]; ok {
			extRef = fmt.Sprintf("%v", v)
			break
		}
	}
	if extRef == "" {
		if v, ok := transferData["externalUniqueId"].(string); ok {
			extRef = v
		}
	}

	return &DtbMpesaCompleteResult{
		ExternalReference: extRef,
		Raw:               resp,
	}, nil
}

func (s *DtbDisbursementService) GetWithdrawalFees(channelID, memberID uint64, amount float64) (map[string]interface{}, error) {
	astra, err := s.clients.Astra(channelID)
	if err != nil {
		return nil, err
	}
	var memberWallet models.MemberWallet
	if err := db.DB.Where("member_id = ?", memberID).Order("id DESC").First(&memberWallet).Error; err != nil {
		return nil, fmt.Errorf("member Astra wallet not found")
	}
	return astra.Get(fmt.Sprintf("wallets/%s/withdrawals/fees", memberWallet.WalletID), map[string]string{
		"amount": fmt.Sprintf("%.2f", amount),
		"type":   "KE_DTB_MPESA",
	})
}
