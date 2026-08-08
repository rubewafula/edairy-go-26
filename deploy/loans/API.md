# Loan Module REST API Specification

All endpoints require JWT authentication unless noted. Pagination: `?page=1&limit=50`. Filtering via query params.

## Loan Products

### `GET /api/loan-module/products`
Query: `active=true`, `code=`

Response:
```json
{
  "data": [{
    "id": 1,
    "product_code": "FARM-001",
    "product_name": "Farm Input Loan",
    "interest_method": "REDUCING_BALANCE",
    "interest_rate": 12.0,
    "min_amount": 5000,
    "max_amount": 500000,
    "repayment_period_months": 12,
    "grace_period_days": 30,
    "processing_fee_rate": 2.0,
    "insurance_fee_rate": 1.0,
    "is_active": true
  }],
  "total": 1
}
```

### `POST /api/loan-module/products`
Body: product fields (see CreateLoanProductRequest).

### `GET /api/loan-module/products/:id`
### `PUT /api/loan-module/products/:id`
### `DELETE /api/loan-module/products/:id` — soft delete; fails if active contracts exist.

---

## Loan Applications

### `POST /api/loan-module/applications`
```json
{
  "member_id": 42,
  "loan_product_id": 1,
  "requested_amount": 100000,
  "requested_term_months": 12,
  "purpose": "Dairy feed purchase",
  "expected_disbursement_date": "2026-08-15",
  "guarantors": [{"member_id": 55, "guaranteed_amount": 50000}]
}
```
Status: `DRAFT`

### `POST /api/loan-module/applications/:id/submit`
Draft → SUBMITTED. Idempotent if already submitted.

### `POST /api/loan-module/applications/:id/approve`
```json
{
  "approval_role": "CREDIT_OFFICER",
  "approved_amount": 80000,
  "approved_term_months": 10,
  "comments": "Approved with reduced amount",
  "conditions": "Must maintain milk delivery"
}
```

### `POST /api/loan-module/applications/:id/reject`
```json
{ "comments": "Insufficient delivery history" }
```

### `GET /api/loan-module/applications`
Query: `status=`, `member_id=`

---

## Loan Contracts

Created automatically when final approval completes, or via `POST /api/loan-module/contracts`.

### `POST /api/loan-module/contracts/:id/disburse`
```json
{
  "amount": 80000,
  "method": "BANK",
  "reference": "DISB-2026-001",
  "disbursement_date": "2026-08-15",
  "post_to_gl": true,
  "create_milk_deduction": true,
  "installment_amount": 8500
}
```
Supports partial/multiple disbursements until `approved_amount` reached.

### `GET /api/loan-module/contracts/:id/schedule`
Returns installment lines with principal, interest, fees, due dates, paid amounts, status.

### `POST /api/loan-module/contracts/:id/repayments`
```json
{
  "amount": 8500,
  "channel": "CASH",
  "payment_date": "2026-09-01",
  "reference": "REP-2026-001",
  "idempotency_key": "uuid",
  "post_to_gl": true
}
```
Allocation follows product priority. Supports partial, overpayment, backdated.

### `GET /api/loan-module/contracts/:id/statement`
Loan statement with opening balance, disbursements, repayments, accruals, closing balance.

### `POST /api/loan-module/contracts/:id/settle`
Query `?quote=true` returns settlement quotation without posting.

### `POST /api/loan-module/contracts/:id/restructure`
```json
{
  "type": "EXTENSION",
  "new_term_months": 18,
  "reason": "Seasonal cash flow",
  "effective_date": "2026-10-01"
}
```

### `POST /api/loan-module/contracts/:id/write-off`
```json
{
  "amount": 15000,
  "reason": "Member exited cooperative",
  "post_to_gl": true
}
```

---

## Interest & Penalties

### `POST /api/loan-module/interest/accrue`
```json
{
  "as_of_date": "2026-07-31",
  "loan_contract_ids": [],
  "post_to_gl": true
}
```
Empty `loan_contract_ids` = all active contracts.

### `POST /api/loan-module/penalties/run`
Batch late penalty calculation for overdue installments.

---

## Reports

### `GET /api/loan-module/reports/portfolio`
Summary: total outstanding, count by status, average rate.

### `GET /api/loan-module/reports/aging`
Buckets: current, 1-30, 31-60, 61-90, 90+ days.

### `GET /api/loan-module/reports/par`
PAR 30/60/90, NPL ratio, recovery rate.

### `GET /api/loan-module/reports/register`
Full loan register with filters.

---

## Idempotency

Use header `Idempotency-Key` or body `idempotency_key` on:
- Disbursements
- Repayments
- Interest accrual batches
- Write-offs

Duplicate keys return the original result (HTTP 200/201).

## Permissions

| Permission | Action |
|------------|--------|
| loan-products.view | List products |
| loan-products.manage | CRUD products |
| loan-applications.view | View applications |
| loan-applications.create | Create/submit |
| loan-applications.approve | Approve/reject |
| loan-contracts.view | View contracts |
| loan-contracts.disburse | Disburse |
| loan-contracts.repay | Record repayments |
| loan-contracts.restructure | Restructure/write-off |
| loan-reports.view | Reports |

## Error Codes

| HTTP | Meaning |
|------|---------|
| 400 | Validation / business rule violation |
| 403 | Missing permission |
| 404 | Resource not found |
| 409 | Invalid status transition / duplicate reference |
| 422 | Unbalanced GL / allocation failure |
