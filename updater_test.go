package main

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   bool
	}{
		{"v2.0.1", "v2.0.0", true},
		{"v2.0.0", "v2.0.1", false},
		{"v2.0.0", "v2.0.0", false},
		{"v1.9.9", "v2.0.0", false},
		{"v2.1.0", "v2.0.9", true},
		{"3.0.0", "v2.0.0", true}, // Test without 'v' prefix
		// New test cases for semantic version parsing
		{"v2.0.10", "v2.0.9", true},  // Multi-digit patch (was broken before)
		{"v2.0.9", "v2.0.10", false}, // Reverse check
		{"v10.0.0", "v9.0.0", true},  // Multi-digit major
		{"v1.10.0", "v1.9.0", true},  // Multi-digit minor
		{"v1.0", "v1.0.0", false},    // Missing patch = 0
		{"v2.0", "v1.9.9", true},     // Short version comparison
	}

	for _, tt := range tests {
		got := CompareVersions(tt.v1, tt.v2)
		if got != tt.want {
			t.Errorf("CompareVersions(%q, %q) = %v; want %v", tt.v1, tt.v2, got, tt.want)
		}
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input string
		want  [3]int
	}{
		{"1.2.3", [3]int{1, 2, 3}},
		{"2.0.0", [3]int{2, 0, 0}},
		{"10.20.30", [3]int{10, 20, 30}},
		{"1.0", [3]int{1, 0, 0}},     // Missing patch
		{"1", [3]int{1, 0, 0}},       // Only major
		{"", [3]int{0, 0, 0}},        // Empty string
		{"abc", [3]int{0, 0, 0}},     // Invalid format
		{"1.2.3.4", [3]int{1, 2, 3}}, // Extra parts ignored
	}

	for _, tt := range tests {
		got := parseVersion(tt.input)
		if got != tt.want {
			t.Errorf("parseVersion(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}
