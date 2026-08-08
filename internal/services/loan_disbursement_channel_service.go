package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DisbursementExecutionInput struct {
	Channel    *models.DisbursementChannel
	NetAmount  float64
	Payload    map[string]interface{}
	ContractNo string
	MemberID   uint64
	Reference  string
}

type DisbursementExecutionResult struct {
	Status            string
	ExternalReference string
	FailureCode       string
	FailureMessage    string
	ProviderReference string
	ChannelPayloadPatch map[string]interface{}
}

type LoanDisbursementChannelService struct {
	dtb     *DtbDisbursementService
	mpesa   *MpesaB2cDisbursementService
	airtel  *AirtelDisbursementService
	equitel *EquitelDisbursementService
}

func NewLoanDisbursementChannelService() *LoanDisbursementChannelService {
	return &LoanDisbursementChannelService{
		dtb:     NewDtbDisbursementService(),
		mpesa:   NewMpesaB2cDisbursementService(),
		airtel:  NewAirtelDisbursementService(),
		equitel: NewEquitelDisbursementService(),
	}
}

func (s *LoanDisbursementChannelService) ListActive() ([]models.DisbursementChannel, error) {
	var channels []models.DisbursementChannel
	err := db.DB.Where("is_active = ?", true).Order("sort_order ASC, channel_name ASC").Find(&channels).Error
	return channels, err
}

func (s *LoanDisbursementChannelService) ListAll() ([]models.DisbursementChannel, error) {
	var channels []models.DisbursementChannel
	err := db.DB.Order("sort_order ASC, channel_name ASC").Find(&channels).Error
	return channels, err
}

func (s *LoanDisbursementChannelService) GetByID(id uint64) (*models.DisbursementChannel, error) {
	var ch models.DisbursementChannel
	if err := db.DB.First(&ch, id).Error; err != nil {
		return nil, fmt.Errorf("disbursement channel not found")
	}
	return &ch, nil
}

func normalizeChannelCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(strings.ReplaceAll(code, " ", "_")))
}

func validChannelType(t string) bool {
	switch strings.ToUpper(strings.TrimSpace(t)) {
	case models.DisbursementChannelCash, models.DisbursementChannelBank, models.DisbursementChannelMobile,
		models.DisbursementChannelItem, models.DisbursementChannelWallet:
		return true
	default:
		return false
	}
}

func parseConfigSchema(raw string) datatypes.JSON {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return datatypes.JSON(raw)
	}
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}

