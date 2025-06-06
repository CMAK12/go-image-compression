package resizer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
	"github.com/nordew/go-errx"
)

type (
	Compressor interface {
		Compress(image image.Image, percent float64) (image.Image, error)
		EncodeImage(img image.Image, format string) (*bytes.Reader, int64, error)
		BuildImageID(fileName string, percent float64) string
		GetImage(file io.Reader) (image.Image, string, error)
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
	baseName := fileName[:len(fileName)-4]

	return fmt.Sprintf("%s_%d", baseName, int(percent*100))
}

func (r *Resizer) GetImage(file io.Reader) (image.Image, string, error) {
	decodedImage, format, err := image.Decode(file)
	if err != nil {
		return nil, "", errx.NewInternal().WithDescriptionAndCause("pkg.resizer.GetImage: ", err)
	}

	return decodedImage, format, nil
}

func (r *Resizer) EncodeImage(img image.Image, format string) (*bytes.Reader, int64, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
			return nil, 0, errx.NewInternal().WithDescriptionAndCause("pkg.resizer.EncodeImage: ", err)
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, 0, errx.NewInternal().WithDescriptionAndCause("pkg.resizer.EncodeImage: ", err)
		}
	default:
		return nil, 0, errx.NewInternal().WithDescription("pkg.resizer.EncodeImage: unsupported format")
	}

	reader := bytes.NewReader(buf.Bytes())

	return reader, reader.Size(), nil
}
