package repository

import "github.com/minio/minio-go/v7"

type Repositories struct {
	ImageRepository ImageRepository
}

func NewRepository(minio *minio.Client) Repositories {
	return Repositories{
		ImageRepository: newImageRepository(minio),
	}
}
