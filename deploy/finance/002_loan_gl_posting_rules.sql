-- Idempotent loan GL accounts + posting rules for production dairies.
-- Uses existing account codes where present; skips rules when accounts are missing.

INSERT IGNORE INTO accounts (account_code, name, account_category_id, is_postable, is_active, created_at, updated_at)
SELECT '4020', 'Loan Penalty Income', ac.id, 1, 1, NOW(), NOW()
FROM accounts a
JOIN account_categories ac ON ac.id = a.account_category_id
WHERE a.account_code = '4010' AND a.deleted_at IS NULL
LIMIT 1;

INSERT IGNORE INTO accounts (account_code, name, account_category_id, is_postable, is_active, created_at, updated_at)
SELECT '4100', 'Loan Processing Fee Income', ac.id, 1, 1, NOW(), NOW()
FROM accounts a
JOIN account_categories ac ON ac.id = a.account_category_id
WHERE a.account_code = '4010' AND a.deleted_at IS NULL
LIMIT 1;

INSERT IGNORE INTO accounts (account_code, name, account_category_id, is_postable, is_active, created_at, updated_at)
SELECT '5105', 'Bad Debt Expense', ac.id, 1, 1, NOW(), NOW()
FROM accounts a
JOIN account_categories ac ON ac.id = a.account_category_id
WHERE a.account_code = '5000' AND a.deleted_at IS NULL
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_DISBURSEMENT', d.id, c.id, 'Loan principal disbursement', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '1200' AND c.account_code = '1000'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_DISBURSEMENT')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_REPAYMENT', d.id, c.id, 'Loan repayment (cash/bank)', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '1000' AND c.account_code = '1200'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_REPAYMENT')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'MEMBER_REPAYMENT_DEDUCTION', d.id, c.id, 'Loan recovery from member milk payment', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '2210' AND c.account_code = '1200'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'MEMBER_REPAYMENT_DEDUCTION')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_PROCESSING_FEE', d.id, c.id, 'Loan processing fee (cash)', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '1000' AND c.account_code = '4100'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_PROCESSING_FEE')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_PROCESSING_FEE_CAPITALIZED', d.id, c.id, 'Processing fee deducted from loan proceeds', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '1200' AND c.account_code = '4100'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_PROCESSING_FEE_CAPITALIZED')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_INTEREST_ACCRUAL', d.id, c.id, 'Interest accrual on outstanding loans', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '1200' AND c.account_code = '4010'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_INTEREST_ACCRUAL')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_PENALTY', d.id, c.id, 'Late payment penalty', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '1200' AND c.account_code = '4020'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_PENALTY')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_PENALTY_WAIVER', d.id, c.id, 'Penalty waiver', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '4020' AND c.account_code = '1200'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_PENALTY_WAIVER')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_WRITE_OFF', d.id, c.id, 'Loan principal write-off', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '5105' AND c.account_code = '1200'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_WRITE_OFF')
LIMIT 1;

INSERT INTO transaction_posting_rules (transaction_type, debit_account_id, credit_account_id, description, created_at, updated_at)
SELECT 'LOAN_RECOVERY_AFTER_WRITE_OFF', d.id, c.id, 'Recovery of previously written-off loan', NOW(), NOW()
FROM accounts d, accounts c
WHERE d.account_code = '1000' AND c.account_code = '5105'
  AND d.deleted_at IS NULL AND c.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM transaction_posting_rules WHERE transaction_type = 'LOAN_RECOVERY_AFTER_WRITE_OFF')
LIMIT 1;

DELETE t1 FROM transaction_posting_rules t1
INNER JOIN transaction_posting_rules t2
  ON t1.transaction_type = t2.transaction_type AND t1.id > t2.id
WHERE t1.transaction_type LIKE 'LOAN%'
   OR t1.transaction_type = 'MEMBER_REPAYMENT_DEDUCTION';
