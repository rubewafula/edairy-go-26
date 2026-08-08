-- Loan Management Module schema (idempotent where possible)
-- Run: mysql -u user -p database < 001_loan_module_schema.sql

CREATE TABLE IF NOT EXISTS loan_products (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_code VARCHAR(32) NOT NULL,
  product_name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  interest_method ENUM('FLAT','REDUCING_BALANCE','EQUAL_INSTALLMENTS','INTEREST_ONLY','BALLOON') NOT NULL DEFAULT 'REDUCING_BALANCE',
  interest_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
  min_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
  max_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
  repayment_period_months INT NOT NULL DEFAULT 12,
  grace_period_days INT NOT NULL DEFAULT 0,
  processing_fee_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
  insurance_fee_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
  fee_collection_method VARCHAR(32) NOT NULL DEFAULT 'DEDUCT_FROM_PROCEEDS',
  penalty_rate_daily DECIMAL(8,4) NOT NULL DEFAULT 0,
  penalty_rate_monthly DECIMAL(8,4) NOT NULL DEFAULT 0,
  penalty_fixed_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
  allocation_priority JSON NULL,
  recovery_priority INT NOT NULL DEFAULT 1,
  eligibility_rules JSON NULL,
  approval_workflow JSON NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL,
  UNIQUE KEY uk_loan_products_code (product_code)
);

CREATE TABLE IF NOT EXISTS loan_applications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  application_no VARCHAR(64) NOT NULL,
  member_id BIGINT UNSIGNED NOT NULL,
  loan_product_id BIGINT UNSIGNED NOT NULL,
  requested_amount DECIMAL(15,2) NOT NULL,
  requested_term_months INT NOT NULL,
  purpose VARCHAR(500) NULL,
  status ENUM('DRAFT','SUBMITTED','UNDER_REVIEW','APPROVED','REJECTED','CANCELLED','DISBURSED') NOT NULL DEFAULT 'DRAFT',
  approved_amount DECIMAL(15,2) NULL,
  approved_term_months INT NULL,
  date_applied DATETIME NOT NULL,
  expected_disbursement_date DATE NULL,
  rejected_reason TEXT NULL,
  loan_contract_id BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL,
  UNIQUE KEY uk_loan_applications_no (application_no),
  KEY idx_loan_applications_member (member_id),
  KEY idx_loan_applications_status (status),
  KEY idx_loan_applications_product (loan_product_id)
);

CREATE TABLE IF NOT EXISTS loan_approvals (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_application_id BIGINT UNSIGNED NOT NULL,
  approval_role ENUM('CREDIT_OFFICER','FINANCE_OFFICER','MANAGER','BOARD') NOT NULL,
  approver_id BIGINT UNSIGNED NOT NULL,
  approval_date DATETIME NOT NULL,
  decision ENUM('APPROVED','REJECTED','RETURNED') NOT NULL,
  approved_amount DECIMAL(15,2) NULL,
  approved_term_months INT NULL,
  comments TEXT NULL,
  conditions TEXT NULL,
  step_order INT NOT NULL DEFAULT 1,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL,
  KEY idx_loan_approvals_application (loan_application_id)
);

