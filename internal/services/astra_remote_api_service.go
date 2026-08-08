package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/config"
)

type astraToken struct {
	HeaderValue     string `json:"headerValue"`
	ExpiresEpochSec int64  `json:"expiresEpochSecs"`
}

type AstraRemoteAPIService struct {
	cfg        config.AstraConfig
	httpClient *http.Client
	tokenMu    sync.Mutex
	token      *astraToken
}

func NewAstraRemoteAPIService() *AstraRemoteAPIService {
	return NewAstraRemoteAPIServiceWithConfig(config.LoadAstraConfig())
}

func NewAstraRemoteAPIServiceWithConfig(cfg config.AstraConfig) *AstraRemoteAPIService {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if !cfg.VerifySSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	return &AstraRemoteAPIService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout:   time.Duration(cfg.RequestTimeoutSec) * time.Second,
			Transport: transport,
		},
	}
}

func (s *AstraRemoteAPIService) Configured() bool {
	return s.cfg.Configured()
}

func (s *AstraRemoteAPIService) getToken() (*astraToken, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	now := time.Now().Unix()
	if s.token != nil && s.token.HeaderValue != "" && now < s.token.ExpiresEpochSec-60 {
		return s.token, nil
	}
	if s.token != nil && s.token.HeaderValue != "" && now < s.token.ExpiresEpochSec-60 {
		if refreshed, err := s.refreshTokenLocked(); err == nil && refreshed != nil {
			return refreshed, nil
		}
	}
	tok, err := s.loginTokenLocked()
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func (s *AstraRemoteAPIService) loginTokenLocked() (*astraToken, error) {
	body := map[string]string{
		"identity": s.cfg.AuthIdentity,
		"password": s.cfg.TenantPassword,
	}
	var resp astraToken
	if err := s.authenticateRequest(http.MethodPost, "login", body, &resp); err != nil {
		return nil, err
	}
	if resp.HeaderValue == "" {
		return nil, fmt.Errorf("astra login returned empty token")
	}
	if resp.ExpiresEpochSec == 0 {
		resp.ExpiresEpochSec = time.Now().Add(30 * time.Minute).Unix()
	}
	s.token = &resp
	return s.token, nil
}

func (s *AstraRemoteAPIService) refreshTokenLocked() (*astraToken, error) {
	if s.token == nil || s.token.HeaderValue == "" {
		return s.loginTokenLocked()
	}
	jwt := strings.TrimPrefix(s.token.HeaderValue, "Bearer ")
	body := map[string]string{"jwt": jwt}
	var resp astraToken
	if err := s.authenticateRequest(http.MethodPost, "renew", body, &resp); err != nil {
		return s.loginTokenLocked()
	}
	if resp.HeaderValue == "" {
		return s.loginTokenLocked()
	}
	if resp.ExpiresEpochSec == 0 {
		resp.ExpiresEpochSec = time.Now().Add(30 * time.Minute).Unix()
	}
	s.token = &resp
	return s.token, nil
}

func (s *AstraRemoteAPIService) authenticateRequest(method, endpoint string, payload interface{}, out interface{}) error {
	url := s.cfg.AuthURL + "/" + strings.TrimLeft(endpoint, "/")
	return s.doJSON(method, url, payload, out, nil, "")
}

func (s *AstraRemoteAPIService) Get(endpoint string, query map[string]string) (map[string]interface{}, error) {
	tok, err := s.getToken()
	if err != nil {
		return nil, err
	}
	url := s.cfg.RemoteURL + "/" + strings.TrimLeft(endpoint, "/")
	if len(query) > 0 {
		parts := make([]string, 0, len(query))
		for k, v := range query {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
		url += "?" + strings.Join(parts, "&")
	}
	var out map[string]interface{}
	err = s.doJSON(http.MethodGet, url, nil, &out, nil, tok.HeaderValue)
	return out, err
}

func (s *AstraRemoteAPIService) Post(endpoint string, payload interface{}) (map[string]interface{}, error) {
	tok, err := s.getToken()
	if err != nil {
		return nil, err
	}
	url := s.cfg.RemoteURL + "/" + strings.TrimLeft(endpoint, "/")
	var out map[string]interface{}
	err = s.doJSON(http.MethodPost, url, payload, &out, nil, tok.HeaderValue)
	return out, err
}

func (s *AstraRemoteAPIService) PostWithAuth(endpoint string, payload interface{}, authHeader string) (map[string]interface{}, error) {
	url := s.cfg.RemoteURL + "/" + strings.TrimLeft(endpoint, "/")
	var out map[string]interface{}
	err := s.doJSON(http.MethodPost, url, payload, &out, nil, authHeader)
	return out, err
}

type astraInitiateWithdrawalResult struct {
	StatusCode int
	Headers    http.Header
	Body       map[string]interface{}
}

func (s *AstraRemoteAPIService) InitiateCashTransfer(endpoint string, payload map[string]interface{}) (*astraInitiateWithdrawalResult, error) {
	tok, err := s.getToken()
	if err != nil {
		return nil, err
	}
	url := s.cfg.RemoteURL + "/" + strings.TrimLeft(endpoint, "/")
	reqBody, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tok.HeaderValue)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("astra withdrawal request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var body map[string]interface{}
	_ = json.Unmarshal(raw, &body)
	log.Printf("[ASTRA] initiate withdrawal %s status=%d body=%s", url, resp.StatusCode, string(raw))

	return &astraInitiateWithdrawalResult{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
	}, nil
}

func (s *AstraRemoteAPIService) ExchangeSCAForJWT(intentID, otp string) (string, error) {
	tok, err := s.getToken()
	if err != nil {
		return "", err
	}
	jwt := strings.TrimPrefix(tok.HeaderValue, "Bearer ")
	body := map[string]interface{}{
		"intentId": intentID,
		"otp":      otp,
		"jwt":      jwt,
	}

	var resp astraToken
	if err := s.authenticateRequest(http.MethodPut, "jwt", body, &resp); err != nil {
		return "", err
	}
	if resp.HeaderValue == "" {
		return "", fmt.Errorf("astra SCA exchange returned empty authorization")
	}
	return resp.HeaderValue, nil
}

func (s *AstraRemoteAPIService) doJSON(method, url string, payload interface{}, out interface{}, extraHeaders map[string]string, authHeader string) error {
	var bodyReader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("astra request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("astra HTTP %d: %s", resp.StatusCode, string(raw))
	}
	if out == nil {
		return nil
	}
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func parseSCAIntentID(headers http.Header) string {
	sca := headers.Get("SCA")
	if sca == "" {
		for k, vals := range headers {
			if strings.EqualFold(k, "SCA") && len(vals) > 0 {
				sca = vals[0]
				break
			}
		}
	}
	if sca == "" {
		return ""
	}
	parts := strings.Split(sca, ";")
	if len(parts) == 0 {
		return ""
	}
	kv := strings.SplitN(parts[0], "=", 2)
	if len(kv) != 2 {
		return ""
	}
	return strings.TrimSpace(kv[1])
}
