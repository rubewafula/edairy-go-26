package services

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/rubewafula/edairy-go-26/internal/config"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type DisbursementChannelConfigService struct{}

func NewDisbursementChannelConfigService() *DisbursementChannelConfigService {
	return &DisbursementChannelConfigService{}
}

func (s *DisbursementChannelConfigService) loadChannel(channelID uint64) (*models.DisbursementChannel, error) {
	var ch models.DisbursementChannel
	if err := db.DB.First(&ch, channelID).Error; err != nil {
		return nil, fmt.Errorf("disbursement channel not found")
	}
	return &ch, nil
}

func (s *DisbursementChannelConfigService) loadDBValues(channelID uint64) (map[string]string, error) {
	var rows []models.DisbursementChannelConfig
	if err := db.DB.Where("disbursement_channel_id = ?", channelID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[row.ConfigKey] = row.ConfigValue
	}
	return out, nil
}

func (s *DisbursementChannelConfigService) GetResolvedValues(channelID uint64) (map[string]string, error) {
	ch, err := s.loadChannel(channelID)
	if err != nil {
		return nil, err
	}
	return s.resolveForProvider(ch.ChannelCode, ch.Provider, channelID)
}

func (s *DisbursementChannelConfigService) resolveForProvider(channelCode, provider string, channelID uint64) (map[string]string, error) {
	template := config.GetProviderConfigTemplateForChannel(channelCode, provider)
	dbValues, err := s.loadDBValues(channelID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(template))
	for _, def := range template {
		if v := strings.TrimSpace(dbValues[def.Key]); v != "" {
			out[def.Key] = v
			continue
		}
		if def.EnvVar != "" {
			if v := strings.TrimSpace(os.Getenv(def.EnvVar)); v != "" {
				out[def.Key] = v
			}
		}
	}
	return out, nil
}

