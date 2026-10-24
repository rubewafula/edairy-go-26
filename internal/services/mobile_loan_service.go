package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	mobileOTPTTL            = 5 * time.Minute
	mobileOTPResendCooldown = 60 * time.Second
	mobileOTPMaxAttempts    = 5
	mobileOTPMaxResends     = 3
	// mobileLoanMinRequestAmount is the smallest cashout a member can request (KES 100).
	// The product's Min Amount is the credit limit needed to borrow at all (see EvaluateBorrowStatus).
	mobileLoanMinRequestAmount = 100.0
	mobileLoanApprovalRole     = "MOBILE_OTP"
)

// ServiceError carries the HTTP status a controller should return.
type ServiceError struct {
	Status  int
	Message string
}

func (e *ServiceError) Error() string { return e.Message }

func serviceErr(status int, format string, args ...interface{}) error {
	return &ServiceError{Status: status, Message: fmt.Sprintf(format, args...)}
}

type MobileLoanService struct {
	loans  *LoanModuleService
	credit *MemberCreditService
}

func NewMobileLoanService() *MobileLoanService {
	return &MobileLoanService{loans: NewLoanModuleService(), credit: NewMemberCreditService()}
}

// memberLoanNames / memberLoanCodes identify the member loan product at setup when none is flagged yet.
var (
	memberLoanNames = []string{"member loan", "member loans"}
	memberLoanCodes = []string{"MEMBER-LOAN", "MEMBER_LOAN", "MEMBER-LOANS", "MEMBER_LOANS", "MEMBERLOAN", "MEMBERLOANS"}
)

// errNoMemberCashoutDefault is what members see when no default member cashout product is set.
var errNoMemberCashoutDefault = fmt.Errorf("Member cashouts are not available right now: the default member cashout product has not been set. Please contact the admin")

// resolveMemberLoanProduct returns the one loan product members can apply for from the mobile app: the
// active product flagged is_default_member_loan (set in Loan Products on the desktop). The default is
// chosen at setup (EnsureMemberLoanProduct); requests never pick one on their own.
func resolveMemberLoanProduct() (*models.LoanProduct, error) {
	var product models.LoanProduct
	if db.DB.Where("is_default_member_loan = ? AND is_active = ?", true, true).
		Order("id ASC").Limit(1).Find(&product).RowsAffected > 0 {
		return &product, nil
	}
	return nil, errNoMemberCashoutDefault
}

// isSellOffProduct matches "Sell Off", "Sell-Off", "selloff", "SELL_OFF" in a product's name or code.
func isSellOffProduct(p models.LoanProduct) bool {
	norm := func(v string) string {
		return strings.NewReplacer(" ", "", "-", "", "_", "").Replace(strings.ToLower(v))
	}
	return strings.Contains(norm(p.ProductName), "selloff") || strings.Contains(norm(p.ProductCode), "selloff")
}

// isMemberLoanProduct matches the product named "Member Loan(s)" or coded MEMBER-LOAN.
func isMemberLoanProduct(p models.LoanProduct) bool {
	name := strings.ToLower(strings.TrimSpace(p.ProductName))
	code := strings.ToUpper(strings.TrimSpace(p.ProductCode))
	for _, n := range memberLoanNames {
		if name == n {
			return true
		}
	}
	for _, c := range memberLoanCodes {
		if code == c {
			return true
		}
	}
	return false
}

// pickDefaultMemberLoanProduct chooses the setup default among active products: the sell-off loan first,
// else the member loan, else MOBILE_LOAN_PRODUCT_CODE (env). Ties go to the lowest id.
func pickDefaultMemberLoanProduct() *models.LoanProduct {
	var active []models.LoanProduct
	db.DB.Where("is_active = ?", true).Order("id ASC").Find(&active)
	for _, match := range []func(models.LoanProduct) bool{isSellOffProduct, isMemberLoanProduct} {
		for i := range active {
			if match(active[i]) {
				return &active[i]
			}
		}
	}
	if code := strings.TrimSpace(os.Getenv("MOBILE_LOAN_PRODUCT_CODE")); code != "" {
		for i := range active {
			if strings.EqualFold(strings.TrimSpace(active[i].ProductCode), code) {
				return &active[i]
			}
		}
	}
	return nil
}

// resolveMobileLoanProduct returns the member loan product. A mobile request may name a product,
// but only the member loan product is accepted; every other product is selected on the desktop.
func resolveMobileLoanProduct(productID uint64) (*models.LoanProduct, error) {
	product, err := resolveMemberLoanProduct()
	if err != nil {
		return nil, err
	}
	if productID > 0 && productID != product.ID {
		return nil, fmt.Errorf("only the %s product can be applied for from the mobile app", product.ProductName)
	}
	return product, nil
}

