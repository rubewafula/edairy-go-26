-- Dev seed data for transporters module
-- Run against dev_dairy_erp:
--   mysql -uedairy -p dev_dairy_erp < deploy/dev/seed-transporters.sql

-- Individual transporter 1
INSERT INTO transporters (transporter_no, category, primary_phone, email_address, status, restricted, created_by, updated_by)
SELECT 'TRP-001', 'INDIVIDUAL', '0711000001', 'john.kamau@dev-transport.test', 'ACTIVE', 0, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM transporters WHERE transporter_no = 'TRP-001' AND deleted_at IS NULL);

SET @t1 := (SELECT id FROM transporters WHERE transporter_no = 'TRP-001' AND deleted_at IS NULL LIMIT 1);

INSERT INTO individual_transporters (transporter_id, first_name, last_name, other_names, gender, date_of_birth, national_id_no, kra_pin, created_by, updated_by)
SELECT @t1, 'John', 'Kamau', 'Mwangi', 'MALE', '1985-03-15', '12345678', 'A001234567B', 1, 1
WHERE @t1 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM individual_transporters WHERE transporter_id = @t1 AND deleted_at IS NULL);

-- Individual transporter 2
INSERT INTO transporters (transporter_no, category, primary_phone, email_address, status, restricted, created_by, updated_by)
SELECT 'TRP-002', 'INDIVIDUAL', '0711000002', 'mary.wanjiku@dev-transport.test', 'ACTIVE', 0, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM transporters WHERE transporter_no = 'TRP-002' AND deleted_at IS NULL);

SET @t2 := (SELECT id FROM transporters WHERE transporter_no = 'TRP-002' AND deleted_at IS NULL LIMIT 1);

INSERT INTO individual_transporters (transporter_id, first_name, last_name, gender, date_of_birth, national_id_no, kra_pin, created_by, updated_by)
SELECT @t2, 'Mary', 'Wanjiku', 'FEMALE', '1990-07-22', '23456789', 'A002345678C', 1, 1
WHERE @t2 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM individual_transporters WHERE transporter_id = @t2 AND deleted_at IS NULL);

-- Company transporter
INSERT INTO transporters (transporter_no, category, primary_phone, email_address, status, restricted, created_by, updated_by)
SELECT 'TRP-C01', 'COMPANY', '0711000003', 'dispatch@devhaul.co.test', 'ACTIVE', 0, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM transporters WHERE transporter_no = 'TRP-C01' AND deleted_at IS NULL);

SET @t3 := (SELECT id FROM transporters WHERE transporter_no = 'TRP-C01' AND deleted_at IS NULL LIMIT 1);

INSERT INTO company_transporters (transporter_id, company_name, registration_no, kra_pin, contact_person_name, contact_person_phone, postal_address, postal_code, town, created_by, updated_by)
SELECT @t3, 'Dev Haul Ltd', 'C.123456', 'P051234567D', 'Peter Ochieng', '0711000099', 'PO Box 100', '00100', 'Nairobi', 1, 1
WHERE @t3 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM company_transporters WHERE transporter_id = @t3 AND deleted_at IS NULL);

-- Route assignments (uses DEV-R01..DEV-R03 from milk journal seed)
INSERT INTO transporter_route_assignments (transporter_id, route_id, start_date, end_date, active, created_by, updated_by)
SELECT @t1, 8038, '2026-01-01', '2026-12-31', 1, 1, 1
WHERE @t1 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_route_assignments WHERE transporter_id = @t1 AND route_id = 8038 AND deleted_at IS NULL);

INSERT INTO transporter_route_assignments (transporter_id, route_id, start_date, end_date, active, created_by, updated_by)
SELECT @t2, 8039, '2026-01-01', '2026-12-31', 1, 1, 1
WHERE @t2 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_route_assignments WHERE transporter_id = @t2 AND route_id = 8039 AND deleted_at IS NULL);

INSERT INTO transporter_route_assignments (transporter_id, route_id, start_date, end_date, active, created_by, updated_by)
SELECT @t3, 8040, '2026-01-01', '2026-12-31', 1, 1, 1
WHERE @t3 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_route_assignments WHERE transporter_id = @t3 AND route_id = 8040 AND deleted_at IS NULL);

-- Vehicles
INSERT INTO transporter_vehicles (transporter_id, route_id, registration_no, vehicle_type, capacity_litres, active, created_by, updated_by)
SELECT @t1, 8038, 'KDA 001A', 'PICKUP', 500.00, 1, 1, 1
WHERE @t1 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_vehicles WHERE registration_no = 'KDA 001A' AND deleted_at IS NULL);

INSERT INTO transporter_vehicles (transporter_id, route_id, registration_no, vehicle_type, capacity_litres, active, created_by, updated_by)
SELECT @t2, 8039, 'KDA 002B', 'VAN', 800.00, 1, 1, 1
WHERE @t2 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_vehicles WHERE registration_no = 'KDA 002B' AND deleted_at IS NULL);

INSERT INTO transporter_vehicles (transporter_id, route_id, registration_no, vehicle_type, capacity_litres, active, created_by, updated_by)
SELECT @t3, 8040, 'KDA 003C', 'TANKER', 5000.00, 1, 1, 1
WHERE @t3 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_vehicles WHERE registration_no = 'KDA 003C' AND deleted_at IS NULL);

-- Drivers
INSERT INTO transporter_drivers (transporter_id, driver_no, first_name, last_name, gender, national_id_no, primary_phone, driving_license_no, driving_license_class, driving_license_expiry, status, created_by, updated_by)
SELECT @t1, 'DRV-001', 'James', 'Mutua', 'MALE', '34567890', '0711000101', 'DL-001234', 'BCE', '2028-12-31', 'ACTIVE', 1, 1
WHERE @t1 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_drivers WHERE driver_no = 'DRV-001' AND deleted_at IS NULL);

INSERT INTO transporter_drivers (transporter_id, driver_no, first_name, last_name, gender, national_id_no, primary_phone, driving_license_no, driving_license_class, driving_license_expiry, status, created_by, updated_by)
SELECT @t3, 'DRV-002', 'Grace', 'Akinyi', 'FEMALE', '45678901', '0711000102', 'DL-002345', 'CE', '2029-06-30', 'ACTIVE', 1, 1
WHERE @t3 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM transporter_drivers WHERE driver_no = 'DRV-002' AND deleted_at IS NULL);
