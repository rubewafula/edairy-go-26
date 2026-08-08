-- Disbursement channels catalog + loan_disbursements status tracking
-- Run after 001_loan_module_schema.sql (or rely on EnsureLoanSchema AutoMigrate)

CREATE TABLE IF NOT EXISTS disbursement_channels (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  channel_code VARCHAR(32) NOT NULL,
  channel_name VARCHAR(128) NOT NULL,
  description VARCHAR(512) NULL,
  channel_type ENUM('CASH','BANK','MOBILE','ITEM','WALLET') NOT NULL,
  provider VARCHAR(64) NULL,
  is_async TINYINT(1) NOT NULL DEFAULT 0,
  requires_external_ref TINYINT(1) NOT NULL DEFAULT 0,
  gl_rule_type VARCHAR(64) NOT NULL DEFAULT 'LOAN_DISBURSEMENT',
  config_schema JSON NULL,
  sort_order INT NOT NULL DEFAULT 0,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL,
  UNIQUE KEY uk_disbursement_channels_code (channel_code),
  KEY idx_disbursement_channels_type (channel_type),
  KEY idx_disbursement_channels_active (is_active)
);

-- Seed default channels (idempotent)
INSERT INTO disbursement_channels
  (channel_code, channel_name, description, channel_type, provider, is_async, requires_external_ref, gl_rule_type, sort_order, is_active, created_at, updated_at)
VALUES
  ('CASH', 'Cash', 'Physical cash paid to member at office', 'CASH', 'INTERNAL', 0, 0, 'LOAN_DISBURSEMENT', 10, 1, NOW(), NOW()),
  ('ITEM_PURCHASE', 'Item Purchase', 'Loan disbursed as goods/services via store (input order)', 'ITEM', 'STORE', 1, 1, 'LOAN_DISBURSEMENT', 20, 1, NOW(), NOW()),
  ('MOBILE_MPESA', 'Mobile - M-Pesa', 'B2C transfer to member M-Pesa wallet', 'MOBILE', 'MPESA', 1, 1, 'LOAN_DISBURSEMENT', 30, 1, NOW(), NOW()),
  ('MOBILE_AIRTEL', 'Mobile - Airtel Money', 'Transfer to member Airtel Money wallet', 'MOBILE', 'AIRTEL', 1, 1, 'LOAN_DISBURSEMENT', 40, 1, NOW(), NOW()),
  ('BANK_COOP', 'Bank - Co-operative Bank', 'Transfer to member Co-op Bank account', 'BANK', 'COOP', 0, 1, 'LOAN_DISBURSEMENT', 50, 1, NOW(), NOW()),
  ('BANK_EQUITY', 'Bank - Equity Bank', 'Transfer to member Equity Bank account', 'BANK', 'EQUITY', 0, 1, 'LOAN_DISBURSEMENT', 60, 1, NOW(), NOW()),
  ('BANK_KCB', 'Bank - KCB', 'Transfer to member KCB account', 'BANK', 'KCB', 0, 1, 'LOAN_DISBURSEMENT', 70, 1, NOW(), NOW()),
  ('WALLET_MEMBER', 'Member Wallet', 'Credit member internal wallet balance', 'WALLET', 'INTERNAL', 0, 0, 'LOAN_DISBURSEMENT', 80, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  channel_name = VALUES(channel_name),
  description = VALUES(description),
  channel_type = VALUES(channel_type),
  provider = VALUES(provider),
  is_async = VALUES(is_async),
  requires_external_ref = VALUES(requires_external_ref),
  updated_at = NOW();
