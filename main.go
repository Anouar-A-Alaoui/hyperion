package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Configuration structure to hold all command line flags
type Config struct {
	Path           string
	ExcludeFolders []string
	ShowFiles      bool
	ExcludeFiles   []string
	ExcludeNames   []string
	MaxDepth       int
	Unicode        bool
	Color          bool
	BgColor        bool
	Compact        bool
	ShowStats      bool
	StatTable      bool
	StatsCount     int
	Chart          bool
}

// Statistics structure to track directory stats
type Stats struct {
	TotalDirs  int
	TotalFiles int
	TotalSize  int64
	FileTypes  map[string]int64
	LargeFiles []FileInfo
}

// FileInfo to track file stats for the largest files
type FileInfo struct {
	Path string
	Size int64
	Type string
}

// TreeChars defines the characters used to draw the tree
type TreeChars struct {
	Line       string
	Branch     string
	LastItem   string
	MiddleItem string
	Indent     string
}

// Main function
func main() {
	// Define and parse command line flags
	var config Config
	var excludeFoldersStr, excludeFilesStr, excludeNamesStr string

	flag.StringVar(&config.Path, "path", ".", "Root directory to scan")
	flag.StringVar(&excludeFoldersStr, "exclude-folders", "node_modules", "Folders to exclude from tree (comma-separated)")
	flag.BoolVar(&config.ShowFiles, "show-files", false, "Whether to show files in output")
	flag.StringVar(&excludeFilesStr, "exclude-files", "", "File extensions to exclude (comma-separated, e.g., '.exe,.dll')")
	flag.StringVar(&excludeNamesStr, "exclude-names", "", "File names to exclude exactly (comma-separated, e.g., 'config.json,README.md')")
	flag.IntVar(&config.MaxDepth, "max-depth", -1, "Maximum depth to recurse (-1 for unlimited)")
	flag.BoolVar(&config.Unicode, "unicode", true, "Use Unicode characters for pretty tree visuals")
	flag.BoolVar(&config.Color, "color", true, "Use colors in output")
	flag.BoolVar(&config.BgColor, "bg-color", false, "Use background color for items")
	flag.BoolVar(&config.Compact, "compact", false, "Enable compact tree layout")
	flag.BoolVar(&config.ShowStats, "show-stats", false, "Show total files, dirs, size")
	flag.BoolVar(&config.StatTable, "stat-table", false, "Show a table of largest files and types")
	flag.IntVar(&config.StatsCount, "stats-count", 10, "Number of top files to show in stats table")
	flag.BoolVar(&config.Chart, "chart", false, "Show a visual chart of file size distribution")
	
	// Check for help flag
	helpFlag    := flag.Bool("help", false, "Show usage and examples")
	aboutFlag   := flag.Bool("about", false, "Show about the software")
	versionFlag := flag.Bool("version", false, "Show version information")

	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	if *aboutFlag {
		showAbout()
		return
	}

	if *versionFlag {
		showVersion()
		return
	}

	// Process comma-separated values into slices
	config.ExcludeFolders = splitCommaString(excludeFoldersStr)
	config.ExcludeFiles   = splitCommaString(excludeFilesStr)
	config.ExcludeNames   = splitCommaString(excludeNamesStr)

	// Auto-detect Unicode support if needed
	if runtime.GOOS == "windows" && config.Unicode {
		unicodeSupported := isTerminalSupportsUnicode()
		if !unicodeSupported {
			fmt.Println("Note: Unicode characters may not display correctly in this terminal.")
			fmt.Println("      Use --unicode=false for ASCII characters instead.")
		}
	}

	// Initialize statistics
	stats := Stats{
		FileTypes:  make(map[string]int64),
		LargeFiles: []FileInfo{},
	}

	// Select tree characters based on Unicode flag
	treeChars := getTreeChars(config.Unicode, config.Compact)

	// Print the root directory
	rootInfo, err := os.Stat(config.Path)
	if err != nil {
		fmt.Printf("Error accessing path %s: %v\n", config.Path, err)
		return
	}

	rootDir := filepath.Base(config.Path)
	
	// Print root directory with appropriate styling
	if config.Color {
		if config.BgColor {
			color.New(color.FgHiWhite, color.BgBlue).Printf("%s\n", rootDir)
		} else {
			color.New(color.FgBlue, color.Bold).Printf("%s\n", rootDir)
		}
	} else {
		fmt.Printf("%s\n", rootDir)
	}

	// Walk the directory tree
	if rootInfo.IsDir() {
		walkDir(config, config.Path, "", "", 0, treeChars, &stats)
	}

	// Show statistics if requested
	if config.ShowStats || config.StatTable || config.Chart {
		printStats(config, stats)
	}
}

