package debuglog

import (
	"encoding/json"
	"os"
	"time"
)

const defaultPath = "/home/reuben/apps/edairy/.cursor/debug-714439.log"

type entry struct {
	SessionID    string         `json:"sessionId"`
	ID           string         `json:"id"`
	Timestamp    int64          `json:"timestamp"`
	Location     string         `json:"location"`
	Message      string         `json:"message"`
	Data         map[string]any `json:"data,omitempty"`
	RunID        string         `json:"runId,omitempty"`
	HypothesisID string         `json:"hypothesisId,omitempty"`
}

func path() string {
	if p := os.Getenv("DEBUG_LOG_PATH"); p != "" {
		return p
	}
	return defaultPath
}

func Log(location, message, hypothesisID, runID string, data map[string]any) {
	e := entry{
		SessionID:    "714439",
		ID:           "log_" + time.Now().Format("150405") + "_" + hypothesisID,
		Timestamp:    time.Now().UnixMilli(),
		Location:     location,
		Message:      message,
		Data:         data,
		RunID:        runID,
		HypothesisID: hypothesisID,
	}
	b, _ := json.Marshal(e)
	b = append(b, '\n')
	if f, err := os.OpenFile(path(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
		_, _ = f.Write(b)
		_ = f.Close()
	}
	_, _ = os.Stderr.Write(b)
}
