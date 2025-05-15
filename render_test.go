package main

import (
	"os"
	"testing"
)

func TestIsTerminalSupportsUnicode(t *testing.T) {
	// Save original TERM value
	originalTerm := os.Getenv("TERM")
	defer os.Setenv("TERM", originalTerm)

	// Test cases
	testCases := []struct {
		term     string
		expected bool
	}{
		{"xterm", true},
		{"xterm-256color", true},
		{"unicode", true},
		{"utf-8", true},
		{"linux", true},
		{"vt100", false},
		{"dumb", false},
		{"", false},
	}

	for _, tc := range testCases {
		os.Setenv("TERM", tc.term)
		result := IsTerminalSupportsUnicode()
		if result != tc.expected {
			t.Errorf("IsTerminalSupportsUnicode with TERM=%s: expected %v, got %v", tc.term, tc.expected, result)
		}
	}
}

func TestShouldExcludeFolder(t *testing.T) {
	excludeFolders := []string{"node_modules", "bin", ".git"}

	tests := []struct {
		name     string
		expected bool
	}{
		{"node_modules", true},
		{"bin", true},
		{".git", true},
		{"src", false},
		{"lib", false},
	}

	for _, test := range tests {
		result := ShouldExcludeFolder(test.name, excludeFolders)
		if result != test.expected {
			t.Errorf("ShouldExcludeFolder(%q, %v): expected %v, got %v", test.name, excludeFolders, test.expected, result)
		}
	}
}

func TestShouldExcludeFile(t *testing.T) {
	excludeFiles := []string{".exe", ".dll"}
	excludeNames := []string{"config.json", "README.md"}

	tests := []struct {
		name     string
		expected bool
	}{
		{"program.exe", true},
		{"library.dll", true},
		{"config.json", true},
		{"README.md", true},
		{"script.js", false},
		{"document.txt", false},
	}

	for _, test := range tests {
		result := ShouldExcludeFile(test.name, excludeFiles, excludeNames)
		if result != test.expected {
			t.Errorf("ShouldExcludeFile(%q, %v, %v): expected %v, got %v", 
				test.name, excludeFiles, excludeNames, test.expected, result)
		}
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"file.txt", ".txt"},
		{"document.PDF", ".pdf"},
		{"script.js", ".js"},
		{"path/to/image.png", ".png"},
		{"noextension", ""},
		{".hidden", ""},
		{"multiple.dots.md", ".md"},
	}

	for _, test := range tests {
		result := GetFileExtension(test.path)
		if result != test.expected {
			t.Errorf("GetFileExtension(%q): expected %q, got %q", test.path, test.expected, result)
		}
	}
}
