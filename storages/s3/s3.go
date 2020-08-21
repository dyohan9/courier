package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/nyaruka/courier"
)

var s3BucketURL = "https://%s.s3.amazonaws.com%s"

func init() {
	courier.RegisterStorage("s3", newStorage)
}

type storage struct {
	client s3iface.S3API
	bucket string
}

func newStorage(config *courier.Config) (courier.Storage, error) {
	s3Session, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.AWSAccessKeyID, config.AWSSecretAccessKey, ""),
		Endpoint:         aws.String(config.S3Endpoint),
		Region:           aws.String(config.S3Region),
		DisableSSL:       aws.Bool(config.S3DisableSSL),
		S3ForcePathStyle: aws.Bool(config.S3ForcePathStyle),
	})
	if err != nil {
		return nil, err
	}
	return &storage{
		client: s3.New(s3Session),
		bucket: config.S3MediaBucket,
	}, nil
}

func (s *storage) Test(ctx context.Context) error {
	params := &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	}
	_, err := s.client.HeadBucket(params)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) PutFile(ctx context.Context, path string, contentType string, content []byte) (string, error) {
	params := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
		ACL:         aws.String(s3.BucketCannedACLPublicRead),
	}
	_, err := s.client.PutObject(params)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(s3BucketURL, s.bucket, path)
	return url, nil
}
