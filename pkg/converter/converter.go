package converter

import (
	"io"
	"mime/multipart"
)

type (
	Converter interface {
		ConvertToBytes(data multipart.File) ([]byte, error)
	}

	сonverter struct{}
)

func NewConverter() Converter {
	return &сonverter{}
}

func (c *сonverter) ConvertToBytes(data multipart.File) ([]byte, error) {
	bytes, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
