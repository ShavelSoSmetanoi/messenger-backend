package s3

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

// Client struct to hold MinIO client and the bucket name
type Client struct {
	Client *minio.Client
	Bucket string
}

// NewS3Client initializes and returns a new MinIO client connected to Yandex Object Storage
func NewS3Client(bucketName string) (*Client, error) {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Initialize the MinIO client with configuration for Yandex Object Storage
	minioClient, err := minio.New(
		os.Getenv("YANDEX_S3_ENDPOINT"),
		&minio.Options{
			Creds:  credentials.NewStaticV4(os.Getenv("YANDEX_S3_ACCESS_KEY"), os.Getenv("YANDEX_S3_SECRET_KEY"), ""),
			Secure: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize MinIO client: %v", err)
	}

	// Create and return the S3 client with the specified bucket name
	return &Client{
		Client: minioClient,
		Bucket: bucketName,
	}, nil
}
