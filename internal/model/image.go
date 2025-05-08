package model

import (
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
)

type Image struct {
	ID          string         `json:"id"`
	File        multipart.File `json:"file"`
	FileSize    int64          `json:"file_size"`
	ContentType string         `json:"content_type"`
	Bucket      string         `json:"bucket"`
}

type ListImageFilter struct {
	ID              string `query:"id"`
	CompressPercent int    `query:"compress_percent"`
}

func NewImage(file multipart.File, fileSize int64, bucket, fileName, contentType string) Image {
	return Image{
		ID:          fmt.Sprintf("%s-%s", uuid.NewString(), fileName),
		File:        file,
		FileSize:    fileSize,
		ContentType: contentType,
		Bucket:      bucket,
	}
}
