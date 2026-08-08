package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type LoanModuleService struct {
	schedule    *LoanScheduleService
	posting     *FinancialPostingService
	channels    *LoanDisbursementChannelService
	dtb         *DtbDisbursementService
	channelCfg  *DisbursementChannelConfigService
	providers   *DisbursementProviderClients
}

func NewLoanModuleService() *LoanModuleService {
	return &LoanModuleService{
		schedule:   NewLoanScheduleService(),
		posting:    NewFinancialPostingService(),
		channels:   NewLoanDisbursementChannelService(),
		dtb:        NewDtbDisbursementService(),
		channelCfg: NewDisbursementChannelConfigService(),
		providers:  GetDisbursementProviderClients(),
	}
}

func (s *LoanModuleService) logAudit(tx *gorm.DB, entityType string, entityID uint64, action string, userID uint64, oldVal, newVal interface{}) {
	oldJSON, _ := json.Marshal(oldVal)
	newJSON, _ := json.Marshal(newVal)
	_ = tx.Create(&models.LoanAuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldValues:  oldJSON,
		NewValues:  newJSON,
		UserID:     &userID,
		CreatedAt:  time.Now(),
	}).Error
}

// --- Products ---

func (s *LoanModuleService) CreateProduct(req dtos.CreateLoanProductRequest, userID uint64) (*models.LoanProduct, error) {
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	alloc := req.AllocationPriority
	if len(alloc) == 0 {
		alloc = strings.Split(models.DefaultAllocationOrder, ",")
	}
	allocJSON, _ := json.Marshal(alloc)
	recoveryPriority := req.RecoveryPriority
	if recoveryPriority <= 0 {
		recoveryPriority = 1
	}
	product := &models.LoanProduct{
		BaseModel:             models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		ProductCode:           req.ProductCode,
		ProductName:           req.ProductName,
		Description:           req.Description,
		InterestMethod:        req.InterestMethod,
		InterestRate:          req.InterestRate,
		MinAmount:             req.MinAmount,
		MaxAmount:             req.MaxAmount,
		RepaymentPeriodMonths: req.RepaymentPeriodMonths,
		GracePeriodDays:       req.GracePeriodDays,
		ProcessingFeeRate:     req.ProcessingFeeRate,
		InsuranceFeeRate:      req.InsuranceFeeRate,
		FeeCollectionMethod:   normalizeFeeCollectionMethod(req.FeeCollectionMethod),
		PenaltyRateDaily:      req.PenaltyRateDaily,
		PenaltyRateMonthly:    req.PenaltyRateMonthly,
		PenaltyFixedAmount:    req.PenaltyFixedAmount,
		AllocationPriority:    allocJSON,
		RecoveryPriority:      recoveryPriority,
		IsActive:              active,
	}
	if err := db.DB.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (s *LoanModuleService) ListProducts(page, limit int, activeOnly bool) ([]models.LoanProduct, int64, error) {
	var items []models.LoanProduct
	var total int64
	q := db.DB.Model(&models.LoanProduct{})
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	q.Count(&total)
	offset := (page - 1) * limit
	err := q.Order("product_code ASC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (s *LoanModuleService) GetProduct(id uint64) (*models.LoanProduct, error) {
	var p models.LoanProduct
	if err := db.DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *LoanModuleService) UpdateProduct(id uint64, req dtos.UpdateLoanProductRequest, userID uint64) (*models.LoanProduct, error) {
	var p models.LoanProduct
	if err := db.DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{"updated_by": userID}
	if req.ProductName != nil {
		updates["product_name"] = *req.ProductName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.InterestMethod != nil {
		updates["interest_method"] = *req.InterestMethod
	}
	if req.InterestRate != nil {
		updates["interest_rate"] = *req.InterestRate
	}
	if req.MinAmount != nil {
		updates["min_amount"] = *req.MinAmount
	}
	if req.MaxAmount != nil {
		updates["max_amount"] = *req.MaxAmount
	}
	if req.RepaymentPeriodMonths != nil {
		updates["repayment_period_months"] = *req.RepaymentPeriodMonths
	}
	if req.GracePeriodDays != nil {
		updates["grace_period_days"] = *req.GracePeriodDays
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.AllocationPriority != nil {
		b, _ := json.Marshal(req.AllocationPriority)
		updates["allocation_priority"] = b
	}
	if req.ProcessingFeeRate != nil {
		updates["processing_fee_rate"] = *req.ProcessingFeeRate
	}
	if req.InsuranceFeeRate != nil {
		updates["insurance_fee_rate"] = *req.InsuranceFeeRate
	}
	if req.FeeCollectionMethod != nil {
		updates["fee_collection_method"] = normalizeFeeCollectionMethod(*req.FeeCollectionMethod)
	}
	if req.PenaltyRateDaily != nil {
		updates["penalty_rate_daily"] = *req.PenaltyRateDaily
	}
	if req.PenaltyRateMonthly != nil {
		updates["penalty_rate_monthly"] = *req.PenaltyRateMonthly
	}
	if req.PenaltyFixedAmount != nil {
		updates["penalty_fixed_amount"] = *req.PenaltyFixedAmount
	}
	if req.RecoveryPriority != nil {
		updates["recovery_priority"] = *req.RecoveryPriority
	}
	if err := db.DB.Model(&p).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetProduct(id)
}

// --- Applications ---

func (s *LoanModuleService) CreateApplication(req dtos.CreateLoanApplicationRequest, userID uint64) (*models.LoanApplication, error) {
	var product models.LoanProduct
	if err := db.DB.First(&product, req.LoanProductID).Error; err != nil {
		return nil, fmt.Errorf("loan product not found")
	}
	if !product.IsActive {
		return nil, fmt.Errorf("loan product is inactive")
	}
	if req.RequestedAmount < product.MinAmount || req.RequestedAmount > product.MaxAmount {
		return nil, fmt.Errorf("requested amount outside product limits (%.2f - %.2f)", product.MinAmount, product.MaxAmount)
	}

	appNo := fmt.Sprintf("LN-APP-%d", time.Now().UnixNano())
	var expected *time.Time
	if req.ExpectedDisbursementDate != "" {
		t := utils.ParseFlexibleDate(req.ExpectedDisbursementDate)
		if !t.IsZero() {
			expected = &t
		}
	}

	var app *models.LoanApplication
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		app = &models.LoanApplication{
			BaseModel:                models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			ApplicationNo:            appNo,
			MemberID:                 req.MemberID,
			LoanProductID:            req.LoanProductID,
			RequestedAmount:          req.RequestedAmount,
			RequestedTermMonths:      req.RequestedTermMonths,
			Purpose:                  req.Purpose,
			Status:                   models.LoanAppDraft,
			DateApplied:              time.Now(),
			ExpectedDisbursementDate: expected,
		}
		if err := tx.Create(app).Error; err != nil {
			return err
		}
		for _, g := range req.Guarantors {
			if err := tx.Create(&models.LoanGuarantor{
				BaseModel:         models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
				LoanApplicationID: app.ID,
				MemberID:          g.MemberID,
				GuaranteedAmount:  g.GuaranteedAmount,
				Relationship:      g.Relationship,
			}).Error; err != nil {
				return err
			}
		}
		s.logAudit(tx, "loan_application", app.ID, "CREATE", userID, nil, app)
		return nil
	})
	return app, err
}

func (s *LoanModuleService) SubmitApplication(id, userID uint64) (*models.LoanApplication, error) {
	var app models.LoanApplication
	if err := db.DB.First(&app, id).Error; err != nil {
		return nil, err
	}
	if app.Status != models.LoanAppDraft {
		return nil, fmt.Errorf("only DRAFT applications can be submitted")
	}
	app.Status = models.LoanAppSubmitted
	app.UpdatedBy = userID
	if err := db.DB.Save(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *LoanModuleService) ApproveApplication(id uint64, req dtos.ApproveLoanApplicationRequest, userID uint64) (*models.LoanApplication, error) {
	var app models.LoanApplication
	if err := db.DB.Preload("Guarantors").First(&app, id).Error; err != nil {
		return nil, err
	}
	if app.Status != models.LoanAppSubmitted && app.Status != models.LoanAppUnderReview {
		return nil, fmt.Errorf("application cannot be approved in status %s", app.Status)
	}

	approvedAmt := req.ApprovedAmount
	if approvedAmt <= 0 {
		approvedAmt = app.RequestedAmount
	}
	approvedTerm := req.ApprovedTermMonths
	if approvedTerm <= 0 {
		approvedTerm = app.RequestedTermMonths
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		approval := &models.LoanApproval{
			BaseModel:          models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			LoanApplicationID: app.ID,
			ApprovalRole:       req.ApprovalRole,
			ApproverID:         userID,
			ApprovalDate:       time.Now(),
			Decision:           "APPROVED",
			ApprovedAmount:     &approvedAmt,
			ApprovedTermMonths: &approvedTerm,
			Comments:           req.Comments,
			Conditions:         req.Conditions,
		}
		if err := tx.Create(approval).Error; err != nil {
			return err
		}

		app.Status = models.LoanAppUnderReview
		app.ApprovedAmount = &approvedAmt
		app.ApprovedTermMonths = &approvedTerm
		app.UpdatedBy = userID

		if req.FinalApproval {
			app.Status = models.LoanAppApproved
			contract, err := s.createContractFromApplication(tx, &app, approvedAmt, approvedTerm, userID)
			if err != nil {
				return err
			}
			app.LoanContractID = &contract.ID
		}
		return tx.Save(&app).Error
	})
	if err != nil {
		return nil, err
	}
	return s.GetApplication(id)
}

func (s *LoanModuleService) RejectApplication(id uint64, req dtos.RejectLoanApplicationRequest, userID uint64) (*models.LoanApplication, error) {
	var app models.LoanApplication
	if err := db.DB.First(&app, id).Error; err != nil {
		return nil, err
	}
	app.Status = models.LoanAppRejected
	app.RejectedReason = req.Comments
	app.UpdatedBy = userID
	if err := db.DB.Save(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *LoanModuleService) GetApplication(id uint64) (*models.LoanApplication, error) {
	var app models.LoanApplication
	if err := db.DB.Preload("LoanProduct").Preload("Guarantors").First(&app, id).Error; err != nil {
		return nil, err
	}
	apps := []models.LoanApplication{app}
	enrichApplicationsWithMembers(apps)
	return &apps[0], nil
}

func (s *LoanModuleService) ListApplications(page, limit int, status string, memberID uint64) ([]models.LoanApplication, int64, error) {
	var items []models.LoanApplication
	var total int64
	q := db.DB.Model(&models.LoanApplication{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if memberID > 0 {
		q = q.Where("member_id = ?", memberID)
	}
	q.Count(&total)
	offset := (page - 1) * limit
	err := q.Preload("LoanProduct").Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return items, total, err
	}
	enrichApplicationsWithMembers(items)
	return items, total, nil
}

func (s *LoanModuleService) createContractFromApplication(tx *gorm.DB, app *models.LoanApplication, amount float64, term int, userID uint64) (*models.LoanContract, error) {
	var product models.LoanProduct
	if err := tx.First(&product, app.LoanProductID).Error; err != nil {
		return nil, err
	}
	contractNo := fmt.Sprintf("LN-%d-%d", app.MemberID, time.Now().Unix())
	contract := &models.LoanContract{
		BaseModel:            models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		ContractNo:           contractNo,
		LoanApplicationID:    &app.ID,
		MemberID:             app.MemberID,
		LoanProductID:        app.LoanProductID,
		PrincipalAmount:      amount,
		ApprovedAmount:       amount,
		InterestRate:         product.InterestRate,
		InterestMethod:       product.InterestMethod,
		TermMonths:           term,
		GracePeriodDays:      product.GracePeriodDays,
		Status:               models.LoanContractPending,
		OutstandingPrincipal: amount,
		MilkDeductionEnabled: true,
		MaxDeductionPercent:  50,
	}
	if err := tx.Create(contract).Error; err != nil {
		return nil, err
	}

	procFee, insFee := computeProductFees(&product, amount)
	scheduleProcFee, scheduleInsFee := procFee, insFee
	if !feesSpreadOverInstallments(product.FeeCollectionMethod) {
		scheduleProcFee, scheduleInsFee = 0, 0
	}
	startDate := time.Now()
	if app.ExpectedDisbursementDate != nil {
		startDate = *app.ExpectedDisbursementDate
	}
	lines := s.schedule.Generate(scheduleInput{
		Principal:       amount,
		AnnualRate:      product.InterestRate,
		TermMonths:      term,
		Method:          product.InterestMethod,
		GracePeriodDays: product.GracePeriodDays,
		StartDate:       startDate,
		ProcessingFee:   scheduleProcFee,
		InsuranceFee:    scheduleInsFee,
	})
	var installmentAmt float64
	if len(lines) > 0 {
		installmentAmt = lines[0].TotalDue
	}
	contract.InstallmentAmount = &installmentAmt
	maturity := lines[len(lines)-1].DueDate
	contract.MaturityDate = &maturity

	for _, line := range lines {
		inst := models.LoanScheduleInstallment{
			BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			LoanContractID: contract.ID,
			InstallmentNo:  line.InstallmentNo,
			DueDate:        line.DueDate,
			PrincipalDue:   line.PrincipalDue,
			InterestDue:    line.InterestDue,
			FeeDue:         line.FeeDue,
			InsuranceDue:   line.InsuranceDue,
			TotalDue:       line.TotalDue,
			BalanceAfter:   line.BalanceAfter,
			Status:         "PENDING",
		}
		if err := tx.Create(&inst).Error; err != nil {
			return nil, err
		}
	}
	if feesSpreadOverInstallments(product.FeeCollectionMethod) && procFee > 0 {
		contract.OutstandingFees = procFee
	}
	return contract, tx.Save(contract).Error
}

func (s *LoanModuleService) GetDisbursementQuote(contractID uint64) (*dtos.LoanDisbursementQuote, error) {
	contract, err := s.GetContract(contractID)
	if err != nil {
		return nil, err
	}
	var product models.LoanProduct
	if err := db.DB.First(&product, contract.LoanProductID).Error; err != nil {
		return nil, err
	}
	procFee, insFee := computeProductFees(&product, contract.ApprovedAmount)
	totalFees := roundMoney(procFee + insFee)
	remainingNet := s.remainingDisbursementNet(contract, &product)
	quote := &dtos.LoanDisbursementQuote{
		ApprovedAmount:   contract.ApprovedAmount,
		ProcessingFee:    procFee,
		InsuranceFee:     insFee,
		TotalFees:        totalFees,
		AlreadyDisbursed: contract.DisbursedAmount,
		FeesDeducted:     contract.FeesDeductedAtDisbursement,
		RemainingNet:     remainingNet,
	}
	if feesDeductedFromProceeds(product.FeeCollectionMethod) {
		if contract.FeesDeductedAtDisbursement == 0 {
			quote.NetDisbursement = roundMoney(contract.ApprovedAmount - totalFees)
		} else {
			quote.NetDisbursement = remainingNet
		}
	} else {
		quote.NetDisbursement = remainingNet
	}
	return quote, nil
}

// --- Contracts ---

func (s *LoanModuleService) GetContract(id uint64) (*models.LoanContract, error) {
	var c models.LoanContract
	if err := db.DB.Preload("LoanProduct").Preload("Installments", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("installment_no ASC")
	}).Preload("Disbursements", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("id DESC")
	}).First(&c, id).Error; err != nil {
		return nil, err
	}
	contracts := []models.LoanContract{c}
	enrichContractsWithMembers(contracts)
	return &contracts[0], nil
}

func (s *LoanModuleService) ListContracts(page, limit int, status string, memberID uint64) ([]models.LoanContract, int64, error) {
	var items []models.LoanContract
	var total int64
	q := db.DB.Model(&models.LoanContract{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if memberID > 0 {
		q = q.Where("member_id = ?", memberID)
	}
	q.Count(&total)
	offset := (page - 1) * limit
	err := q.Preload("LoanProduct").Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return items, total, err
	}
	enrichContractsWithMembers(items)
	return items, total, nil
}

func (s *LoanModuleService) ListDisbursementChannels(all bool) ([]models.DisbursementChannel, error) {
	if all {
		return s.channels.ListAll()
	}
	return s.channels.ListActive()
}

func (s *LoanModuleService) GetDisbursementChannel(id uint64) (*models.DisbursementChannel, error) {
	return s.channels.GetByID(id)
}

func (s *LoanModuleService) CreateDisbursementChannel(req dtos.CreateDisbursementChannelRequest, userID uint64) (*models.DisbursementChannel, error) {
	return s.channels.Create(req, userID)
}

func (s *LoanModuleService) UpdateDisbursementChannel(id uint64, req dtos.UpdateDisbursementChannelRequest, userID uint64) (*models.DisbursementChannel, error) {
	return s.channels.Update(id, req, userID)
}

func (s *LoanModuleService) DeleteDisbursementChannel(id uint64) error {
	return s.channels.Delete(id)
}

func (s *LoanModuleService) DisburseContract(id uint64, req dtos.DisburseLoanRequest, userID uint64) (*models.LoanDisbursement, error) {
	if req.IdempotencyKey != "" {
		var existing models.LoanDisbursement
		if err := db.DB.Where("idempotency_key = ?", req.IdempotencyKey).First(&existing).Error; err == nil {
			return &existing, nil
		}
	}

	channelCode := s.channels.ResolveChannelCode(req.ChannelCode, req.Method)
	channel, err := s.channels.GetByCode(channelCode)
	if err != nil {
		return nil, err
	}

	var disbursement *models.LoanDisbursement
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var contract models.LoanContract
		if err := tx.First(&contract, id).Error; err != nil {
			return err
		}
		var product models.LoanProduct
		if err := tx.First(&product, contract.LoanProductID).Error; err != nil {
			return err
		}
		if contract.Status != models.LoanContractPending && contract.Status != models.LoanContractActive {
			return fmt.Errorf("contract cannot be disbursed in status %s", contract.Status)
		}

		procFee, insFee := computeProductFees(&product, contract.ApprovedAmount)
		totalFees := roundMoney(procFee + insFee)
		remainingNet := s.remainingDisbursementNet(&contract, &product)
		if remainingNet <= 0 {
			return fmt.Errorf("contract is already fully disbursed")
		}

		cashAmount := req.Amount
		feeMethod := normalizeFeeCollectionMethod(product.FeeCollectionMethod)
		var feesThisDisbursement float64

		if feeMethod == models.LoanFeeDeductFromProceeds {
			if contract.FeesDeductedAtDisbursement == 0 {
				feesThisDisbursement = totalFees
			}
			expectedNet := roundMoney(remainingNet)
			if cashAmount <= 0 || cashAmount >= contract.ApprovedAmount-contract.DisbursedAmount-0.01 {
				cashAmount = expectedNet
			}
			if cashAmount > expectedNet+0.01 {
				return fmt.Errorf("net disbursement exceeds remaining amount (%.2f)", expectedNet)
			}
		} else {
			if cashAmount <= 0 {
				cashAmount = remainingNet
			}
			if cashAmount > remainingNet+0.01 {
				return fmt.Errorf("disbursement exceeds remaining approved amount (%.2f)", remainingNet)
			}
		}
		if cashAmount <= 0 {
			return fmt.Errorf("disbursement amount must be positive")
		}

		disbDate := utils.ParseFlexibleDate(req.DisbursementDate)
		if disbDate.IsZero() {
			disbDate = time.Now()
		}

		grossAmount := roundMoney(cashAmount + feesThisDisbursement)
		disbursement = &models.LoanDisbursement{
			BaseModel:             models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			LoanContractID:        contract.ID,
			DisbursementChannelID: &channel.ID,
			ChannelCode:           channel.ChannelCode,
			Status:                models.DisbursementStatusPending,
			GrossAmount:           grossAmount,
			FeesDeducted:          feesThisDisbursement,
			Amount:                cashAmount,
			Method:                s.channels.LegacyMethod(channel.ChannelType),
			Reference:             req.Reference,
			IdempotencyKey:        req.IdempotencyKey,
			ChannelPayload:        encodeChannelPayload(req.ChannelPayload),
			DisbursementDate:      disbDate,
			Notes:                 req.Notes,
		}
		if err := tx.Create(disbursement).Error; err != nil {
			return err
		}

		result := s.channels.Execute(DisbursementExecutionInput{
			Channel:    channel,
			NetAmount:  cashAmount,
			Payload:    req.ChannelPayload,
			ContractNo: contract.ContractNo,
			MemberID:   contract.MemberID,
			Reference:  req.Reference,
		})
		applyExecutionResult(disbursement, result)

		switch result.Status {
		case models.DisbursementStatusSuccess:
			if err := s.finalizeSuccessfulDisbursement(tx, &contract, &product, disbursement,
				cashAmount, feesThisDisbursement, procFee, insFee, disbDate,
				postDisburseFinalizeRequest{
					Reference: req.Reference, IdempotencyKey: req.IdempotencyKey,
					PostToGL: req.PostToGL, CreateMilkDeduction: req.CreateMilkDeduction,
					InstallmentAmount: req.InstallmentAmount,
				}, userID); err != nil {
				disbursement.Status = models.DisbursementStatusFailed
				disbursement.FailureCode = "FINALIZE_FAILED"
				disbursement.FailureMessage = err.Error()
				tx.Save(disbursement)
				return err
			}
		case models.DisbursementStatusFailed:
			return fmt.Errorf("%s: %s", result.FailureCode, result.FailureMessage)
		}

		return tx.Save(disbursement).Error
	})
	return disbursement, err
}

func (s *LoanModuleService) ConfirmDisbursement(disbursementID uint64, req dtos.ConfirmLoanDisbursementRequest, userID uint64) (*models.LoanDisbursement, error) {
	var disbursement models.LoanDisbursement
	if err := db.DB.First(&disbursement, disbursementID).Error; err != nil {
		return nil, err
	}
	if disbursement.Status == models.DisbursementStatusSuccess {
		return &disbursement, nil
	}
	if disbursement.Status != models.DisbursementStatusProcessing && disbursement.Status != models.DisbursementStatusPending {
		return nil, fmt.Errorf("disbursement cannot be confirmed in status %s", disbursement.Status)
	}

	var otpPayloadPatch []byte
	if disbursement.ChannelCode == models.DisbursementCodeMobileDtbMpesa && strings.TrimSpace(req.OTP) != "" {
		payload := decodeChannelPayload([]byte(disbursement.ChannelPayload))
		scaIntentID := strings.TrimSpace(req.ScaIntentID)
		if scaIntentID == "" {
			scaIntentID = payloadStringFromJSON(payload, "sca_intent_id")
		}
		transferData := extractTransferData(payload)
		var contract models.LoanContract
		if err := db.DB.First(&contract, disbursement.LoanContractID).Error; err != nil {
			return nil, err
		}
		channelID, err := s.resolveDisbursementChannelID(&disbursement)
		if err != nil {
			return nil, err
		}
		result, err := s.dtb.CompleteMpesaTransfer(channelID, contract.MemberID, scaIntentID, strings.TrimSpace(req.OTP), transferData)
		if err != nil {
			return nil, err
		}
		req.Success = true
		if req.ExternalReference == "" {
			req.ExternalReference = result.ExternalReference
		}
		patch := map[string]interface{}{
			"awaiting_otp":     false,
			"otp_confirmed_at": time.Now().Format(time.RFC3339),
		}
		if result.Raw != nil {
			patch["provider_response"] = result.Raw
		}
		otpPayloadPatch = mergeChannelPayload([]byte(disbursement.ChannelPayload), patch)
	} else if disbursement.ChannelCode == models.DisbursementCodeMobileDtbMpesa && !req.Success && strings.TrimSpace(req.OTP) == "" {
		return nil, fmt.Errorf("DTB M-Pesa disbursement requires OTP confirmation")
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&disbursement, disbursementID).Error; err != nil {
			return err
		}
		if len(otpPayloadPatch) > 0 {
			disbursement.ChannelPayload = otpPayloadPatch
		}
		var contract models.LoanContract
		if err := tx.First(&contract, disbursement.LoanContractID).Error; err != nil {
			return err
		}
		var product models.LoanProduct
		if err := tx.First(&product, contract.LoanProductID).Error; err != nil {
			return err
		}

		now := time.Now()
		disbursement.ProcessedAt = &now
		if !req.Success {
			disbursement.Status = models.DisbursementStatusFailed
			disbursement.FailureCode = req.FailureCode
			if disbursement.FailureCode == "" {
				disbursement.FailureCode = "PROVIDER_FAILED"
			}
			disbursement.FailureMessage = req.FailureMessage
			return tx.Save(&disbursement).Error
		}

		if req.ExternalReference != "" {
			disbursement.ExternalReference = req.ExternalReference
		}
		disbursement.Status = models.DisbursementStatusSuccess
		disbursement.CompletedAt = &now
		disbursement.FailureCode = ""
		disbursement.FailureMessage = ""

		procFee, insFee := computeProductFees(&product, contract.ApprovedAmount)
		feesThisDisbursement := disbursement.FeesDeducted

		postGL := req.PostToGL
		if !postGL && disbursement.GLTransactionID == nil {
			postGL = true
		}
		createMilk := req.CreateMilkDeduction

		if err := s.finalizeSuccessfulDisbursement(tx, &contract, &product, &disbursement,
			disbursement.Amount, feesThisDisbursement, procFee, insFee, disbursement.DisbursementDate,
			postDisburseFinalizeRequest{
				Reference: disbursement.Reference, IdempotencyKey: disbursement.IdempotencyKey,
				PostToGL: postGL, CreateMilkDeduction: createMilk,
			}, userID); err != nil {
			return err
		}
		return tx.Save(&disbursement).Error
	})
	if err != nil {
		return nil, err
	}
	return &disbursement, nil
}