// EnsureMemberLoanProduct runs at setup (API start). When no active product is the default member loan,
// it flags the sell-off loan, else the member loan (see pickDefaultMemberLoanProduct). Without a default,
// member cashouts are refused with errNoMemberCashoutDefault until an admin sets one in Loan Products.
func EnsureMemberLoanProduct() {
	if product, err := resolveMemberLoanProduct(); err == nil {
		log.Printf("[MemberLoan] mobile loans use product %d %q (min KES %.2f, max KES %.2f)",
			product.ID, product.ProductName, product.MinAmount, product.MaxAmount)
		return
	}
	p := pickDefaultMemberLoanProduct()
	if p == nil {
		log.Printf("[MemberLoan] WARNING: no default member cashout product is set and no sell-off or member loan product was found; member cashouts are refused until one is set in Loan Products")
		return
	}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := clearDefaultMemberLoan(tx, p.ID); err != nil {
			return err
		}
		return tx.Model(&models.LoanProduct{}).Where("id = ?", p.ID).Update("is_default_member_loan", true).Error
	})
	if err != nil {
		log.Printf("[MemberLoan] WARNING: could not flag %q as the default member loan: %v", p.ProductName, err)
		return
	}
	log.Printf("[MemberLoan] flagged product %d %q as the default member loan; mobile loans use it (min KES %.2f, max KES %.2f)",
		p.ID, p.ProductName, p.MinAmount, p.MaxAmount)
}

func productTermMonths(product *models.LoanProduct) int {
	if product.RepaymentPeriodMonths > 0 {
		return product.RepaymentPeriodMonths
	}
	return 1
}

// projectedRepayable is everything the member will have deducted from milk for a new loan of this amount.
func (s *MobileLoanService) projectedRepayable(product *models.LoanProduct, memberID uint64, amount float64, term int) float64 {
	procFee, insFee := computeProductFees(product, amount)
	if !feesAddedToRepayment(product) {
		procFee, insFee = 0, 0
	}
	in := scheduleInput{
		Principal:       amount,
		AnnualRate:      product.InterestRate,
		TermMonths:      term,
		Method:          product.InterestMethod,
		GracePeriodDays: product.GracePeriodDays,
		StartDate:       time.Now(),
		ProcessingFee:   procFee,
		InsuranceFee:    insFee,
		FeesUpfront:     product.IsDefaultMemberLoan,
	}
	if interest, spread := spreadInterest(db.DB, product, memberID, amount); spread {
		in.FixedInterest, in.AnnualRate, in.Method = interest, 0, models.LoanInterestFlat
	}
	lines := s.loans.schedule.Generate(in)
	total := 0.0
	for _, line := range lines {
		total += line.TotalDue
	}
	if total < amount {
		total = amount
	}
	return roundMoney(total)
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func hashOTP(applicationID uint64, otp string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not configured")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "loan-otp:%d:%s", applicationID, otp)
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func otpMatches(applicationID uint64, otp, storedHash string) bool {
	candidate, err := hashOTP(applicationID, strings.TrimSpace(otp))
	if err != nil {
		return false
	}
	return hmac.Equal([]byte(candidate), []byte(storedHash))
}

func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:4] + strings.Repeat("*", len(phone)-7) + phone[len(phone)-3:]
}

func sendLoanOTP(memberID uint64, otp string, amount float64, applicationID, userID uint64) error {
	if os.Getenv("ENV") == "development" {
		log.Printf("[MobileLoan][development] OTP for loan application %d: %s", applicationID, otp)
		writeDevOTPLog(memberID, otp, amount, applicationID)
	}
	message := fmt.Sprintf("Your eDairy loan OTP is %s for KES %s. Valid %d minutes. Do not share this code.",
		otp, formatKES(amount), int(mobileOTPTTL.Minutes()))
	return notifyMemberSMS(memberID, message, userID)
}

