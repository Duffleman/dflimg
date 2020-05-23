package storageproviders

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"dflimg"
	"dflimg/dflerr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	pkgerr "github.com/pkg/errors"
)

const (
	// S3Region is the region to upload to in S3
	S3Region = "eu-west-1"
	// S3Bucket is the name of the bucket to upload to in S3
	S3Bucket = "s3.duffleman.co.uk"
	// S3RootKey is the folder that stores the images inside the bucket
	S3RootKey = "i.dfl.mn"
)

// AWS is a storage provider for AWS S3
type AWS struct {
	driver *session.Session
}

// NewAWSProvider returns a new provider for AWS, because 'driver' is private
func NewAWSProvider(driver *session.Session) StorageProvider {
	return &AWS{driver: driver}
}

// NewAWSProviderFromEnv builds an AWS driver from ENV variables
func NewAWSProviderFromEnv() (StorageProvider, error) {
	awsDriver, err := session.NewSession(&aws.Config{Region: aws.String(S3Region)})
	if err != nil {
		return nil, err
	}

	aws := NewAWSProvider(awsDriver)

	return aws, nil
}

// CheckEnvVariables checks for the minimum required AWS variables
func (a *AWS) CheckEnvVariables() error {
	requiredEnvVars := []string{
		"AWS_ACCESS_KEY",
		"AWS_SECRET_KEY",
	}

	for _, k := range requiredEnvVars {
		if val := os.Getenv(k); val == "" {
			return fmt.Errorf("missing env variable (%s)", k)
		}
	}

	return nil
}

// SupportsTwoStage returns whether this provider supports URL signing
func (a *AWS) SupportsTwoStage() bool {
	return true
}

// GenerateKey returns a file key for uploading to AWS S3
func (a *AWS) GenerateKey(fileID string) string {
	return fmt.Sprintf("%s/%s", S3RootKey, fileID)
}

// PrepareUpload to AWS S3
func (a *AWS) PrepareUpload(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	s3req, _ := s3.New(a.driver).PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(S3Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})

	url, err := s3req.Presign(expiry)
	if err != nil {
		return "", pkgerr.Wrap(err, "unable to create presigned s3 url")
	}

	return url, nil
}

// Get a file from AWS S3
func (a *AWS) Get(ctx context.Context, resource *dflimg.Resource) ([]byte, *time.Time, error) {
	s3item, err := s3.New(a.driver).GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(resource.Link),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, nil, dflerr.ErrNotFound
		}

		return nil, nil, err
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, s3item.Body)
	if err != nil {
		return nil, nil, err
	}

	bytes := buf.Bytes()

	return bytes, s3item.LastModified, nil
}

// Upload a file directly to AWS S3
func (a *AWS) Upload(ctx context.Context, key, contentType string, file bytes.Buffer) error {
	_, err := s3.New(a.driver).PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(S3Bucket),
		Key:           aws.String(key),
		ACL:           aws.String("private"),
		Body:          bytes.NewReader(file.Bytes()),
		ContentLength: aws.Int64(int64(file.Len())),
		ContentType:   aws.String(contentType),
	})

	return err
}