func extractTransferData(payload map[string]interface{}) map[string]interface{} {
	if raw, ok := payload["transfer_data"]; ok {
		switch v := raw.(type) {
		case map[string]interface{}:
			return v
		case string:
			var out map[string]interface{}
			if err := json.Unmarshal([]byte(v), &out); err == nil {
				return out
			}
		case json.RawMessage:
			var out map[string]interface{}
			if err := json.Unmarshal(v, &out); err == nil {
				return out
			}
		}
	}
	return payload
}

func (s *LoanModuleService) ProcessWithdrawalCallback(data map[string]interface{}) error {
	extID := ""
	for _, key := range []string{"externalUniqueId", "external_unique_id", "externalUniqueID"} {
		if v, ok := data[key]; ok {
			extID = fmt.Sprintf("%v", v)
			break
		}
	}
	if extID == "" {
		return fmt.Errorf("withdrawal callback missing externalUniqueId")
	}

	var disbursements []models.LoanDisbursement
	if err := db.DB.Where(
		"channel_code = ? AND status IN ? AND JSON_UNQUOTE(JSON_EXTRACT(channel_payload, '$.external_unique_id')) = ?",
		models.DisbursementCodeMobileDtbMpesa,
		[]string{models.DisbursementStatusProcessing, models.DisbursementStatusPending},
		extID,
	).Find(&disbursements).Error; err != nil {
		return err
	}
	if len(disbursements) == 0 {
		return fmt.Errorf("no matching loan disbursement for externalUniqueId %s", extID)
	}

	status := strings.ToUpper(fmt.Sprintf("%v", data["status"]))
	success := status == "SUCCESS" || status == "SUCCESSFUL" || status == "COMPLETED"
	failure := status == "FAILED" || status == "ERROR" || status == "CANCELLED"

	for _, d := range disbursements {
		patch := map[string]interface{}{"callback": data}
		_ = db.DB.Model(&models.LoanDisbursement{}).Where("id = ?", d.ID).
			Update("channel_payload", mergeChannelPayload([]byte(d.ChannelPayload), patch))

		if success {
			extRef := extID
			if v, ok := data["withdrawalId"]; ok {
				extRef = fmt.Sprintf("%v", v)
			}
			_, err := s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
				Success:           true,
				ExternalReference: extRef,
				PostToGL:          true,
				CreateMilkDeduction: true,
			}, d.CreatedBy)
			if err != nil {
				return err
			}
		} else if failure {
			_, _ = s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
				Success:        false,
				FailureCode:    "PROVIDER_FAILED",
				FailureMessage: fmt.Sprintf("withdrawal callback status: %s", status),
			}, d.CreatedBy)
		}
	}
	return nil
}