// RequestLoan checks eligibility, creates a DRAFT mobile application and sends an OTP to the member's registered phone.
func (s *MobileLoanService) RequestLoan(req dtos.MobileLoanRequest, userID uint64) (*dtos.MobileLoanRequestResponse, error) {
	amount := roundMoney(req.Amount)
	if amount <= 0 {
		return nil, serviceErr(http.StatusBadRequest, "amount must be greater than zero")
	}

	var member models.Member
	if err := db.DB.Select("id, status, primary_phone").First(&member, req.MemberID).Error; err != nil {
		return nil, serviceErr(http.StatusNotFound, "member not found")
	}
	if !memberCanBorrow(member.Status) {
		return nil, serviceErr(http.StatusForbidden, "member is not allowed to borrow (status %s)", member.Status)
	}
	phone := utils.NormalizePhone(member.PrimaryPhone)
	if phone == "" {
		return nil, serviceErr(http.StatusUnprocessableEntity, "member has no valid registered phone number")
	}

	product, err := resolveMobileLoanProduct(req.LoanProductID)
	if err != nil {
		return nil, serviceErr(http.StatusUnprocessableEntity, "%s", err.Error())
	}
	if amount < mobileLoanMinRequestAmount {
		return nil, serviceErr(http.StatusBadRequest, "minimum cashout is KES %s", formatKES(mobileLoanMinRequestAmount))
	}
	term := productTermMonths(product)
	repayable := s.projectedRepayable(product, member.ID, amount, term)

	if err := s.expireStaleRequests(member.ID); err != nil {
		return nil, err
	}
	if err := s.ensureNoOpenRequest(member.ID); err != nil {
		return nil, err
	}

	// Advances need an advance rate for the member (Default Milk Rates -> Advance Rate).
	if MemberAdvanceRate(db.DB, member.ID) <= 0 {
		return nil, serviceErr(http.StatusUnprocessableEntity,
			"advances are not available yet: no advance rate is set for this member's route")
	}

	// The member loan is borrowed strictly against the member's existing credit limit
	// (member_registrations.credit_limit, the value GET /member-credit-limit returns).
	storedLimit, err := storedCreditLimit(db.DB, member.ID)
	if err != nil {
		return nil, err
	}
	if need := s.creditNeededFor(product, member.ID, amount, term); need > storedLimit+0.005 {
		return nil, serviceErr(http.StatusUnprocessableEntity,
			"maximum cashout is KES %s (credit limit KES %s after charges)",
			formatKES(s.maxRequestFor(product, member.ID, math.Max(0, storedLimit), term)), formatKES(math.Max(0, storedLimit)))
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	app := &models.LoanApplication{
		BaseModel:           models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		ApplicationNo:       fmt.Sprintf("LN-APP-%d", now.UnixNano()),
		MemberID:            member.ID,
		LoanProductID:       product.ID,
		RequestedAmount:     amount,
		RequestedTermMonths: term,
		Purpose:             strings.TrimSpace(req.Purpose),
		Status:              models.LoanAppDraft,
		DateApplied:         now,
		Source:              models.LoanAppSourceMobile,
	}
	otpRecord := &models.LoanApplicationOTP{
		BaseModel:   models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		MemberID:    member.ID,
		Phone:       phone,
		ExpiresAt:   now.Add(mobileOTPTTL),
		MaxAttempts: mobileOTPMaxAttempts,
		LastSentAt:  now,
		Status:      models.LoanOTPStatusPending,
	}
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(app).Error; err != nil {
			return err
		}
		hash, err := hashOTP(app.ID, otp)
		if err != nil {
			return err
		}
		otpRecord.LoanApplicationID = app.ID
		otpRecord.OTPHash = hash
		if err := tx.Create(otpRecord).Error; err != nil {
			return err
		}
		s.loans.logAudit(tx, "loan_application", app.ID, "MOBILE_REQUEST", userID, nil, app)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := sendLoanOTP(member.ID, otp, amount, app.ID, userID); err != nil {
		s.cancelRequest(app.ID, otpRecord.ID, models.LoanOTPStatusExpired, "OTP could not be sent: "+err.Error(), userID)
		return nil, serviceErr(http.StatusBadGateway, "could not send OTP to the member's phone: %s", err.Error())
	}
	s.credit.RefreshBorrowStatusAsync(member.ID, "loan_application", app.ID, userID)

	return &dtos.MobileLoanRequestResponse{
		ApplicationID:     app.ID,
		ApplicationNo:     app.ApplicationNo,
		MemberID:          member.ID,
		Amount:            amount,
		TotalRepayable:    repayable,
		TermMonths:        term,
		LoanProductID:     product.ID,
		MaskedPhone:       maskPhone(phone),
		OTPExpiresAt:      otpRecord.ExpiresAt,
		ResendAvailableAt: now.Add(mobileOTPResendCooldown),
		CreditLimit:       storedLimit,
	}, nil
}

// expireStaleRequests cancels mobile drafts whose OTP has expired so they no longer block a new request.
func (s *MobileLoanService) expireStaleRequests(memberID uint64) error {
	var stale []models.LoanApplicationOTP
	if err := db.DB.Where("member_id = ? AND status = ? AND expires_at < ?", memberID, models.LoanOTPStatusPending, time.Now()).
		Find(&stale).Error; err != nil {
		return err
	}
	for _, o := range stale {
		s.cancelRequest(o.LoanApplicationID, o.ID, models.LoanOTPStatusExpired, "OTP expired", 0)
	}
	return nil
}

// ensureNoOpenRequest uses the same rule as can_borrow: the member's last application must be processed.
func (s *MobileLoanService) ensureNoOpenRequest(memberID uint64) error {
	open := lastPendingLoanApplication(db.DB, memberID)
	if open == nil {
		return nil
	}
	return serviceErr(http.StatusConflict, "%s", pendingApplicationReason(open))
}

func (s *MobileLoanService) cancelRequest(applicationID, otpID uint64, otpStatus, reason string, userID uint64) {
	db.DB.Model(&models.LoanApplicationOTP{}).Where("id = ? AND status = ?", otpID, models.LoanOTPStatusPending).
		Updates(map[string]interface{}{"status": otpStatus, "updated_by": userID})
	db.DB.Model(&models.LoanApplication{}).Where("id = ? AND status = ?", applicationID, models.LoanAppDraft).
		Updates(map[string]interface{}{"status": models.LoanAppCancelled, "rejected_reason": reason, "updated_by": userID})
	var app models.LoanApplication
	if db.DB.Select("id, member_id").First(&app, applicationID).Error == nil {
		s.credit.RefreshBorrowStatusAsync(app.MemberID, "loan_application", applicationID, userID)
	}
}

func (s *MobileLoanService) loadMobileDraft(tx *gorm.DB, applicationID uint64, lock bool) (*models.LoanApplication, *models.LoanApplicationOTP, error) {
	q := tx
	if lock {
		q = tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var app models.LoanApplication
	if err := q.First(&app, applicationID).Error; err != nil {
		return nil, nil, serviceErr(http.StatusNotFound, "loan request not found")
	}
	if app.Source != models.LoanAppSourceMobile {
		return nil, nil, serviceErr(http.StatusNotFound, "loan request not found")
	}
	var otp models.LoanApplicationOTP
	q2 := tx
	if lock {
		q2 = tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q2.Where("loan_application_id = ?", app.ID).Order("id DESC").First(&otp).Error; err != nil {
		return &app, nil, serviceErr(http.StatusNotFound, "no OTP found for this loan request")
	}
	return &app, &otp, nil
}

// ResendOTP issues a fresh code within the resend limit and cooldown.
func (s *MobileLoanService) ResendOTP(applicationID, userID uint64) (*dtos.MobileLoanRequestResponse, error) {
	var app *models.LoanApplication
	var record *models.LoanApplicationOTP
	var otp string
	var resendRejected error
	now := time.Now()
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		app, record, err = s.loadMobileDraft(tx, applicationID, true)
		if err != nil {
			return err
		}
		// A request still awaiting OTP confirmation can be reinitiated, including after its code expired.
		if app.Status != models.LoanAppDraft ||
			(record.Status != models.LoanOTPStatusPending && record.Status != models.LoanOTPStatusExpired) {
			return serviceErr(http.StatusConflict, "this loan request is no longer awaiting OTP (status %s)", app.Status)
		}
		if record.ResendCount >= mobileOTPMaxResends {
			return serviceErr(http.StatusTooManyRequests, "OTP resend limit reached; start a new loan request")
		}
		// Re-validate before sending a new code, so the member is not asked to confirm a loan they no longer qualify for.
		if tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").First(&models.Member{}, app.MemberID).Error == nil {
			if reason := s.revalidateLoanRequest(tx, app); reason != "" {
				app.Status = models.LoanAppRejected
				app.RejectedReason = reason
				record.Status = models.LoanOTPStatusExpired
				if err := tx.Save(app).Error; err != nil {
					return err
				}
				if err := tx.Save(record).Error; err != nil {
					return err
				}
				resendRejected = serviceErr(http.StatusUnprocessableEntity, "not eligible: %s", reason)
				return nil
			}
		}
		if wait := record.LastSentAt.Add(mobileOTPResendCooldown).Sub(now); wait > 0 {
			return serviceErr(http.StatusTooManyRequests, "wait %d seconds before requesting another OTP", int(math.Ceil(wait.Seconds())))
		}
		if otp, err = generateOTP(); err != nil {
			return err
		}
		hash, err := hashOTP(app.ID, otp)
		if err != nil {
			return err
		}
		record.OTPHash = hash
		record.Status = models.LoanOTPStatusPending
		record.ExpiresAt = now.Add(mobileOTPTTL)
		record.Attempts = 0
		record.ResendCount++
		record.LastSentAt = now
		record.UpdatedBy = userID
		return tx.Save(record).Error
	})
	if err != nil {
		return nil, err
	}
	if resendRejected != nil {
		s.credit.RefreshBorrowStatusAsync(app.MemberID, "loan_application", app.ID, userID)
		return nil, resendRejected
	}
	if err := sendLoanOTP(app.MemberID, otp, app.RequestedAmount, app.ID, userID); err != nil {
		return nil, serviceErr(http.StatusBadGateway, "could not send OTP to the member's phone: %s", err.Error())
	}
	return &dtos.MobileLoanRequestResponse{
		ApplicationID:     app.ID,
		ApplicationNo:     app.ApplicationNo,
		MemberID:          app.MemberID,
		Amount:            app.RequestedAmount,
		TermMonths:        app.RequestedTermMonths,
		LoanProductID:     app.LoanProductID,
		MaskedPhone:       maskPhone(record.Phone),
		OTPExpiresAt:      record.ExpiresAt,
		ResendAvailableAt: now.Add(mobileOTPResendCooldown),
	}, nil
}

