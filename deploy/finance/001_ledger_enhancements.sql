-- P0: Ledger integrity enhancements (idempotent, safe to re-run)
-- Run: mysql -u user -p database < 001_ledger_enhancements.sql

-- Unique transaction reference (skip if duplicates exist — run cleanup-coa.sql first)
-- ALTER TABLE transactions ADD UNIQUE INDEX idx_transactions_reference (reference);

-- MySQL 8 does not support ADD COLUMN IF NOT EXISTS; use procedures for idempotency.
DROP PROCEDURE IF EXISTS edairy_add_ledger_columns;
DELIMITER //
CREATE PROCEDURE edairy_add_ledger_columns()
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'transactions' AND column_name = 'idempotency_key'
  ) THEN
    ALTER TABLE transactions ADD COLUMN idempotency_key VARCHAR(64) NULL;
  END IF;
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'transactions' AND column_name = 'reversal_of_id'
  ) THEN
    ALTER TABLE transactions ADD COLUMN reversal_of_id BIGINT UNSIGNED NULL;
  END IF;
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'general_ledger_entries' AND column_name = 'posted_at'
  ) THEN
    ALTER TABLE general_ledger_entries ADD COLUMN posted_at DATETIME NULL;
  END IF;
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'general_ledger_entries' AND column_name = 'is_reversal'
  ) THEN
    ALTER TABLE general_ledger_entries ADD COLUMN is_reversal TINYINT(1) NOT NULL DEFAULT 0;
  END IF;
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'general_ledger_entries' AND column_name = 'reverses_entry_id'
  ) THEN
    ALTER TABLE general_ledger_entries ADD COLUMN reverses_entry_id BIGINT UNSIGNED NULL;
  END IF;
END//
DELIMITER ;
CALL edairy_add_ledger_columns();
DROP PROCEDURE edairy_add_ledger_columns;

-- Normalize GORM-created longtext idempotency_key before unique index.
ALTER TABLE transactions MODIFY idempotency_key VARCHAR(64) NULL;

SET @idx_gle := (
  SELECT COUNT(*) FROM information_schema.statistics
  WHERE table_schema = DATABASE() AND table_name = 'general_ledger_entries' AND index_name = 'idx_gle_account_date'
);
SET @sql_gle := IF(@idx_gle = 0, 'CREATE INDEX idx_gle_account_date ON general_ledger_entries (account_id, transaction_date)', 'SELECT 1');
PREPARE stmt FROM @sql_gle; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx_idem := (
  SELECT COUNT(*) FROM information_schema.statistics
  WHERE table_schema = DATABASE() AND table_name = 'transactions' AND index_name = 'idx_transactions_idempotency'
);
SET @sql_idem := IF(@idx_idem = 0, 'CREATE UNIQUE INDEX idx_transactions_idempotency ON transactions (idempotency_key)', 'SELECT 1');
PREPARE stmt FROM @sql_idem; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS financial_periods (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'OPEN',
  closed_at DATETIME NULL,
  closed_by BIGINT UNSIGNED NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL
);

CREATE TABLE IF NOT EXISTS bank_accounts (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  account_no VARCHAR(64) NULL,
  bank_name VARCHAR(255) NULL,
  ledger_account_id BIGINT UNSIGNED NOT NULL,
  is_default TINYINT(1) NOT NULL DEFAULT 0,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NULL,
  updated_by BIGINT UNSIGNED NULL
);