type mobileCallbackOutcome struct {
	Success        bool
	ExternalRef    string
	FailureCode    string
	FailureMessage string
	Raw            map[string]interface{}
}

func (s *LoanModuleService) processMobileDisbursementCallback(channelCode, matchReference string, outcome mobileCallbackOutcome) error {
	if matchReference == "" {
		return fmt.Errorf("callback missing match reference")
	}
	var disbursements []models.LoanDisbursement
	if err := db.DB.Where(
		`channel_code = ? AND status IN ? AND (
			reference = ? OR
			JSON_UNQUOTE(JSON_EXTRACT(channel_payload, '$.client_reference')) = ? OR
			JSON_UNQUOTE(JSON_EXTRACT(channel_payload, '$.provider_transaction_id')) = ?
		)`,
		channelCode,
		[]string{models.DisbursementStatusProcessing, models.DisbursementStatusPending},
		matchReference, matchReference, matchReference,
	).Find(&disbursements).Error; err != nil {
		return err
	}
	if len(disbursements) == 0 {
		return fmt.Errorf("no matching loan disbursement for reference %s (channel %s)", matchReference, channelCode)
	}

	for _, d := range disbursements {
		if d.Status == models.DisbursementStatusSuccess {
			continue
		}
		patch := map[string]interface{}{"callback_result": outcome.Raw}
		_ = db.DB.Model(&models.LoanDisbursement{}).Where("id = ?", d.ID).
			Update("channel_payload", mergeChannelPayload([]byte(d.ChannelPayload), patch))

		if outcome.Success {
			_, err := s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
				Success:             true,
				ExternalReference:   outcome.ExternalRef,
				PostToGL:            true,
				CreateMilkDeduction: true,
			}, d.CreatedBy)
			if err != nil {
				return err
			}
		} else {
			_, _ = s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
				Success:        false,
				FailureCode:    outcome.FailureCode,
				FailureMessage: outcome.FailureMessage,
			}, d.CreatedBy)
		}
	}
	return nil
}

