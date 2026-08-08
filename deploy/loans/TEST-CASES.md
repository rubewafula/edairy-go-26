# Loan Module Test Cases

## Unit Tests

### LoanScheduleService
- FLAT: 100k @ 12% flat, 12 months → total interest = 12k, 12 equal installments
- REDUCING_BALANCE: verify sum(principal) = original principal
- EQUAL_INSTALLMENTS: PMT matches Excel PMT()
- INTEREST_ONLY: all principal in final installment
- BALLOON: balloon = remaining principal
- Grace period: first due date offset correctly

### LoanRepaymentService.AllocatePayment
- Full payment covers penalty → fee → interest → principal in order
- Partial payment stops when amount exhausted
- Overpayment reduces principal ahead of schedule
- Zero/negative amount rejected

### LoanInterestService
- Daily accrual on outstanding principal
- No double accrual for same date
- Skips SETTLED/WRITTEN_OFF contracts

### LoanApplicationService
- Invalid status transitions rejected (DRAFT→DISBURSED)
- Multi-step approval requires correct role order
- Approved amount ≤ product max and ≥ min

## Integration Tests

### Disbursement → GL
1. Create product, application, approve
2. Disburse 50k with post_to_gl=true
3. Assert: transaction POSTED, DR 1200 / CR 1000 balanced
4. Assert: contract status ACTIVE, schedule generated

### Milk Deduction
1. Disburse with create_milk_deduction=true
2. Assert recurrent_deduction created with correct reference
3. Generate member payroll with gross > deduction
4. Approve payroll
5. Assert loan_repayment recorded, schedule updated, GL MEMBER_REPAYMENT_DEDUCTION

### Manual Repayment
1. POST repayment CASH 8500
2. Assert allocation rows, GL LOAN_REPAYMENT
3. Idempotent retry with same key returns same repayment

### Write-off
1. Active contract with 10k outstanding
2. Write-off 10k
3. GL LOAN_WRITE_OFF, status WRITTEN_OFF

### Reversal
1. Reverse disbursement transaction via financial-transactions reverse API
2. Contract returns to PENDING_DISBURSEMENT

### Phone normalization (mobile disbursement)
- `0712345678` → `254712345678`
- Empty phone with no fallback → error

## Mobile disbursement integration

### MOBILE_MPESA (Safaricom Daraja B2C)
1. Disburse via `MOBILE_MPESA` with valid Safaricom MSISDN
2. Assert `loan_disbursements.status = PROCESSING`, `channel_payload.provider = SAFARICOM_B2C`
3. POST sandbox `/api/mpesa-b2c-result` with `Result.ResultCode = 0`
4. Assert disbursement `SUCCESS`, contract GL posted, balances updated
5. Duplicate callback → idempotent (no double GL)

### MOBILE_AIRTEL
1. Disburse via `MOBILE_AIRTEL` sandbox
2. Assert `PROCESSING` with Airtel transaction id in payload
3. POST `/api/airtel-disbursement-callback` with `status_code: TS`
4. Assert `SUCCESS`

### MOBILE_EQUITEL (Jenga Finserve)
1. Seed channel via `007_equitel_channel.sql`
2. Disburse KES 100+ via `MOBILE_EQUITEL`
3. POST `/api/jenga-mobile-callback` with `ResponseCode: 0`
4. Assert `SUCCESS`

### Validation
- Invalid phone → immediate `FAILED`, no provider API call
- Equitel amount < KES 100 → validation error before API call
- M-Pesa amount < KES 10 → validation error

### Provider status
- `GET /loan-module/disbursements/:id/provider-status` for stuck M-Pesa disbursement

## UAT Scenarios

| # | Scenario | Expected |
|---|----------|----------|
| 1 | Member applies for farm loan | Application DRAFT→SUBMITTED |
| 2 | Credit officer approves 80k | Status APPROVED, audit log |
| 3 | Finance disburses to bank | Member receives funds, GL posted |
| 4 | Milk payment deducts installment | Payslip shows deduction, loan balance reduced |
| 5 | Member pays cash top-up | Manual repayment, receipt reference |
| 6 | Loan 45 days overdue | Penalty calculated, aging bucket 31-60 |
| 7 | Early settlement quote | Total payoff amount displayed |
| 8 | Board write-off NPL | Portfolio PAR decreases |
| 9 | Loan restructure extension | New schedule, history preserved |
| 10 | Portfolio report for board | Outstanding, PAR30, recovery rate |
| 11 | Disburse via MOBILE_MPESA | Member receives M-Pesa, GL posted after callback |
| 12 | Disburse via MOBILE_AIRTEL | Airtel wallet credited after callback |
| 13 | Disburse via MOBILE_EQUITEL | Equitel wallet credited after Jenga callback |
| 14 | Invalid mobile phone on disburse | Failed before provider call |

## Performance

- Schedule generation for 10k contracts < 30s batch
- Portfolio report with 50k contracts paginated < 2s per page
- Indexes on member_id, status, due_date verified via EXPLAIN
