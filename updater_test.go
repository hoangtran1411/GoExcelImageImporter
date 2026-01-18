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
	}

	for _, tt := range tests {
		got := CompareVersions(tt.v1, tt.v2)
		if got != tt.want {
			t.Errorf("CompareVersions(%q, %q) = %v; want %v", tt.v1, tt.v2, got, tt.want)
		}
	}
}
