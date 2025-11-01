package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"restaurant-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// S3Service handles S3 operations for tenant isolation
type S3Service struct {
	client     *s3.Client
	bucketName string
}

// NewS3Service creates a new S3Service instance
func NewS3Service(cfg *config.Config) (*S3Service, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// If credentials are provided via environment, set them explicitly
	if cfg.AWSAccessKeyID != "" && cfg.AWSSecretAccessKey != "" {
		// Note: In production, use IAM roles instead of explicit credentials
		// This is for development/testing purposes
	}

	return &S3Service{
		client:     s3.NewFromConfig(awsCfg),
		bucketName: cfg.S3BucketName,
	}, nil
}

// UploadFile uploads a file to S3 with tenant-specific prefix
func (s *S3Service) UploadFile(ctx context.Context, restaurantID uint, fileName string, fileType string, fileReader io.Reader) (string, error) {
	// Generate unique key with tenant prefix
	fileExtension := getFileExtension(fileName)
	uniqueID := uuid.New().String()
	key := fmt.Sprintf("restaurant-%d/menu-items/%s%s", restaurantID, uniqueID, fileExtension)

	// Upload file
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        fileReader,
		ContentType: aws.String(fileType),
		ACL:         types.ObjectCannedACLPrivate, // Private by default
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return key, nil
}

// GeneratePresignedURL generates a presigned URL for accessing an S3 object
func (s *S3Service) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// getFileExtension extracts the file extension from a filename
func getFileExtension(fileName string) string {
	extension := ""
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			extension = fileName[i:]
			break
		}
	}
	// Default to .jpg if no extension found
	if extension == "" {
		extension = ".jpg"
	}
	return extension
}

