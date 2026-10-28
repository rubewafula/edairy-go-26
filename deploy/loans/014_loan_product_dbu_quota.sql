-- DBU allocation quota per loan product (the API also adds these columns on startup).
ALTER TABLE loan_products
  ADD COLUMN dbu_min_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
  ADD COLUMN dbu_max_amount DECIMAL(15,2) NOT NULL DEFAULT 500000;
