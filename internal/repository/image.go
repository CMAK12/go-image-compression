package repository

import (
	"context"
	"go-image-compression/internal/model"

	"github.com/minio/minio-go/v7"
)

type (
	ImageRepository interface {
		List(ctx context.Context, filter model.ListImageFilter) ([]model.Image, error)
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

func (r *imageRepository) List(ctx context.Context, filter model.ListImageFilter) ([]model.Image, error) {
	return nil, nil
}

func (r *imageRepository) Create(ctx context.Context, image model.Image) error {
	return nil
}