CREATE TABLE IF NOT EXISTS loan_contracts (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  contract_no VARCHAR(64) NOT NULL,
  loan_application_id BIGINT UNSIGNED NULL,
  member_id BIGINT UNSIGNED NOT NULL,
  loan_product_id BIGINT UNSIGNED NOT NULL,
  principal_amount DECIMAL(15,2) NOT NULL,
  approved_amount DECIMAL(15,2) NOT NULL,
  disbursed_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
  fees_deducted_at_disbursement DECIMAL(15,2) NOT NULL DEFAULT 0,
  interest_rate DECIMAL(8,4) NOT NULL,
  interest_method VARCHAR(32) NOT NULL,
  term_months INT NOT NULL,
  grace_period_days INT NOT NULL DEFAULT 0,
  status ENUM('PENDING_DISBURSEMENT','ACTIVE','OVERDUE','RESTRUCTURED','SETTLED','WRITTEN_OFF','CLOSED') NOT NULL DEFAULT 'PENDING_DISBURSEMENT',
  outstanding_principal DECIMAL(15,2) NOT NULL DEFAULT 0,
  outstanding_interest DECIMAL(15,2) NOT NULL DEFAULT 0,
  outstanding_fees DECIMAL(15,2) NOT NULL DEFAULT 0,
  outstanding_penalties DECIMAL(15,2) NOT NULL DEFAULT 0,
  total_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
  days_in_arrears INT NOT NULL DEFAULT 0,
  disbursement_date DATE NULL,
  maturity_date DATE NULL,
  last_payment_date DATE NULL,
  installment_amount DECIMAL(15,2) NULL,
  milk_deduction_enabled TINYINT(1) NOT NULL DEFAULT 1,
  max_deduction_percent DECIMAL(5,2) NOT NULL DEFAULT 50.00,
  recurrent_deduction_id BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL,
  UNIQUE KEY uk_loan_contracts_no (contract_no),
  KEY idx_loan_contracts_member (member_id),
  KEY idx_loan_contracts_status (status),
  KEY idx_loan_contracts_product (loan_product_id)
);

CREATE TABLE IF NOT EXISTS loan_disbursements (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  method ENUM('CASH','BANK','MOBILE_MONEY','WALLET') NOT NULL DEFAULT 'BANK',
  reference VARCHAR(128) NOT NULL,
  idempotency_key VARCHAR(64) NULL,
  disbursement_date DATE NOT NULL,
  gl_transaction_id BIGINT UNSIGNED NULL,
  notes TEXT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL,
  UNIQUE KEY uk_loan_disbursements_ref (reference),
  UNIQUE KEY uk_loan_disbursements_idem (idempotency_key),
  KEY idx_loan_disbursements_contract (loan_contract_id)
);

CREATE TABLE IF NOT EXISTS loan_schedule_installments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  installment_no INT NOT NULL,
  due_date DATE NOT NULL,
  principal_due DECIMAL(15,2) NOT NULL DEFAULT 0,
  interest_due DECIMAL(15,2) NOT NULL DEFAULT 0,
  fee_due DECIMAL(15,2) NOT NULL DEFAULT 0,
  insurance_due DECIMAL(15,2) NOT NULL DEFAULT 0,
  penalty_due DECIMAL(15,2) NOT NULL DEFAULT 0,
  total_due DECIMAL(15,2) NOT NULL DEFAULT 0,
  principal_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
  interest_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
  fee_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
  insurance_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
  penalty_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
  balance_after DECIMAL(15,2) NOT NULL DEFAULT 0,
  status ENUM('PENDING','PARTIAL','PAID','OVERDUE','WAIVED') NOT NULL DEFAULT 'PENDING',
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  KEY idx_loan_schedule_contract (loan_contract_id),
  KEY idx_loan_schedule_due (due_date),
  KEY idx_loan_schedule_status (status)
);

CREATE TABLE IF NOT EXISTS loan_repayments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  channel ENUM('MILK_DEDUCTION','CASH','BANK','MOBILE_MONEY','JOURNAL') NOT NULL,
  reference VARCHAR(128) NOT NULL,
  idempotency_key VARCHAR(64) NULL,
  payment_date DATE NOT NULL,
  member_payroll_id BIGINT UNSIGNED NULL,
  member_payslip_id BIGINT UNSIGNED NULL,
  gl_transaction_id BIGINT UNSIGNED NULL,
  notes TEXT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL,
  UNIQUE KEY uk_loan_repayments_ref (reference),
  UNIQUE KEY uk_loan_repayments_idem (idempotency_key),
  KEY idx_loan_repayments_contract (loan_contract_id),
  KEY idx_loan_repayments_date (payment_date)
);

