package storage

import (
	"io"
	"os"
	"path/filepath"

	"gitlab.com/xerofenix/csd-career/user-service/internal/config"
)

type Storage struct {
	UploadDir string
}

func New(cfg *config.Config) (*Storage, error) {
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		return nil, err
	}
	return &Storage{UploadDir: cfg.UploadDir}, nil
}

func (s *Storage) SaveResume(src io.Reader, filename string) (string, error) {
	filepath := filepath.Join(s.UploadDir, filename)
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filepath, nil
}