func (s *LoanModuleService) ProcessMpesaB2CResult(data map[string]interface{}) error {
	originatorID, transactionID, resultCode, resultDesc := parseMpesaResultCallback(data)
	if originatorID == "" {
		return fmt.Errorf("mpesa callback missing OriginatorConversationID")
	}
	success := resultCode == "0"
	outcome := mobileCallbackOutcome{
		Success:        success,
		ExternalRef:    transactionID,
		FailureCode:    resultCode,
		FailureMessage: resultDesc,
		Raw:            data,
	}
	if !success && outcome.FailureCode == "" {
		outcome.FailureCode = "MPESA_FAILED"
	}
	if !success && outcome.FailureMessage == "" {
		outcome.FailureMessage = "M-Pesa B2C disbursement failed"
	}
	return s.processMobileDisbursementCallback(models.DisbursementCodeMobileMpesa, originatorID, outcome)
}

func (s *LoanModuleService) ProcessMpesaB2CTimeout(data map[string]interface{}) error {
	originatorID, _, _, _ := parseMpesaResultCallback(data)
	if originatorID == "" {
		if v, ok := data["OriginatorConversationID"]; ok {
			originatorID = fmt.Sprintf("%v", v)
		}
	}
	if originatorID == "" {
		return fmt.Errorf("mpesa timeout callback missing reference")
	}
	return s.processMobileDisbursementCallback(models.DisbursementCodeMobileMpesa, originatorID, mobileCallbackOutcome{
		Success:        false,
		FailureCode:    "QUEUE_TIMEOUT",
		FailureMessage: "M-Pesa B2C queue timeout",
		Raw:            data,
	})
}

func (s *LoanModuleService) ProcessAirtelDisbursementCallback(data map[string]interface{}) error {
	txnID, statusCode, airtelMoneyID, success := parseAirtelCallback(data)
	matchRef := txnID
	if matchRef == "" || matchRef == "<nil>" {
		if nested, ok := data["transaction"].(map[string]interface{}); ok {
			matchRef = fmt.Sprintf("%v", nested["reference"])
		}
	}
	extRef := airtelMoneyID
	if extRef == "" || extRef == "<nil>" {
		extRef = txnID
	}
	return s.processMobileDisbursementCallback(models.DisbursementCodeMobileAirtel, matchRef, mobileCallbackOutcome{
		Success:        success,
		ExternalRef:    extRef,
		FailureCode:    statusCode,
		FailureMessage: fmt.Sprintf("Airtel disbursement status: %s", statusCode),
		Raw:            data,
	})
}

func (s *LoanModuleService) ProcessJengaMobileCallback(data map[string]interface{}) error {
	reference, txnRef, responseCode, success := parseJengaCallback(data)
	matchRef := reference
	if matchRef == "" || matchRef == "<nil>" {
		matchRef = txnRef
	}
	extRef := txnRef
	if extRef == "" || extRef == "<nil>" {
		extRef = reference
	}
	return s.processMobileDisbursementCallback(models.DisbursementCodeMobileEquitel, matchRef, mobileCallbackOutcome{
		Success:        success,
		ExternalRef:    extRef,
		FailureCode:    responseCode,
		FailureMessage: fmt.Sprintf("Jenga Equitel status: %s", responseCode),
		Raw:            data,
	})
}

func (s *LoanModuleService) resolveDisbursementChannelID(d *models.LoanDisbursement) (uint64, error) {
	if d.DisbursementChannelID != nil && *d.DisbursementChannelID > 0 {
		return *d.DisbursementChannelID, nil
	}
	if d.ChannelCode == "" {
		return 0, fmt.Errorf("disbursement channel is missing")
	}
	ch, err := s.channels.GetByCode(d.ChannelCode)
	if err != nil {
		return 0, err
	}
	return ch.ID, nil
}

func (s *LoanModuleService) GetDisbursementChannelConfig(channelID uint64) ([]dtos.DisbursementChannelConfigItem, error) {
	return s.channelCfg.ListForAdmin(channelID)
}

func (s *LoanModuleService) UpdateDisbursementChannelConfig(channelID uint64, items []dtos.DisbursementChannelConfigUpdateItem) error {
	return s.channelCfg.Upsert(channelID, items)
}

