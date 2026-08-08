# eDairy Loan Management Module

Production-ready SACCO-style lending for dairy cooperatives. Integrates with Members, Finance/GL, Milk Payments, Savings, Shares, Cash/Bank, and Notifications.

## Architecture

```mermaid
flowchart TB
  subgraph UI [React Admin]
    LP[Loan Products]
    LA[Applications]
    LD[Loan Dashboard]
    LR[Repayments]
  end

  subgraph API [Go API Layer]
    LC[LoanModuleController]
  end

  subgraph Domain [Domain Services]
    LProd[LoanProductService]
    LApp[LoanApplicationService]
    LCon[LoanContractService]
    LSch[LoanScheduleService]
    LRep[LoanRepaymentService]
    LInt[LoanInterestService]
    LMon[LoanMonitoringService]
    LPost[FinancialPostingService]
  end

  subgraph Integrations
    MP[MemberPayrollService]
    Led[LedgerService]
    RD[RecurrentDeductions]
  end

  UI --> LC
  LC --> Domain
  LCon --> LSch
  LCon --> LPost --> Led
  LRep --> LPost
  LRep --> RD
  MP --> LRep
```

## Entity Relationship (core)

```mermaid
erDiagram
  LOAN_PRODUCTS ||--o{ LOAN_APPLICATIONS : offers
  MEMBERS ||--o{ LOAN_APPLICATIONS : applies
  LOAN_APPLICATIONS ||--o{ LOAN_APPROVALS : has
  LOAN_APPLICATIONS ||--o| LOAN_CONTRACTS : becomes
  LOAN_CONTRACTS ||--o{ LOAN_DISBURSEMENTS : funded_by
  LOAN_CONTRACTS ||--o{ LOAN_SCHEDULE_INSTALLMENTS : schedule
  LOAN_CONTRACTS ||--o{ LOAN_REPAYMENTS : receives
  LOAN_REPAYMENTS ||--o{ LOAN_REPAYMENT_ALLOCATIONS : allocates
  LOAN_CONTRACTS ||--o{ LOAN_INTEREST_ACCRUALS : accrues
  LOAN_CONTRACTS ||--o{ LOAN_PENALTIES : penalized
  LOAN_CONTRACTS ||--o{ LOAN_GUARANTORS : secured_by
  LOAN_CONTRACTS ||--o{ LOAN_DOCUMENTS : documented
  LOAN_CONTRACTS ||--o{ LOAN_RESTRUCTURINGS : restructured
  LOAN_CONTRACTS ||--o{ LOAN_WRITE_OFFS : written_off
  LOAN_CONTRACTS ||--o{ LOAN_AUDIT_LOGS : audited
```

## Business Workflows

### Application → Disbursement

```mermaid
sequenceDiagram
  participant M as Member
  participant App as ApplicationService
  participant Apr as ApprovalWorkflow
  participant Con as ContractService
  participant Sch as ScheduleService
  participant GL as LedgerService

  M->>App: Submit application
  App->>Apr: Route for approval
  Apr->>App: APPROVED
  App->>Con: Create contract
  Con->>Sch: Generate schedule
  Con->>GL: Post fees (optional)
  Con->>Con: Disburse (cash/bank/wallet)
  Con->>GL: LOAN_DISBURSEMENT
  Con->>Con: Create milk deduction plan
```

### Milk Payment Deduction

```mermaid
sequenceDiagram
  participant Pay as MemberPayroll
  participant RD as RecurrentDeduction
  participant Rep as RepaymentService
  participant Sch as Schedule
  participant GL as LedgerService

  Pay->>RD: Load unsettled loan deductions
  Pay->>Rep: Allocate from net milk pay
  Rep->>Sch: Apply to installments (penalty/fee/int/principal)
  Rep->>GL: LOAN_REPAYMENT + MEMBER_REPAYMENT_DEDUCTION
  Rep->>RD: Update paid_amount / settled
```

## Application Statuses

| Status | Description |
|--------|-------------|
| DRAFT | Editable by member/staff |
| SUBMITTED | Awaiting review |
| UNDER_REVIEW | Credit assessment in progress |
| APPROVED | Approved, pending contract/disbursement |
| REJECTED | Declined with reason |
| CANCELLED | Withdrawn |
| DISBURSED | Contract created and funded |

## Contract Statuses

