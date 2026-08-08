package services

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/config"
)

type airtelOAuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AirtelDisbursementAPI struct {
	cfg        config.AirtelConfig
	httpClient *http.Client
	tokenMu    sync.Mutex
	token      string
	tokenExp   time.Time
}

func NewAirtelDisbursementAPI() *AirtelDisbursementAPI {
	return NewAirtelDisbursementAPIWithConfig(config.LoadAirtelConfig())
}

func NewAirtelDisbursementAPIWithConfig(cfg config.AirtelConfig) *AirtelDisbursementAPI {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if !cfg.VerifySSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	return &AirtelDisbursementAPI{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout:   time.Duration(cfg.RequestTimeoutSec) * time.Second,
			Transport: transport,
		},
	}
}

func (s *AirtelDisbursementAPI) Configured() bool {
	return s.cfg.Configured()
}

func (s *AirtelDisbursementAPI) getAccessToken() (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()
	if s.token != "" && time.Now().Before(s.tokenExp.Add(-30*time.Second)) {
		return s.token, nil
	}
	url := s.cfg.BaseURL + "/auth/oauth2/token"
	body := map[string]string{
		"client_id":     s.cfg.ClientID,
		"client_secret": s.cfg.ClientSecret,
		"grant_type":    "client_credentials",
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("airtel oauth failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("airtel oauth HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var tok airtelOAuthToken
	if err := json.Unmarshal(raw, &tok); err != nil {
		return "", err
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("airtel oauth returned empty token")
	}
	s.token = tok.AccessToken
	expires := tok.ExpiresIn
	if expires <= 0 {
		expires = 170
	}
	s.tokenExp = time.Now().Add(time.Duration(expires) * time.Second)
	return s.token, nil
}

func (s *AirtelDisbursementAPI) loadPublicKey() (*rsa.PublicKey, error) {
	keyPEM, err := os.ReadFile(s.cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read airtel public key: %w", err)
	}
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		// try as raw base64 cert
		cert, err := x509.ParseCertificate(keyPEM)
		if err == nil {
			pub, ok := cert.PublicKey.(*rsa.PublicKey)
			if ok {
				return pub, nil
			}
		}
		return nil, fmt.Errorf("invalid airtel public key PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		cert, err2 := x509.ParseCertificate(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse airtel public key: %w", err)
		}
		pub = cert.PublicKey
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("airtel public key is not RSA")
	}
	return rsaPub, nil
}

func (s *AirtelDisbursementAPI) encryptPIN(pin string) (string, error) {
	pub, err := s.loadPublicKey()
	if err != nil {
		return "", err
	}
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(pin))
	if err != nil {
		return "", fmt.Errorf("encrypt airtel PIN: %w", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

type AirtelDisbursementRequest struct {
	Payee struct {
		Msisdn string `json:"msisdn"`
	} `json:"payee"`
	Reference string `json:"reference"`
	Pin       string `json:"pin"`
	Transaction struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
		ID       string `json:"id"`
	} `json:"transaction"`
}

type AirtelDisbursementResponse struct {
	Data struct {
		Transaction struct {
			ID         string `json:"id"`
			Status     string `json:"status"`
			AirtelMoneyID string `json:"airtel_money_id"`
		} `json:"transaction"`
	} `json:"data"`
	Status struct {
		Code        string `json:"code"`
		Message     string `json:"message"`
		Description string `json:"description"`
	} `json:"status"`
}

func (s *AirtelDisbursementAPI) airtelMSISDN(phone254 string) string {
	if strings.HasPrefix(phone254, "254") && len(phone254) == 12 {
		return phone254[3:]
	}
	return phone254
}

func (s *AirtelDisbursementAPI) InitiateDisbursement(phone254, clientReference string, amount int) (*AirtelDisbursementResponse, error) {
	if !s.Configured() {
		return nil, fmt.Errorf("Airtel disbursement is not configured")
	}
	token, err := s.getAccessToken()
	if err != nil {
		return nil, err
	}
	encPIN, err := s.encryptPIN(s.cfg.DisbursementPIN)
	if err != nil {
		return nil, err
	}

	reqBody := AirtelDisbursementRequest{
		Reference: clientReference,
		Pin:       encPIN,
	}
	reqBody.Payee.Msisdn = s.airtelMSISDN(phone254)
	reqBody.Transaction.Amount = amount
	reqBody.Transaction.Currency = s.cfg.Currency
	reqBody.Transaction.ID = clientReference

	b, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/standard/v2/disbursements/", strings.TrimRight(s.cfg.BaseURL, "/"))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Country", s.cfg.Country)
	req.Header.Set("X-Currency", s.cfg.Currency)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("airtel disbursement failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	log.Printf("[AIRTEL] disbursement status=%d body=%s", resp.StatusCode, string(raw))
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("airtel disbursement HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var out AirtelDisbursementResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.Status.Code != "" && out.Status.Code != "200" && out.Status.Code != "TS" {
		msg := out.Status.Message
		if msg == "" {
			msg = out.Status.Description
		}
		return nil, fmt.Errorf("airtel disbursement rejected: %s", msg)
	}
	return &out, nil
}

func parseAirtelCallback(data map[string]interface{}) (txnID, statusCode, airtelMoneyID string, success bool) {
	txn, ok := data["transaction"].(map[string]interface{})
	if !ok {
		return "", "", "", false
	}
	txnID = fmt.Sprintf("%v", txn["id"])
	statusCode = strings.ToUpper(fmt.Sprintf("%v", txn["status_code"]))
	airtelMoneyID = fmt.Sprintf("%v", txn["airtel_money_id"])
	success = statusCode == "TS" || statusCode == "SUCCESS" || statusCode == "SUCCESSFUL"
	return txnID, statusCode, airtelMoneyID, success
}
