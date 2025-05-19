package resizer

import (
	"fmt"
	"image"
	"mime/multipart"

	"github.com/nfnt/resize"
)

type (
	Compressor interface {
		Compress(image image.Image, percent float64) (image.Image, error)
		BuildImageID(fileName string, percent float64) string
		GetImage(file multipart.File) (image.Image, error)
	}

	Resizer struct{}
)

func NewResizer() Compressor {
	return &Resizer{}
}

func (r *Resizer) Compress(image image.Image, percent float64) (image.Image, error) {
	resizedWidth := uint(float64(image.Bounds().Size().X) * percent)
	resizedHeight := uint(float64(image.Bounds().Size().Y) * percent)

	resizedImage := resize.Resize(resizedWidth, resizedHeight, image, resize.Lanczos3)

	return resizedImage, nil
}

func (r *Resizer) BuildImageID(fileName string, percent float64) string {
	return fmt.Sprintf("%s_%d", fileName, int(percent*100))
}

func (r *Resizer) GetImage(file multipart.File) (image.Image, error) {
	decodedImage, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return decodedImage, nil
}
