-- Dev seed data for milk journal CSV import (POST /api/milk-journals/import)
-- Run against dev_dairy_erp:
--   mysql -uedairy -p dev_dairy_erp < deploy/dev/seed-milk-journal-import.sql

-- 1. Delivery shifts (required — import falls back to first shift when lookup fails)
INSERT INTO milk_delivery_shifts (name, description, created_by, updated_by)
SELECT * FROM (
  SELECT 'Morning' AS name, 'Morning milk collection shift' AS description, 1 AS created_by, 1 AS updated_by
  UNION ALL
  SELECT 'Evening', 'Evening milk collection shift', 1, 1
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM milk_delivery_shifts WHERE name = seed.name AND deleted_at IS NULL);

-- 2. Distinct route codes/names (import matches route_name OR route_code; duplicates break routing)
UPDATE routes SET route_code = 'DEV-R01', route_name = 'Dev Route 01' WHERE id = 8038;
UPDATE routes SET route_code = 'DEV-R02', route_name = 'Dev Route 02' WHERE id = 8039;
UPDATE routes SET route_code = 'DEV-R03', route_name = 'Dev Route 03' WHERE id = 8040;
UPDATE routes SET route_code = 'DEV-R04', route_name = 'Dev Route 04' WHERE id = 8041;
UPDATE routes SET route_code = 'DEV-R05', route_name = 'Dev Route 05' WHERE id = 8042;
UPDATE routes SET route_code = 'DEV-R06', route_name = 'Dev Route 06' WHERE id = 8043;
UPDATE routes SET route_code = 'DEV-R07', route_name = 'Dev Route 07' WHERE id = 8044;
UPDATE routes SET route_code = 'DEV-R08', route_name = 'Dev Route 08' WHERE id = 8045;
UPDATE routes SET route_code = 'DEV-R09', route_name = 'Dev Route 09' WHERE id = 8046;
UPDATE routes SET route_code = 'DEV-R10', route_name = 'Dev Route 10' WHERE id = 8047;

-- 3. Sample milk cans (optional column in import template)
INSERT INTO milk_cans (can_id, can_type, can_size, units, route_id, tare_weight, created_by, updated_by)
SELECT * FROM (
  SELECT 'CAN-001' AS can_id, 'Standard' AS can_type, 50.00 AS can_size, 'L' AS units, 8039 AS route_id, 2.5 AS tare_weight, 1 AS created_by, 1 AS updated_by
  UNION ALL
  SELECT 'CAN-002', 'Standard', 50.00, 'L', 8038, 2.5, 1, 1
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM milk_cans WHERE can_id = seed.can_id AND deleted_at IS NULL);
