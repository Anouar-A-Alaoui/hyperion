package main

import (
	"testing"
	"os"
	"path/filepath"
	"strings"
)

func TestSplitCommaString(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"node_modules", []string{"node_modules"}},
		{"node_modules,bin,obj", []string{"node_modules", "bin", "obj"}},
		{"node_modules, bin, obj", []string{"node_modules", "bin", "obj"}},
	}

	for _, test := range tests {
		result := splitCommaString(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("splitCommaString(%q): expected length %d, got %d", test.input, len(test.expected), len(result))
			continue
		}
		for i, v := range result {
			if v != test.expected[i] {
				t.Errorf("splitCommaString(%q)[%d]: expected %q, got %q", test.input, i, test.expected[i], v)
			}
		}
	}
}

func TestGetTreeChars(t *testing.T) {
	tests := []struct {
		unicode bool
		compact bool
		branch  string
	}{
		{true, false, "├── "},
		{true, true, "├"},
		{false, false, "|-- "},
		{false, true, "+"},
	}

	for _, test := range tests {
		result := getTreeChars(test.unicode, test.compact)
		if !strings.Contains(result.Branch, test.branch) && !strings.Contains(result.MiddleItem, test.branch) {
			t.Errorf("getTreeChars(%v, %v): expected branch containing %q, got %q or %q", 
				test.unicode, test.compact, test.branch, result.Branch, result.MiddleItem)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1500, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{2 * 1024 * 1024, "2.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
		{1024 * 1024 * 1024 * 1024, "1.0 TB"},
	}

	for _, test := range tests {
		result := formatSize(test.size)
		if result != test.expected {
			t.Errorf("formatSize(%d): expected %q, got %q", test.size, test.expected, result)
		}
	}
}

// TestWalkDir is a more complex test that requires setting up a temporary directory structure
func TestWalkDir(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "hyperion-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	dirs := []string{
		"dir1",
		"dir1/subdir1",
		"dir1/subdir2",
		"dir2",
		"node_modules", // This should be excluded by default
	}

	files := []string{
		"file1.txt",
		"file2.exe", // This could be excluded by extension
		"dir1/file3.txt",
		"dir1/subdir1/file4.txt",
		"dir2/README.md", // This could be excluded by name
	}

	// Create directories
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create files with some content
	for _, file := range files {
		content := []byte("Test content for " + file)
		err := os.WriteFile(filepath.Join(tempDir, file), content, 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	// Run a simple test that doesn't depend on stdout capture
	// We'll just verify that stats are collected correctly
	config := Config{
		Path:           tempDir,
		ExcludeFolders: []string{"node_modules"},
		ShowFiles:      true,
		ExcludeFiles:   []string{".exe"},
		ExcludeNames:   []string{"README.md"},
		MaxDepth:       -1,
		Unicode:        true,
		Color:          false,
		StatTable:      true,
	}

	stats := Stats{
		FileTypes:  make(map[string]int64),
		LargeFiles: []FileInfo{},
	}

	treeChars := getTreeChars(config.Unicode, config.Compact)

	// Call walkDir (but we won't check the output, just the stats)
	walkDir(config, config.Path, "", "", 0, treeChars, &stats)

	// Verify statistics
	expectedDirs := 4  // tempDir, dir1, dir1/subdir1, dir1/subdir2, dir2 (excluding node_modules)
	if stats.TotalDirs != expectedDirs {
		t.Errorf("Expected %d directories, got %d", expectedDirs, stats.TotalDirs)
	}

	expectedFiles := 3  // file1.txt, dir1/file3.txt, dir1/subdir1/file4.txt (excluding .exe and README.md)
	if stats.TotalFiles != expectedFiles {
		t.Errorf("Expected %d files, got %d", expectedFiles, stats.TotalFiles)
	}

	// Check that the correct file types were counted
	if _, exists := stats.FileTypes[".txt"]; !exists {
		t.Error("Expected .txt in file types, but it wasn't found")
	}

	if _, exists := stats.FileTypes[".exe"]; exists {
		t.Error("Found .exe in file types, but it should have been excluded")
	}
}