// Function for version display
func showVersion() {
	fmt.Println("Hyperion - Advanced Directory Tree Visualizer")
	fmt.Println("Version: 1.0.0")
}

// Print the help message and examples
func showHelp() {
	helpText := `
	🔧 hyperion - Directory Tree Visualizer

	Usage:
	hyperion [flags]

	Flags:
	--path string             Root directory to scan (default ".")
	--exclude-folders string  Folders to exclude from tree (default "node_modules")
	--show-files              Whether to show files in output (default false)
	--exclude-files string    File extensions to exclude (e.g., ".exe,.dll")
	--exclude-names string    File names to exclude exactly (e.g., "config.json,README.md")
	--max-depth int           Maximum depth to recurse (-1 for unlimited) (default -1)
	--unicode                 Use Unicode characters for pretty tree visuals (default true)
	--color                   Use colors in output (default true)
	--bg-color                Use background color for items (default false)
	--compact                 Enable compact tree layout (default false)
	--show-stats              Show total files, dirs, size (default false)
	--stat-table              Show a table of largest files and types (default false)
	--stats-count int         Number of top files to show in stats table (default 10)
	--chart                   Show a visual chart of file size distribution (default false)
	--help                    Show usage and examples
	--about                   Show about
	--version                 Show version

	Examples:
	# Basic usage (folders only)
	hyperion

	# Exclude folders and file types
	hyperion --show-files --exclude-folders "bin,obj" --exclude-files ".exe,.dll"

	# Show stats with Unicode and color
	hyperion --show-files --unicode --color --show-stats

	# Show top 15 largest files with chart
	hyperion --show-files --stat-table --stats-count 15 --chart

	# Compact view with background color
	hyperion --show-files --compact --bg-color
	`
	fmt.Println(helpText)
}

// About
func showAbout() {
	helpText := `
	______  __                           _____              
	___  / / /____  ________________________(_)____________ 
	__  /_/ /__  / / /__  __ \  _ \_  ___/_  /_  __ \_  __ \
	_  __  / _  /_/ /__  /_/ /  __/  /   _  / / /_/ /  / / /
	/_/ /_/  _\__, / _  .___/\___//_/    /_/  \____//_/ /_/ 
	         /____/  /_/ 
	
	--------------------------------------------------------

	Hyperion - Advanced Directory Tree Visualizer
	Version: 1.0.0
	License: MIT

	A powerful command-line tool for visualizing directory structures with:
	- Customizable tree display with Unicode/ASCII characters
	- Colorized output with file type differentiation
	- Comprehensive filtering options
	- Detailed statistics and analytics
	- Interactive charts and tables

	Features:
	✓ Beautiful tree visualization with configurable characters
	✓ Smart filtering of files and directories
	✓ File type statistics and size analysis
	✓ Largest files identification
	✓ Visual charts of file distribution
	✓ Cross-platform support

	Author        : Anouar AL ECHEIKH EL ALAOUI 
	Repository    : https://github.com/Anouar-A-Alaoui/hyperion
	Documentation : https://github.com/Anouar-A-Alaoui/hyperion

	Use 'hyperion --help' for usage instructions.	
	`
	fmt.Println(helpText)
}

