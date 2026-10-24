package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// appLogDir is api\logs: APP_LOG_DIR, else the folder next to bin\ that holds server.exe, else ./logs.
func appLogDir() string {
	if dir := strings.TrimSpace(os.Getenv("APP_LOG_DIR")); dir != "" {
		return dir
	}
	if exe, err := os.Executable(); err == nil {
		binDir := filepath.Dir(exe)
		if strings.EqualFold(filepath.Base(binDir), "bin") {
			return filepath.Join(filepath.Dir(binDir), "logs")
		}
	}
	return "logs"
}

func appendAppLog(file, line string) {
	dir := appLogDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(filepath.Join(dir, file), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s | %s\n", time.Now().Format("2006-01-02 15:04:05"), line)
}

// writeCreditDiagnosticLog records every credit recalculation (logs\credit-changes.log) with the
// numbers behind it, so "the limit did not change" can be traced to a reason.
func writeCreditDiagnosticLog(memberID uint64, change CreditChange, b *MemberCreditBreakdown,
	computed float64, prevComputed *float64, delta, oldLimit, newLimit float64) {
	prev := "none (first run: baseline only)"
	if prevComputed != nil {
		prev = fmt.Sprintf("%.2f", *prevComputed)
	}
	why := ""
	switch {
	case b.ZeroRateKilos > 0:
		why = fmt.Sprintf(" | WARNING %.2f kg have NO milk rate (rate 0) - set a default/route milk rate", b.ZeroRateKilos)
	case delta == 0 && strings.HasPrefix(change.Action, "MILK"):
		why = " | no change: milk on confirmed journals unchanged since last run, or the journal date is in an approved payroll period"
	}
	appendAppLog("credit-changes.log", fmt.Sprintf(
		"member %d | %s %s#%d %s | confirmed unpaid milk %.2f kg = KES %.2f, rejects %.2f, loans %.2f, deductions %.2f | formula %.2f (prev %s) | delta %+.2f | limit %.2f -> %.2f%s",
		memberID, change.Action, change.SourceType, change.SourceID, change.Reference,
		b.MilkKilos, b.MilkValueUnpaid, b.RejectValue, b.LoanObligations, b.OtherDeductions,
		computed, prev, delta, oldLimit, newLimit, why))
}

// LogAPIStart writes which server.exe started (logs\api-start.log).
func LogAPIStart() {
	exe, _ := os.Executable()
	built := ""
	if st, err := os.Stat(exe); err == nil {
		built = st.ModTime().Format("2006-01-02 15:04:05")
	}
	appendAppLog("api-start.log", fmt.Sprintf("API started: %s (built %s), credit = confirmed milk", exe, built))
	log.Printf("[API] logs in %s", appLogDir())
}