// VerifyOTP confirms the code, re-checks the milk balance under the member lock, approves the loan and sends it to M-Pesa.
func (s *MobileLoanService) VerifyOTP(applicationID uint64, code string, userID uint64) (*dtos.MobileLoanStatusResponse, error) {
	var verifyErr error
	var contractID uint64
	var app *models.LoanApplication

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var record *models.LoanApplicationOTP
		var err error
		app, record, err = s.loadMobileDraft(tx, applicationID, true)
		if err != nil {
			return err
		}
		if app.Status != models.LoanAppDraft {
			return nil
		}
		now := time.Now()
		if record.Status != models.LoanOTPStatusPending {
			verifyErr = serviceErr(http.StatusConflict, "OTP is %s; start a new loan request", strings.ToLower(record.Status))
			return nil
		}
		if now.After(record.ExpiresAt) {
			record.Status = models.LoanOTPStatusExpired
			verifyErr = serviceErr(http.StatusGone, "OTP has expired; request a new one")
			return tx.Save(record).Error
		}
		if !otpMatches(app.ID, code, record.OTPHash) {
			record.Attempts++
			if record.Attempts >= record.MaxAttempts {
				record.Status = models.LoanOTPStatusLocked
				app.Status = models.LoanAppCancelled
				app.RejectedReason = "Too many invalid OTP attempts"
				if err := tx.Save(app).Error; err != nil {
					return err
				}
				verifyErr = serviceErr(http.StatusForbidden, "too many invalid OTP attempts; the loan request has been cancelled")
			} else {
				verifyErr = serviceErr(http.StatusUnauthorized, "invalid OTP; %d attempt(s) left", record.MaxAttempts-record.Attempts)
			}
			return tx.Save(record).Error
		}

		record.Status = models.LoanOTPStatusVerified
		record.VerifiedAt = &now
		record.UpdatedBy = userID
		if err := tx.Save(record).Error; err != nil {
			return err
		}

		// Lock the member before any non-locking read: under REPEATABLE READ the snapshot is taken at the first
		// plain read, so it must come after a concurrent loan for this member has committed.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").First(&models.Member{}, app.MemberID).Error; err != nil {
			return err
		}

		// Re-validate now (the OTP may be confirmed long after the request): member status, any other
		// unprocessed loan, credit limit vs minimum borrowable amount, and amount vs limit.
		if reason := s.revalidateLoanRequest(tx, app); reason != "" {
			app.Status = models.LoanAppRejected
			app.RejectedReason = reason
			verifyErr = serviceErr(http.StatusUnprocessableEntity, "not eligible: %s", reason)
			return tx.Save(app).Error
		}

		amount := app.RequestedAmount
		term := app.RequestedTermMonths
		if err := tx.Create(&models.LoanApproval{
			BaseModel:          models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			LoanApplicationID:  app.ID,
			ApprovalRole:       mobileLoanApprovalRole,
			ApproverID:         userID,
			ApprovalDate:       now,
			Decision:           "APPROVED",
			ApprovedAmount:     &amount,
			ApprovedTermMonths: &term,
			Comments:           "Approved automatically after member OTP verification",
		}).Error; err != nil {
			return err
		}
		app.Status = models.LoanAppApproved
		app.ApprovedAmount = &amount
		app.ApprovedTermMonths = &term
		app.UpdatedBy = userID
		// Record the credit baseline before the contract exists, so the new loan is deducted exactly once.
		if _, err := s.credit.Recalculate(tx, app.MemberID, CreditChange{
			Action:     models.CreditActionRecalculated,
			SourceType: "mobile_loan_verify",
			SourceID:   app.ID,
			Reference:  app.ApplicationNo,
			UserID:     userID,
		}); err != nil {
			return err
		}
		contract, err := s.loans.createContractFromApplication(tx, app, amount, term, userID)
		if err != nil {
			return err
		}
		app.LoanContractID = &contract.ID
		contractID = contract.ID
		if err := tx.Save(app).Error; err != nil {
			return err
		}
		// The new contract's repayable amount comes off the stored credit limit (Recalculate applies the change).
		if _, err := s.credit.Recalculate(tx, app.MemberID, CreditChange{
			Action:     models.CreditActionLoanReserved,
			SourceType: "loan_contract",
			SourceID:   contract.ID,
			Reference:  contract.ContractNo,
			Narration:  "Mobile loan approved after OTP verification",
			UserID:     userID,
		}); err != nil {
			return err
		}
		s.loans.logAudit(tx, "loan_application", app.ID, "MOBILE_OTP_APPROVED", userID, nil, app)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if verifyErr != nil {
		if app != nil {
			s.credit.RefreshBorrowStatusAsync(app.MemberID, "loan_application", app.ID, userID)
		}
		return nil, verifyErr
	}
	if contractID == 0 {
		return s.GetRequest(applicationID)
	}

	var member models.Member
	if err := db.DB.Select("id, primary_phone").First(&member, app.MemberID).Error; err != nil {
		return nil, err
	}
	// Queue the M-Pesa payout for the disbursement worker. The loan is already approved, deducted from
	// the credit limit, and its repayment deduction is added now; the worker sends the money.
	_, queueErr := s.loans.DisburseContract(contractID, dtos.DisburseLoanRequest{
		ChannelCode:         models.DisbursementCodeMobileMpesa,
		ChannelPayload:      map[string]interface{}{"phone_number": member.PrimaryPhone},
		Reference:           "MLOAN-" + app.ApplicationNo,
		IdempotencyKey:      fmt.Sprintf("MLOAN-%d", app.ID),
		PostToGL:            true,
		CreateMilkDeduction: true,
		Notes:               "Mobile loan after OTP verification (queued for M-Pesa)",
		QueueOnly:           true,
	}, userID)
	if queueErr != nil {
		if err := s.loans.releaseMobileContract(contractID, queueErr.Error(), userID); err != nil {
			log.Printf("[MobileLoan] release contract %d after failed queueing: %v", contractID, err)
		}
		status, _ := s.GetRequest(applicationID)
		if status != nil {
			status.FailureMessage = queueErr.Error()
		}
		return status, serviceErr(http.StatusInternalServerError, "could not queue the loan for disbursement: %s", queueErr.Error())
	}

	received := fmt.Sprintf("Dear %s, your loan request of KES %s has been received and is being processed. You will receive the money on %s shortly. Amount to be deducted from your milk: KES %s. Your milk balance after this deduction: KES %s.",
		memberFirstName(app.MemberID), formatKES(app.RequestedAmount), maskPhone(utils.NormalizePhone(member.PrimaryPhone)),
		formatKES(contractRemainingRepayable(db.DB, contractID)), formatKES(memberActualBalance(app.MemberID)))
	if err := notifyMemberSMS(app.MemberID, received, userID); err != nil {
		log.Printf("[MobileLoan] 'request received' SMS for member %d: %v", app.MemberID, err)
	}
	return s.GetRequest(applicationID)
}

