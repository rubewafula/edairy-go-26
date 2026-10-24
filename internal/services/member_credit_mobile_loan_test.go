package services

import (
	"regexp"
	"strings"
	"testing"

	"github.com/rubewafula/edairy-go-26/internal/models"
)

func ptrU64(v uint64) *uint64 { return &v }

func TestMilkRateResolverFollowsPayrollPrecedence(t *testing.T) {
	r := newMilkRateResolver()
	const period, member, route = uint64(4), uint64(819), uint64(14)

	r.addDefault(nil, nil, 40)
	if got := r.Rate(period, member, route); got != 40 {
		t.Fatalf("global default: got %v, want 40", got)
	}
	r.addDefault(nil, ptrU64(route), 42)
	if got := r.Rate(period, member, route); got != 42 {
		t.Fatalf("route default: got %v, want 42", got)
	}
	r.addSpecial(period, nil, ptrU64(route), 44)
	if got := r.Rate(period, member, route); got != 44 {
		t.Fatalf("route special for period: got %v, want 44", got)
	}
	r.addDefault(ptrU64(member), nil, 46)
	if got := r.Rate(period, member, route); got != 46 {
		t.Fatalf("member default beats route special: got %v, want 46", got)
	}
	r.addSpecial(period, ptrU64(member), nil, 48)
	if got := r.Rate(period, member, route); got != 48 {
		t.Fatalf("member special for period: got %v, want 48", got)
	}
	if got := r.Rate(period+1, member, route); got != 46 {
		t.Fatalf("special rate must not leak into another period: got %v, want 46", got)
	}
	if got := r.Rate(0, 999, 999); got != 40 {
		t.Fatalf("unknown member and route fall back to global: got %v, want 40", got)
	}
}

func TestGenerateOTPIsSixDigits(t *testing.T) {
	re := regexp.MustCompile(`^\d{6}$`)
	for i := 0; i < 200; i++ {
		otp, err := generateOTP()
		if err != nil {
			t.Fatal(err)
		}
		if !re.MatchString(otp) {
			t.Fatalf("OTP %q is not 6 digits", otp)
		}
	}
}

func TestOTPHashBindsCodeToApplication(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	hash, err := hashOTP(42, "123456")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(hash, "123456") {
		t.Fatal("hash must not contain the plain OTP")
	}
	if !otpMatches(42, "123456", hash) {
		t.Fatal("correct OTP should match")
	}
	if !otpMatches(42, " 123456 ", hash) {
		t.Fatal("surrounding whitespace should be ignored")
	}
	if otpMatches(42, "654321", hash) {
		t.Fatal("wrong OTP must not match")
	}
	if otpMatches(43, "123456", hash) {
		t.Fatal("OTP must not be reusable for another application")
	}

	t.Setenv("JWT_SECRET", "")
	if _, err := hashOTP(42, "123456"); err == nil {
		t.Fatal("hashing without JWT_SECRET must fail")
	}
}

func TestParseMpesaB2CResult(t *testing.T) {
	data := map[string]interface{}{
		"Result": map[string]interface{}{
			"ResultCode":               float64(0),
			"ResultDesc":               "The service request is processed successfully.",
			"OriginatorConversationID": "MLOAN-LN-APP-1",
			"ConversationID":           "AG_20261018_1",
			"TransactionID":            "SJI1ABC123",
			"ResultParameters": map[string]interface{}{
				"ResultParameter": []interface{}{
					map[string]interface{}{"Key": "TransactionAmount", "Value": float64(5000)},
					map[string]interface{}{"Key": "TransactionReceipt", "Value": "SJI1ABC123"},
					map[string]interface{}{"Key": "ReceiverPartyPublicName", "Value": "254712345678 - Jane Doe"},
				},
			},
		},
	}
	res := parseMpesaB2CResult(data)
	if res.ResultCode != "0" || res.OriginatorID != "MLOAN-LN-APP-1" || res.TransactionID != "SJI1ABC123" {
		t.Fatalf("unexpected parse result: %+v", res)
	}
	if !res.HasAmount || res.Amount != 5000 {
		t.Fatalf("amount not parsed: %+v", res)
	}

	empty := parseMpesaB2CResult(map[string]interface{}{"Result": map[string]interface{}{"ResultCode": float64(1)}})
	if empty.TransactionID != "" || empty.OriginatorID != "" {
		t.Fatalf("missing fields must stay empty, got %+v", empty)
	}
}

func TestMpesaResultIssueRejectsInconsistentSuccess(t *testing.T) {
	d := models.LoanDisbursement{Amount: 5000}
	d.ID = 7

	if issue := mpesaResultIssue(d, mpesaB2CResult{Amount: 5000, HasAmount: true}); !strings.Contains(issue, "TransactionID") {
		t.Fatalf("missing TransactionID should be flagged, got %q", issue)
	}
	if issue := mpesaResultIssue(d, mpesaB2CResult{TransactionID: "X"}); !strings.Contains(issue, "TransactionAmount") {
		t.Fatalf("missing amount should be flagged, got %q", issue)
	}
	if issue := mpesaResultIssue(d, mpesaB2CResult{TransactionID: "X", Amount: 50000, HasAmount: true}); !strings.Contains(issue, "amount mismatch") {
		t.Fatalf("amount mismatch should be flagged, got %q", issue)
	}
}

func TestProductionPayoutBlockedInDevelopment(t *testing.T) {
	t.Setenv("ENV", "development")
	t.Setenv("MPESA_ALLOW_PRODUCTION_IN_DEV", "")
	if !productionPayoutBlocked("production") {
		t.Fatal("production payouts must be blocked in development")
	}
	if productionPayoutBlocked("sandbox") {
		t.Fatal("sandbox payouts must be allowed in development")
	}
	t.Setenv("MPESA_ALLOW_PRODUCTION_IN_DEV", "true")
	if productionPayoutBlocked("production") {
		t.Fatal("explicit override should allow production payouts")
	}
	t.Setenv("ENV", "production")
	t.Setenv("MPESA_ALLOW_PRODUCTION_IN_DEV", "")
	if productionPayoutBlocked("production") {
		t.Fatal("production environment must not be blocked")
	}
}

func TestMemberCanBorrow(t *testing.T) {
	for _, status := range []string{"pending", "ACTIVE", ""} {
		if !memberCanBorrow(status) {
			t.Fatalf("status %q should be allowed to borrow", status)
		}
	}
	for _, status := range []string{"SUSPENDED", "suspended", "Exited"} {
		if memberCanBorrow(status) {
			t.Fatalf("status %q must not be allowed to borrow", status)
		}
	}
}

func TestFormatKESAndMaskPhone(t *testing.T) {
	cases := map[float64]string{0: "0.00", 950: "950.00", 12500: "12,500.00", 1234567.891: "1,234,567.89", -4500.5: "-4,500.50"}
	for in, want := range cases {
		if got := formatKES(in); got != want {
			t.Fatalf("formatKES(%v) = %q, want %q", in, got, want)
		}
	}
	if got := maskPhone("254712345678"); got != "2547*****678" {
		t.Fatalf("maskPhone = %q", got)
	}
}

func TestUniqueMemberIDs(t *testing.T) {
	got := uniqueMemberIDs([]uint64{3, 0, 3, 5, 5, 1})
	if len(got) != 3 || got[0] != 3 || got[1] != 5 || got[2] != 1 {
		t.Fatalf("unexpected ids %v", got)
	}
}
