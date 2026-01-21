package main

import (
	"time"
)

type FileInfo struct {
	Path        string
	SizeBytes   int64
	CreatedAt   time.Time
	ModifiedAt  time.Time
	AccessedAt  time.Time
	IsFile      bool
	IsDirectory bool
	MimeType    string   // e.g., application/pdf, image/jpeg etc.
	FileType    string  // e.g., pdf, doc, img, code, archive
}