func (s *LoanModuleService) GetDisbursementProviderStatus(disbursementID uint64) (map[string]interface{}, error) {
	var d models.LoanDisbursement
	if err := db.DB.First(&d, disbursementID).Error; err != nil {
		return nil, err
	}
	if d.ChannelCode != models.DisbursementCodeMobileMpesa {
		return map[string]interface{}{
			"disbursement_id": d.ID,
			"channel_code":    d.ChannelCode,
			"status":          d.Status,
			"message":         "provider status query only supported for MOBILE_MPESA",
		}, nil
	}
	payload := decodeChannelPayload([]byte(d.ChannelPayload))
	txnID := payloadStringFromJSON(payload, "provider_transaction_id")
	originatorID := payloadStringFromJSON(payload, "client_reference", "originator_conversation_id")
	if originatorID == "" {
		originatorID = d.Reference
	}
	channelID, err := s.resolveDisbursementChannelID(&d)
	if err != nil {
		return nil, err
	}
	mpesaDaraja, err := s.providers.Mpesa(channelID)
	if err != nil {
		return nil, err
	}
	if !mpesaDaraja.Configured() {
		return nil, fmt.Errorf("M-Pesa Daraja is not configured")
	}
	result, err := mpesaDaraja.QueryTransactionStatus(originatorID, txnID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"disbursement_id": d.ID,
		"channel_code":    d.ChannelCode,
		"status":          d.Status,
		"provider_result": result,
	}, nil
}

func (s *LoanModuleService) ensureMilkDeduction(tx *gorm.DB, contract *models.LoanContract, installment float64, userID uint64) (uint64, error) {
	dedTypeID, err := s.loanDeductionTypeID(tx)
	if err != nil {
		return 0, err
	}
	ref := fmt.Sprintf("LOAN-%s", contract.ContractNo)
	var existing models.RecurrentDeduction
	if err := tx.Where("reference = ? AND customer_id = ? AND settled = 0", ref, contract.MemberID).First(&existing).Error; err == nil {
		return existing.ID, nil
	}
	rd := &models.RecurrentDeduction{
		BaseModel:       models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		CustomerID:      contract.MemberID,
		TotalAmount:     contract.ApprovedAmount,
		PaidAmount:      contract.TotalPaid,
		RecurrentAmount: installment,
		DeductionTypeID: dedTypeID,
		Reference:       ref,
		CustomerType:    "member",
		Settled:         0,
		PrincipalAmount: contract.OutstandingPrincipal,
		TransactionDate: time.Now(),
	}
	if err := tx.Create(rd).Error; err != nil {
		return 0, err
	}
	return rd.ID, nil
}

func (s *LoanModuleService) loanDeductionTypeID(tx *gorm.DB) (uint64, error) {
	var dt models.DeductionType
	err := tx.Where("code = ?", "LOAN").First(&dt).Error
	if err == nil {
		return dt.ID, nil
	}
	dt = models.DeductionType{Code: "LOAN", Description: "Loan Repayment", Status: "active"}
	if err := tx.Create(&dt).Error; err != nil {
		return 0, err
	}
	return dt.ID, nil
}

// --- Repayments ---

func (s *LoanModuleService) RecordRepayment(contractID uint64, req dtos.RecordLoanRepaymentRequest, userID uint64) (*models.LoanRepaymentRecord, error) {
	if req.IdempotencyKey != "" {
		var existing models.LoanRepaymentRecord
		if err := db.DB.Where("idempotency_key = ?", req.IdempotencyKey).First(&existing).Error; err == nil {
			return &existing, nil
		}
	}

	var repayment *models.LoanRepaymentRecord
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var contract models.LoanContract
		if err := tx.Preload("LoanProduct").First(&contract, contractID).Error; err != nil {
			return err
		}
		if contract.Status == models.LoanContractSettled || contract.Status == models.LoanContractWrittenOff {
			return fmt.Errorf("loan is closed")
		}

		payDate := utils.ParseFlexibleDate(req.PaymentDate)
		if payDate.IsZero() {
			payDate = time.Now()
		}

		repayment = &models.LoanRepaymentRecord{
			BaseModel:       models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			LoanContractID:  contract.ID,
			Amount:          req.Amount,
			Channel:         req.Channel,
			Reference:       req.Reference,
			IdempotencyKey:  req.IdempotencyKey,
			PaymentDate:     payDate,
			MemberPayrollID: req.MemberPayrollID,
			MemberPayslipID: req.MemberPayslipID,
			Notes:           req.Notes,
		}
		if err := tx.Create(repayment).Error; err != nil {
			return err
		}

		allocs, err := s.allocatePayment(tx, &contract, repayment, userID)
		if err != nil {
			return err
		}
		for i := range allocs {
			allocs[i].LoanRepaymentID = repayment.ID
			if err := tx.Create(&allocs[i]).Error; err != nil {
				return err
			}
		}

		if req.PostToGL {
			idem := req.IdempotencyKey
			if idem == "" {
				idem = req.Reference
			}
			ruleType := "LOAN_REPAYMENT"
			if req.Channel == models.LoanChannelMilkDeduction {
				ruleType = "MEMBER_REPAYMENT_DEDUCTION"
			}
			result, err := s.posting.postRule(DomainPostRequest{
				UserID: userID, Reference: req.Reference, IdempotencyKey: idem,
				Amount: req.Amount, TransactionDate: payDate,
				Description: fmt.Sprintf("Loan repayment %s via %s", contract.ContractNo, req.Channel),
				HeaderType: "LOAN", RuleType: ruleType,
			})
			if err != nil {
				return fmt.Errorf("GL posting failed: %w", err)
			}
			repayment.GLTransactionID = &result.Transaction.ID
			tx.Save(repayment)
		}

		s.syncRecurrentDeduction(tx, &contract)
		s.updateContractBalances(tx, &contract)
		contract.LastPaymentDate = &payDate
		contract.TotalPaid = roundMoney(contract.TotalPaid + req.Amount)
		contract.UpdatedBy = userID
		if contract.OutstandingPrincipal <= 0 && contract.OutstandingInterest <= 0 &&
			contract.OutstandingFees <= 0 && contract.OutstandingPenalties <= 0 {
			contract.Status = models.LoanContractSettled
		}
		return tx.Save(&contract).Error
	})
	return repayment, err
}

func (s *LoanModuleService) allocatePayment(tx *gorm.DB, contract *models.LoanContract, repayment *models.LoanRepaymentRecord, userID uint64) ([]models.LoanRepaymentAllocation, error) {
	var installments []models.LoanScheduleInstallment
	if err := tx.Where("loan_contract_id = ? AND status IN ?", contract.ID, []string{"PENDING", "PARTIAL", "OVERDUE"}).
		Order("due_date ASC, installment_no ASC").Find(&installments).Error; err != nil {
		return nil, err
	}

	order := defaultAllocationOrder(contract.LoanProduct)
	remaining := repayment.Amount
	var allocs []models.LoanRepaymentAllocation

	type bucket struct {
		kind  string
		due   *float64
		paid  *float64
	}
	for remaining > 0.001 && len(installments) > 0 {
		inst := &installments[0]
		buckets := []bucket{
			{"PENALTY", &inst.PenaltyDue, &inst.PenaltyPaid},
			{"FEE", &inst.FeeDue, &inst.FeePaid},
			{"INTEREST", &inst.InterestDue, &inst.InterestPaid},
			{"PRINCIPAL", &inst.PrincipalDue, &inst.PrincipalPaid},
			{"INSURANCE", &inst.InsuranceDue, &inst.InsurancePaid},
		}
		// reorder by product priority
		ordered := make([]bucket, 0, len(buckets))
		for _, key := range order {
			for _, b := range buckets {
				if strings.EqualFold(b.kind, key) {
					ordered = append(ordered, b)
				}
			}
		}

		allocatedThisRound := false
		for _, b := range ordered {
			outstanding := roundMoney(*b.due - *b.paid)
			if outstanding <= 0 {
				continue
			}
			pay := outstanding
			if pay > remaining {
				pay = remaining
			}
			*b.paid = roundMoney(*b.paid + pay)
			remaining = roundMoney(remaining - pay)
			instID := inst.ID
			allocs = append(allocs, models.LoanRepaymentAllocation{
				BaseModel:                 models.BaseModel{CreatedBy: userID},
				LoanScheduleInstallmentID: &instID,
				AllocationType:            b.kind,
				Amount:                    pay,
			})
			allocatedThisRound = true
			if remaining <= 0.001 {
				break
			}
		}

		totalDue := roundMoney(inst.PrincipalDue + inst.InterestDue + inst.FeeDue + inst.InsuranceDue + inst.PenaltyDue)
		totalPaid := roundMoney(inst.PrincipalPaid + inst.InterestPaid + inst.FeePaid + inst.InsurancePaid + inst.PenaltyPaid)
		switch {
		case totalPaid >= totalDue-0.01:
			inst.Status = "PAID"
		case totalPaid > 0:
			inst.Status = "PARTIAL"
		}
		if err := tx.Save(inst).Error; err != nil {
			return nil, err
		}
		if !allocatedThisRound || inst.Status == "PAID" {
			installments = installments[1:]
		}
		if !allocatedThisRound {
			break
		}
	}
	return allocs, nil
}

