package model

// type Image struct {
// 	ID          string         `json:"id"`
// 	File        multipart.File `json:"file"`
// 	FileSize    int64          `json:"file_size"`
// 	ContentType string         `json:"content_type"`
// }

type ListImageFilter struct {
	ID string `query:"id"`
}

// func NewImage(file multipart.File, fileSize int64, fileName, contentType string) Image {
// 	return Image{
// 		ID:          fmt.Sprintf("%s_%s_100", uuid.NewString(), fileName),
// 		File:        file,
// 		FileSize:    fileSize,
// 		ContentType: contentType,
// 	}
// }

// func NewImageWithID(file multipart.File, fileSize int64, imageID, contentType string) Image {
// 	return Image{
// 		ID:          imageID,
// 		File:        file,
// 		FileSize:    fileSize,
// 		ContentType: contentType,
// 	}
// }
