package repository

import (
	"context"
	"fmt"
	"io"

	"go-image-compression/internal/consts"
	"go-image-compression/internal/model"

	"github.com/minio/minio-go/v7"
	"github.com/nordew/go-errx"
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
	image, err := r.minio.GetObject(ctx, consts.BucketName, filter.ID, minio.GetObjectOptions{})
	if err != nil {
		return nil, errx.NewInternal().WithDescriptionAndCause(
			fmt.Sprintf("%s: %s", codepath, err.Error()),
			err,
		)
	}

	return image, nil
}

func (r *imageRepository) Create(ctx context.Context, image model.Image) error {
	_, err := r.minio.PutObject(ctx, image.Bucket, image.ID, image.File, image.FileSize,
		minio.PutObjectOptions{ContentType: image.ContentType},
	)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(
			fmt.Sprintf("%s: %s", codepath, err.Error()),
			err,
		)
	}

	return nil
}