// Split a comma-separated string into a slice
func splitCommaString(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// Get tree characters based on Unicode flag and compact mode
func getTreeChars(useUnicode bool, compact bool) TreeChars {
	if useUnicode {
		if compact {
			return TreeChars{
				Line:       "│",
				Branch:     "├",
				LastItem:   "└",
				MiddleItem: "├",
				Indent:     " ",
			}
		}
		return TreeChars{
			Line:       "│   ",
			Branch:     "├── ",
			LastItem:   "└── ",
			MiddleItem: "├── ",
			Indent:     "    ",
		}
	}
	if compact {
		return TreeChars{
			Line:       "|",
			Branch:     "|",
			LastItem:   "`",
			MiddleItem: "+",
			Indent:     " ",
		}
	}
	return TreeChars{
		Line:       "|   ",
		Branch:     "|-- ",
		LastItem:   "`-- ",
		MiddleItem: "+-- ",
		Indent:     "    ",
	}
}

// Walk directory tree recursively
func walkDir(config Config, path string, prefix string, lastPrefix string, depth int, treeChars TreeChars, stats *Stats) {
	// Check max depth
	if config.MaxDepth != -1 && depth > config.MaxDepth {
		return
	}

	// Read directory entries
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", path, err)
		return
	}

	// Filter and sort entries
	var dirs, files []fs.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			// Check if folder should be excluded
			if !shouldExcludeFolder(name, config.ExcludeFolders) {
				dirs = append(dirs, entry)
			}
		} else if config.ShowFiles {
			// Check if file should be excluded by extension or name
			if !shouldExcludeFile(name, config.ExcludeFiles, config.ExcludeNames) {
				files = append(files, entry)
			}
		}
	}

	// Process directories
	for i, entry := range dirs {
		isLast := i == len(dirs)-1 && len(files) == 0

		entryPath := filepath.Join(path, entry.Name())
		newPrefix, newLastPrefix := renderDir(entry.Name(), isLast, prefix, config, treeChars)

		stats.TotalDirs++
		walkDir(config, entryPath, newPrefix, newLastPrefix, depth+1, treeChars, stats)
	}

	// Process files
	for i, entry := range files {
		isLast := i == len(files)-1

		entryPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", entryPath, err)
			continue
		}

		fileSize := info.Size()
		fileExt := getFileExtension(entry.Name())
		
		// Update statistics
		stats.TotalFiles++
		stats.TotalSize += fileSize
		stats.FileTypes[fileExt] += fileSize
		
		// Track large files for stat table
		if config.StatTable {
			stats.LargeFiles = append(stats.LargeFiles, FileInfo{
				Path: entryPath,
				Size: fileSize,
				Type: fileExt,
			})
		}

		// Determine if the file is a symlink
		isSymlink := info.Mode()&os.ModeSymlink != 0
		
		// Render the file
		renderFile(entry.Name(), isLast, prefix, isSymlink, config, treeChars)
	}
}

// Render directory with appropriate styling
func renderDir(name string, isLast bool, prefix string, config Config, treeChars TreeChars) (string, string) {
	var newPrefix, newLastPrefix string
	
	if isLast {
		fmt.Print(prefix)
		fmt.Print(treeChars.LastItem)
		newPrefix = prefix + treeChars.Indent
		newLastPrefix = prefix + treeChars.Indent
	} else {
		fmt.Print(prefix)
		fmt.Print(treeChars.MiddleItem)
		newPrefix = prefix + treeChars.Line
		newLastPrefix = prefix + treeChars.Indent
	}

	if config.Color {
		if config.BgColor {
			color.New(color.FgHiWhite, color.BgBlue).Printf("%s\n", name)
		} else {
			color.New(color.FgBlue, color.Bold).Printf("%s\n", name)
		}
	} else {
		fmt.Printf("%s\n", name)
	}
	
	return newPrefix, newLastPrefix
}

