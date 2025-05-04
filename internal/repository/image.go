package repository

import (
	"context"
	"fmt"
	"io"

	"go-image-compression/internal/consts"
	"go-image-compression/internal/model"

	"github.com/minio/minio-go/v7"
)

type (
	ImageRepository interface {
		Get(ctx context.Context, filter model.ListImageFilter) (io.Reader, error)
		Create(ctx context.Context, image model.Image) error
	}

	imageRepository struct {
		minio *minio.Client
	}
)

func newImageRepository(minio *minio.Client) ImageRepository {
	return &imageRepository{
		minio: minio,
	}
}

const codepath = "repository/image.go"

func (r *imageRepository) Get(ctx context.Context, filter model.ListImageFilter) (io.Reader, error) {
	bucketName := findImageBucketName(filter.CompressPercent)

	image, err := r.minio.GetObject(ctx, bucketName, filter.ID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", codepath, err)
	}

	return image, nil
}

func (r *imageRepository) Create(ctx context.Context, image model.Image) error {
	_, err := r.minio.PutObject(ctx, image.Bucket, image.ID, image.File, image.FileSize, minio.PutObjectOptions{
		ContentType: image.ContentType,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	return nil
}

func findImageBucketName(compressPercent int) string {
	switch compressPercent {
	case 100:
		return consts.FullImageBucket
	case 75:
		return consts.QuarterImageBucket
	case 50:
		return consts.HalfImageBucket
	case 25:
		return consts.QuarterImageBucket
	default:
		return ""
	}
}
