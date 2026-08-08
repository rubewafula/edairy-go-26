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

type mpesaOAuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

type MpesaDarajaService struct {
	cfg        config.MpesaConfig
	httpClient *http.Client
	tokenMu    sync.Mutex
	token      string
	tokenExp   time.Time
}

func NewMpesaDarajaService() *MpesaDarajaService {
	return NewMpesaDarajaServiceWithConfig(config.LoadMpesaConfig())
}

func NewMpesaDarajaServiceWithConfig(cfg config.MpesaConfig) *MpesaDarajaService {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if !cfg.VerifySSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	return &MpesaDarajaService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout:   time.Duration(cfg.RequestTimeoutSec) * time.Second,
			Transport: transport,
		},
	}
}

func (s *MpesaDarajaService) Configured() bool {
	return s.cfg.Configured()
}

func (s *MpesaDarajaService) getAccessToken() (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()
	if s.token != "" && time.Now().Before(s.tokenExp.Add(-60*time.Second)) {
		return s.token, nil
	}
	req, err := http.NewRequest(http.MethodGet, s.cfg.OAuthURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(s.cfg.ConsumerKey, s.cfg.ConsumerSecret)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("mpesa oauth failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("mpesa oauth HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var tok mpesaOAuthToken
	if err := json.Unmarshal(raw, &tok); err != nil {
		return "", err
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("mpesa oauth returned empty token")
	}
	s.token = tok.AccessToken
	s.tokenExp = time.Now().Add(3500 * time.Second)
	return s.token, nil
}

func (s *MpesaDarajaService) encryptSecurityCredential(password string) (string, error) {
	certPEM, err := os.ReadFile(s.cfg.CertPath)
	if err != nil {
		return "", fmt.Errorf("read mpesa cert: %w", err)
	}
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return "", fmt.Errorf("invalid mpesa certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse mpesa cert: %w", err)
	}
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("mpesa cert public key is not RSA")
	}
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(password))
	if err != nil {
		return "", fmt.Errorf("encrypt security credential: %w", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

type MpesaB2CRequest struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	InitiatorName            string `json:"InitiatorName"`
	SecurityCredential       string `json:"SecurityCredential"`
	CommandID                string `json:"CommandID"`
	Amount                   string `json:"Amount"`
	PartyA                   string `json:"PartyA"`
	PartyB                   string `json:"PartyB"`
	Remarks                  string `json:"Remarks"`
	QueueTimeOutURL          string `json:"QueueTimeOutURL"`
	ResultURL                string `json:"ResultURL"`
	Occassion                string `json:"Occassion"`
}

type MpesaB2CResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

func (s *MpesaDarajaService) InitiateB2C(originatorID, partyB string, amount int, remarks string) (*MpesaB2CResponse, error) {
	if !s.Configured() {
		return nil, fmt.Errorf("M-Pesa Daraja is not configured")
	}
	token, err := s.getAccessToken()
	if err != nil {
		return nil, err
	}
	secCred, err := s.encryptSecurityCredential(s.cfg.InitiatorPassword)
	if err != nil {
		return nil, err
	}
	body := MpesaB2CRequest{
		OriginatorConversationID: originatorID,
		InitiatorName:            s.cfg.InitiatorName,
		SecurityCredential:       secCred,
		CommandID:                s.cfg.B2CCommandID,
		Amount:                   fmt.Sprintf("%d", amount),
		PartyA:                   s.cfg.B2CShortcode,
		PartyB:                   partyB,
		Remarks:                  remarks,
		QueueTimeOutURL:          s.cfg.B2CTimeoutURL,
		ResultURL:                s.cfg.B2CResultURL,
		Occassion:                remarks,
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, s.cfg.B2CURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mpesa b2c request failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	log.Printf("[MPESA] B2C response status=%d body=%s", resp.StatusCode, string(raw))
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("mpesa b2c HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var out MpesaB2CResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ResponseCode != "" && out.ResponseCode != "0" {
		return nil, fmt.Errorf("mpesa b2c rejected: %s", out.ResponseDescription)
	}
	return &out, nil
}

func (s *MpesaDarajaService) QueryTransactionStatus(originatorID, transactionID string) (map[string]interface{}, error) {
	token, err := s.getAccessToken()
	if err != nil {
		return nil, err
	}
	secCred, err := s.encryptSecurityCredential(s.cfg.InitiatorPassword)
	if err != nil {
		return nil, err
	}
	payload := map[string]string{
		"Initiator":              s.cfg.InitiatorName,
		"SecurityCredential":   secCred,
		"CommandID":            "TransactionStatusQuery",
		"TransactionID":        transactionID,
		"OriginatorConversationID": originatorID,
		"PartyA":               s.cfg.B2CShortcode,
		"IdentifierType":       "4",
		"ResultURL":            s.cfg.B2CResultURL,
		"QueueTimeOutURL":      s.cfg.B2CTimeoutURL,
		"Remarks":              "Loan disbursement status query",
		"Occasion":             "Query",
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, s.cfg.StatusQueryURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var out map[string]interface{}
	_ = json.Unmarshal(raw, &out)
	if resp.StatusCode >= 400 {
		return out, fmt.Errorf("mpesa status query HTTP %d: %s", resp.StatusCode, string(raw))
	}
	return out, nil
}

func parseMpesaResultCallback(data map[string]interface{}) (originatorID, transactionID, resultCode, resultDesc string) {
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", "", "", ""
	}
	originatorID = fmt.Sprintf("%v", result["OriginatorConversationID"])
	transactionID = fmt.Sprintf("%v", result["TransactionID"])
	resultCode = fmt.Sprintf("%v", result["ResultCode"])
	resultDesc = fmt.Sprintf("%v", result["ResultDesc"])
	return originatorID, transactionID, resultCode, strings.TrimSpace(resultDesc)
}
