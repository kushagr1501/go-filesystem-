package main

import (
	"fmt"
	"sort"
	"strings"
)

type FileTypeStats struct {
	Count          int
	TotalSizeBytes int64
	FileType       string
}

func GenerateStats(files []FileExport) []FileTypeStats {
	statsMap := make(map[string]*FileTypeStats)
	for _, path := range files {
		info, err := getFileInfoNative(path.Path)
		if err != nil {
			continue
		}

		if !ShouldInclude(path.Path, filterConfig) {
			continue
		}

		if !ShouldIncludeSize(path.Path, filterConfig, info.SizeBytes) {
			continue
		}

		ftype := info.FileType
		if ftype == "Unknown" {
			ftype = "Other"
		}
		if _, exists := statsMap[ftype]; !exists {
			statsMap[ftype] = &FileTypeStats{
				Count:          0,
				TotalSizeBytes: 0,
				FileType:       ftype,
			}
		}

		statsMap[ftype].Count++
		statsMap[ftype].TotalSizeBytes += info.SizeBytes
		statsMap[ftype].FileType = ftype
	}
	var statsSlice []FileTypeStats
	for _, stat := range statsMap {
		statsSlice = append(statsSlice, *stat)
	}

	// Sort by count (most files first)
	sort.Slice(statsSlice, func(i, j int) bool {
		return statsSlice[i].Count > statsSlice[j].Count
	})

	return statsSlice

}

// DisplayTypeDistribution shows file type statistics with visual bars
func DisplayTypeDistribution(stats []FileTypeStats, grandTotal int, grandSize int64) {
	if len(stats) == 0 {
		PrintInfo("No files found matching filters")
		return
	}

	PrintDivider()
	fmt.Printf("%sFile Type Distribution (%d files, %s total)%s\n",
		ColorYellow+ColorBold, grandTotal, formatFileSize(grandSize), ColorReset)
	PrintDivider()

	// Find maximum values for bar scaling
	maxCount := 0
	maxSize := int64(0)
	for _, stat := range stats {
		if stat.Count > maxCount {
			maxCount = stat.Count
		}
		if stat.TotalSizeBytes > maxSize {
			maxSize = stat.TotalSizeBytes
		}
	}

	// Display each type with bars
	for _, stat := range stats {
		countBarLength := 0
		sizeBarLength := 0

		if maxCount > 0 {
			countBarLength = int(float64(stat.Count) / float64(maxCount) * 40)
		}
		if maxSize > 0 {
			sizeBarLength = int(float64(stat.TotalSizeBytes) / float64(maxSize) * 40)
		}

		countBar := ColorGreen + strings.Repeat("█", countBarLength) + ColorDim + strings.Repeat("░", 40-countBarLength) + ColorReset
		sizeBar := ColorYellow + strings.Repeat("█", sizeBarLength) + ColorDim + strings.Repeat("░", 40-sizeBarLength) + ColorReset

		fmt.Printf("\n%s%-12s%s ", ColorCyan+ColorBold, stat.FileType, ColorReset)
		fmt.Printf("%5d files  ", stat.Count)
		fmt.Printf("%10s\n", formatFileSize(stat.TotalSizeBytes))
		fmt.Printf("  Count: %s\n", countBar)
		fmt.Printf("  Size:  %s\n", sizeBar)
	}

	PrintDivider()

	// Show summary
	fmt.Printf("\n%sSummary:%s\n", ColorBold, ColorReset)
	fmt.Printf("  Total Files: %s%d%s\n", ColorYellow+ColorBold, grandTotal, ColorReset)
	fmt.Printf("  Total Size:  %s\n",
		ColorYellow+ColorBold+formatFileSize(grandSize)+ColorReset)
	fmt.Printf("  File Types:  %s%d%s\n", ColorYellow+ColorBold, len(stats), ColorReset)
	PrintDivider()
}

// GetDistributionTotals calculates grand totals
func GetDistributionTotals(files []FileExport) (int, int64) {
	grandTotal := 0
	grandSize := int64(0)

	for _, path := range files {
		info, err := getFileInfoNative(path.Path)
		if err != nil {
			continue
		}

		// Apply filters
		if !ShouldInclude(path.Path, filterConfig) {
			continue
		}

		if !ShouldIncludeSize(path.Path, filterConfig, info.SizeBytes) {
			continue
		}

		grandTotal++
		grandSize += info.SizeBytes
	}

	return grandTotal, grandSize
}
