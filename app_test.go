package main

import (
	"testing"
)

func TestExtractState(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Standard comma", "Sweetwater, TN", "TN"},
		{"Comma with space", "Jonesborough, NC", "NC"},
		{"No state", "Central Station North", ""},
		{"Multiple commas", "Station, Area, GA", "GA"},
		{"Trailing punctuation", "Amherst, OH!", ""}, // Regex looks for exactly [A-Z]{2}$
		{"Lowercase", "Sweetwater, tn", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractState(tt.input)
			if got != tt.expected {
				t.Errorf("extractState(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