// GetRequest reports the current state of a mobile loan request for the app to poll.
func (s *MobileLoanService) GetRequest(applicationID uint64) (*dtos.MobileLoanStatusResponse, error) {
	app, record, err := s.loadMobileDraft(db.DB, applicationID, false)
	if err != nil && app == nil {
		return nil, err
	}
	resp := &dtos.MobileLoanStatusResponse{
		ApplicationID:     app.ID,
		ApplicationNo:     app.ApplicationNo,
		MemberID:          app.MemberID,
		Amount:            app.RequestedAmount,
		ApplicationStatus: app.Status,
		ContractID:        app.LoanContractID,
	}
	if app.Status == models.LoanAppCancelled || app.Status == models.LoanAppRejected {
		resp.FailureMessage = app.RejectedReason
	}
	if record != nil {
		resp.OTPStatus = record.Status
		if record.Status == models.LoanOTPStatusPending {
			expires := record.ExpiresAt
			resp.OTPExpiresAt = &expires
			resp.OTPAttemptsLeft = record.MaxAttempts - record.Attempts
		}
	}
	if app.LoanContractID != nil {
		var contract models.LoanContract
		if err := db.DB.Select("id, contract_no, status").First(&contract, *app.LoanContractID).Error; err == nil {
			resp.ContractNo = contract.ContractNo
			resp.ContractStatus = contract.Status
			resp.TotalRepayable = contractRemainingRepayable(db.DB, contract.ID)
		}
		var d models.LoanDisbursement
		if err := db.DB.Where("loan_contract_id = ?", *app.LoanContractID).Order("id DESC").First(&d).Error; err == nil {
			id := d.ID
			resp.DisbursementID = &id
			resp.DisbursementStatus = d.Status
			if d.Status == models.DisbursementStatusMpesaFailed || d.Status == models.DisbursementStatusMpesaError {
				// Members are not shown M-Pesa failures; the admin retries the payout.
				resp.DisbursementStatus = models.DisbursementStatusProcessing
			}
			resp.DisbursementReference = d.Reference
			if d.Status == models.DisbursementStatusSuccess {
				resp.MpesaReference = d.ExternalReference
			}
			if d.Status == models.DisbursementStatusFailed && resp.FailureMessage == "" {
				resp.FailureMessage = d.FailureMessage
			}
		}
	}
	var member models.Member
	if err := db.DB.Select("id, credit_limit").First(&member, app.MemberID).Error; err == nil {
		resp.CreditLimit = member.CreditLimit
	}
	return resp, nil
}

