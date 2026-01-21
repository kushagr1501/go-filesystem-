package main

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
	"time"
)

// listDirectoryNative scans a directory recursively using native Go
func listDirectoryNative(rootPath string) ([]string, error) {
	var files []string
	
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil  // Skip errors, keep going
		}
		
		if !d.IsDir() {
			files = append(files, path)
		}
		
		return nil
	})
	
	return files, err
}
// getFileInfoNative gets file metadata using native Go + Windows API
func getFileInfoNative(path string) (*FileInfo, error) {
	// Get basic file info
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	
	result := &FileInfo{
		Path:        path,
		SizeBytes:   info.Size(),
		IsFile:      !info.IsDir(),
		IsDirectory: info.IsDir(),
	}
	
	// Get Windows-specific times (created, accessed)
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err == nil {
		var attrs syscall.Win32FileAttributeData
		if syscall.GetFileAttributesEx(pathPtr, syscall.GetFileExInfoStandard, (*byte)(unsafe.Pointer(&attrs))) == nil {
			result.CreatedAt = time.Unix(0, attrs.CreationTime.Nanoseconds())
			result.AccessedAt = time.Unix(0, attrs.LastAccessTime.Nanoseconds())
		}
	}
	
	// Modified time from os.Stat()
	result.ModifiedAt = info.ModTime()
	
	// Get file type from extension
	result.FileType = getFileType(path)
	
	return result, nil
}