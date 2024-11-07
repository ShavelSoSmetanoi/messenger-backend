package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Client struct {
	Client *s3.Client
	Bucket string
}

func NewS3Client(bucketName string) (*Client, error) {
	// Загрузить переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Создаем конфигурацию SDK
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("YANDEX_S3_REGION")),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
				os.Getenv("YANDEX_S3_ACCESS_KEY"),
				os.Getenv("YANDEX_S3_SECRET_KEY"),
				"",
			)),
		),
		// Устанавливаем endpoint для Yandex Object Storage
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           os.Getenv("YANDEX_S3_ENDPOINT"), // Например: "https://storage.yandexcloud.net"
				SigningRegion: region,
			}, nil
		})),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	// Инициализируем S3 клиент
	return &Client{
		Client: s3.NewFromConfig(cfg),
		Bucket: bucketName,
	}, nil
}
