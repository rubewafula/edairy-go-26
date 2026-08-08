-- Per-channel provider configuration (EAV). Values override .env when set.
CREATE TABLE IF NOT EXISTS disbursement_channel_config (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  disbursement_channel_id BIGINT UNSIGNED NOT NULL,
  config_key VARCHAR(128) NOT NULL,
  config_value TEXT NULL,
  is_secret TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  UNIQUE KEY uk_channel_config (disbursement_channel_id, config_key),
  KEY idx_channel_config_channel (disbursement_channel_id),
  CONSTRAINT fk_channel_config_channel
    FOREIGN KEY (disbursement_channel_id) REFERENCES disbursement_channels(id)
);
