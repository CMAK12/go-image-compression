package model

type Payload struct {
	ImageID  string `json:"image_id"`
	MIMEType string `json:"mime_type"`
}

func NewPayload(imageID, mimeType string) Payload {
	return Payload{
		ImageID:  imageID,
		MIMEType: mimeType,
	}
}
