package repository

import (
	"context"
	"fmt"
	"mime/multipart"

	"go-image-compression/internal/model"
	"go-image-compression/pkg/db"

	"github.com/nordew/go-errx"
)

type (
	ImageRepository interface {
		Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error)
		Create(ctx context.Context, image model.Image) error
	}

	imageRepository struct {
		db db.Storage
	}
)

func newImageRepository(db db.Storage) ImageRepository {
	return &imageRepository{
		db: db,
	}
}

const codepath = "repository/image.go"
const BucketName = "images"

func (r *imageRepository) Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error) {
	image, err := r.db.Download(ctx, db.GetObjectOptions{
		Bucket: BucketName,
		Object: filter.ID,
	})
	if err != nil {
		return nil, errx.NewInternal().WithDescriptionAndCause(
			fmt.Sprintf("%s: %s", codepath, err.Error()),
			err,
		)
	}

	return image, nil
}

func (r *imageRepository) Create(ctx context.Context, image model.Image) error {
	err := r.db.Upload(ctx, db.PutObjectOptions{
		Bucket:      BucketName,
		ObjectName:  image.ID,
		Data:        image.File,
		Size:        image.FileSize,
		ContentType: image.ContentType,
	})
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(
			fmt.Sprintf("%s: %s", codepath, err.Error()),
			err,
		)
	}

	return nil
}