// ListRequests returns the member's mobile loan requests, newest first, with a display status.
func (s *MobileLoanService) ListRequests(memberID uint64, limit int) ([]dtos.MobileLoanListItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var apps []models.LoanApplication
	if err := db.DB.Select("id, status, rejected_reason, date_applied, created_at").
		Where("member_id = ? AND source = ?", memberID, models.LoanAppSourceMobile).
		Order("id DESC").Limit(limit).Find(&apps).Error; err != nil {
		return nil, err
	}
	items := make([]dtos.MobileLoanListItem, 0, len(apps))
	for _, a := range apps {
		st, err := s.GetRequest(a.ID)
		if err != nil || st == nil {
			continue
		}
		item := dtos.MobileLoanListItem{MobileLoanStatusResponse: *st, CreatedAt: a.CreatedAt}
		if item.CreatedAt.IsZero() {
			item.CreatedAt = a.DateApplied
		}
		item.DisplayStatus, item.StatusMessage = mobileLoanDisplayStatus(st)
		items = append(items, item)
	}
	return items, nil
}

func mobileLoanDisplayStatus(st *dtos.MobileLoanStatusResponse) (string, string) {
	switch st.ApplicationStatus {
	case models.LoanAppDraft:
		if st.OTPStatus == models.LoanOTPStatusPending && st.OTPExpiresAt != nil && time.Now().Before(*st.OTPExpiresAt) {
			return "Awaiting OTP", "Waiting for the member to confirm the OTP sent to their phone."
		}
		return "Expired", "The OTP was not confirmed in time. Request a new cashout."
	case models.LoanAppSubmitted, models.LoanAppUnderReview:
		return "Under review", "The request is being reviewed."
	case models.LoanAppRejected:
		return "Rejected", st.FailureMessage
	case models.LoanAppCancelled:
		msg := st.FailureMessage
		if msg == "" {
			msg = "The request was cancelled."
		}
		return "Cancelled", msg
	case models.LoanAppDisbursed:
		return "Disbursed", disbursedMessage(st)
	case models.LoanAppApproved:
		switch st.DisbursementStatus {
		case models.DisbursementStatusSuccess:
			return "Disbursed", disbursedMessage(st)
		case models.DisbursementStatusProcessing:
			return "Sending", "Sent to M-Pesa; waiting for Safaricom to confirm."
		case models.DisbursementStatusPending:
			return "Queued", "Approved and queued; the money will be sent via M-Pesa shortly."
		case models.DisbursementStatusFailed:
			return "Failed", st.FailureMessage
		}
		return "Approved", "Approved; waiting to be sent."
	}
	return st.ApplicationStatus, st.FailureMessage
}

