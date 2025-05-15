package main

import (
	"fmt"
	"sort"
	"strings"
)

// chartRenderer represents an interface for rendering charts
type chartRenderer interface {
	renderChart(stats Stats, config Config) string
}

// unicodeChartRenderer implements chartRenderer using Unicode block characters
type unicodeChartRenderer struct{}

func (r *unicodeChartRenderer) renderChart(stats Stats, config Config) string {
	var sb strings.Builder
	sb.WriteString("\n📊 File Size Distribution Chart:\n")
	
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
		
		fmt.Fprintf(&sb, "  %-15s [%-50s] %6.1f%% (%s)\n", 
			info.Ext, 
			bar, 
			percentage, 
			formatSize(info.Size))
	}
	
	return sb.String()
}

// asciiChartRenderer implements chartRenderer using ASCII characters
type asciiChartRenderer struct{}

func (r *asciiChartRenderer) renderChart(stats Stats, config Config) string {
	var sb strings.Builder
	sb.WriteString("\n# File Size Distribution Chart:\n")
	
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
		
		fmt.Fprintf(&sb, "  %-15s [%-50s] %6.1f%% (%s)\n", 
			info.Ext, 
			bar, 
			percentage, 
			formatSize(info.Size))
	}
	
	return sb.String()
}

// sparklineChartRenderer implements chartRenderer using sparkline characters
type sparklineChartRenderer struct{}

func (r *sparklineChartRenderer) renderChart(stats Stats, config Config) string {
	var sb strings.Builder
	sb.WriteString("\n⚡ File Size Distribution Sparklines:\n")
	
	// Sparkline characters for different levels
	// from lowest to highest: ▁ ▂ ▃ ▄ ▅ ▆ ▇ █
	sparkChars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	
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
	
	// Print the chart
	for _, info := range typeInfos {
		// Calculate the level (0-7)
		level := int((float64(info.Size) / float64(maxSize)) * float64(len(sparkChars)-1))
		if level < 0 {
			level = 0
		} else if level >= len(sparkChars) {
			level = len(sparkChars) - 1
		}
		
		// Generate a 20-character sparkline where only the relevant part is filled
		sparkline := strings.Repeat(" ", 20)
		sparkChar := string(sparkChars[level])
		
		percentage := float64(info.Size) / float64(stats.TotalSize) * 100
		
		fmt.Fprintf(&sb, "  %-15s %s %6.1f%% (%s)\n", 
			info.Ext, 
			sparkChar, 
			percentage, 
			formatSize(info.Size))
	}
	
	return sb.String()
}

// getChartRenderer returns the appropriate chart renderer based on config
func getChartRenderer(config Config) chartRenderer {
	if config.Unicode {
		return &unicodeChartRenderer{}
	}
	return &asciiChartRenderer{}
}
