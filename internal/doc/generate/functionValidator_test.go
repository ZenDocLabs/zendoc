package generate

import (
	"testing"
)

func TestIsPrivateFunction(t *testing.T) {
	tests := []struct {
		name     string
		function string
		expected bool
	}{
		{"Empty string", "", false},
		{"Uppercase start", "PublicFunction", false},
		{"Lowercase start", "privateFunction", true},
		{"Single character uppercase", "A", false},
		{"Single character lowercase", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrivateFunction(tt.function)
			if result != tt.expected {
				t.Errorf("IsPrivateFunction(%v) = %v; want %v", tt.function, result, tt.expected)
			}
		})
	}
}
