package reports

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var s3Client *minio.Client
var bucketName string

func InitS3() error {
	endpoint := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	if endpoint == "" {
		log.Println("S3 not configured: S3_ENDPOINT not set, skipping S3 initialization")
		return nil
	}

	log.Printf("S3 initialization: endpoint=%s, bucket=%s, useSSL=%s", 
		endpoint, os.Getenv("S3_BUCKET"), os.Getenv("S3_USE_SSL"))

	// Strip inline comments starting with '#'
	if idx := strings.Index(endpoint, "#"); idx != -1 {
		endpoint = strings.TrimSpace(endpoint[:idx])
	}

	// Trim URL schemes and trailing slashes if present (MinIO SDK requires host:port only)
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimSuffix(endpoint, "/")
	endpoint = strings.TrimSpace(endpoint)

	accessKey := strings.TrimSpace(os.Getenv("S3_ACCESS_KEY"))
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := strings.TrimSpace(os.Getenv("S3_SECRET_KEY"))
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	bucketName = strings.TrimSpace(os.Getenv("S3_BUCKET"))
	if bucketName == "" {
		bucketName = "arishem-reports"
	}

	useSSL := os.Getenv("S3_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	s3Client = client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := s3Client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		if err := s3Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	log.Printf("S3 initialized with bucket: %s", bucketName)
	return nil
}

func GetS3Client() *minio.Client {
	return s3Client
}

func GetBucketName() string {
	return bucketName
}

func GetSignedURL(ctx context.Context, orgID, reportID, format string) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	storageKey := fmt.Sprintf("reports/%s.%s", reportID, format)

	presignedURL, err := s3Client.PresignedGetObject(ctx, bucketName, storageKey, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

func UploadReport(ctx context.Context, orgID, reportID, format, contentType, content string) error {
	if s3Client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	storageKey := fmt.Sprintf("reports/%s.%s", reportID, format)

	_, err := s3Client.PutObject(ctx, bucketName, storageKey, strings.NewReader(content), int64(len(content)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload report: %w", err)
	}

	log.Printf("Uploaded report: %s to s3://%s/%s", reportID, bucketName, storageKey)
	return nil
}