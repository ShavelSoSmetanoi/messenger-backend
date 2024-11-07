package file

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
)

type S3FileService struct {
	client *s3.Client
	bucket string
}

func NewS3FileService(client *s3.Client, bucket string) *S3FileService {
	return &S3FileService{client: client, bucket: bucket}
}

// UploadFile uploads a file to S3 and returns its ID or key
func (s *S3FileService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate a unique file ID or key (could use UUID or filename)
	fileKey := fileHeader.Filename

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileKey),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	return fileKey, nil
}

// DownloadFile retrieves a file from S3
func (s *S3FileService) DownloadFile(ctx context.Context, fileID string) (*File, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, err
	}

	return &File{
		Path:     fileID,
		Content:  resp.Body,
		FileType: *resp.ContentType,
	}, nil
}

// DeleteFile removes a file from S3
func (s *S3FileService) DeleteFile(ctx context.Context, fileID string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileID),
	})
	return err
}

// GetFileInfo retrieves metadata for a file stored in S3
func (s *S3FileService) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	head, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		ID:       fileID,
		Name:     fileID, // Or retrieve original name if stored separately
		Size:     *head.ContentLength,
		Uploaded: *head.LastModified,
	}, nil
}
