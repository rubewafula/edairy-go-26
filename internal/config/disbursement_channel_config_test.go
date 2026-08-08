package config

import "testing"

func TestMaskSecretValue(t *testing.T) {
	if got := MaskSecretValue(""); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
	if got := MaskSecretValue("abc"); got != "****" {
		t.Fatalf("expected ****, got %q", got)
	}
	if got := MaskSecretValue("secret1234"); got != "******1234" {
		t.Fatalf("expected ******1234, got %q", got)
	}
}

func TestBuildMpesaConfigFromResolvedOverridesEnv(t *testing.T) {
	values := map[string]string{
		"env":            "sandbox",
		"consumer_key":   "db-key",
		"consumer_secret": "db-secret",
		"b2c_shortcode":  "600000",
	}
	cfg := BuildMpesaConfig(values)
	if cfg.ConsumerKey != "db-key" {
		t.Fatalf("consumer_key = %q", cfg.ConsumerKey)
	}
	if cfg.B2CShortcode != "600000" {
		t.Fatalf("b2c_shortcode = %q", cfg.B2CShortcode)
	}
	if cfg.OAuthURL == "" || cfg.B2CURL == "" {
		t.Fatal("expected default sandbox URLs")
	}
}

func TestGetProviderConfigTemplate(t *testing.T) {
	mpesa := GetProviderConfigTemplate("MPESA")
	if len(mpesa) == 0 {
		t.Fatal("expected MPESA template keys")
	}
	if GetProviderConfigTemplate("UNKNOWN") != nil {
		t.Fatal("expected nil for unknown provider")
	}
}

func TestBuildAirtelConfigDefaults(t *testing.T) {
	cfg := BuildAirtelConfig(map[string]string{"env": "sandbox"})
	if cfg.Country != "KE" || cfg.Currency != "KES" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.BaseURL == "" {
		t.Fatal("expected default base URL")
	}
}