func (s *LoanDisbursementChannelService) Create(req dtos.CreateDisbursementChannelRequest, userID uint64) (*models.DisbursementChannel, error) {
	code := normalizeChannelCode(req.ChannelCode)
	if code == "" {
		return nil, fmt.Errorf("channel_code is required")
	}
	channelType := strings.ToUpper(strings.TrimSpace(req.ChannelType))
	if !validChannelType(channelType) {
		return nil, fmt.Errorf("invalid channel_type %q", req.ChannelType)
	}
	var existing models.DisbursementChannel
	if err := db.DB.Unscoped().Where("channel_code = ?", code).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("channel_code %q already exists", code)
	}
	isAsync := false
	if req.IsAsync != nil {
		isAsync = *req.IsAsync
	}
	requiresRef := false
	if req.RequiresExternalRef != nil {
		requiresRef = *req.RequiresExternalRef
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	glRule := strings.TrimSpace(req.GLRuleType)
	if glRule == "" {
		glRule = "LOAN_DISBURSEMENT"
	}
	ch := models.DisbursementChannel{
		BaseModel:           models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		ChannelCode:         code,
		ChannelName:         strings.TrimSpace(req.ChannelName),
		Description:         strings.TrimSpace(req.Description),
		ChannelType:         channelType,
		Provider:            strings.TrimSpace(req.Provider),
		IsAsync:             isAsync,
		RequiresExternalRef: requiresRef,
		GLRuleType:          glRule,
		ConfigSchema:        parseConfigSchema(req.ConfigSchema),
		SortOrder:           sortOrder,
		IsActive:            isActive,
	}
	if ch.ChannelName == "" {
		return nil, fmt.Errorf("channel_name is required")
	}
	if err := db.DB.Create(&ch).Error; err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *LoanDisbursementChannelService) Update(id uint64, req dtos.UpdateDisbursementChannelRequest, userID uint64) (*models.DisbursementChannel, error) {
	var ch models.DisbursementChannel
	if err := db.DB.First(&ch, id).Error; err != nil {
		return nil, fmt.Errorf("disbursement channel not found")
	}
	if req.ChannelName != nil {
		name := strings.TrimSpace(*req.ChannelName)
		if name == "" {
			return nil, fmt.Errorf("channel_name cannot be empty")
		}
		ch.ChannelName = name
	}
	if req.Description != nil {
		ch.Description = strings.TrimSpace(*req.Description)
	}
	if req.ChannelType != nil {
		channelType := strings.ToUpper(strings.TrimSpace(*req.ChannelType))
		if !validChannelType(channelType) {
			return nil, fmt.Errorf("invalid channel_type %q", *req.ChannelType)
		}
		ch.ChannelType = channelType
	}
	if req.Provider != nil {
		ch.Provider = strings.TrimSpace(*req.Provider)
	}
	if req.IsAsync != nil {
		ch.IsAsync = *req.IsAsync
	}
	if req.RequiresExternalRef != nil {
		ch.RequiresExternalRef = *req.RequiresExternalRef
	}
	if req.GLRuleType != nil {
		glRule := strings.TrimSpace(*req.GLRuleType)
		if glRule != "" {
			ch.GLRuleType = glRule
		}
	}
	if req.ConfigSchema != nil {
		ch.ConfigSchema = parseConfigSchema(*req.ConfigSchema)
	}
	if req.SortOrder != nil {
		ch.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		ch.IsActive = *req.IsActive
	}
	ch.UpdatedBy = userID
	if err := db.DB.Save(&ch).Error; err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *LoanDisbursementChannelService) Delete(id uint64) error {
	var ch models.DisbursementChannel
	if err := db.DB.First(&ch, id).Error; err != nil {
		return fmt.Errorf("disbursement channel not found")
	}
	var count int64
	db.DB.Model(&models.LoanDisbursement{}).
		Where("disbursement_channel_id = ? OR channel_code = ?", id, ch.ChannelCode).
		Count(&count)
	if count > 0 {
		return fmt.Errorf("channel has %d disbursement(s); deactivate it instead of deleting", count)
	}
	return db.DB.Delete(&ch).Error
}

func (s *LoanDisbursementChannelService) GetByCode(code string) (*models.DisbursementChannel, error) {
	var ch models.DisbursementChannel
	err := db.DB.Where("channel_code = ? AND is_active = ?", strings.ToUpper(strings.TrimSpace(code)), true).First(&ch).Error
	if err != nil {
		return nil, fmt.Errorf("disbursement channel %q not found or inactive", code)
	}
	return &ch, nil
}

func (s *LoanDisbursementChannelService) ResolveChannelCode(channelCode, legacyMethod string) string {
	if c := strings.ToUpper(strings.TrimSpace(channelCode)); c != "" {
		return c
	}
	switch strings.ToUpper(strings.TrimSpace(legacyMethod)) {
	case "CASH":
		return models.DisbursementCodeCash
	case "MOBILE_MONEY", "MOBILE":
		return models.DisbursementCodeMobileMpesa
	case "WALLET":
		return models.DisbursementCodeWalletMember
	default:
		return models.DisbursementCodeBankCoop
	}
}

func (s *LoanDisbursementChannelService) LegacyMethod(channelType string) string {
	switch channelType {
	case models.DisbursementChannelCash:
		return models.LoanChannelCash
	case models.DisbursementChannelMobile:
		return models.LoanChannelMobileMoney
	case models.DisbursementChannelWallet:
		return "WALLET"
	default:
		return models.LoanChannelBank
	}
}

