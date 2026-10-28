-- Demo seed for dev_dairy_erp (idempotent where possible).
-- Usage on edairy-db:
--   mysql dev_dairy_erp < seed-dev-demo-data.sql

-- SMS campaigns (dev had none; outbox API expects campaign ids)
INSERT IGNORE INTO sms_campaigns (
  campaign_code, campaign_name, description, sms_message_id,
  total_recipients, total_sent, total_failed, status, created_by, updated_by
) VALUES
  ('DEV-SMS-001', 'Loan disbursement notices', 'Notify members when loan funds are sent', 63, 0, 0, 0, 'draft', 1, 1),
  ('DEV-SMS-002', 'Milk payment alerts', 'Payment confirmation after milk delivery', 63, 0, 0, 0, 'draft', 1, 1),
  ('DEV-SMS-003', 'Member welcome', 'Welcome SMS for new registrations', 58, 0, 0, 0, 'draft', 1, 1),
  ('DEV-SMS-004', 'Arrears reminder', 'Loan repayment reminders', 63, 0, 0, 0, 'scheduled', 1, 1),
  ('DEV-SMS-005', 'General broadcast', 'Co-op announcements and events', 59, 0, 0, 0, 'draft', 1, 1);

-- Loan products (expand catalog for UI/testing)
INSERT INTO loan_products (
  created_at, updated_at, created_by, updated_by,
  product_code, product_name, description, interest_method, interest_rate,
  min_amount, max_amount, repayment_period_months, grace_period_days,
  processing_fee_rate, is_active, fee_collection_method
)
SELECT NOW(3), NOW(3), 1, 1, product_code, product_name, description, interest_method, interest_rate,
       min_amount, max_amount, repayment_period_months, grace_period_days,
       processing_fee_rate, 1, 'DEDUCT_FROM_PROCEEDS'
FROM (
  SELECT 'DEV-DAIRY-02' AS product_code, 'Dairy Working Capital' AS product_name,
         'Short-term working capital for milk handlers' AS description,
         'REDUCING_BALANCE' AS interest_method, 14.0000 AS interest_rate,
         5000.00 AS min_amount, 150000.00 AS max_amount, 6 AS repayment_period_months, 14 AS grace_period_days,
         1.5000 AS processing_fee_rate
  UNION ALL SELECT 'DEV-EQUIP-03', 'Equipment Lease Loan', 'Tractors, coolers, and dairy equipment', 'EQUAL_INSTALLMENTS', 11.5000,
         20000.00, 500000.00, 24, 30, 2.0000
  UNION ALL SELECT 'DEV-FEED-04', 'Feed & Fodder Credit', 'Seasonal input financing', 'FLAT', 10.0000,
         3000.00, 80000.00, 9, 7, 1.0000
  UNION ALL SELECT 'DEV-EMERG-05', 'Emergency Member Loan', 'Quick approval for urgent needs', 'REDUCING_BALANCE', 15.0000,
         1000.00, 50000.00, 3, 0, 0.5000
  UNION ALL SELECT 'DEV-SCHOOL-06', 'School Fees Bridge', 'Term-based school fee support', 'EQUAL_INSTALLMENTS', 12.0000,
         5000.00, 120000.00, 12, 30, 1.2500
) AS products
WHERE NOT EXISTS (SELECT 1 FROM loan_products lp WHERE lp.product_code = products.product_code);