func (s *LoanModuleService) syncRecurrentDeduction(tx *gorm.DB, contract *models.LoanContract) {
	if contract.RecurrentDeductionID == nil {
		return
	}
	var rd models.RecurrentDeduction
	if err := tx.First(&rd, *contract.RecurrentDeductionID).Error; err != nil {
		return
	}
	rd.PaidAmount = contract.TotalPaid
	rd.PrincipalAmount = contract.OutstandingPrincipal
	if contract.InstallmentAmount != nil {
		rd.RecurrentAmount = *contract.InstallmentAmount
	}
	if contract.OutstandingPrincipal <= 0 && contract.OutstandingInterest <= 0 {
		rd.Settled = 1
	}
	tx.Save(&rd)
}

func (s *LoanModuleService) updateContractBalances(tx *gorm.DB, contract *models.LoanContract) {
	var insts []models.LoanScheduleInstallment
	tx.Where("loan_contract_id = ?", contract.ID).Find(&insts)
	var op, oi, of, opn float64
	for _, inst := range insts {
		op += roundMoney(inst.PrincipalDue - inst.PrincipalPaid)
		oi += roundMoney(inst.InterestDue - inst.InterestPaid)
		of += roundMoney(inst.FeeDue - inst.FeePaid)
		opn += roundMoney(inst.PenaltyDue - inst.PenaltyPaid)
	}
	contract.OutstandingPrincipal = roundMoney(op)
	contract.OutstandingInterest = roundMoney(oi)
	contract.OutstandingFees = roundMoney(of)
	contract.OutstandingPenalties = roundMoney(opn)
}

// RecordRepaymentFromPayroll updates loan balances after milk payroll GL is already posted.
func (s *LoanModuleService) RecordRepaymentFromPayroll(contractID uint64, amount float64, payslipID, payrollID, userID uint64, reference string) error {
	_, err := s.RecordRepayment(contractID, dtos.RecordLoanRepaymentRequest{
		Amount:          amount,
		Channel:         models.LoanChannelMilkDeduction,
		Reference:       reference + "-SCHED",
		PostToGL:        false,
		MemberPayrollID: &payrollID,
		MemberPayslipID: &payslipID,
	}, userID)
	return err
}

// FindActiveContractByMemberReference resolves LOAN-{contractNo} references from recurrent deductions.
func (s *LoanModuleService) FindContractByDeductionReference(ref string) (*models.LoanContract, error) {
	ref = strings.TrimPrefix(ref, "LOAN-")
	var c models.LoanContract
	if err := db.DB.Where("contract_no = ? AND status IN ?", ref, []string{models.LoanContractActive, models.LoanContractOverdue}).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// --- Interest ---

func (s *LoanModuleService) AccrueInterest(req dtos.AccrueInterestRequest, userID uint64) (int, error) {
	asOf := utils.ParseFlexibleDate(req.AsOfDate)
	if asOf.IsZero() {
		asOf = time.Now()
	}
	q := db.DB.Where("status IN ?", []string{models.LoanContractActive, models.LoanContractOverdue})
	if len(req.LoanContractIDs) > 0 {
		q = q.Where("id IN ?", req.LoanContractIDs)
	}
	var contracts []models.LoanContract
	if err := q.Find(&contracts).Error; err != nil {
		return 0, err
	}
	count := 0
	for _, c := range contracts {
		var existing models.LoanInterestAccrual
		if err := db.DB.Where("loan_contract_id = ? AND accrual_date = ?", c.ID, asOf.Format("2006-01-02")).First(&existing).Error; err == nil {
			continue
		}
		dailyRate := c.InterestRate / 100 / 365
		interest := roundMoney(c.OutstandingPrincipal * dailyRate)
		if interest <= 0 {
			continue
		}
		accrual := &models.LoanInterestAccrual{
			BaseModel:      models.BaseModel{CreatedBy: userID},
			LoanContractID: c.ID,
			AccrualDate:    asOf,
			PrincipalBase:  c.OutstandingPrincipal,
			InterestAmount: interest,
		}
		if err := db.DB.Create(accrual).Error; err != nil {
			return count, err
		}
		if req.PostToGL {
			ref := fmt.Sprintf("LOAN-INT-%d-%s", c.ID, asOf.Format("20060102"))
			result, err := s.posting.postRule(DomainPostRequest{
				UserID: userID, Reference: ref, Amount: interest, TransactionDate: asOf,
				Description: fmt.Sprintf("Interest accrual %s", c.ContractNo),
				HeaderType: "LOAN", RuleType: "LOAN_INTEREST_ACCRUAL",
			})
			if err == nil {
				accrual.Posted = true
				accrual.GLTransactionID = &result.Transaction.ID
				db.DB.Save(accrual)
			}
		}
		count++
	}
	return count, nil
}

// --- Write-off & Restructure ---

func (s *LoanModuleService) WriteOffContract(id uint64, req dtos.WriteOffLoanRequest, userID uint64) (*models.LoanWriteOffRecord, error) {
	var record *models.LoanWriteOffRecord
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var contract models.LoanContract
		if err := tx.First(&contract, id).Error; err != nil {
			return err
		}
		woDate := utils.ParseFlexibleDate(req.WriteOffDate)
		if woDate.IsZero() {
			woDate = time.Now()
		}
		record = &models.LoanWriteOffRecord{
			BaseModel:      models.BaseModel{CreatedBy: userID},
			LoanContractID: contract.ID,
			Amount:         req.Amount,
			Reason:         req.Reason,
			WriteOffDate:   woDate,
			ApprovedBy:     &userID,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		if req.PostToGL {
			ref := fmt.Sprintf("LOAN-WO-%d-%d", contract.ID, time.Now().Unix())
			result, err := s.posting.postRule(DomainPostRequest{
				UserID: userID, Reference: ref, Amount: req.Amount, TransactionDate: woDate,
				Description: fmt.Sprintf("Write-off %s", contract.ContractNo),
				HeaderType: "LOAN", RuleType: "LOAN_WRITE_OFF",
			})
			if err != nil {
				return err
			}
			record.GLTransactionID = &result.Transaction.ID
			tx.Save(record)
		}
		contract.Status = models.LoanContractWrittenOff
		contract.OutstandingPrincipal = roundMoney(contract.OutstandingPrincipal - req.Amount)
		if contract.OutstandingPrincipal < 0 {
			contract.OutstandingPrincipal = 0
		}
		contract.UpdatedBy = userID
		return tx.Save(&contract).Error
	})
	return record, err
}

