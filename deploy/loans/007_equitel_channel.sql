-- Equitel mobile disbursement channel (Equity Finserve/Jenga)
INSERT INTO disbursement_channels
  (channel_code, channel_name, description, channel_type, provider, is_async, requires_external_ref, gl_rule_type, sort_order, is_active, config_schema, created_at, updated_at)
VALUES
  ('MOBILE_EQUITEL', 'Mobile - Equitel', 'Transfer to member Equitel wallet via Equity Finserve', 'MOBILE', 'EQUITY_JENGA', 1, 1, 'LOAN_DISBURSEMENT', 45, 1,
   JSON_OBJECT('required', JSON_ARRAY('phone_number')), NOW(), NOW())
ON DUPLICATE KEY UPDATE
  channel_name = VALUES(channel_name),
  description = VALUES(description),
  provider = VALUES(provider),
  is_async = VALUES(is_async),
  requires_external_ref = VALUES(requires_external_ref),
  config_schema = VALUES(config_schema),
  is_active = VALUES(is_active),
  updated_at = NOW();
