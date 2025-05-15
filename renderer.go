package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// Renderer represents a tree node renderer
type Renderer interface {
	RenderDir(name string, isLast bool, prefix string) (string, string)
	RenderFile(name string, isLast bool, prefix string, isSymlink bool)
}

// ColorRenderer renders the tree with colors
type ColorRenderer struct {
	config Config
	treeChars TreeChars
}

// Create a new color renderer
func NewColorRenderer(config Config, treeChars TreeChars) *ColorRenderer {
	return &ColorRenderer{
		config: config,
		treeChars: treeChars,
	}
}

// RenderDir renders a directory with color
func (r *ColorRenderer) RenderDir(name string, isLast bool, prefix string) (string, string) {
	var newPrefix, newLastPrefix string
	
	if isLast {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.LastItem)
		newPrefix = prefix + r.treeChars.Indent
		newLastPrefix = prefix + r.treeChars.Indent
	} else {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.MiddleItem)
		newPrefix = prefix + r.treeChars.Line
		newLastPrefix = prefix + r.treeChars.Indent
	}

	if r.config.BgColor {
		color.New(color.FgHiWhite, color.BgBlue).Printf("%s\n", name)
	} else {
		color.New(color.FgBlue, color.Bold).Printf("%s\n", name)
	}
	
	return newPrefix, newLastPrefix
}

// RenderFile renders a file with color
func (r *ColorRenderer) RenderFile(name string, isLast bool, prefix string, isSymlink bool) {
	if isLast {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.LastItem)
	} else {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.MiddleItem)
	}

	if isSymlink {
		if r.config.BgColor {
			color.New(color.FgHiWhite, color.BgMagenta).Printf("%s\n", name)
		} else {
			color.New(color.FgMagenta).Printf("%s\n", name)
		}
	} else {
		if r.config.BgColor {
			color.New(color.FgHiWhite, color.BgGreen).Printf("%s\n", name)
		} else {
			color.New(color.FgGreen).Printf("%s\n", name)
		}
	}
}

// BasicRenderer renders the tree without colors
type BasicRenderer struct {
	treeChars TreeChars
}

// Create a new basic renderer
func NewBasicRenderer(treeChars TreeChars) *BasicRenderer {
	return &BasicRenderer{
		treeChars: treeChars,
	}
}

// RenderDir renders a directory without color
func (r *BasicRenderer) RenderDir(name string, isLast bool, prefix string) (string, string) {
	var newPrefix, newLastPrefix string
	
	if isLast {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.LastItem)
		newPrefix = prefix + r.treeChars.Indent
		newLastPrefix = prefix + r.treeChars.Indent
	} else {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.MiddleItem)
		newPrefix = prefix + r.treeChars.Line
		newLastPrefix = prefix + r.treeChars.Indent
	}

	fmt.Printf("%s\n", name)
	
	return newPrefix, newLastPrefix
}

// RenderFile renders a file without color
func (r *BasicRenderer) RenderFile(name string, isLast bool, prefix string, isSymlink bool) {
	if isLast {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.LastItem)
	} else {
		fmt.Print(prefix)
		fmt.Print(r.treeChars.MiddleItem)
	}

	fmt.Printf("%s\n", name)
}

// Helper functions for tree rendering

// GetRenderer returns the appropriate renderer based on config
func GetRenderer(config Config, treeChars TreeChars) Renderer {
	if config.Color {
		return NewColorRenderer(config, treeChars)
	}
	return NewBasicRenderer(treeChars)
}

// ShouldExcludeFolder checks if a folder should be excluded
func ShouldExcludeFolder(name string, excludeFolders []string) bool {
	for _, excludeFolder := range excludeFolders {
		if name == excludeFolder {
			return true
		}
	}
	return false
}

// ShouldExcludeFile checks if a file should be excluded by extension or name
func ShouldExcludeFile(name string, excludeFiles, excludeNames []string) bool {
	// Check exclusion by exact name
	for _, excludeName := range excludeNames {
		if name == excludeName {
			return true
		}
	}
	
	// Check exclusion by extension
	for _, excludeExt := range excludeFiles {
		if strings.HasSuffix(strings.ToLower(name), strings.ToLower(excludeExt)) {
			return true
		}
	}
	
	return false
}

// GetFileExtension gets the file extension in lowercase (with dot)
func GetFileExtension(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	return ext
}

// IsTerminalSupportsUnicode checks if the terminal supports Unicode
func IsTerminalSupportsUnicode() bool {
	if runtime.GOOS == "windows" {
		// Additional checks could be added for Windows terminals
		return false
	}
	
	// Crude check for Unicode support via environment variables
	term := strings.ToLower(os.Getenv("TERM"))
	if strings.Contains(term, "xterm") || strings.Contains(term, "unicode") || 
	   strings.Contains(term, "utf") || strings.Contains(term, "linux") {
		return true
	}
	
	return false
}
