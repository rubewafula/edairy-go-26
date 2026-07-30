# Finance Deployment Template

Standard accounting setup for new eDairy deployments.

## Quick Start

```bash
# 1. Apply schema enhancements
mysql -u $DB_USER -p$DB_PASSWORD $DB_NAME < deploy/finance/001_ledger_enhancements.sql

# 2. Seed chart of accounts + posting rules (new dairy)
cd edairy-go && go run ./cmd/seed-finance/main.go

# 3. Cleanup existing migrated database (Arithi-style duplicates)
mysql -u $DB_USER -p$DB_PASSWORD $DB_NAME < deploy/finance/cleanup-coa.sql
```

## What's Included

| Asset | Purpose |
|-------|---------|
| `chart-of-accounts.yaml` | Default COA, posting rules, checklist |
| `001_ledger_enhancements.sql` | GL immutability columns, periods, bank accounts |
| `cleanup-coa.sql` | Dedupe rules, fix 444 accounts, add missing rules |
| `cmd/seed-finance` | Go seed command for new dairies |

## API Endpoints (new)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/financial-transactions` | Post any configured rule type to GL |
| POST | `/api/financial-transactions/:id/reverse` | Reverse posted transaction |
| GET | `/api/general-ledger?from=&to=` | General ledger report |
| GET | `/api/cash-flow-statement?from=&to=` | Cash flow statement |
| GET | `/api/member-statements/:member_id` | Member financial statement |
| GET | `/api/loan-statements/:loan_id` | Loan statement |
| GET | `/api/financial-periods` | List accounting periods |
| POST | `/api/financial-periods/:id/close` | Period-end close |
| GET | `/api/trial-balance?from=&to=` | Trial balance (date filtered) |

## Idempotency

Pass `Idempotency-Key` header or `idempotency_key` in body on:
- Share payments
- Cash transactions (`post_to_gl: true`)
- Loan transactions (`post_to_gl: true`)
- Financial transaction POST

## Opening Checklist

See `chart-of-accounts.yaml` → `opening_checklist`