func (s *LoanDisbursementChannelService) Execute(input DisbursementExecutionInput) DisbursementExecutionResult {
	if input.Channel == nil {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "CHANNEL_MISSING",
			FailureMessage: "disbursement channel is required",
		}
	}

	switch input.Channel.ChannelCode {
	case models.DisbursementCodeCash:
		return s.executeCash(input)
	case models.DisbursementCodeBankCoop, models.DisbursementCodeBankEquity, models.DisbursementCodeBankKCB:
		return s.executeBankTransfer(input)
	case models.DisbursementCodeMobileMpesa:
		return s.executeMobileMpesa(input)
	case models.DisbursementCodeMobileDtbMpesa:
		return s.executeDtbMpesa(input)
	case models.DisbursementCodeMobileAirtel:
		return s.executeMobileAirtel(input)
	case models.DisbursementCodeMobileEquitel:
		return s.executeMobileEquitel(input)
	case models.DisbursementCodeItemPurchase:
		return s.executeItemPurchase(input)
	case models.DisbursementCodeWalletMember:
		return s.executeMemberWallet(input)
	default:
		switch input.Channel.ChannelType {
		case models.DisbursementChannelCash:
			return s.executeCash(input)
		case models.DisbursementChannelBank:
			return s.executeBankTransfer(input)
		case models.DisbursementChannelMobile:
			return s.executeMobileMpesa(input)
		case models.DisbursementChannelItem:
			return s.executeItemPurchase(input)
		case models.DisbursementChannelWallet:
			return s.executeMemberWallet(input)
		default:
			return DisbursementExecutionResult{
				Status:         models.DisbursementStatusFailed,
				FailureCode:    "CHANNEL_UNSUPPORTED",
				FailureMessage: fmt.Sprintf("no handler for channel %s (type %s)", input.Channel.ChannelCode, input.Channel.ChannelType),
			}
		}
	}
}

func payloadString(payload map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if payload == nil {
			return ""
		}
		if v, ok := payload[key]; ok {
			if str, ok := v.(string); ok && strings.TrimSpace(str) != "" {
				return strings.TrimSpace(str)
			}
		}
	}
	return ""
}

func (s *LoanDisbursementChannelService) executeCash(_ DisbursementExecutionInput) DisbursementExecutionResult {
	return DisbursementExecutionResult{Status: models.DisbursementStatusSuccess}
}

func (s *LoanDisbursementChannelService) executeBankTransfer(input DisbursementExecutionInput) DisbursementExecutionResult {
	account := payloadString(input.Payload, "account_number", "bank_account")
	extRef := payloadString(input.Payload, "external_reference", "bank_reference", "transaction_ref")
	if input.Channel.RequiresExternalRef && extRef == "" {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "BANK_REF_REQUIRED",
			FailureMessage: "bank_reference (or external_reference) is required for bank disbursement",
		}
	}
	if account == "" {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "ACCOUNT_REQUIRED",
			FailureMessage: "account_number is required for bank disbursement",
		}
	}
	return DisbursementExecutionResult{
		Status:            models.DisbursementStatusSuccess,
		ExternalReference: extRef,
	}
}

func (s *LoanDisbursementChannelService) executeDtbMpesa(input DisbursementExecutionInput) DisbursementExecutionResult {
	phone := payloadString(input.Payload, "phone_number", "msisdn", "phone")
	result, err := s.dtb.InitiateMpesaTransfer(input.Channel.ID, input.MemberID, input.NetAmount, phone, input.ContractNo)
	if err != nil {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "DTB_INIT_FAILED",
			FailureMessage: err.Error(),
		}
	}
	transferJSON, _ := json.Marshal(result.TransferData)
	return DisbursementExecutionResult{
		Status:            models.DisbursementStatusProcessing,
		ProviderReference: result.ProviderReference,
		FailureMessage:    "DTB M-Pesa OTP sent — enter OTP to complete disbursement",
		ChannelPayloadPatch: map[string]interface{}{
			"phone_number":        result.PhoneNumber,
			"sca_intent_id":       result.SCAIntentID,
			"external_unique_id":  result.ExternalUniqueID,
			"transfer_data":       json.RawMessage(transferJSON),
			"dtb_provider":        "KE_DTB_MPESA",
			"awaiting_otp":        true,
		},
	}
}