func disbursedMessage(st *dtos.MobileLoanStatusResponse) string {
	if st.MpesaReference != "" {
		return "Money sent via M-Pesa (ref " + st.MpesaReference + ")."
	}
	return "Money sent."
}

// revalidateLoanRequest re-checks a request still awaiting OTP against the member's CURRENT stored
// credit limit and borrow status. Returns "" when it may proceed, else the reason it may not.
// The caller must hold the member row lock.
func (s *MobileLoanService) revalidateLoanRequest(tx *gorm.DB, app *models.LoanApplication) string {
	// The default member loan may have changed since the request was made: move the request onto the
	// current default so it is approved on today's terms (term, fees, interest, limits).
	if cur, err := resolveMobileLoanProduct(0); err == nil {
		if cur.ID != app.LoanProductID {
			term := productTermMonths(cur)
			if err := tx.Model(&models.LoanApplication{}).Where("id = ?", app.ID).Updates(map[string]interface{}{
				"loan_product_id": cur.ID, "requested_term_months": term,
			}).Error; err != nil {
				return "could not move the request to the current loan product"
			}
			log.Printf("[MobileLoan] application %d moved from product %d to the default member loan %d (%s)",
				app.ID, app.LoanProductID, cur.ID, cur.ProductName)
			app.LoanProductID, app.RequestedTermMonths = cur.ID, term
		}
		if cur.MinAmount > 0 && app.RequestedAmount < cur.MinAmount-0.005 {
			return fmt.Sprintf("The minimum advance is KES %s", formatKES(cur.MinAmount))
		}
		if cur.MaxAmount > 0 && app.RequestedAmount > cur.MaxAmount+0.005 {
			return fmt.Sprintf("The maximum advance is KES %s", formatKES(cur.MaxAmount))
		}
	} else {
		return err.Error()
	}
	var member models.Member
	if err := tx.Select("id, status, credit_limit").First(&member, app.MemberID).Error; err != nil {
		return "member not found"
	}
	borrow := s.credit.evaluateBorrowStatus(tx, member.ID, member.Status, member.CreditLimit, app.ID)
	if !borrow.CanBorrow {
		if borrow.Reason != nil {
			return *borrow.Reason
		}
		return "member cannot borrow right now"
	}
	need := app.RequestedAmount
	var product models.LoanProduct
	if err := tx.First(&product, app.LoanProductID).Error; err == nil {
		term := app.RequestedTermMonths
		if term <= 0 {
			term = productTermMonths(&product)
		}
		need = s.creditNeededFor(&product, member.ID, app.RequestedAmount, term)
	}
	if need > member.CreditLimit+0.005 {
		return fmt.Sprintf("Credit limit KES %s is below the KES %s this advance of KES %s needs",
			formatKES(math.Max(0, member.CreditLimit)), formatKES(need), formatKES(app.RequestedAmount))
	}
	return ""
}

// storedCreditLimit reads member_registrations.credit_limit without recalculating it.
func storedCreditLimit(tx *gorm.DB, memberID uint64) (float64, error) {
	var member models.Member
	if err := tx.Select("id, credit_limit").First(&member, memberID).Error; err != nil {
		return 0, fmt.Errorf("member %d not found", memberID)
	}
	return member.CreditLimit, nil
}

