package repository

import (
	"go-image-compression/pkg/db"
)

type Repositories struct {
	ImageRepository ImageRepository
}

func NewRepository(db db.Storage) Repositories {
	return Repositories{
		ImageRepository: newImageRepository(db),
	}
}
