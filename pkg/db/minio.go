package db

import (
	"context"
	"go-image-compression/internal/config"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nordew/go-errx"
)

type minioStorage struct {
	client *minio.Client
}

func NewMinioStorage(cfg config.MinioConfig) (Storage, error) {
	mc, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.SSLMode,
	})
	if err != nil {
		return nil, err
	}

	return &minioStorage{client: mc}, nil
}

func (s *minioStorage) Upload(ctx context.Context, options PutObjectOptions) error {
	_, err := s.client.PutObject(ctx, options.Bucket, options.ObjectName, options.Data, options.Size, minio.PutObjectOptions{
		ContentType: options.ContentType,
	})
	return err
}

func (s *minioStorage) Download(ctx context.Context, options GetObjectOptions) (multipart.File, error) {
	obj, err := s.client.GetObject(ctx, options.Bucket, options.Object, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *minioStorage) Delete(ctx context.Context, bucket, object string) error {
	return s.client.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
}

func (s *minioStorage) BucketExists(ctx context.Context, bucket string) (bool, error) {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *minioStorage) CreateBucket(ctx context.Context, bucket string) error {
	if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		if exists, err := s.client.BucketExists(ctx, bucket); err != nil || !exists {
			return err
		}

		return err
	}
	return nil
}

func GetFileStat(header *multipart.FileHeader) (multipart.File, int64, string, string, error) {
	file, err := header.Open()
	if err != nil {
		return nil, 0, "", "", err
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		return nil, 0, "", "", err
	}

	fileName := header.Filename
	if fileName == "" {
		return nil, 0, "", "", errx.NewInternal().WithDescriptionAndCause("pkg.minio.GetFileStat: ", err)
	}

	return file, header.Size, fileName, contentType, nil
}
