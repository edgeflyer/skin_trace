package storage

import (
	"fmt"
	"io"
	"os"
)

type Local struct {
	dir string
}

func NewLocal(dir string) (*Local, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create upload directory: %w", err)
	}

	return &Local{dir: dir}, nil
}

// src是照片内容，ext是扩展名
func (s *Local) Save(src io.Reader, ext string) (string, error) {
	file, err := os.CreateTemp(s.dir, "photo-*"+ext)
	if err != nil {
		return "", fmt.Errorf("create photo file: %w", err)
	}

	path := file.Name()

	if _, err := io.Copy(file, src); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return "", fmt.Errorf("write photo file: %w", err)
	}

	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("close photo file: %w", err)
	}

	return path, nil
}

func (s *Local) Delete(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete photo file: %w", err)
	}

	return nil
}

func (s *Local) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read photo file: %w", err)
	}

	return data, nil
}
