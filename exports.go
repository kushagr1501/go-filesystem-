package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileExport struct {
	Path          string    `json:"path"`
	Name          string    `json:"name"`
	SizeBytes      int64     `json:"size_bytes"`
	SizeMB         float64   `json:"size_mb"`
	FileType      string    `json:"file_type"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
	AccessedAt    time.Time `json:"accessed_at"`
	IsUnused      bool      `json:"is_unused"`
	IsZeroByte    bool      `json:"is_zero_byte"`
}

// ScanExport represents full scan results
type ScanExport struct {
	ScanPath      string       `json:"scan_path"`
	ScanTime      time.Time     `json:"scan_time"`
	FileCount     int           `json:"file_count"`
	TotalSize     int64         `json:"total_size"`
	FilterConfig  FilterConfig  `json:"filter_config"`
	Files         []FileExport  `json:"files"`
}

func ExportToJSON(export ScanExport, filePath string) error {
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}
	
	return nil
}

func ExportToCSV(export ScanExport, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write CSV header
	header := []string{
		"Path", "Name", "Size_Bytes", "Size_MB", "File_Type",
		"Created_At", "Modified_At", "Accessed_At",
		"Is_Unused", "Is_Zero_Byte",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}
	
	// Write file rows
	for _, file := range export.Files {
		row := []string{
			file.Path,
			file.Name,
			strconv.FormatInt(file.SizeBytes, 10),
			fmt.Sprintf("%.2f", file.SizeMB),
			file.FileType,
			file.CreatedAt.Format(time.RFC3339),
			file.ModifiedAt.Format(time.RFC3339),
			file.AccessedAt.Format(time.RFC3339),
			strconv.FormatBool(file.IsUnused),
			strconv.FormatBool(file.IsZeroByte),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %v", err)
		}
	}
	
	return nil
}

//  converts scanned files to export format
func PrepareExportData(files []string, filterConfig FilterConfig, scanPath string) (ScanExport, error) {
	export := ScanExport{
		ScanPath:     scanPath,
		ScanTime:     time.Now(),
		FileCount:    len(files),
		FilterConfig: filterConfig,
		Files:        []FileExport{},
	}
	
	var totalSize int64
	
	for _, path := range files {
		info, err := getFileInfoNative(path)
		if err != nil {
			continue
		}
		
		totalSize += info.SizeBytes
		
		fileExport := FileExport{
			Path:       info.Path,
			Name:       getFileName(info.Path),
			SizeBytes:   info.SizeBytes,
			SizeMB:      float64(info.SizeBytes) / (1024 * 1024),
			FileType:    info.FileType,
			CreatedAt:   info.CreatedAt,
			ModifiedAt:  info.ModifiedAt,
			AccessedAt:  info.AccessedAt,
			IsUnused:    ExplainUnused(info, 60) != nil,
			IsZeroByte:  ExplainZeroByte(info) != nil,
		}
		
		export.Files = append(export.Files, fileExport)
	}
	
	export.TotalSize = totalSize
	return export, nil
}