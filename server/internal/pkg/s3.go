package pkg

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// To resolve acl error when uploading of objects to the storage
type S3Client struct {
	client     *s3.Client
	bucketName string
	region     string
}

// Database model for storing file references
type FileRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	URL          string    `gorm:"not null" json:"url"`
	Filename     string    `gorm:"not null" json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	CreatedAt    time.Time `json:"created_at"`
}

type UploadResponse struct {
	ID          uint   `json:"id"`
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

func (s3c *S3Client) NewS3Client(ctx context.Context) (*S3Client, error) {
	accessKey := os.Getenv("S3_ACCESS_KEY_ID")
	secretKey := os.Getenv("S3_ACCESS_KEY")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || region == "" || bucketName == "" {
		return nil, fmt.Errorf("missing required environment variables: S3_ACCESS_KEY_ID, S3_ACCESS_KEY, S3_REGION, S3_BUCKET_NAME")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

func (s3c *S3Client) UploadFile(ctx context.Context, file io.Reader, originalFilename string, contentType string, fileSize int64) (*UploadResponse, error) {
	ext := filepath.Ext(originalFilename)

	if queryIndex := strings.Index(ext, "?"); queryIndex != -1 {
		ext = ext[:queryIndex]
	}

	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uniqueID, ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload parameters
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s3c.bucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	_, err := s3c.client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Format: https://bucket-name.s3.region.amazonaws.com/filename
	permanentURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s3c.bucketName, s3c.region, filename)

	return &UploadResponse{
		URL:         permanentURL,
		Filename:    filename,
		Size:        fileSize,
		ContentType: contentType,
	}, nil
}

func (s3c *S3Client) DeleteFile(ctx context.Context, filename string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3c.bucketName),
		Key:    aws.String(filename),
	}

	_, err := s3c.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}
	return nil
}

func (s3c *S3Client) FileExists(ctx context.Context, filename string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s3c.bucketName),
		Key:    aws.String(filename),
	}

	_, err := s3c.client.HeadObject(ctx, input)
	if err != nil {
		// var notFound *types.NotFound
		if err := err.(*types.NotFound); err != nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Generate presigned URL for temporary access (optional feature)
func (s3c *S3Client) GetPresignedURL(ctx context.Context, filename string, duration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3c.client)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s3c.bucketName),
		Key:    aws.String(filename),
	}

	presignedURL, err := presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.URL, nil
}