func (s *DisbursementChannelConfigService) ListForAdmin(channelID uint64) ([]dtos.DisbursementChannelConfigItem, error) {
	ch, err := s.loadChannel(channelID)
	if err != nil {
		return nil, err
	}
	template := config.GetProviderConfigTemplateForChannel(ch.ChannelCode, ch.Provider)
	dbValues, err := s.loadDBValues(channelID)
	if err != nil {
		return nil, err
	}

	items := make([]dtos.DisbursementChannelConfigItem, 0, len(template))
	for _, def := range template {
		item := dtos.DisbursementChannelConfigItem{
			Key:       def.Key,
			Label:     def.Label,
			InputType: def.InputType,
			IsSecret:  def.IsSecret,
			Source:    config.ConfigSourceUnset,
		}
		if dbVal := strings.TrimSpace(dbValues[def.Key]); dbVal != "" {
			item.Source = config.ConfigSourceDB
			item.HasValue = true
			if def.IsSecret {
				item.Value = config.MaskSecretValue(dbVal)
			} else {
				item.Value = dbVal
			}
		} else if def.EnvVar != "" {
			if envVal := strings.TrimSpace(os.Getenv(def.EnvVar)); envVal != "" {
				item.Source = config.ConfigSourceEnv
				item.HasValue = true
				if def.IsSecret {
					item.Value = config.MaskSecretValue(envVal)
				} else {
					item.Value = envVal
				}
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *DisbursementChannelConfigService) Upsert(channelID uint64, items []dtos.DisbursementChannelConfigUpdateItem) error {
	ch, err := s.loadChannel(channelID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("no config items provided")
	}

	for _, item := range items {
		key := strings.ToLower(strings.TrimSpace(item.Key))
		if key == "" {
			continue
		}
		def, ok := config.LookupConfigKeyDefForChannel(ch.ChannelCode, ch.Provider, key)
		if !ok {
			return fmt.Errorf("unknown config key %q for provider %q", key, ch.Provider)
		}

		value := strings.TrimSpace(item.Value)
		if value == "" {
			if err := db.DB.Where("disbursement_channel_id = ? AND config_key = ?", channelID, def.Key).
				Delete(&models.DisbursementChannelConfig{}).Error; err != nil {
				return err
			}
			continue
		}

		if def.IsSecret && strings.Contains(value, "*") {
			continue
		}

		var existing models.DisbursementChannelConfig
		err := db.DB.Where("disbursement_channel_id = ? AND config_key = ?", channelID, def.Key).First(&existing).Error
		if err != nil {
			row := models.DisbursementChannelConfig{
				DisbursementChannelID: channelID,
				ConfigKey:             def.Key,
				ConfigValue:           value,
				IsSecret:              def.IsSecret,
			}
			if err := db.DB.Create(&row).Error; err != nil {
				return err
			}
			continue
		}
		existing.ConfigValue = value
		existing.IsSecret = def.IsSecret
		if err := db.DB.Save(&existing).Error; err != nil {
			return err
		}
	}

	InvalidateDisbursementProviderClients(channelID)
	return nil
}

func (s *DisbursementChannelConfigService) BuildMpesaConfig(channelID uint64) (config.MpesaConfig, error) {
	values, err := s.GetResolvedValues(channelID)
	if err != nil {
		return config.MpesaConfig{}, err
	}
	return config.BuildMpesaConfig(values), nil
}

func (s *DisbursementChannelConfigService) BuildAirtelConfig(channelID uint64) (config.AirtelConfig, error) {
	values, err := s.GetResolvedValues(channelID)
	if err != nil {
		return config.AirtelConfig{}, err
	}
	return config.BuildAirtelConfig(values), nil
}

func (s *DisbursementChannelConfigService) BuildJengaConfig(channelID uint64) (config.JengaConfig, error) {
	values, err := s.GetResolvedValues(channelID)
	if err != nil {
		return config.JengaConfig{}, err
	}
	return config.BuildJengaConfig(values), nil
}

func (s *DisbursementChannelConfigService) BuildAstraConfig(channelID uint64) (config.AstraConfig, error) {
	values, err := s.GetResolvedValues(channelID)
	if err != nil {
		return config.AstraConfig{}, err
	}
	return config.BuildAstraConfig(values), nil
}

// DisbursementProviderClients caches provider HTTP clients per channel.
type DisbursementProviderClients struct {
	configSvc *DisbursementChannelConfigService
	mu        sync.Mutex
	mpesa     map[uint64]*MpesaDarajaService
	airtel    map[uint64]*AirtelDisbursementAPI
	jenga     map[uint64]*JengaFinserveService
	astra     map[uint64]*AstraRemoteAPIService
}

var globalProviderClients *DisbursementProviderClients
var providerClientsOnce sync.Once

func GetDisbursementProviderClients() *DisbursementProviderClients {
	providerClientsOnce.Do(func() {
		globalProviderClients = &DisbursementProviderClients{
			configSvc: NewDisbursementChannelConfigService(),
			mpesa:     make(map[uint64]*MpesaDarajaService),
			airtel:    make(map[uint64]*AirtelDisbursementAPI),
			jenga:     make(map[uint64]*JengaFinserveService),
			astra:     make(map[uint64]*AstraRemoteAPIService),
		}
	})
	return globalProviderClients
}

func InvalidateDisbursementProviderClients(channelID uint64) {
	clients := GetDisbursementProviderClients()
	clients.mu.Lock()
	defer clients.mu.Unlock()
	delete(clients.mpesa, channelID)
	delete(clients.airtel, channelID)
	delete(clients.jenga, channelID)
	delete(clients.astra, channelID)
}

func (p *DisbursementProviderClients) Mpesa(channelID uint64) (*MpesaDarajaService, error) {
	p.mu.Lock()
	if svc, ok := p.mpesa[channelID]; ok {
		p.mu.Unlock()
		return svc, nil
	}
	p.mu.Unlock()

	cfg, err := p.configSvc.BuildMpesaConfig(channelID)
	if err != nil {
		return nil, err
	}
	svc := NewMpesaDarajaServiceWithConfig(cfg)

	p.mu.Lock()
	p.mpesa[channelID] = svc
	p.mu.Unlock()
	return svc, nil
}

func (p *DisbursementProviderClients) Airtel(channelID uint64) (*AirtelDisbursementAPI, error) {
	p.mu.Lock()
	if svc, ok := p.airtel[channelID]; ok {
		p.mu.Unlock()
		return svc, nil
	}
	p.mu.Unlock()

	cfg, err := p.configSvc.BuildAirtelConfig(channelID)
	if err != nil {
		return nil, err
	}
	svc := NewAirtelDisbursementAPIWithConfig(cfg)

	p.mu.Lock()
	p.airtel[channelID] = svc
	p.mu.Unlock()
	return svc, nil
}

func (p *DisbursementProviderClients) Jenga(channelID uint64) (*JengaFinserveService, error) {
	p.mu.Lock()
	if svc, ok := p.jenga[channelID]; ok {
		p.mu.Unlock()
		return svc, nil
	}
	p.mu.Unlock()

	cfg, err := p.configSvc.BuildJengaConfig(channelID)
	if err != nil {
		return nil, err
	}
	svc := NewJengaFinserveServiceWithConfig(cfg)

	p.mu.Lock()
	p.jenga[channelID] = svc
	p.mu.Unlock()
	return svc, nil
}

func (p *DisbursementProviderClients) Astra(channelID uint64) (*AstraRemoteAPIService, error) {
	p.mu.Lock()
	if svc, ok := p.astra[channelID]; ok {
		p.mu.Unlock()
		return svc, nil
	}
	p.mu.Unlock()

	cfg, err := p.configSvc.BuildAstraConfig(channelID)
	if err != nil {
		return nil, err
	}
	svc := NewAstraRemoteAPIServiceWithConfig(cfg)

	p.mu.Lock()
	p.astra[channelID] = svc
	p.mu.Unlock()
	return svc, nil
}
