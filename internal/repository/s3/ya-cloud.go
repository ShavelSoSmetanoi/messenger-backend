package s3

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

type Client struct {
	Client *minio.Client
	Bucket string
}

func NewS3Client(bucketName string) (*Client, error) {
	// Загрузить переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Инициализируем MinIO клиент с настройками для Yandex Object Storage
	minioClient, err := minio.New(
		os.Getenv("YANDEX_S3_ENDPOINT"), // Указываем endpoint для Yandex Object Storage
		&minio.Options{
			Creds:  credentials.NewStaticV4(os.Getenv("YANDEX_S3_ACCESS_KEY"), os.Getenv("YANDEX_S3_SECRET_KEY"), ""),
			Secure: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize MinIO client: %v", err)
		return nil, err
	}

	// Создаем клиент для работы с S3
	return &Client{
		Client: minioClient,
		Bucket: bucketName,
	}, nil
}
