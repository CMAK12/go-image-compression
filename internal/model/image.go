package model

type Image struct {
	ID string `json:"id"`
}

type ListImageFilter struct {
	ID              string `json:"id"`
	CompressPercent int    `json:"compress_percent"`
}
