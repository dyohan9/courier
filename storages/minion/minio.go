package minion

import (
	"bytes"
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nyaruka/courier"
	"strings"
)

func init() {
	courier.RegisterStorage("minio", newStorage)
}

type storage struct {
	client *minio.Client
	bucket string
}

func newStorage(config *courier.Config) (courier.Storage, error) {
	minioClient, err := minio.New(
		config.MinioEndpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.MinioAccessKeyID, config.MinioSecretAccessKey, ""),
			Region: config.MinioRegion,
			Secure: config.MinioSecure,
		},
	)
	if err != nil {
		return nil, err
	}
	return &storage{
		client: minioClient,
		bucket: config.MinioMediaBucket,
	}, nil
}

func (s *storage) Test(ctx context.Context, ) error {
	found, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("bucket not found")
	}
	return nil
}

func (s *storage) PutFile(ctx context.Context, path string, contentType string, content []byte) (string, error) {
	path = strings.TrimPrefix(path, "/")
	info, err := s.client.PutObject(
		ctx,
		s.bucket,
		path,
		bytes.NewReader(content),
		-1,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}
	return info.Location, nil
}
