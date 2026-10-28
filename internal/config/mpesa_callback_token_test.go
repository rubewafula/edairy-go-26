package config

import "testing"

func TestWithMpesaCallbackToken(t *testing.T) {
	t.Setenv("MPESA_CALLBACK_TOKEN", "")
	if got := WithMpesaCallbackToken("https://api.example/api/mpesa-b2c-result"); got != "https://api.example/api/mpesa-b2c-result" {
		t.Fatalf("no token configured should leave URL unchanged, got %q", got)
	}

	t.Setenv("MPESA_CALLBACK_TOKEN", "s3cr et")
	if got := WithMpesaCallbackToken("https://api.example/api/mpesa-b2c-result"); got != "https://api.example/api/mpesa-b2c-result?token=s3cr+et" {
		t.Fatalf("unexpected URL %q", got)
	}
	if got := WithMpesaCallbackToken("https://api.example/cb?x=1"); got != "https://api.example/cb?x=1&token=s3cr+et" {
		t.Fatalf("unexpected URL with query %q", got)
	}
	if got := WithMpesaCallbackToken("https://api.example/cb?token=abc"); got != "https://api.example/cb?token=abc" {
		t.Fatalf("existing token must not be duplicated, got %q", got)
	}
}

func TestMpesaCallbackTokenValid(t *testing.T) {
	t.Setenv("MPESA_CALLBACK_TOKEN", "")
	if !MpesaCallbackTokenValid("") {
		t.Fatal("callbacks are accepted when no token is configured")
	}
	t.Setenv("MPESA_CALLBACK_TOKEN", "expected")
	if !MpesaCallbackTokenValid("expected") {
		t.Fatal("matching token should be valid")
	}
	if MpesaCallbackTokenValid("") || MpesaCallbackTokenValid("wrong") {
		t.Fatal("missing or wrong token must be rejected")
	}
}