// Render file with appropriate styling
func renderFile(name string, isLast bool, prefix string, isSymlink bool, config Config, treeChars TreeChars) {
	if isLast {
		fmt.Print(prefix)
		fmt.Print(treeChars.LastItem)
	} else {
		fmt.Print(prefix)
		fmt.Print(treeChars.MiddleItem)
	}

	if config.Color {
		if isSymlink {
			if config.BgColor {
				color.New(color.FgHiWhite, color.BgMagenta).Printf("%s\n", name)
			} else {
				color.New(color.FgMagenta).Printf("%s\n", name)
			}
		} else {
			if config.BgColor {
				color.New(color.FgHiWhite, color.BgGreen).Printf("%s\n", name)
			} else {
				color.New(color.FgGreen).Printf("%s\n", name)
			}
		}
	} else {
		fmt.Printf("%s\n", name)
	}
}

// Check if a folder should be excluded
func shouldExcludeFolder(name string, excludeFolders []string) bool {
	for _, excludeFolder := range excludeFolders {
		if name == excludeFolder {
			return true
		}
	}
	return false
}

// Check if a file should be excluded by extension or name
func shouldExcludeFile(name string, excludeFiles, excludeNames []string) bool {
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

// Get file extension in lowercase (with dot)
func getFileExtension(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	return ext
}

// Check if terminal supports Unicode
func isTerminalSupportsUnicode() bool {
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

// Print statistics information
func printStats(config Config, stats Stats) {
	fmt.Println("\n📊 Statistics:")
	fmt.Printf("  - Total Directories: %d\n", stats.TotalDirs)
	fmt.Printf("  - Total Files: %d\n", stats.TotalFiles)
	fmt.Printf("  - Total Size: %s\n", formatSize(stats.TotalSize))

	// Print top largest files table if requested
	if config.StatTable && len(stats.LargeFiles) > 0 {
		fmt.Println("\n📈 Largest Files:")
		
		// Sort files by size (descending)
		sort.Slice(stats.LargeFiles, func(i, j int) bool {
			return stats.LargeFiles[i].Size > stats.LargeFiles[j].Size
		})
		
		// Create a tabwriter for aligned columns
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "  Size\tPath\tType\t")
		fmt.Fprintln(w, "  ----\t----\t----\t")
		
		// Limit to the requested count
		limit := config.StatsCount
		if limit > len(stats.LargeFiles) {
			limit = len(stats.LargeFiles)
		}
		
		for i := 0; i < limit; i++ {
			file := stats.LargeFiles[i]
			relativePath, _ := filepath.Rel(config.Path, file.Path)
			fmt.Fprintf(w, "  %s\t%s\t%s\t\n", 
				formatSize(file.Size), 
				relativePath, 
				file.Type)
		}
		w.Flush()
	}

	// Print file type distribution
	if len(stats.FileTypes) > 0 {
		fmt.Println("\n🗂️ File Type Distribution:")
		
		// Convert map to slice for sorting
		type TypeInfo struct {
			Ext  string
			Size int64
		}
		
		typeInfos := make([]TypeInfo, 0, len(stats.FileTypes))
		for ext, size := range stats.FileTypes {
			if ext == "" {
				ext = "(no extension)"
			}
			typeInfos = append(typeInfos, TypeInfo{ext, size})
		}
		
		// Sort by size
		sort.Slice(typeInfos, func(i, j int) bool {
			return typeInfos[i].Size > typeInfos[j].Size
		})
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "  Type\tSize\tPercentage\t")
		fmt.Fprintln(w, "  ----\t----\t----------\t")
		
		for _, info := range typeInfos {
			percentage := float64(info.Size) / float64(stats.TotalSize) * 100
			fmt.Fprintf(w, "  %s\t%s\t%.1f%%\t\n", 
				info.Ext, 
				formatSize(info.Size), 
				percentage)
		}
		w.Flush()
	}

	// Print chart if requested
	if config.Chart && len(stats.FileTypes) > 0 {
		printChart(stats, config)
	}
}

