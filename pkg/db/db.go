package db

import (
	"context"
	"io"
	"mime/multipart"
)

type GetObjectOptions struct {
	Bucket string
	Object string
}

type PutObjectOptions struct {
	Bucket      string
	ObjectName  string
	Data        io.Reader
	Size        int64
	ContentType string
}

type Storage interface {
	Upload(ctx context.Context, options PutObjectOptions) error
	Download(ctx context.Context, options GetObjectOptions) (multipart.File, error)
	Delete(ctx context.Context, bucket, object string) error

	BucketExists(ctx context.Context, bucket string) (bool, error)
	CreateBucket(ctx context.Context, bucket string) error
}