CREATE TABLE IF NOT EXISTS loan_repayment_allocations (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_repayment_id BIGINT UNSIGNED NOT NULL,
  loan_schedule_installment_id BIGINT UNSIGNED NULL,
  allocation_type ENUM('PENALTY','FEE','INTEREST','PRINCIPAL','INSURANCE') NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  created_at DATETIME NULL,
  KEY idx_loan_alloc_repayment (loan_repayment_id)
);

CREATE TABLE IF NOT EXISTS loan_charges (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  charge_type ENUM('PROCESSING','INSURANCE','LEGAL','OTHER') NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  capitalized TINYINT(1) NOT NULL DEFAULT 0,
  charged_date DATE NOT NULL,
  gl_transaction_id BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  KEY idx_loan_charges_contract (loan_contract_id)
);

CREATE TABLE IF NOT EXISTS loan_penalties (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  loan_schedule_installment_id BIGINT UNSIGNED NULL,
  penalty_type ENUM('LATE_PAYMENT','DAILY','MONTHLY','FIXED') NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  penalty_date DATE NOT NULL,
  waived TINYINT(1) NOT NULL DEFAULT 0,
  waived_by BIGINT UNSIGNED NULL,
  waived_at DATETIME NULL,
  gl_transaction_id BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  KEY idx_loan_penalties_contract (loan_contract_id)
);

CREATE TABLE IF NOT EXISTS loan_interest_accruals (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  accrual_date DATE NOT NULL,
  principal_base DECIMAL(15,2) NOT NULL,
  interest_amount DECIMAL(15,2) NOT NULL,
  posted TINYINT(1) NOT NULL DEFAULT 0,
  gl_transaction_id BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  UNIQUE KEY uk_loan_accrual_contract_date (loan_contract_id, accrual_date),
  KEY idx_loan_accruals_date (accrual_date)
);

CREATE TABLE IF NOT EXISTS loan_guarantors (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_application_id BIGINT UNSIGNED NOT NULL,
  member_id BIGINT UNSIGNED NOT NULL,
  guaranteed_amount DECIMAL(15,2) NOT NULL,
  relationship VARCHAR(128) NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  KEY idx_loan_guarantors_application (loan_application_id)
);

CREATE TABLE IF NOT EXISTS loan_documents (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_application_id BIGINT UNSIGNED NULL,
  loan_contract_id BIGINT UNSIGNED NULL,
  document_type VARCHAR(64) NOT NULL,
  file_path VARCHAR(512) NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  uploaded_by BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  KEY idx_loan_documents_application (loan_application_id),
  KEY idx_loan_documents_contract (loan_contract_id)
);

CREATE TABLE IF NOT EXISTS loan_restructurings (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  restructure_type ENUM('EXTENSION','REFINANCING','RESCHEDULING','TOP_UP','WAIVER') NOT NULL,
  previous_term_months INT NULL,
  new_term_months INT NULL,
  previous_rate DECIMAL(8,4) NULL,
  new_rate DECIMAL(8,4) NULL,
  reason TEXT NULL,
  effective_date DATE NOT NULL,
  created_by BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  KEY idx_loan_restructurings_contract (loan_contract_id)
);

CREATE TABLE IF NOT EXISTS loan_write_offs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  loan_contract_id BIGINT UNSIGNED NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  reason TEXT NULL,
  write_off_date DATE NOT NULL,
  gl_transaction_id BIGINT UNSIGNED NULL,
  approved_by BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  KEY idx_loan_write_offs_contract (loan_contract_id)
);

CREATE TABLE IF NOT EXISTS loan_audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  entity_type VARCHAR(64) NOT NULL,
  entity_id BIGINT UNSIGNED NOT NULL,
  action VARCHAR(64) NOT NULL,
  old_values JSON NULL,
  new_values JSON NULL,
  user_id BIGINT UNSIGNED NULL,
  ip_address VARCHAR(45) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_loan_audit_entity (entity_type, entity_id),
  KEY idx_loan_audit_created (created_at)
);