// Format file size in human-readable format
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// Print ASCII/Unicode chart of file type distribution
func printChart(stats Stats, config Config) {
	if config.Unicode {
		printUnicodeChart(stats, config)
	} else {
		printAsciiChart(stats, config)
	}
}

// Print chart with Unicode block characters
func printUnicodeChart(stats Stats, config Config) {
	fmt.Println("\n📊 File Size Distribution Chart:")
	
	// Convert map to slice for sorting
	type TypeInfo struct {
		Ext  string
		Size int64
	}
	
	typeInfos := make([]TypeInfo, 0, len(stats.FileTypes))
	for ext, size := range stats.FileTypes {
		if ext == "" {
			ext = "(no extension)"
		}
		typeInfos = append(typeInfos, TypeInfo{ext, size})
	}
	
	// Sort by size
	sort.Slice(typeInfos, func(i, j int) bool {
		return typeInfos[i].Size > typeInfos[j].Size
	})
	
	// Limit to top items for the chart
	limit := config.StatsCount
	if limit > len(typeInfos) {
		limit = len(typeInfos)
	}
	
	typeInfos = typeInfos[:limit]
	
	// Find the maximum size for scaling
	maxSize := int64(0)
	for _, info := range typeInfos {
		if info.Size > maxSize {
			maxSize = info.Size
		}
	}
	
	// Maximum bar width
	const maxWidth = 50
	
	// Print the chart
	for _, info := range typeInfos {
		// Calculate bar width proportional to size
		width := int(float64(info.Size) / float64(maxSize) * float64(maxWidth))
		if width < 1 {
			width = 1
		}
		
		// Print the bar
		bar := strings.Repeat("█", width)
		percentage := float64(info.Size) / float64(stats.TotalSize) * 100
		
		fmt.Printf("  %-15s [%-50s] %6.1f%% (%s)\n", 
			info.Ext, 
			bar, 
			percentage, 
			formatSize(info.Size))
	}
}

// Print chart with ASCII characters
func printAsciiChart(stats Stats, config Config) {
	fmt.Println("\n# File Size Distribution Chart:")
	
	// Convert map to slice for sorting
	type TypeInfo struct {
		Ext  string
		Size int64
	}
	
	typeInfos := make([]TypeInfo, 0, len(stats.FileTypes))
	for ext, size := range stats.FileTypes {
		if ext == "" {
			ext = "(no extension)"
		}
		typeInfos = append(typeInfos, TypeInfo{ext, size})
	}
	
	// Sort by size
	sort.Slice(typeInfos, func(i, j int) bool {
		return typeInfos[i].Size > typeInfos[j].Size
	})
	
	// Limit to top items for the chart
	limit := config.StatsCount
	if limit > len(typeInfos) {
		limit = len(typeInfos)
	}
	
	typeInfos = typeInfos[:limit]
	
	// Find the maximum size for scaling
	maxSize := int64(0)
	for _, info := range typeInfos {
		if info.Size > maxSize {
			maxSize = info.Size
		}
	}
	
	// Maximum bar width
	const maxWidth = 50
	
	// Print the chart
	for _, info := range typeInfos {
		// Calculate bar width proportional to size
		width := int(float64(info.Size) / float64(maxSize) * float64(maxWidth))
		if width < 1 {
			width = 1
		}
		
		// Print the bar
		bar := strings.Repeat("#", width)
		percentage := float64(info.Size) / float64(stats.TotalSize) * 100
		
		fmt.Printf("  %-15s [%-50s] %6.1f%% (%s)\n", 
			info.Ext, 
			bar, 
			percentage, 
			formatSize(info.Size))
	}
}