// writeDevOTPLog appends each OTP to a dedicated file in development so it can be read without the
// server console: DEV_OTP_LOG_FILE if set, else logs/dev-otp.log under the working directory.
// Never enabled outside ENV=development.
func writeDevOTPLog(memberID uint64, otp string, amount float64, applicationID uint64) {
	path := strings.TrimSpace(os.Getenv("DEV_OTP_LOG_FILE"))
	if path == "" {
		path = filepath.Join("logs", "dev-otp.log")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Printf("[MobileLoan][development] cannot create OTP log dir: %v", err)
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("[MobileLoan][development] cannot open OTP log: %v", err)
		return
	}
	defer f.Close()
	var member models.Member
	db.DB.Select("id, member_no, primary_phone").First(&member, memberID)
	fmt.Fprintf(f, "%s | application %d | member %d (%s) | phone %s | KES %s | OTP %s\n",
		time.Now().Format("2006-01-02 15:04:05"), applicationID, memberID, member.MemberNo,
		utils.NormalizePhone(member.PrimaryPhone), formatKES(amount), otp)
}

// CashoutQuote is what a member pays back for an advance: the "advance price".
type CashoutQuote struct {
	Amount         float64 `json:"amount"`
	ProcessingFee  float64 `json:"processing_fee"`  // from the product's fee bands
	YouReceive     float64 `json:"you_receive"`     // sent to M-Pesa (amount - fees taken upfront)
	Repayable      float64 `json:"repayable"`       // total deducted from milk = advance price
	Cost           float64 `json:"cost"`            // repayable - you_receive (spread interest + fees)
	CostPercent    float64 `json:"cost_percent"`    // cost as % of the amount
	CreditUsed     float64 `json:"credit_used"`     // how much of the credit limit this takes
	TermMonths     int     `json:"term_months"`
	InterestRate   float64 `json:"interest_rate"`
	AdvanceRate    float64 `json:"advance_rate"`    // credit rate, KES per litre
	CreditRate     float64 `json:"credit_rate"`
	DefaultRate    float64 `json:"default_rate"`
	ChargePerLitre float64 `json:"charge_per_litre"`
}

// creditNeededFor is how much credit an advance of amount uses: its Total Deductible (repayable),
// the same figure that comes off the member balance and is recovered through milk deductions.
func (s *MobileLoanService) creditNeededFor(product *models.LoanProduct, memberID uint64, amount float64, term int) float64 {
	return s.projectedRepayable(product, memberID, amount, term)
}

// maxRequestFor is the largest advance whose credit use fits in limit (and the product maximum).
func (s *MobileLoanService) maxRequestFor(product *models.LoanProduct, memberID uint64, limit float64, term int) float64 {
	hi := math.Floor(limit*100) / 100
	if product.MaxAmount > 0 && hi > product.MaxAmount {
		hi = product.MaxAmount
	}
	if hi <= 0 {
		return 0
	}
	if s.creditNeededFor(product, memberID, hi, term) <= limit+0.005 {
		return hi
	}
	lo := 0.0
	for i := 0; i < 40 && hi-lo > 0.01; i++ {
		mid := math.Floor((lo+hi)/2*100) / 100
		if s.creditNeededFor(product, memberID, mid, term) <= limit+0.005 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return roundMoney(lo)
}

// upfrontFees are the fees taken from the payout (fee collection "deduct from proceeds").
func upfrontFees(product *models.LoanProduct, amount float64) (processing, total float64) {
	proc, ins := computeProductFees(product, amount)
	if feesOffPayout(product) {
		return proc, roundMoney(proc + ins)
	}
	return proc, 0
}

// QuoteCashout prices an advance of `amount` on the member loan product.
func (s *MobileLoanService) QuoteCashout(memberID uint64, amount float64) (*CashoutQuote, error) {
	amount = roundMoney(amount)
	if amount <= 0 {
		return nil, serviceErr(http.StatusBadRequest, "amount must be greater than zero")
	}
	product, err := resolveMobileLoanProduct(0)
	if err != nil {
		return nil, serviceErr(http.StatusUnprocessableEntity, "%s", err.Error())
	}
	term := productTermMonths(product)
	repayable := s.projectedRepayable(product, memberID, amount, term)
	procFee, upfront := upfrontFees(product, amount)
	receive := roundMoney(amount - upfront)
	q := &CashoutQuote{Amount: amount, ProcessingFee: procFee, YouReceive: receive, Repayable: repayable,
		Cost: roundMoney(repayable - receive), TermMonths: term, InterestRate: product.InterestRate,
		ChargePerLitre: product.ChargePerLitre}
	if amount > 0 {
		q.CostPercent = roundMoney((repayable - receive) / amount * 100)
	}
	if memberID > 0 {
		mr := GetMemberRates(db.DB, memberID)
		q.CreditRate, q.DefaultRate, q.AdvanceRate = mr.CreditRate, mr.DefaultRate, mr.CreditRate
		q.CreditUsed = s.creditNeededFor(product, memberID, amount, term)
	}
	return q, nil
}
