package file

import (
	"context"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/messageDB"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"log"
	"mime/multipart"
)

// S3FileService provides file management functionality for S3
type S3FileService struct {
	messageRepo messageDB.MessageRepository
	client      *minio.Client
	bucket      string
}

// NewS3FileService creates a new instance of S3FileService
func NewS3FileService(client *minio.Client, bucket string, repo messageDB.MessageRepository) *S3FileService {
	return &S3FileService{
		messageRepo: repo,
		client:      client,
		bucket:      bucket,
	}
}

// UploadFile uploads a file to S3 and returns its unique ID or key
func (s *S3FileService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate a unique file ID or key using UUID
	fileID := uuid.New().String() + "_" + fileHeader.Filename

	// Upload file to S3
	_, err = s.client.PutObject(ctx, s.bucket, fileID, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: "application/octet-stream", // Default content type
	})
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		return "", err
	}

	return fileID, nil
}

// DownloadFile retrieves a file from S3
func (s *S3FileService) DownloadFile(ctx context.Context, fileID string) (*File, error) {
	// Get file metadata
	fileInfo, err := s.client.StatObject(ctx, s.bucket, fileID, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("Error retrieving file metadata: %v", err)
		return nil, err
	}

	// Retrieve the file itself
	resp, err := s.client.GetObject(ctx, s.bucket, fileID, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error retrieving file from S3: %v", err)
		return nil, err
	}
	defer resp.Close()

	return &File{
		Path:     fileID,
		Content:  resp,                 // This is the file content
		FileType: fileInfo.ContentType, // Get content type from metadata
	}, nil
}

// DeleteFile removes a file from S3
func (s *S3FileService) DeleteFile(ctx context.Context, fileID string) error {
	err := s.client.RemoveObject(ctx, s.bucket, fileID, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Error deleting file from S3: %v", err)
	}
	return err
}

// GetFileInfo retrieves metadata for a file stored in S3
func (s *S3FileService) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	// Get object metadata
	head, err := s.client.StatObject(ctx, s.bucket, fileID, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("Error retrieving file metadata: %v", err)
		return nil, err
	}

	return &FileInfo{
		ID:       fileID,
		Name:     fileID, // Or retrieve original name if stored separately
		Size:     head.Size,
		Uploaded: head.LastModified,
		FileType: head.ContentType,
	}, nil
}
