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

func TestUpdateInfoDefaults(t *testing.T) {
	info := UpdateInfo{
		Available:  false,
		CurrentVer: "v1.0.0",
	}

	if info.Available {
		t.Error("Expected Available to be false")
	}
	if info.CurrentVer != "v1.0.0" {
		t.Errorf("Expected CurrentVer to be v1.0.0, got %s", info.CurrentVer)
	}
	if info.LatestVer != "" {
		t.Errorf("Expected LatestVer to be empty, got %s", info.LatestVer)
	}
	if info.DownloadURL != "" {
		t.Errorf("Expected DownloadURL to be empty, got %s", info.DownloadURL)
	}
}

func TestGitHubReleaseStruct(t *testing.T) {
	release := GitHubRelease{
		TagName: "v2.0.0",
		HTMLURL: "https://github.com/test/repo/releases/v2.0.0",
	}

	if release.TagName != "v2.0.0" {
		t.Errorf("Expected TagName v2.0.0, got %s", release.TagName)
	}
	if len(release.Assets) != 0 {
		t.Errorf("Expected 0 assets, got %d", len(release.Assets))
	}
}

func TestCurrentVersion(t *testing.T) {
	// Verify CurrentVersion is set and starts with 'v'
	if CurrentVersion == "" {
		t.Error("CurrentVersion should not be empty")
	}
	if CurrentVersion[0] != 'v' {
		t.Errorf("CurrentVersion should start with 'v', got %s", CurrentVersion)
	}
}

func TestGitHubConstants(t *testing.T) {
	if GitHubOwner == "" {
		t.Error("GitHubOwner should not be empty")
	}
	if GitHubRepo == "" {
		t.Error("GitHubRepo should not be empty")
	}
}
