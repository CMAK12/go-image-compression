package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"go-image-compression/internal/model"
	"go-image-compression/pkg/db"

	"github.com/nordew/go-errx"
)

type (
	ImageRepository struct {
		db db.Storage
	}
)

func newImageRepository(db db.Storage) ImageRepository {
	return ImageRepository{
		db: db,
	}
}

const BucketName = "images"

func (r *ImageRepository) Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error) {
	image, err := r.db.Download(ctx, db.GetObjectOptions{
		Bucket: BucketName,
		Object: filter.ID,
	})
	if err != nil {
		return nil, errx.NewInternal().WithDescriptionAndCause(
			fmt.Sprintf("repository.image.Get: %s", err.Error()),
			err,
		)
	}

	return image, nil
}

func (r *ImageRepository) Create(ctx context.Context, img io.Reader, size int64, imageID, contentType string) error {
	err := r.db.Upload(ctx, db.PutObjectOptions{
		Bucket:      BucketName,
		ObjectName:  imageID,
		Data:        img,
		Size:        size,
		ContentType: contentType,
	})
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(
			fmt.Sprintf("repository.image.Create: %s", err.Error()),
			err,
		)
	}

	return nil
}
