package file

import (
	"context"
	"io"
	"mime/multipart"
	"time"
)

// FileService defines methods for handling file operations
type FileService interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
	DownloadFile(ctx context.Context, fileID string) (*File, error)
	DeleteFile(ctx context.Context, fileID string) error
	GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error)
}

// File represents the file structure for download operation
type File struct {
	Path     string
	Content  io.ReadCloser
	FileType string
}

// FileInfo represents metadata for a file
type FileInfo struct {
	ID       string
	Name     string
	Size     int64
	Uploaded time.Time
	FileType string
}