func (s *LoanDisbursementChannelService) executeMobileMpesa(input DisbursementExecutionInput) DisbursementExecutionResult {
	phone := payloadString(input.Payload, "phone_number", "msisdn", "phone")
	ref := input.Reference
	if ref == "" {
		ref = fmt.Sprintf("MPESA-%s-%d", input.ContractNo, time.Now().UnixNano())
	}
	result, err := s.mpesa.InitiateB2C(input.Channel.ID, input.MemberID, input.NetAmount, phone, input.ContractNo, ref)
	if err != nil {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "MPESA_INIT_FAILED",
			FailureMessage: err.Error(),
		}
	}
	return DisbursementExecutionResult{
		Status:              models.DisbursementStatusProcessing,
		ProviderReference:   result.ProviderReference,
		FailureMessage:      "Safaricom M-Pesa B2C submitted — awaiting provider confirmation",
		ChannelPayloadPatch: mpesaB2CPayloadPatch(result),
	}
}

func (s *LoanDisbursementChannelService) executeMobileAirtel(input DisbursementExecutionInput) DisbursementExecutionResult {
	phone := payloadString(input.Payload, "phone_number", "msisdn", "phone")
	ref := input.Reference
	if ref == "" {
		ref = fmt.Sprintf("AIRTEL-%s-%d", input.ContractNo, time.Now().UnixNano())
	}
	result, err := s.airtel.InitiateDisbursement(input.Channel.ID, input.MemberID, input.NetAmount, phone, ref)
	if err != nil {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "AIRTEL_INIT_FAILED",
			FailureMessage: err.Error(),
		}
	}
	return DisbursementExecutionResult{
		Status:              models.DisbursementStatusProcessing,
		ProviderReference:   result.ProviderReference,
		FailureMessage:      "Airtel Money payout submitted — awaiting provider confirmation",
		ChannelPayloadPatch: airtelPayloadPatch(result),
	}
}

func (s *LoanDisbursementChannelService) executeMobileEquitel(input DisbursementExecutionInput) DisbursementExecutionResult {
	phone := payloadString(input.Payload, "phone_number", "msisdn", "phone")
	ref := input.Reference
	if ref == "" {
		ref = fmt.Sprintf("EQUITEL-%s-%d", input.ContractNo, time.Now().UnixNano())
	}
	result, err := s.equitel.InitiatePayout(input.Channel.ID, input.MemberID, input.NetAmount, phone, input.ContractNo, ref)
	if err != nil {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "EQUITEL_INIT_FAILED",
			FailureMessage: err.Error(),
		}
	}
	return DisbursementExecutionResult{
		Status:              models.DisbursementStatusProcessing,
		ProviderReference:   result.ProviderReference,
		FailureMessage:      "Equitel payout submitted — awaiting provider confirmation",
		ChannelPayloadPatch: equitelPayloadPatch(result),
	}
}

func (s *LoanDisbursementChannelService) executeItemPurchase(input DisbursementExecutionInput) DisbursementExecutionResult {
	orderRef := payloadString(input.Payload, "store_order_ref", "order_reference", "external_reference")
	if orderRef == "" {
		return DisbursementExecutionResult{
			Status:         models.DisbursementStatusFailed,
			FailureCode:    "ORDER_REF_REQUIRED",
			FailureMessage: "store_order_ref is required for item purchase disbursement",
		}
	}
	return DisbursementExecutionResult{
		Status:            models.DisbursementStatusProcessing,
		ExternalReference: orderRef,
		FailureMessage:    "Item purchase disbursement pending store fulfillment",
	}
}

func (s *LoanDisbursementChannelService) executeMemberWallet(_ DisbursementExecutionInput) DisbursementExecutionResult {
	return DisbursementExecutionResult{
		Status:         models.DisbursementStatusFailed,
		FailureCode:    "NOT_IMPLEMENTED",
		FailureMessage: "member wallet disbursement is not implemented yet",
	}
}

