package services

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

type jengaAuthResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	IssuedAt    string `json:"issuedAt"`
}

type JengaFinserveService struct {
	cfg        config.JengaConfig
	httpClient *http.Client
	tokenMu    sync.Mutex
	token      string
	tokenExp   time.Time
}

func NewJengaFinserveService() *JengaFinserveService {
	return NewJengaFinserveServiceWithConfig(config.LoadJengaConfig())
}

func NewJengaFinserveServiceWithConfig(cfg config.JengaConfig) *JengaFinserveService {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if !cfg.VerifySSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	return &JengaFinserveService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout:   time.Duration(cfg.RequestTimeoutSec) * time.Second,
			Transport: transport,
		},
	}
}

func (s *JengaFinserveService) Configured() bool {
	return s.cfg.Configured()
}

func (s *JengaFinserveService) getAccessToken() (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()
	if s.token != "" && time.Now().Before(s.tokenExp.Add(-60*time.Second)) {
		return s.token, nil
	}
	url := s.cfg.BaseURL + "/authentication/api/v3/authenticate/merchant"
	body := map[string]string{
		"merchantCode":    s.cfg.MerchantCode,
		"consumerSecret":  s.cfg.ConsumerSecret,
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", s.cfg.APIKey)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("jenga auth failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("jenga auth HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var out jengaAuthResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", err
	}
	if out.AccessToken == "" {
		return "", fmt.Errorf("jenga auth returned empty token")
	}
	s.token = out.AccessToken
	expires := out.ExpiresIn
	if expires <= 0 {
		expires = 3600
	}
	s.tokenExp = time.Now().Add(time.Duration(expires) * time.Second)
	return s.token, nil
}

func (s *JengaFinserveService) loadPrivateKey() (*rsa.PrivateKey, error) {
	keyPEM, err := os.ReadFile(s.cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read jenga private key: %w", err)
	}
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("invalid jenga private key PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse jenga private key: %w", err)
		}
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("jenga private key is not RSA")
	}
	return rsaKey, nil
}

func (s *JengaFinserveService) signMobileTransfer(amount, currency, reference, accountNumber string) (string, error) {
	priv, err := s.loadPrivateKey()
	if err != nil {
		return "", err
	}
	signStr := amount + currency + reference + accountNumber
	hash := sha256.Sum256([]byte(signStr))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("jenga sign: %w", err)
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

type JengaSendMobileRequest struct {
	Source      jengaSource      `json:"source"`
	Destination jengaDestination `json:"destination"`
	Transfer    jengaTransfer    `json:"transfer"`
}

type jengaSource struct {
	CountryCode   string `json:"countryCode"`
	Name          string `json:"name"`
	AccountNumber string `json:"accountNumber"`
}

type jengaDestination struct {
	Type         string `json:"type"`
	CountryCode  string `json:"countryCode"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
	WalletName   string `json:"walletName"`
}

type jengaTransfer struct {
	Type         string `json:"type"`
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
	Reference    string `json:"reference"`
	Date         string `json:"date"`
	Description  string `json:"description"`
	CallbackURL  string `json:"callbackUrl"`
}

type JengaSendMobileResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TransactionReference string `json:"transactionReference"`
	} `json:"data"`
}

func (s *JengaFinserveService) SendMobileEquitel(recipientName, phone254, amountStr, reference, description string) (*JengaSendMobileResponse, error) {
	if !s.Configured() {
		return nil, fmt.Errorf("Jenga/Equitel is not configured")
	}
	token, err := s.getAccessToken()
	if err != nil {
		return nil, err
	}
	signature, err := s.signMobileTransfer(amountStr, "KES", reference, s.cfg.SourceAccount)
	if err != nil {
		return nil, err
	}

	sourceName := s.cfg.SourceAccountName
	if sourceName == "" {
		sourceName = "EDAIRY"
	}
	body := JengaSendMobileRequest{
		Source: jengaSource{
			CountryCode:   s.cfg.SourceCountryCode,
			Name:          sourceName,
			AccountNumber: s.cfg.SourceAccount,
		},
		Destination: jengaDestination{
			Type:         "mobile",
			CountryCode:  "KE",
			Name:         recipientName,
			MobileNumber: phone254,
			WalletName:   "Equitel",
		},
		Transfer: jengaTransfer{
			Type:         "MobileWallet",
			Amount:       amountStr,
			CurrencyCode: "KES",
			Reference:    reference,
			Date:         time.Now().Format("2006-01-02"),
			Description:  description,
			CallbackURL:  s.cfg.CallbackURL,
		},
	}

	b, _ := json.Marshal(body)
	url := s.cfg.BaseURL + "/v3-apis/transaction-api/v3.0/remittance/sendmobile"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("signature", signature)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jenga sendmobile failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	log.Printf("[JENGA] sendmobile status=%d body=%s", resp.StatusCode, string(raw))
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("jenga sendmobile HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var out JengaSendMobileResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if strings.EqualFold(out.Status, "false") || (out.Code != "" && out.Code != "0") {
		msg := out.Message
		if msg == "" {
			msg = "jenga sendmobile rejected"
		}
		return nil, fmt.Errorf(msg)
	}
	return &out, nil
}

func parseJengaCallback(data map[string]interface{}) (reference, txnRef, responseCode string, success bool) {
	txnRef = fmt.Sprintf("%v", data["transactionReference"])
	if nested, ok := data["data"].(map[string]interface{}); ok {
		responseCode = fmt.Sprintf("%v", nested["ResponseCode"])
		if ref := fmt.Sprintf("%v", nested["Reference"]); ref != "" && ref != "<nil>" {
			reference = ref
		}
	}
	if responseCode == "" {
		responseCode = fmt.Sprintf("%v", data["code"])
	}
	success = responseCode == "0" || strings.EqualFold(fmt.Sprintf("%v", data["status"]), "true")
	return reference, txnRef, responseCode, success
}