| Status | Description |
|--------|-------------|
| PENDING_DISBURSEMENT | Approved, not yet funded |
| ACTIVE | Performing loan |
| OVERDUE | Past due installments |
| RESTRUCTURED | Terms modified |
| SETTLED | Fully paid |
| WRITTEN_OFF | Bad debt write-off |
| CLOSED | Administrative closure |

## Interest Methods

| Method | Description |
|--------|-------------|
| FLAT | Total interest = P × rate × periods; equal installments |
| REDUCING_BALANCE | Interest on outstanding principal each period |
| EQUAL_INSTALLMENTS | Standard PMT amortization |
| INTEREST_ONLY | Interest each period; principal at maturity |
| BALLOON | Periodic interest/small principal; large final payment |

## Repayment Allocation Priority (default)

1. Penalties  
2. Fees (processing, insurance)  
3. Accrued interest  
4. Principal  

Configurable per product via `allocation_priority` JSON array.

## Double-Entry Posting Matrix

See [posting-matrix.yaml](./posting-matrix.yaml). All entries posted via `LedgerService` (immutable GL, reversals only).

| Event | Debit | Credit |
|-------|-------|--------|
| Loan disbursement | 1200 Member Loans | 1000 Cash at Bank |
| Loan repayment (cash) | 1000 Cash at Bank | 1200 Member Loans |
| Milk deduction | 2215 Member Payables | 1200 Member Loans |
| Processing fee (cash) | 1000 Cash | 4100 Loan Fee Income |
| Processing fee (capitalized) | 1200 Member Loans | 4100 Loan Fee Income |
| Interest accrual | 1200 Member Loans | 4010 Interest Income |
| Penalty charged | 1200 Member Loans | 4020 Penalty Income |
| Penalty waiver | 4020 Penalty Income | 1200 Member Loans |
| Write-off | 5105 Bad Debt Expense | 1200 Member Loans |
| Recovery after write-off | 1000 Cash | 5105 Bad Debt Expense |

## REST API Summary

Base path: `/api/loan-module`

| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/products` | Loan products CRUD |
| GET/PUT/DELETE | `/products/:id` | Single product |
| GET/POST | `/applications` | Applications |
| POST | `/applications/:id/submit` | Submit for review |
| POST | `/applications/:id/approve` | Approval step |
| POST | `/applications/:id/reject` | Reject |
| GET/POST | `/contracts` | Active loans |
| POST | `/contracts/:id/disburse` | Disburse funds |
| GET | `/contracts/:id/schedule` | Repayment schedule |
| POST | `/contracts/:id/repayments` | Manual repayment |
| GET | `/contracts/:id/statement` | Loan statement |
| POST | `/contracts/:id/settle` | Early settlement quote/execute |
| POST | `/contracts/:id/restructure` | Restructure |
| POST | `/contracts/:id/write-off` | Write-off |
| POST | `/interest/accrue` | Batch interest accrual |
| GET | `/reports/portfolio` | Portfolio summary |
| GET | `/reports/aging` | Aging report |
| GET | `/reports/par` | PAR / NPL metrics |

Full specification: [API.md](./API.md)

## Deployment

1. Run `001_loan_module_schema.sql` (or rely on `EnsureLoanSchema()` on API startup).
2. Seed posting rules from `posting-matrix.yaml` via finance seed or SQL.
3. Grant permissions: `loan-products.*`, `loan-applications.*`, `loan-contracts.*`, `loan-reports.view`.
4. Configure default loan products per dairy.

## Migration from Legacy

- Legacy `member_loans` / `loans` CRUD remains for backward compatibility.
- New SACCO workflows use `loan_contracts` and related tables.
- Optional data migration script maps approved legacy loans → contracts (future).

## Test Cases

See [TEST-CASES.md](./TEST-CASES.md) for unit, integration, and UAT scenarios.

## Administrator Guide

1. **Configure products** — set rates, limits, grace, fees, allocation priority.
2. **Approval workflow** — assign roles (credit, finance, manager, board).
3. **Disburse** — choose channel; GL posts automatically when `post_to_gl=true`.
4. **Milk recovery** — on disbursement a `recurrent_deduction` is created; payroll approval recovers automatically.
5. **Interest accrual** — schedule cron or manual `POST /loan-module/interest/accrue`.
6. **Monitoring** — use portfolio, aging, and PAR reports for board meetings.