func EnsureDefaultDisbursementChannels() error {
	if err := db.DB.AutoMigrate(&models.DisbursementChannel{}); err != nil {
		return err
	}
	defaults := []models.DisbursementChannel{
		{ChannelCode: models.DisbursementCodeCash, ChannelName: "Cash", Description: "Physical cash paid to member at office", ChannelType: models.DisbursementChannelCash, Provider: "INTERNAL", GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 10, IsActive: true},
		{ChannelCode: models.DisbursementCodeItemPurchase, ChannelName: "Item Purchase", Description: "Loan disbursed as goods/services via store", ChannelType: models.DisbursementChannelItem, Provider: "STORE", IsAsync: true, RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 20, IsActive: true},
		{ChannelCode: models.DisbursementCodeMobileMpesa, ChannelName: "Mobile - M-Pesa", Description: "B2C transfer to member M-Pesa wallet", ChannelType: models.DisbursementChannelMobile, Provider: "MPESA", IsAsync: true, RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 30, IsActive: true, ConfigSchema: datatypes.JSON(`{"required":["phone_number"]}`)},
		{ChannelCode: models.DisbursementCodeMobileDtbMpesa, ChannelName: "Mobile - DTB M-Pesa", Description: "DTB/Astra M-Pesa payout (KE_DTB_MPESA) with OTP confirmation", ChannelType: models.DisbursementChannelMobile, Provider: "DTB", IsAsync: true, RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 35, IsActive: true, ConfigSchema: datatypes.JSON(`{"required":["phone_number"]}`)},
		{ChannelCode: models.DisbursementCodeMobileAirtel, ChannelName: "Mobile - Airtel Money", Description: "Transfer to member Airtel Money wallet", ChannelType: models.DisbursementChannelMobile, Provider: "AIRTEL", IsAsync: true, RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 40, IsActive: true, ConfigSchema: datatypes.JSON(`{"required":["phone_number"]}`)},
		{ChannelCode: models.DisbursementCodeMobileEquitel, ChannelName: "Mobile - Equitel", Description: "Transfer to member Equitel wallet via Equity Finserve", ChannelType: models.DisbursementChannelMobile, Provider: "EQUITY_JENGA", IsAsync: true, RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 45, IsActive: true, ConfigSchema: datatypes.JSON(`{"required":["phone_number"]}`)},
		{ChannelCode: models.DisbursementCodeBankCoop, ChannelName: "Bank - Co-operative Bank", Description: "Transfer to member Co-op Bank account", ChannelType: models.DisbursementChannelBank, Provider: "COOP", RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 50, IsActive: true, ConfigSchema: datatypes.JSON(`{"required":["account_number","bank_reference"]}`)},
		{ChannelCode: models.DisbursementCodeBankEquity, ChannelName: "Bank - Equity Bank", Description: "Transfer to member Equity Bank account", ChannelType: models.DisbursementChannelBank, Provider: "EQUITY", RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 60, IsActive: true, ConfigSchema: datatypes.JSON(`{"required":["account_number","bank_reference"]}`)},
		{ChannelCode: models.DisbursementCodeBankKCB, ChannelName: "Bank - KCB", Description: "Transfer to member KCB account", ChannelType: models.DisbursementChannelBank, Provider: "KCB", RequiresExternalRef: true, GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 70, IsActive: true, ConfigSchema: datatypes.JSON(`{"required":["account_number","bank_reference"]}`)},
		{ChannelCode: models.DisbursementCodeWalletMember, ChannelName: "Member Wallet", Description: "Credit member internal wallet balance", ChannelType: models.DisbursementChannelWallet, Provider: "INTERNAL", GLRuleType: "LOAN_DISBURSEMENT", SortOrder: 80, IsActive: false},
	}
	for _, ch := range defaults {
		var existing models.DisbursementChannel
		if err := db.DB.Where("channel_code = ?", ch.ChannelCode).First(&existing).Error; err == nil {
			continue
		}
		if err := db.DB.Create(&ch).Error; err != nil {
			return err
		}
	}
	return nil
}

func encodeChannelPayload(payload map[string]interface{}) datatypes.JSON {
	if len(payload) == 0 {
		return nil
	}
	b, _ := json.Marshal(payload)
	return datatypes.JSON(b)
}

func applyExecutionResult(disbursement *models.LoanDisbursement, result DisbursementExecutionResult) {
	now := time.Now()
	disbursement.ProcessedAt = &now
	disbursement.Status = result.Status
	disbursement.ExternalReference = result.ExternalReference
	if disbursement.ExternalReference == "" {
		disbursement.ExternalReference = result.ProviderReference
	}
	disbursement.FailureCode = result.FailureCode
	disbursement.FailureMessage = result.FailureMessage
	if len(result.ChannelPayloadPatch) > 0 {
		disbursement.ChannelPayload = mergeChannelPayload([]byte(disbursement.ChannelPayload), result.ChannelPayloadPatch)
	}
	if result.Status == models.DisbursementStatusSuccess {
		disbursement.CompletedAt = &now
	}
}

