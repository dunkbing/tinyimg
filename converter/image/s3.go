package image

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	svc *s3.Client
}

var bucketName = "optipic"

var s3Endpoint = os.Getenv("S3_ENDPOINT")
var accessKeyId = os.Getenv("S3_ACCESS_KEY")
var accessKeySecret = os.Getenv("S3_SECRET_KEY")

func NewS3Client() (*S3Client, error) {
	// Create a new AWS session
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s3Endpoint,
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{
		svc: client,
	}, nil
}

func (c *S3Client) UploadFile(key, filePath string) error {
	logger.Info("Uploading file to S3", "key", key, "filePath", filePath)
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	uploader := manager.NewUploader(c.svc)
	chunks := strings.Split(key, "/")
	filename := chunks[len(chunks)-1]
	contentDisposition := fmt.Sprintf(`attachment; filename="%s"`, filename)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
		ACL:    "public-read",
		Metadata: map[string]string{
			"Content-Disposition": contentDisposition,
		},
	})

	return err
}

func (c *S3Client) GetFileUrl(key string) (string, error) {
	presignClient := s3.NewPresignClient(c.svc)

	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}
