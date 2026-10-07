package deepseek

import (
	"encoding/base64"
	"net/http"
)

func NewImagePart(data []byte) ContentPart {
	mimeType := http.DetectContentType(data)
	encoded := base64.StdEncoding.EncodeToString(data)

	return ContentPart{
		Type: "image_url",
		ImageURL: &ImageURL{
			URL: "data:" + mimeType + ";base64," + encoded,
		},
	}
}
