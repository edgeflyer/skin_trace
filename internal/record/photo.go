package record

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

var ErrInvalidPhoto = errors.New("photo must be JPEG or PNG")

func detectPhotoExtension(photo io.ReadSeeker) (string, error) {
	_, format, err := image.DecodeConfig(photo)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidPhoto, err)
	}

	// 识别格式时读取了部分内容，保存前要回到文件开头。
	if _, err := photo.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("rewind photo: %w", err)
	}

	switch format {
	case "jpeg":
		return ".jpg", nil
	case "png":
		return ".png", nil
	default:
		return "", ErrInvalidPhoto
	}
}
