package main

import (
	"testing"
)

func TestChartRenderers(t *testing.T) {
	// Create a simple stats structure for testing
	stats := Stats{
		TotalDirs:  5,
		TotalFiles: 10,
		TotalSize:  102400, // 100 KB
		FileTypes: map[string]int64{
			".txt":  51200,  // 50 KB
			".go":   30720,  // 30 KB
			".md":   15360,  // 15 KB
			".json": 5120,   // 5 KB
		},
	}
	
	config := Config{
		StatsCount: 5,
		Unicode:    true,
	}
	
	// Test Unicode chart renderer
	unicodeRenderer := &unicodeChartRenderer{}
	unicodeOutput := unicodeRenderer.renderChart(stats, config)
	
	if unicodeOutput == "" {
		t.Errorf("Unicode chart renderer returned empty output")
	}
	
	if len(unicodeOutput) < 100 {
		t.Errorf("Unicode chart output suspiciously short: %d chars", len(unicodeOutput))
	}
	
	// Test ASCII chart renderer
	asciiRenderer := &asciiChartRenderer{}
	asciiOutput := asciiRenderer.renderChart(stats, config)
	
	if asciiOutput == "" {
		t.Errorf("ASCII chart renderer returned empty output")
	}
	
	if len(asciiOutput) < 100 {
		t.Errorf("ASCII chart output suspiciously short: %d chars", len(asciiOutput))
	}
	
	// Test sparkline chart renderer
	sparklineRenderer := &sparklineChartRenderer{}
	sparklineOutput := sparklineRenderer.renderChart(stats, config)
	
	if sparklineOutput == "" {
		t.Errorf("Sparkline chart renderer returned empty output")
	}
	
	if len(sparklineOutput) < 100 {
		t.Errorf("Sparkline chart output suspiciously short: %d chars", len(sparklineOutput))
	}
}

func TestGetChartRenderer(t *testing.T) {
	// Test with Unicode enabled
	config := Config{
		Unicode: true,
	}
	
	renderer := getChartRenderer(config)
	_, isUnicode := renderer.(*unicodeChartRenderer)
	
	if !isUnicode {
		t.Errorf("getChartRenderer with Unicode=true should return unicodeChartRenderer")
	}
	
	// Test with Unicode disabled
	config.Unicode = false
	
	renderer = getChartRenderer(config)
	_, isAscii := renderer.(*asciiChartRenderer)
	
	if !isAscii {
		t.Errorf("getChartRenderer with Unicode=false should return asciiChartRenderer")
	}
}
