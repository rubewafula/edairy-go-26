package services

import (
	"encoding/json"
	"strings"
)

func decodeChannelPayload(raw []byte) map[string]interface{} {
	if len(raw) == 0 {
		return map[string]interface{}{}
	}
	var out map[string]interface{}
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return map[string]interface{}{}
	}
	return out
}

func mergeChannelPayload(existing []byte, patch map[string]interface{}) []byte {
	base := decodeChannelPayload(existing)
	for k, v := range patch {
		base[k] = v
	}
	b, _ := json.Marshal(base)
	return b
}

func payloadStringFromJSON(payload map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := payload[key]; ok {
			if str, ok := v.(string); ok && strings.TrimSpace(str) != "" {
				return strings.TrimSpace(str)
			}
		}
	}
	return ""
}
