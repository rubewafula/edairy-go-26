package utils

import "testing"

func TestNormalizeKenyaPhone(t *testing.T) {
	tests := []struct {
		override string
		fallback string
		want     string
		wantErr  bool
	}{
		{"0712345678", "", "254712345678", false},
		{"254712345678", "", "254712345678", false},
		{"712345678", "", "254712345678", false},
		{"", "0712345678", "254712345678", false},
		{"", "", "", true},
		{"123", "", "", true},
	}
	for _, tc := range tests {
		got, err := NormalizeKenyaPhone(tc.override, tc.fallback)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected error for override=%q fallback=%q", tc.override, tc.fallback)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != tc.want {
			t.Fatalf("NormalizeKenyaPhone(%q, %q) = %q, want %q", tc.override, tc.fallback, got, tc.want)
		}
	}
}
