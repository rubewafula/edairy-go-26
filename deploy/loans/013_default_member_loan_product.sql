-- Default member loan product: the product used for mobile loan applications.
-- EnsureLoanSchema (AutoMigrate) also adds this column on API startup; this script is for manual deploys.
-- MySQL 8.0.29+ supports IF NOT EXISTS on ADD COLUMN; drop it on older servers if the column is missing.
ALTER TABLE loan_products
  ADD COLUMN IF NOT EXISTS is_default_member_loan TINYINT(1) NOT NULL DEFAULT 0;

CREATE INDEX idx_loan_products_is_default_member_loan ON loan_products (is_default_member_loan);

-- Optional: mark an existing product as the default (only one should be flagged)
-- UPDATE loan_products SET is_default_member_loan = (product_code = 'MOBILE') WHERE deleted_at IS NULL;