type postDisburseFinalizeRequest struct {
	Reference           string
	IdempotencyKey      string
	PostToGL            bool
	CreateMilkDeduction bool
	InstallmentAmount   float64
}

func (s *LoanModuleService) finalizeSuccessfulDisbursement(
	tx *gorm.DB,
	contract *models.LoanContract,
	product *models.LoanProduct,
	disbursement *models.LoanDisbursement,
	cashAmount, feesThisDisbursement, procFee, insFee float64,
	disbDate time.Time,
	req postDisburseFinalizeRequest,
	userID uint64,
) error {
	contract.DisbursedAmount = roundMoney(contract.DisbursedAmount + cashAmount)
	if feesThisDisbursement > 0 {
		contract.FeesDeductedAtDisbursement = roundMoney(contract.FeesDeductedAtDisbursement + feesThisDisbursement)
		if procFee > 0 {
			charge := &models.LoanChargeRecord{
				BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
				LoanContractID: contract.ID,
				ChargeType:     "PROCESSING",
				Amount:         procFee,
				Capitalized:    false,
				ChargedDate:    disbDate,
			}
			if err := tx.Create(charge).Error; err != nil {
				return err
			}
		}
		if insFee > 0 {
			charge := &models.LoanChargeRecord{
				BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
				LoanContractID: contract.ID,
				ChargeType:     "INSURANCE",
				Amount:         insFee,
				Capitalized:    false,
				ChargedDate:    disbDate,
			}
			if err := tx.Create(charge).Error; err != nil {
				return err
			}
		}
	}
	contract.Status = models.LoanContractActive
	if contract.DisbursementDate == nil {
		contract.DisbursementDate = &disbDate
	}
	contract.UpdatedBy = userID

	if req.PostToGL {
		idem := req.IdempotencyKey
		if idem == "" {
			idem = req.Reference
		}
		feeMethod := normalizeFeeCollectionMethod(product.FeeCollectionMethod)
		var result *PostFromRuleResult
		var err error
		if feeMethod == models.LoanFeeDeductFromProceeds && feesThisDisbursement > 0 {
			result, err = s.posting.PostLoanNetDisbursement(userID, req.Reference, idem, contract.ContractNo,
				LoanNetDisbursementAmounts{
					GrossPrincipal: roundMoney(cashAmount + feesThisDisbursement),
					NetCash:        cashAmount,
					ProcessingFee:  procFee,
					InsuranceFee:   insFee,
				}, disbDate)
		} else {
			result, err = s.posting.postRule(DomainPostRequest{
				UserID: userID, Reference: req.Reference, IdempotencyKey: idem,
				Amount: cashAmount, TransactionDate: disbDate,
				Description: fmt.Sprintf("Loan disbursement %s via %s", contract.ContractNo, disbursement.ChannelCode),
				HeaderType: "LOAN", RuleType: "LOAN_DISBURSEMENT",
			})
		}
		if err != nil {
			return fmt.Errorf("GL posting failed: %w", err)
		}
		disbursement.GLTransactionID = &result.Transaction.ID
		if feesThisDisbursement > 0 {
			tx.Model(&models.LoanChargeRecord{}).
				Where("loan_contract_id = ? AND charged_date = ? AND gl_transaction_id IS NULL", contract.ID, disbDate).
				Update("gl_transaction_id", result.Transaction.ID)
		}
	}

	if req.CreateMilkDeduction && contract.MilkDeductionEnabled {
		installment := req.InstallmentAmount
		if installment <= 0 && contract.InstallmentAmount != nil {
			installment = *contract.InstallmentAmount
		}
		if installment <= 0 {
			installment = roundMoney(contract.ApprovedAmount / float64(contract.TermMonths))
		}
		rdID, err := s.ensureMilkDeduction(tx, contract, installment, userID)
		if err != nil {
			return err
		}
		contract.RecurrentDeductionID = &rdID
	}

	if s.isContractFullyDisbursed(contract, product) {
		if appID := contract.LoanApplicationID; appID != nil {
			tx.Model(&models.LoanApplication{}).Where("id = ?", *appID).Updates(map[string]interface{}{
				"status": models.LoanAppDisbursed, "updated_by": userID,
			})
		}
	}
	return tx.Save(contract).Error
}