func (s *LoanModuleService) RestructureContract(id uint64, req dtos.RestructureLoanRequest, userID uint64) (*models.LoanRestructuring, error) {
	var record *models.LoanRestructuring
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var contract models.LoanContract
		if err := tx.Preload("LoanProduct").First(&contract, id).Error; err != nil {
			return err
		}
		effective := utils.ParseFlexibleDate(req.EffectiveDate)
		prevTerm := contract.TermMonths
		prevRate := contract.InterestRate
		record = &models.LoanRestructuring{
			BaseModel:          models.BaseModel{CreatedBy: userID},
			LoanContractID:     contract.ID,
			RestructureType:    req.RestructureType,
			PreviousTermMonths: &prevTerm,
			PreviousRate:       &prevRate,
			NewTermMonths:      req.NewTermMonths,
			NewRate:            req.NewRate,
			Reason:             req.Reason,
			EffectiveDate:      effective,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		if req.NewTermMonths != nil {
			contract.TermMonths = *req.NewTermMonths
		}
		if req.NewRate != nil {
			contract.InterestRate = *req.NewRate
		}
		contract.Status = models.LoanContractRestructured
		contract.UpdatedBy = userID
		if err := tx.Save(&contract).Error; err != nil {
			return err
		}
		// Regenerate remaining schedule from outstanding
		tx.Where("loan_contract_id = ? AND status IN ?", contract.ID, []string{"PENDING", "PARTIAL", "OVERDUE"}).Delete(&models.LoanScheduleInstallment{})
		lines := s.schedule.Generate(scheduleInput{
			Principal:  contract.OutstandingPrincipal,
			AnnualRate: contract.InterestRate,
			TermMonths: contract.TermMonths,
			Method:     contract.InterestMethod,
			StartDate:  effective,
		})
		for _, line := range lines {
			if err := tx.Create(&models.LoanScheduleInstallment{
				BaseModel: models.BaseModel{CreatedBy: userID}, LoanContractID: contract.ID,
				InstallmentNo: line.InstallmentNo, DueDate: line.DueDate,
				PrincipalDue: line.PrincipalDue, InterestDue: line.InterestDue, TotalDue: line.TotalDue,
				BalanceAfter: line.BalanceAfter, Status: "PENDING",
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return record, err
}

// --- Settlement ---

func (s *LoanModuleService) SettleContract(id uint64, req dtos.SettleLoanRequest, userID uint64) (*dtos.LoanSettlementQuote, error) {
	contract, err := s.GetContract(id)
	if err != nil {
		return nil, err
	}
	settleDate := utils.ParseFlexibleDate(req.SettlementDate)
	if settleDate.IsZero() {
		settleDate = time.Now()
	}
	quote := &dtos.LoanSettlementQuote{
		ContractID:           contract.ID,
		SettlementDate:       settleDate,
		OutstandingPrincipal: contract.OutstandingPrincipal,
		OutstandingInterest:  contract.OutstandingInterest,
		OutstandingFees:      contract.OutstandingFees,
		OutstandingPenalties: contract.OutstandingPenalties,
		TotalSettlement: roundMoney(contract.OutstandingPrincipal + contract.OutstandingInterest +
			contract.OutstandingFees + contract.OutstandingPenalties),
	}
	if req.QuoteOnly {
		return quote, nil
	}
	_, err = s.RecordRepayment(id, dtos.RecordLoanRepaymentRequest{
		Amount:      quote.TotalSettlement,
		Channel:     models.LoanChannelCash,
		Reference:   fmt.Sprintf("SETTLE-%s-%d", contract.ContractNo, time.Now().Unix()),
		PaymentDate: settleDate.Format("2006-01-02"),
		PostToGL:    req.PostToGL,
		Notes:       "Full loan settlement",
	}, userID)
	return quote, err
}

// --- Statement ---

func (s *LoanModuleService) GetStatement(contractID uint64) (*dtos.LoanModuleStatement, error) {
	contract, err := s.GetContract(contractID)
	if err != nil {
		return nil, err
	}
	stmt := &dtos.LoanModuleStatement{
		ContractID: contract.ID,
		ContractNo: contract.ContractNo,
		MemberID:   contract.MemberID,
	}
	var disbs []models.LoanDisbursement
	db.DB.Where("loan_contract_id = ?", contractID).Order("disbursement_date ASC").Find(&disbs)
	balance := 0.0
	for _, d := range disbs {
		balance += d.Amount
		stmt.Disbursements = append(stmt.Disbursements, dtos.LoanModuleStatementLine{
			Date: d.DisbursementDate, Reference: d.Reference, Description: "Disbursement",
			Debit: d.Amount, Balance: balance,
		})
	}
	var reps []models.LoanRepaymentRecord
	db.DB.Where("loan_contract_id = ?", contractID).Order("payment_date ASC").Find(&reps)
	for _, r := range reps {
		balance -= r.Amount
		stmt.Repayments = append(stmt.Repayments, dtos.LoanModuleStatementLine{
			Date: r.PaymentDate, Reference: r.Reference, Description: "Repayment " + r.Channel,
			Credit: r.Amount, Balance: balance,
		})
	}
	stmt.ClosingBalance = contract.OutstandingPrincipal + contract.OutstandingInterest +
		contract.OutstandingFees + contract.OutstandingPenalties
	for _, inst := range contract.Installments {
		stmt.Schedule = append(stmt.Schedule, dtos.LoanScheduleLine{
			InstallmentNo: inst.InstallmentNo, DueDate: inst.DueDate,
			PrincipalDue: inst.PrincipalDue, InterestDue: inst.InterestDue, TotalDue: inst.TotalDue,
			TotalPaid: roundMoney(inst.PrincipalPaid + inst.InterestPaid + inst.FeePaid + inst.InsurancePaid + inst.PenaltyPaid),
			Status: inst.Status,
		})
	}
	return stmt, nil
}

// --- Monitoring & Reports ---

func (s *LoanModuleService) PortfolioReport() (*dtos.LoanPortfolioReport, error) {
	r := &dtos.LoanPortfolioReport{}
	db.DB.Model(&models.LoanContract{}).Count(&r.TotalContracts)
	db.DB.Model(&models.LoanContract{}).Where("status = ?", models.LoanContractActive).Count(&r.ActiveContracts)
	db.DB.Model(&models.LoanContract{}).Where("status = ?", models.LoanContractOverdue).Count(&r.OverdueContracts)
	db.DB.Model(&models.LoanContract{}).Select("COALESCE(SUM(outstanding_principal + outstanding_interest + outstanding_fees + outstanding_penalties),0)").Scan(&r.TotalOutstanding)
	db.DB.Model(&models.LoanContract{}).Select("COALESCE(SUM(outstanding_principal),0)").Scan(&r.TotalPrincipal)
	db.DB.Model(&models.LoanContract{}).Select("COALESCE(SUM(outstanding_interest),0)").Scan(&r.TotalInterestDue)
	db.DB.Model(&models.LoanContract{}).Select("COALESCE(SUM(disbursed_amount),0)").Scan(&r.TotalDisbursed)
	db.DB.Model(&models.LoanContract{}).Select("COALESCE(SUM(total_paid),0)").Scan(&r.TotalRepaid)
	db.DB.Model(&models.LoanContract{}).Select("COALESCE(AVG(interest_rate),0)").Scan(&r.AverageInterestRate)
	return r, nil
}

func (s *LoanModuleService) AgingReport() (*dtos.LoanAgingReport, error) {
	asOf := time.Now()
	report := &dtos.LoanAgingReport{AsOfDate: asOf}
	buckets := []struct {
		name  string
		minD  int
		maxD  int
	}{
		{"current", 0, 0},
		{"1-30", 1, 30},
		{"31-60", 31, 60},
		{"61-90", 61, 90},
		{"90+", 91, 9999},
	}
	for _, b := range buckets {
		var count int64
		var outstanding float64
		q := db.DB.Model(&models.LoanScheduleInstallment{}).
			Joins("JOIN loan_contracts ON loan_contracts.id = loan_schedule_installments.loan_contract_id").
			Where("loan_schedule_installments.status IN ?", []string{"PENDING", "PARTIAL", "OVERDUE"})
		if b.name == "current" {
			q = q.Where("loan_schedule_installments.due_date >= ?", asOf.Format("2006-01-02"))
		} else {
			q = q.Where("DATEDIFF(?, loan_schedule_installments.due_date) BETWEEN ? AND ?", asOf.Format("2006-01-02"), b.minD, b.maxD)
		}
		q.Count(&count)
		q.Select("COALESCE(SUM(loan_schedule_installments.total_due - loan_schedule_installments.principal_paid - loan_schedule_installments.interest_paid - loan_schedule_installments.fee_paid - loan_schedule_installments.insurance_paid - loan_schedule_installments.penalty_paid),0)").Scan(&outstanding)
		report.Buckets = append(report.Buckets, dtos.LoanAgingBucket{Bucket: b.name, Count: count, Outstanding: outstanding})
	}
	return report, nil
}

func (s *LoanModuleService) PARReport() (*dtos.LoanPARReport, error) {
	portfolio, _ := s.PortfolioReport()
	r := &dtos.LoanPARReport{AsOfDate: time.Now(), TotalPortfolio: portfolio.TotalOutstanding}
	if portfolio.TotalOutstanding > 0 {
		var par30, par60, par90 float64
		asOf := time.Now().Format("2006-01-02")
		for _, item := range []struct {
			d   int
			ptr *float64
		}{{30, &par30}, {60, &par60}, {90, &par90}} {
			db.DB.Model(&models.LoanScheduleInstallment{}).
				Joins("JOIN loan_contracts ON loan_contracts.id = loan_schedule_installments.loan_contract_id").
				Where("loan_schedule_installments.status IN ? AND DATEDIFF(?, loan_schedule_installments.due_date) >= ?",
					[]string{"PENDING", "PARTIAL", "OVERDUE"}, asOf, item.d).
				Select("COALESCE(SUM(loan_schedule_installments.principal_due - loan_schedule_installments.principal_paid),0)").Scan(item.ptr)
		}
		r.PAR30 = roundMoney(par30 / portfolio.TotalOutstanding * 100)
		r.PAR60 = roundMoney(par60 / portfolio.TotalOutstanding * 100)
		r.PAR90 = roundMoney(par90 / portfolio.TotalOutstanding * 100)
		r.NPLRatio = r.PAR90
	}
	if portfolio.TotalDisbursed > 0 {
		r.RecoveryRate = roundMoney(portfolio.TotalRepaid / portfolio.TotalDisbursed * 100)
	}
	return r, nil
}

func (s *LoanModuleService) RegisterReport(status string, memberID uint64) (*dtos.LoanRegisterReport, error) {
	var contracts []models.LoanContract
	q := db.DB.Model(&models.LoanContract{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if memberID > 0 {
		q = q.Where("member_id = ?", memberID)
	}
	if err := q.Preload("LoanProduct").Order("id DESC").Find(&contracts).Error; err != nil {
		return nil, err
	}
	enrichContractsWithMembers(contracts)
	rows := make([]dtos.LoanRegisterRow, 0, len(contracts))
	for _, c := range contracts {
		productName := ""
		if c.LoanProduct != nil {
			productName = c.LoanProduct.ProductName
		}
		totalOutstanding := c.OutstandingPrincipal + c.OutstandingInterest + c.OutstandingFees + c.OutstandingPenalties
		rows = append(rows, dtos.LoanRegisterRow{
			ContractID:           c.ID,
			ContractNo:           c.ContractNo,
			MemberID:             c.MemberID,
			MemberNo:             c.MemberNo,
			MemberName:           c.MemberName,
			ProductName:          productName,
			ApprovedAmount:       c.ApprovedAmount,
			DisbursedAmount:      c.DisbursedAmount,
			OutstandingPrincipal: c.OutstandingPrincipal,
			OutstandingInterest:  c.OutstandingInterest,
			OutstandingFees:      c.OutstandingFees,
			OutstandingPenalties: c.OutstandingPenalties,
			TotalOutstanding:     roundMoney(totalOutstanding),
			Status:               c.Status,
			InterestRate:         c.InterestRate,
			TermMonths:           c.TermMonths,
			DaysInArrears:        c.DaysInArrears,
			DisbursementDate:     c.DisbursementDate,
		})
	}
	return &dtos.LoanRegisterReport{
		AsOfDate: time.Now(),
		Rows:     rows,
		Total:    int64(len(rows)),
	}, nil
}

func (s *LoanModuleService) ApplicationsPipelineReport() (*dtos.LoanApplicationsPipelineReport, error) {
	type row struct {
		Status         string
		Count          int64
		TotalRequested float64
	}
	var rows []row
	if err := db.DB.Model(&models.LoanApplication{}).
		Select("status, COUNT(*) AS count, COALESCE(SUM(requested_amount),0) AS total_requested").
		Group("status").
		Order("status ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]dtos.LoanApplicationStatusSummary, 0, len(rows))
	var totalApps int64
	var totalRequested float64
	for _, r := range rows {
		items = append(items, dtos.LoanApplicationStatusSummary{
			Status:         r.Status,
			Count:          r.Count,
			TotalRequested: r.TotalRequested,
		})
		totalApps += r.Count
		totalRequested += r.TotalRequested
	}
	return &dtos.LoanApplicationsPipelineReport{
		AsOfDate:          time.Now(),
		Items:             items,
		TotalApplications: totalApps,
		TotalRequested:    roundMoney(totalRequested),
	}, nil
}

func (s *LoanModuleService) DisbursementsReport(status, channelCode string) (*dtos.LoanDisbursementsReport, error) {
	var disbursements []models.LoanDisbursement
	q := db.DB.Model(&models.LoanDisbursement{}).Preload("DisbursementChannel")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if channelCode != "" {
		q = q.Where("channel_code = ?", channelCode)
	}
	if err := q.Order("disbursement_date DESC, id DESC").Find(&disbursements).Error; err != nil {
		return nil, err
	}
	contractIDs := make([]uint64, 0, len(disbursements))
	seen := make(map[uint64]struct{})
	for _, d := range disbursements {
		if _, ok := seen[d.LoanContractID]; !ok {
			seen[d.LoanContractID] = struct{}{}
			contractIDs = append(contractIDs, d.LoanContractID)
		}
	}
	contractMap := make(map[uint64]models.LoanContract)
	if len(contractIDs) > 0 {
		var contracts []models.LoanContract
		db.DB.Where("id IN ?", contractIDs).Find(&contracts)
		enrichContractsWithMembers(contracts)
		for _, c := range contracts {
			contractMap[c.ID] = c
		}
	}
	rows := make([]dtos.LoanDisbursementReportRow, 0, len(disbursements))
	var totalAmount float64
	for _, d := range disbursements {
		contract := contractMap[d.LoanContractID]
		rows = append(rows, dtos.LoanDisbursementReportRow{
			ID:                d.ID,
			ContractNo:        contract.ContractNo,
			MemberNo:          contract.MemberNo,
			MemberName:        contract.MemberName,
			ChannelCode:       d.ChannelCode,
			Status:            d.Status,
			Amount:            d.Amount,
			GrossAmount:       d.GrossAmount,
			FeesDeducted:      d.FeesDeducted,
			DisbursementDate:  d.DisbursementDate,
			ExternalReference: d.ExternalReference,
			Reference:         d.Reference,
		})
		totalAmount += d.Amount
	}
	return &dtos.LoanDisbursementsReport{
		AsOfDate:    time.Now(),
		Rows:        rows,
		Total:       int64(len(rows)),
		TotalAmount: roundMoney(totalAmount),
	}, nil
}

func (s *LoanModuleService) RunPenalties(asOf time.Time, userID uint64) (int, error) {
	if asOf.IsZero() {
		asOf = time.Now()
	}
	var overdue []models.LoanScheduleInstallment
	db.DB.Where("status IN ? AND due_date < ?", []string{"PENDING", "PARTIAL", "OVERDUE"}, asOf.Format("2006-01-02")).Find(&overdue)
	count := 0
	for _, inst := range overdue {
		var contract models.LoanContract
		if err := db.DB.Preload("LoanProduct").First(&contract, inst.LoanContractID).Error; err != nil {
			continue
		}
		product := contract.LoanProduct
		if product == nil {
			continue
		}
		daysLate := int(asOf.Sub(inst.DueDate).Hours() / 24)
		penalty := product.PenaltyFixedAmount
		if product.PenaltyRateDaily > 0 {
			outstanding := roundMoney(inst.TotalDue - inst.PrincipalPaid - inst.InterestPaid - inst.FeePaid - inst.InsurancePaid - inst.PenaltyPaid)
			penalty += roundMoney(outstanding * product.PenaltyRateDaily / 100 * float64(daysLate))
		}
		if penalty <= 0 {
			continue
		}
		inst.PenaltyDue = roundMoney(inst.PenaltyDue + penalty)
		inst.TotalDue = roundMoney(inst.TotalDue + penalty)
		inst.Status = "OVERDUE"
		db.DB.Save(&inst)
		db.DB.Create(&models.LoanPenaltyRecord{
			BaseModel: models.BaseModel{CreatedBy: userID},
			LoanContractID: inst.LoanContractID, LoanScheduleInstallmentID: &inst.ID,
			PenaltyType: "LATE_PAYMENT", Amount: penalty, PenaltyDate: asOf,
		})
		contract.OutstandingPenalties = roundMoney(contract.OutstandingPenalties + penalty)
		contract.Status = models.LoanContractOverdue
		contract.DaysInArrears = daysLate
		db.DB.Save(&contract)
		count++
	}
	return count, nil
}

func memberDisplayName(m models.Member) string {
	return strings.TrimSpace(m.FirstName + " " + m.LastName)
}

func loadMembersByIDs(ids []uint64) map[uint64]models.Member {
	out := make(map[uint64]models.Member)
	if len(ids) == 0 {
		return out
	}
	unique := make([]uint64, 0, len(ids))
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	if len(unique) == 0 {
		return out
	}
	var members []models.Member
	if err := db.DB.Where("id IN ?", unique).Find(&members).Error; err != nil {
		return out
	}
	for _, m := range members {
		out[m.ID] = m
	}
	return out
}

func enrichApplicationsWithMembers(items []models.LoanApplication) {
	if len(items) == 0 {
		return
	}
	ids := make([]uint64, len(items))
	for i := range items {
		ids[i] = items[i].MemberID
	}
	memberMap := loadMembersByIDs(ids)
	for i := range items {
		if m, ok := memberMap[items[i].MemberID]; ok {
			items[i].MemberNo = m.MemberNo
			items[i].MemberName = memberDisplayName(m)
		}
	}
}

func enrichContractsWithMembers(items []models.LoanContract) {
	if len(items) == 0 {
		return
	}
	ids := make([]uint64, len(items))
	for i := range items {
		ids[i] = items[i].MemberID
	}
	memberMap := loadMembersByIDs(ids)
	for i := range items {
		if m, ok := memberMap[items[i].MemberID]; ok {
			items[i].MemberNo = m.MemberNo
			items[i].MemberName = memberDisplayName(m)
		}
	}
}

// EnsureLoanSchema migrates loan module tables.
func EnsureLoanSchema() error {
	if err := db.DB.AutoMigrate(
		&models.LoanProduct{},
		&models.LoanApplication{},
		&models.LoanApproval{},
		&models.LoanGuarantor{},
		&models.LoanDocument{},
		&models.LoanContract{},
		&models.DisbursementChannel{},
		&models.DisbursementChannelConfig{},
		&models.LoanDisbursement{},
		&models.LoanScheduleInstallment{},
		&models.LoanRepaymentRecord{},
		&models.LoanRepaymentAllocation{},
		&models.LoanChargeRecord{},
		&models.LoanPenaltyRecord{},
		&models.LoanInterestAccrual{},
		&models.LoanRestructuring{},
		&models.LoanWriteOffRecord{},
		&models.LoanAuditLog{},
	); err != nil {
		return err
	}
	return EnsureDefaultDisbursementChannels()
}

